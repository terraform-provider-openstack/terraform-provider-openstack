package openstack

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
)

const subnetInUseBody = `{"NeutronError":{"type":"SubnetInUse","message":"One or more ports have an IP allocation from this subnet.","detail":""}}`

func TestNetworkingSubnetV2DeleteConflictFailsAfterGracePeriod(t *testing.T) {
	const subnetID = "subnet-in-use"

	var deletes atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		deletes.Add(1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(subnetInUseBody))
	}))

	defer server.Close()

	refresh := networkingSubnetV2StateRefreshFuncDeleteWithGrace(context.Background(), subnetTestClient(server.URL), subnetID, 0)

	_, state, err := refresh()
	if err == nil {
		t.Fatal("expected a conflict error, got nil")
	}

	if state != "" {
		t.Fatalf("expected refresh to stop (empty state), got %q", state)
	}

	if got := deletes.Load(); got != 1 {
		t.Fatalf("expected a single delete attempt, got %d", got)
	}

	msg := err.Error()
	for _, want := range []string{"SubnetInUse", "One or more ports have an IP allocation", subnetID} {
		if !strings.Contains(msg, want) {
			t.Errorf("error %q does not contain %q", msg, want)
		}
	}
}

func TestNetworkingSubnetV2DeleteConflictRetriesWithinGracePeriod(t *testing.T) {
	const subnetID = "subnet-in-use"

	var deletes atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		deletes.Add(1)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(subnetInUseBody))
	}))

	defer server.Close()

	refresh := networkingSubnetV2StateRefreshFuncDeleteWithGrace(context.Background(), subnetTestClient(server.URL), subnetID, time.Hour)

	for attempt := 1; attempt <= 3; attempt++ {
		_, state, err := refresh()
		if err != nil {
			t.Fatalf("attempt %d: expected a transient conflict to keep retrying, got %v", attempt, err)
		}

		if state != "ACTIVE" {
			t.Fatalf("attempt %d: state = %q, want ACTIVE", attempt, state)
		}
	}

	if got := deletes.Load(); got != 3 {
		t.Fatalf("expected 3 delete attempts, got %d", got)
	}
}

func TestNetworkingSubnetV2DeleteConflictClearsWhenDependentsGoAway(t *testing.T) {
	const subnetID = "subnet-draining"

	var deletes atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			// First attempt conflicts while ports drain, later attempts succeed.
			if deletes.Add(1) == 1 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusConflict)
				_, _ = w.Write([]byte(subnetInUseBody))

				return
			}

			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))

	defer server.Close()

	refresh := networkingSubnetV2StateRefreshFuncDeleteWithGrace(context.Background(), subnetTestClient(server.URL), subnetID, time.Hour)

	_, state, err := refresh()
	if err != nil || state != "ACTIVE" {
		t.Fatalf("first attempt: state = %q, err = %v, want ACTIVE and no error", state, err)
	}

	_, state, err = refresh()
	if err != nil {
		t.Fatalf("second attempt: expected success, got %v", err)
	}

	if state != "DELETED" {
		t.Fatalf("second attempt: state = %q, want DELETED", state)
	}
}

func TestNetworkingSubnetV2DeleteNotFoundIsDeleted(t *testing.T) {
	const subnetID = "subnet-gone"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	defer server.Close()

	refresh := networkingSubnetV2StateRefreshFuncDelete(context.Background(), subnetTestClient(server.URL), subnetID)

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

	refresh := networkingSubnetV2StateRefreshFuncDelete(context.Background(), subnetTestClient(server.URL), subnetID)

	_, state, err := refresh()
	if err != nil {
		t.Fatalf("expected an in-progress delete to keep retrying, got %v", err)
	}

	if state != "ACTIVE" {
		t.Fatalf("state = %q, want ACTIVE", state)
	}
}

func subnetTestClient(endpoint string) *gophercloud.ServiceClient {
	return &gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       endpoint + "/",
	}
}
