package openstack

import (
	"reflect"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/servergroups"
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
	raw := []interface{}{
		"affinity",
	}

	expected := []string{
		"affinity",
	}

	actual := expandComputeServerGroupV2Policies(raw)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Results differ. Want: #%v, but got %#v", expected, actual)
	}
}
