package openstack

import (
	"reflect"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/servergroups"
	th "github.com/gophercloud/gophercloud/testhelper"
	thclient "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestComputeServerGroupV2CreateOpts(t *testing.T) {
	createOpts := ComputeServerGroupV2CreateOpts{
		servergroups.CreateOpts{
			Name:     "foo",
			Policies: []string{"affinity"},
		},
		map[string]string{
			"foo": "bar",
		},
	}

	expected := map[string]interface{}{
		"server_group": map[string]interface{}{
			"name":     "foo",
			"policies": []interface{}{"affinity"},
			"foo":      "bar",
		},
	}

	actual, err := createOpts.ToServerGroupCreateMap()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Maps differ. Want: %#v, but got: %#v", expected, actual)
	}
}

func TestExpandComputeServerGroupV2Policies(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	raw := []interface{}{
		"affinity",
	}

	expected := []string{
		"affinity",
	}

	client := thclient.ServiceClient()
	actual := expandComputeServerGroupV2Policies(client, raw)

	if client.Microversion != "" {
		t.Fatalf("Expected no microversion in client, but got %s", client.Microversion)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Results differ. Want: #%v, but got %#v", expected, actual)
	}
}

func TestExpandComputeServerGroupV2PoliciesMicroversions(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	raw := []interface{}{
		"affinity",
		"soft-anti-affinity",
		"soft-affinity",
	}

	expected := []string{
		"affinity",
		"soft-anti-affinity",
		"soft-affinity",
	}

	client := thclient.ServiceClient()
	actual := expandComputeServerGroupV2Policies(client, raw)

	if client.Microversion != "2.15" {
		t.Fatalf("Expected 2.15 microversion in client, but got %s", client.Microversion)
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Results differ. Want: #%v, but got %#v", expected, actual)
	}
}
