package openstack

import (
	"reflect"
	"testing"
)

func testResourceNetworkingV2PortStateDataV0() map[string]any {
	return map[string]any{
		"name": "test",
		"all_fixed_ips": []any{
			"192.168.0.10",
			"192.168.0.11",
		},
	}
}

func testResourceNetworkingV2PortStateDataV1() map[string]any {
	return map[string]any{
		"name": "test",
		"all_fixed_ips": []map[string]any{
			{
				"ip_address": "192.168.0.10",
				"subnet_id":  "",
			},
			{
				"ip_address": "192.168.0.11",
				"subnet_id":  "",
			},
		},
	}
}

func TestAccNetworkingV2PortStateUpgradeV0(t *testing.T) {
	expected := testResourceNetworkingV2PortStateDataV1()

	actual, err := upgradeNetworkingPortV2StateV0toV1(t.Context(), testResourceNetworkingV2PortStateDataV0(), nil)
	if err != nil {
		t.Fatalf("error migrating state: %s", err)
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("\n\nexpected:\n\n%#v\n\ngot:\n\n%#v\n\n", expected, actual)
	}
}
