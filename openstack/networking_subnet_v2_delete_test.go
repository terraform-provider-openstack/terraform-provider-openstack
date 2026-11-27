package openstack

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

const subnetInUseBody = `{"NeutronError":{"type":"SubnetInUse","message":"One or more ports have an IP allocation from this subnet.","detail":""}}`

func TestNetworkingSubnetV2DeleteConflictRetriesAndRecordsDetails(t *testing.T) {
	const subnetID = "subnet-in-use"

	var deletes atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		deletes.Add(1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(subnetInUseBody))
	}))

	defer server.Close()

	var lastConflict error

	refresh := networkingSubnetV2StateRefreshFuncDelete(context.Background(), subnetTestClient(server.URL), subnetID, &lastConflict)

	for attempt := 1; attempt <= 3; attempt++ {
		_, state, err := refresh()
		if err != nil {
			t.Fatalf("attempt %d: a 409 must keep retrying, got error %v", attempt, err)
		}

		if state != "ACTIVE" {
			t.Fatalf("attempt %d: state = %q, want ACTIVE", attempt, state)
		}
	}

	if got := deletes.Load(); got != 3 {
		t.Fatalf("expected 3 delete attempts, got %d", got)
	}

	if lastConflict == nil {
		t.Fatal("expected the last conflict to be recorded")
	}

	for _, want := range []string{"SubnetInUse", "One or more ports have an IP allocation"} {
		if !strings.Contains(lastConflict.Error(), want) {
			t.Errorf("conflict %q does not contain %q", lastConflict, want)
		}
	}
}

func TestNetworkingSubnetV2DeleteConflictClearsWhenDependentsGoAway(t *testing.T) {
	const subnetID = "subnet-draining"

	var deletes atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusNotFound)

			return
		}

		// The first attempt conflicts while ports drain, later attempts succeed.
		if deletes.Add(1) == 1 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write([]byte(subnetInUseBody))

			return
		}

		w.WriteHeader(http.StatusNoContent)
	}))

	defer server.Close()

	var lastConflict error

	refresh := networkingSubnetV2StateRefreshFuncDelete(context.Background(), subnetTestClient(server.URL), subnetID, &lastConflict)

	_, state, err := refresh()
	if err != nil || state != "ACTIVE" {
		t.Fatalf("first attempt: state = %q, err = %v, want ACTIVE and no error", state, err)
	}

	if lastConflict == nil {
		t.Fatal("first attempt: expected the conflict to be recorded")
	}

	_, state, err = refresh()
	if err != nil {
		t.Fatalf("second attempt: expected success, got %v", err)
	}

	if state != "DELETED" {
		t.Fatalf("second attempt: state = %q, want DELETED", state)
	}

	if lastConflict != nil {
		t.Fatalf("expected the conflict to be cleared once the delete succeeded, got %v", lastConflict)
	}
}

func TestNetworkingSubnetV2DeleteTimeoutReportsLastConflict(t *testing.T) {
	const subnetID = "subnet-in-use"

	var deletes atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		deletes.Add(1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(subnetInUseBody))
	}))

	defer server.Close()

	var lastConflict error

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    networkingSubnetV2StateRefreshFuncDelete(context.Background(), subnetTestClient(server.URL), subnetID, &lastConflict),
		Timeout:    500 * time.Millisecond,
		MinTimeout: 10 * time.Millisecond,
	}

	_, waitErr := stateConf.WaitForStateContext(context.Background())
	if waitErr == nil {
		t.Fatal("expected the wait to time out while the subnet is still in use")
	}

	if got := deletes.Load(); got < 2 {
		t.Fatalf("expected the delete to be retried until the timeout, got %d attempt(s)", got)
	}

	msg := subnetDeleteWaitError(waitErr, lastConflict).Error()
	for _, want := range []string{"timeout while waiting", "SubnetInUse", "One or more ports have an IP allocation"} {
		if !strings.Contains(msg, want) {
			t.Errorf("error %q does not contain %q", msg, want)
		}
	}
}

func TestNetworkingSubnetV2DeleteWaitErrorWithoutConflictIsUnchanged(t *testing.T) {
	waitErr := errors.New("boom")

	if got := subnetDeleteWaitError(waitErr, nil); !errors.Is(got, waitErr) || got.Error() != "boom" {
		t.Fatalf("expected the wait error to pass through unchanged, got %v", got)
	}
}

func TestNetworkingSubnetV2DeleteConflictWithoutNeutronBody(t *testing.T) {
	const subnetID = "subnet-in-use"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte("conflict"))
	}))

	defer server.Close()

	var lastConflict error

	refresh := networkingSubnetV2StateRefreshFuncDelete(context.Background(), subnetTestClient(server.URL), subnetID, &lastConflict)

	_, state, err := refresh()
	if err != nil || state != "ACTIVE" {
		t.Fatalf("state = %q, err = %v, want ACTIVE and no error", state, err)
	}

	if lastConflict == nil || !strings.Contains(lastConflict.Error(), "cannot be deleted while dependent objects still use it") {
		t.Fatalf("expected a generic conflict error, got %v", lastConflict)
	}
}

func TestNetworkingSubnetV2DeleteNotFoundIsDeleted(t *testing.T) {
	const subnetID = "subnet-gone"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	defer server.Close()

	refresh := networkingSubnetV2StateRefreshFuncDelete(context.Background(), subnetTestClient(server.URL), subnetID, nil)

	_, state, err := refresh()
	if err != nil {
		t.Fatalf("expected already-deleted subnet to succeed, got %v", err)
	}

	if state != "DELETED" {
		t.Fatalf("state = %q, want DELETED", state)
	}
}

func TestNetworkingSubnetV2DeleteAcceptedStillRetries(t *testing.T) {
	const subnetID = "subnet-deleting"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"subnet":{"id":"` + subnetID + `","cidr":"10.0.0.0/24","ip_version":4}}`))
		default:
			http.NotFound(w, r)
		}
	}))

	defer server.Close()

	var lastConflict error

	refresh := networkingSubnetV2StateRefreshFuncDelete(context.Background(), subnetTestClient(server.URL), subnetID, &lastConflict)

	_, state, err := refresh()
	if err != nil {
		t.Fatalf("expected an in-progress delete to keep retrying, got %v", err)
	}

	if state != "ACTIVE" {
		t.Fatalf("state = %q, want ACTIVE", state)
	}

	if lastConflict != nil {
		t.Fatalf("expected no conflict to be recorded, got %v", lastConflict)
	}
}

func subnetTestClient(endpoint string) *gophercloud.ServiceClient {
	return &gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       endpoint + "/",
	}
}
