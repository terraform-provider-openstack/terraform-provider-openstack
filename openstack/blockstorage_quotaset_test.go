package openstack

import (
	"reflect"
	"testing"
)

func TestBlockStorageVolumeTypeQuotaConversion(t *testing.T) {
	raw := map[string]interface{}{
		"foo": "42",
		"bar": "43",
	}

	expected := map[string]interface{}{
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

func TestBlockStorageVolumeTypeQuotaConversion_err(t *testing.T) {
	raw := map[string]interface{}{
		"foo": 100,
		"bar": 200,
	}

	_, err := blockStorageQuotasetVolTypeQuotaToInt(raw)

	if err == nil {
		t.Fatal("Expected error in asserting string")
	}
}

func TestBlockStorageVolumeTypeQuotaConversion_err2(t *testing.T) {
	raw := map[string]interface{}{
		"foo": "foo",
		"bar": "bar",
	}

	_, err := blockStorageQuotasetVolTypeQuotaToInt(raw)

	if err == nil {
		t.Fatal("Expected error in converting to int")
	}
}
