package openstack

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
)

func TestNetworkingSubnetV2DeleteConflictFailsImmediately(t *testing.T) {
	const subnetID = "subnet-in-use"
	body := `{"NeutronError":{"type":"SubnetInUse","message":"One or more ports have an IP allocation from this subnet.","detail":""}}`

	var deletes int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete || !strings.HasSuffix(r.URL.Path, "/subnets/"+subnetID) {
			http.NotFound(w, r)
			return
		}
		deletes++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(body))
	}))
	defer server.Close()

	refresh := networkingSubnetV2StateRefreshFuncDelete(context.Background(), subnetTestClient(server.URL), subnetID)
	_, state, err := refresh()
	if err == nil {
		t.Fatal("expected a conflict error, got nil")
	}
	if state != "" {
		t.Fatalf("expected refresh to stop (empty state), got %q", state)
	}
	if deletes != 1 {
		t.Fatalf("expected a single delete attempt, got %d", deletes)
	}
	msg := err.Error()
	for _, want := range []string{"SubnetInUse", "One or more ports have an IP allocation", subnetID} {
		if !strings.Contains(msg, want) {
			t.Errorf("error %q does not contain %q", msg, want)
		}
	}
}

func TestNetworkingSubnetV2DeleteNotFoundIsDeleted(t *testing.T) {
	const subnetID = "subnet-gone"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
