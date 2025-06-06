package openstack

import (
	"reflect"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/keypairs"
)

func TestUnitComputeKeyPairV2CreateOpts(t *testing.T) {
	createOpts := ComputeKeyPairV2CreateOpts{
		keypairs.CreateOpts{
			Name: "kp_1",
		},
		map[string]string{
			"foo": "bar",
		},
	}

	expected := map[string]any{
		"keypair": map[string]any{
			"name": "kp_1",
			"foo":  "bar",
		},
	}

	actual, err := createOpts.ToKeyPairCreateMap()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Maps differ. Want: %#v, but got: %#v", expected, actual)
	}
}
