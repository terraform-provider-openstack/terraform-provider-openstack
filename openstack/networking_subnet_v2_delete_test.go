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
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}
		deletes++
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(body))
	}))
	defer server.Close()

	client := &gophercloud.ServiceClient{
		ProviderClient: &gophercloud.ProviderClient{},
		Endpoint:       server.URL + "/",
	}

	refresh := networkingSubnetV2StateRefreshFuncDelete(context.Background(), client, subnetID)
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
	for _, want := range []string{"SubnetInUse", "dependent objects", subnetID} {
		if !strings.Contains(msg, want) {
			t.Errorf("error %q does not contain %q", msg, want)
		}
	}
}
