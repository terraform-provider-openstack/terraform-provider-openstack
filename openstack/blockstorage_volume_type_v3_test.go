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

func TestUnitExpandBlockStorageVolumeTypeV3Encryption(t *testing.T) {
	raw := []any{
		map[string]any{
			"provider":         "luks",
			"cipher":           "aes-xts-plain64",
			"key_size":         256,
			"control_location": "front-end",
		},
	}

	expected := volumetypes.CreateEncryptionOpts{
		Provider:        "luks",
		Cipher:          "aes-xts-plain64",
		KeySize:         256,
		ControlLocation: "front-end",
	}

	actual := expandBlockStorageVolumeTypeV3Encryption(raw)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Results differ. Want: %#v, but got %#v", expected, actual)
	}
}

func TestUnitFlattenBlockStorageVolumeTypeV3Encryption(t *testing.T) {
	enc := &volumetypes.GetEncryptionType{
		EncryptionID:    "81e069c6-7394-4856-8df7-3b237ca61f74",
		Provider:        "luks",
		Cipher:          "aes-xts-plain64",
		KeySize:         256,
		ControlLocation: "front-end",
	}

	expected := []map[string]any{
		{
			"provider":         "luks",
			"cipher":           "aes-xts-plain64",
			"key_size":         256,
			"control_location": "front-end",
		},
	}

	actual := flattenBlockStorageVolumeTypeV3Encryption(enc)

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Results differ. Want: %#v, but got %#v", expected, actual)
	}

	if actual := flattenBlockStorageVolumeTypeV3Encryption(&volumetypes.GetEncryptionType{}); actual != nil {
		t.Fatalf("Want nil for a type without encryption, but got %#v", actual)
	}
}
