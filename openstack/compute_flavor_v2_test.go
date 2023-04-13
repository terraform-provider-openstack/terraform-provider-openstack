package openstack

import (
	"reflect"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
)

func TestUnitExpandComputeFlavorV2ExtraSpecs(t *testing.T) {
	raw := map[string]interface{}{
		"foo": "bar",
		"bar": "baz",
	}

	expected := flavors.ExtraSpecsOpts{
		"foo": "bar",
		"bar": "baz",
	}

	actual := expandComputeFlavorV2ExtraSpecs(raw)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Results differ. Want: %#v, but got %#v", expected, actual)
	}
}
