package openstack

import (
	"reflect"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumetypes"
)

func TestUnitExpandBlockStorageVolumeTypeV3ExtraSpecs(t *testing.T) {
	raw := map[string]any{
		"foo": "foo",
		"bar": "bar",
	}

	expected := volumetypes.ExtraSpecsOpts{
		"foo": "foo",
		"bar": "bar",
	}

	actual := expandBlockStorageVolumeTypeV3ExtraSpecs(raw)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Results differ. Want: %#v, but got %#v", expected, actual)
	}
}
