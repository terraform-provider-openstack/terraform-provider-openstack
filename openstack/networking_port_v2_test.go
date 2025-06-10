package openstack

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/extradhcpopts"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestUnitExpandNetworkingPortDHCPOptsV2Create(t *testing.T) {
	r := resourceNetworkingPortV2()
	d := r.TestResourceData()
	d.SetId("1")

	dhcpOpts1 := map[string]any{
		"ip_version": 4,
		"name":       "A",
		"value":      "true",
	}
	dhcpOpts2 := map[string]any{
		"ip_version": 6,
		"name":       "B",
		"value":      "false",
	}
	extraDHCPOpts := []map[string]any{dhcpOpts1, dhcpOpts2}
	d.Set("extra_dhcp_option", extraDHCPOpts)

	expectedDHCPOptions := []extradhcpopts.CreateExtraDHCPOpt{
		{
			OptName:   "B",
			OptValue:  "false",
			IPVersion: gophercloud.IPVersion(6),
		},
		{
			OptName:   "A",
			OptValue:  "true",
			IPVersion: gophercloud.IPVersion(4),
		},
	}

	actualDHCPOptions := expandNetworkingPortDHCPOptsV2Create(d.Get("extra_dhcp_option").(*schema.Set))

	assert.ElementsMatch(t, expectedDHCPOptions, actualDHCPOptions)
}

func TestUnitExpandNetworkingPortDHCPOptsEmptyV2Create(t *testing.T) {
	r := resourceNetworkingPortV2()
	d := r.TestResourceData()
	d.SetId("1")

	expectedDHCPOptions := []extradhcpopts.CreateExtraDHCPOpt{}

	actualDHCPOptions := expandNetworkingPortDHCPOptsV2Create(d.Get("extra_dhcp_option").(*schema.Set))

	assert.ElementsMatch(t, expectedDHCPOptions, actualDHCPOptions)
}

func TestUnitExpandNetworkingPortDHCPOptsV2Update(t *testing.T) {
	r := resourceNetworkingPortV2()
	d := r.TestResourceData()
	d.SetId("1")

	dhcpOpts1 := map[string]any{
		"ip_version": 4,
		"name":       "A",
		"value":      "true",
	}
	dhcpOpts2 := map[string]any{
		"ip_version": 6,
		"name":       "B",
		"value":      "false",
	}
	extraDHCPOpts := []map[string]any{dhcpOpts1, dhcpOpts2}
	d.Set("extra_dhcp_option", extraDHCPOpts)

	optsValueTrue := "true"
	optsValueFalse := "false"
	expectedDHCPOptions := []extradhcpopts.UpdateExtraDHCPOpt{
		{
			OptName:   "B",
			OptValue:  &optsValueFalse,
			IPVersion: gophercloud.IPVersion(6),
		},
		{
			OptName:   "A",
			OptValue:  &optsValueTrue,
			IPVersion: gophercloud.IPVersion(4),
		},
	}

	actualDHCPOptions := expandNetworkingPortDHCPOptsV2Update(nil, d.Get("extra_dhcp_option").(*schema.Set))

	assert.ElementsMatch(t, expectedDHCPOptions, actualDHCPOptions)
}

func TestUnitExpandNetworkingPortDHCPOptsEmptyV2Update(t *testing.T) {
	r := resourceNetworkingPortV2()
	d := r.TestResourceData()
	d.SetId("1")

	var expectedDHCPOptions []extradhcpopts.UpdateExtraDHCPOpt

	actualDHCPOptions := expandNetworkingPortDHCPOptsV2Update(nil, d.Get("extra_dhcp_option").(*schema.Set))

	assert.ElementsMatch(t, expectedDHCPOptions, actualDHCPOptions)
}

func TestUnitExpandNetworkingPortDHCPOptsV2Delete(t *testing.T) {
	r := resourceNetworkingPortV2()
	d := r.TestResourceData()
	d.SetId("1")

	dhcpOpts1 := map[string]any{
		"ip_version": 4,
		"name":       "A",
		"value":      "true",
	}
	dhcpOpts2 := map[string]any{
		"ip_version": 6,
		"name":       "B",
		"value":      "false",
	}
	extraDHCPOpts := []map[string]any{dhcpOpts1, dhcpOpts2}
	d.Set("extra_dhcp_option", extraDHCPOpts)

	expectedDHCPOptions := []extradhcpopts.UpdateExtraDHCPOpt{
		{
			OptName: "B",
		},
		{
			OptName: "A",
		},
	}

	actualDHCPOptions := expandNetworkingPortDHCPOptsV2Update(d.Get("extra_dhcp_option").(*schema.Set), nil)

	assert.ElementsMatch(t, expectedDHCPOptions, actualDHCPOptions)
}

func TestUnitFlattenNetworkingPort2DHCPOptionsV2(t *testing.T) {
	dhcpOptions := extradhcpopts.ExtraDHCPOptsExt{
		ExtraDHCPOpts: []extradhcpopts.ExtraDHCPOpt{
			{
				OptName:   "A",
				OptValue:  "true",
				IPVersion: 4,
			},
			{
				OptName:   "B",
				OptValue:  "false",
				IPVersion: 6,
			},
		},
	}

	expectedDHCPOptions := []map[string]any{
		{
			"ip_version": 4,
			"name":       "A",
			"value":      "true",
		},
		{
			"ip_version": 6,
			"name":       "B",
			"value":      "false",
		},
	}

	actualDHCPOptions := flattenNetworkingPortDHCPOptsV2(dhcpOptions)

	assert.ElementsMatch(t, expectedDHCPOptions, actualDHCPOptions)
}

func TestUnitExpandNetworkingPortAllowedAddressPairsV2(t *testing.T) {
	r := resourceNetworkingPortV2()
	d := r.TestResourceData()
	d.SetId("1")

	addressPairs1 := map[string]any{
		"ip_address":  "192.0.2.1",
		"mac_address": "mac1",
	}
	addressPairs2 := map[string]any{
		"ip_address":  "198.51.100.1",
		"mac_address": "mac2",
	}
	allowedAddressPairs := []map[string]any{addressPairs1, addressPairs2}
	d.Set("allowed_address_pairs", allowedAddressPairs)

	expectedAllowedAddressPairs := []ports.AddressPair{
		{
			IPAddress:  "192.0.2.1",
			MACAddress: "mac1",
		},
		{
			IPAddress:  "198.51.100.1",
			MACAddress: "mac2",
		},
	}

	actualAllowedAddressPairs := expandNetworkingPortAllowedAddressPairsV2(d.Get("allowed_address_pairs").(*schema.Set))

	assert.ElementsMatch(t, expectedAllowedAddressPairs, actualAllowedAddressPairs)
}

func TestUnitFlattenNetworkingPortAllowedAddressPairsV2(t *testing.T) {
	allowedAddressPairs := []ports.AddressPair{
		{
			IPAddress:  "192.0.2.1",
			MACAddress: "mac1",
		},
		{
			IPAddress:  "198.51.100.1",
			MACAddress: "mac2",
		},
	}
	mac := "mac3"

	expectedAllowedAddressPairs := []map[string]any{
		{
			"ip_address":  "192.0.2.1",
			"mac_address": "mac1",
		},
		{
			"ip_address":  "198.51.100.1",
			"mac_address": "mac2",
		},
	}

	actualAllowedAddressPairs := flattenNetworkingPortAllowedAddressPairsV2(mac, allowedAddressPairs)

	assert.ElementsMatch(t, expectedAllowedAddressPairs, actualAllowedAddressPairs)
}

func TestUnitExpandNetworkingPortFixedIPV2NoFixedIPs(t *testing.T) {
	r := resourceNetworkingPortV2()
	d := r.TestResourceData()
	d.SetId("1")
	d.Set("no_fixed_ip", true)

	actualFixedIP := expandNetworkingPortFixedIPV2(d)

	assert.Empty(t, actualFixedIP)
}

func TestUnitExpandNetworkingPortFixedIPV2SomeFixedIPs(t *testing.T) {
	r := resourceNetworkingPortV2()
	d := r.TestResourceData()
	d.SetId("1")

	fixedIP1 := map[string]any{
		"subnet_id":  "aaa",
		"ip_address": "192.0.201.101",
	}
	fixedIP2 := map[string]any{
		"subnet_id":  "bbb",
		"ip_address": "192.0.202.102",
	}
	fixedIP := []map[string]any{fixedIP1, fixedIP2}
	d.Set("fixed_ip", fixedIP)

	expectedFixedIP := []ports.IP{
		{
			SubnetID:  "aaa",
			IPAddress: "192.0.201.101",
		},
		{
			SubnetID:  "bbb",
			IPAddress: "192.0.202.102",
		},
	}

	actualFixedIP := expandNetworkingPortFixedIPV2(d)

	assert.ElementsMatch(t, expectedFixedIP, actualFixedIP)
}
