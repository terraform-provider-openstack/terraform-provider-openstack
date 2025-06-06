package openstack

import (
	"reflect"
	"testing"
)

func testResourceObjectStorageV1ContainerStateDataV0() map[string]any {
	return map[string]any{
		"name": "test",
		"versioning": []any{
			map[string]any{
				"type":     "versions",
				"location": "test",
			},
		},
	}
}

func testResourceObjectStorageV1ContainerStateDataV1() map[string]any {
	v0 := testResourceObjectStorageV1ContainerStateDataV0()

	return map[string]any{
		"name":              v0["name"],
		"versioning":        false,
		"versioning_legacy": v0["versioning"],
	}
}

func TestAccObjectStorageV1ContainerStateUpgradeV0(t *testing.T) {
	expected := testResourceObjectStorageV1ContainerStateDataV1()

	actual, err := resourceObjectStorageContainerStateUpgradeV0(t.Context(), testResourceObjectStorageV1ContainerStateDataV0(), nil)
	if err != nil {
		t.Fatalf("error migrating state: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("\n\nexpected:\n\n%#v\n\ngot:\n\n%#v\n\n", expected, actual)
	}
}
