package openstack

import (
	"reflect"
	"testing"
)

func TestUnitBlockStorageVolumeTypeQuotaConversion(t *testing.T) {
	raw := map[string]any{
		"foo": "42",
		"bar": "43",
	}

	expected := map[string]any{
		"foo": 42,
		"bar": 43,
	}

	actual, err := blockStorageQuotasetVolTypeQuotaToInt(raw)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("Results differ. Want: %#v, but got %#v", expected, actual)
	}
}

func TestUnitBlockStorageVolumeTypeQuotaConversion_err(t *testing.T) {
	raw := map[string]any{
		"foo": 100,
		"bar": 200,
	}

	_, err := blockStorageQuotasetVolTypeQuotaToInt(raw)

	if err == nil {
		t.Fatal("Expected error in asserting string")
	}
}

func TestUnitBlockStorageVolumeTypeQuotaConversion_err2(t *testing.T) {
	raw := map[string]any{
		"foo": "foo",
		"bar": "bar",
	}

	_, err := blockStorageQuotasetVolTypeQuotaToInt(raw)

	if err == nil {
		t.Fatal("Expected error in converting to int")
	}
}
