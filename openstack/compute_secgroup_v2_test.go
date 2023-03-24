package openstack

import (
	"reflect"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/secgroups"
)

func TestUnitExpandComputeSecGroupV2CreateRules(t *testing.T) {
	r := resourceComputeSecGroupV2()
	d := r.TestResourceData()
	d.SetId("1")

	rule1 := map[string]interface{}{
		"from_port":   22,
		"to_port":     22,
		"ip_protocol": "tcp",
		"cidr":        "0.0.0.0/0",
	}

	rule2 := map[string]interface{}{
		"from_port":   1,
		"to_port":     65536,
		"ip_protocol": "udp",
		"cidr":        "0.0.0.0/0",
	}

	rule3 := map[string]interface{}{
		"from_port":   -1,
		"to_port":     -1,
		"ip_protocol": "icmp",
		"cidr":        "0.0.0.0/0",
	}

	rules := []map[string]interface{}{rule1, rule2, rule3}
	d.Set("rule", rules)

	expectedRules := []secgroups.CreateRuleOpts{
		{
			ParentGroupID: "1",
			FromPort:      -1,
			ToPort:        -1,
			IPProtocol:    "icmp",
			CIDR:          "0.0.0.0/0",
		},

		{
			ParentGroupID: "1",
			FromPort:      1,
			ToPort:        65536,
			IPProtocol:    "udp",
			CIDR:          "0.0.0.0/0",
		},

		{
			ParentGroupID: "1",
			FromPort:      22,
			ToPort:        22,
			IPProtocol:    "tcp",
			CIDR:          "0.0.0.0/0",
		},
	}

	actualRules := expandComputeSecGroupV2CreateRules(d)

	if !reflect.DeepEqual(expectedRules, actualRules) {
		t.Fatalf("Rules differ. Want: %#v, but got: %#v", expectedRules, actualRules)
	}
}

func TestUnitExpandComputeSecGroupV2Rule(t *testing.T) {
	r := resourceComputeSecGroupV2()
	d := r.TestResourceData()
	d.SetId("1")

	rule1 := map[string]interface{}{
		"id":          "2",
		"from_port":   22,
		"to_port":     22,
		"ip_protocol": "tcp",
		"cidr":        "0.0.0.0/0",
	}

	expectedRules := secgroups.Rule{
		ParentGroupID: "1",
		ID:            "2",
		FromPort:      22,
		ToPort:        22,
		IPProtocol:    "tcp",
		IPRange:       secgroups.IPRange{CIDR: "0.0.0.0/0"},
	}

	actualRules := expandComputeSecGroupV2Rule(d, rule1)

	if !reflect.DeepEqual(expectedRules, actualRules) {
		t.Fatalf("Results differ. Want: %#v, but got: %#v", expectedRules, actualRules)
	}
}
