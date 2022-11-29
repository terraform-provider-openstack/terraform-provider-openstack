package openstack

import (
	"reflect"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/rules"
)

func TestExpandNetworkingSecgroupV2CreateRules(t *testing.T) {
	r := resourceNetworkingSecGroupV2()
	d := r.TestResourceData()
	d.SetId("1")

	rule1 := map[string]interface{}{
		"port_range_min":   22,
		"port_range_max":   22,
		"protocol":         "tcp",
		"remote_ip_prefix": "0.0.0.0/0",
	}

	rule2 := map[string]interface{}{
		"port_range_min":   1,
		"port_range_max":   65536,
		"protocol":         "udp",
		"remote_ip_prefix": "0.0.0.0/0",
	}

	rule3 := map[string]interface{}{
		"port_range_min":   -1,
		"port_range_max":   -1,
		"protocol":         "icmp",
		"remote_ip_prefix": "0.0.0.0/0",
	}

	allRules := []map[string]interface{}{rule1, rule2, rule3}
	d.Set("rule", allRules)

	expectedRules := []rules.CreateOpts{
		{
			SecGroupID:     "1",
			PortRangeMin:   -1,
			PortRangeMax:   -1,
			Protocol:       "icmp",
			RemoteIPPrefix: "0.0.0.0/0",
		},

		{
			SecGroupID:     "1",
			PortRangeMin:   1,
			PortRangeMax:   65536,
			Protocol:       "udp",
			RemoteIPPrefix: "0.0.0.0/0",
		},

		{
			SecGroupID:     "1",
			PortRangeMin:   22,
			PortRangeMax:   22,
			Protocol:       "tcp",
			RemoteIPPrefix: "0.0.0.0/0",
		},
	}

	actualRules, err := expandNetworkingSecgroupV2CreateRules(d)

	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(expectedRules, actualRules) {
		t.Fatalf("Rules differ. Want: %#v, but got: %#v", expectedRules, actualRules)
	}
}

func TestExpandNetworkingSecgroupV2Rule(t *testing.T) {
	r := resourceNetworkingSecGroupV2()
	d := r.TestResourceData()
	d.SetId("1")

	rule1 := map[string]interface{}{
		"id":               "2",
		"port_range_min":   22,
		"port_range_max":   22,
		"protocol":         "tcp",
		"remote_ip_prefix": "0.0.0.0/0",
	}

	expectedRules := rules.SecGroupRule{
		SecGroupID:     "1",
		ID:             "2",
		PortRangeMin:   22,
		PortRangeMax:   22,
		Protocol:       "tcp",
		RemoteIPPrefix: "0.0.0.0/0",
	}

	actualRules := expandNetworkingSecgroupV2Rule(d, rule1)

	if !reflect.DeepEqual(expectedRules, actualRules) {
		t.Fatalf("Results differ. Want: %#v, but got: %#v", expectedRules, actualRules)
	}
}
