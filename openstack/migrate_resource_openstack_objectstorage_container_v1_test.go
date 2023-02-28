package openstack

import (
	"context"
	"reflect"
	"testing"
)

func testResourceObjectStorageV1ContainerStateDataV0() map[string]interface{} {
	return map[string]interface{}{
		"name": "test",
		"versioning": []interface{}{
			map[string]interface{}{
				"type":     "versions",
				"location": "test",
			},
		},
	}
}

func testResourceObjectStorageV1ContainerStateDataV1() map[string]interface{} {
	v0 := testResourceObjectStorageV1ContainerStateDataV0()
	return map[string]interface{}{
		"name":              v0["name"],
		"versioning":        false,
		"versioning_legacy": v0["versioning"],
	}
}

func TestAccObjectStorageV1ContainerStateUpgradeV0(t *testing.T) {
	expected := testResourceObjectStorageV1ContainerStateDataV1()
	actual, err := resourceObjectStorageContainerStateUpgradeV0(context.Background(), testResourceObjectStorageV1ContainerStateDataV0(), nil)
	if err != nil {
		t.Fatalf("error migrating state: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("\n\nexpected:\n\n%#v\n\ngot:\n\n%#v\n\n", expected, actual)
	}
}
