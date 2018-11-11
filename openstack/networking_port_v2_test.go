package openstack

import (
	"reflect"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/extradhcpopts"
	"github.com/hashicorp/terraform/helper/schema"
)

func TestExpandNetworkingPortDHCPOptsV2Create(t *testing.T) {
	r := resourceNetworkingPortV2()
	d := r.TestResourceData()
	d.SetId("1")
	dhcpOpts1 := map[string]interface{}{
		"ip_version": 4,
		"opt_name":   "A",
		"opt_value":  "true",
	}
	dhcpOpts2 := map[string]interface{}{
		"ip_version": 6,
		"opt_name":   "B",
		"opt_value":  "false",
	}
	extraDHCPOpts := []map[string]interface{}{dhcpOpts1, dhcpOpts2}
	d.Set("extra_dhcp_opts", extraDHCPOpts)

	expectedDHCPOptions := []extradhcpopts.CreateExtraDHCPOpt{
		{
			OptName:   "A",
			OptValue:  "true",
			IPVersion: gophercloud.IPVersion(4),
		},
		{
			OptName:   "B",
			OptValue:  "false",
			IPVersion: gophercloud.IPVersion(6),
		},
	}

	actualDHCPOptions := expandNetworkingPortDHCPOptsV2Create(d.Get("extra_dhcp_opts").(*schema.Set))

	if !reflect.DeepEqual(expectedDHCPOptions, actualDHCPOptions) {
		t.Fatalf("DHCP options differs, want: %+v, but got: %+v",
			expectedDHCPOptions, actualDHCPOptions)
	}
}

func TestExpandNetworkingPortDHCPOptsEmptyV2Create(t *testing.T) {
	r := resourceNetworkingPortV2()
	d := r.TestResourceData()
	d.SetId("1")

	expectedDHCPOptions := []extradhcpopts.CreateExtraDHCPOpt{}

	actualDHCPOptions := expandNetworkingPortDHCPOptsV2Create(d.Get("extra_dhcp_opts").(*schema.Set))

	if !reflect.DeepEqual(expectedDHCPOptions, actualDHCPOptions) {
		t.Fatalf("DHCP options differs, want: %+v, but got: %+v",
			expectedDHCPOptions, actualDHCPOptions)
	}
}

func TestExpandNetworkingPortDHCPOptsV2Update(t *testing.T) {
	r := resourceNetworkingPortV2()
	d := r.TestResourceData()
	d.SetId("1")
	dhcpOpts1 := map[string]interface{}{
		"ip_version": 4,
		"opt_name":   "A",
		"opt_value":  "true",
	}
	dhcpOpts2 := map[string]interface{}{
		"ip_version": 6,
		"opt_name":   "B",
		"opt_value":  "false",
	}
	extraDHCPOpts := []map[string]interface{}{dhcpOpts1, dhcpOpts2}
	d.Set("extra_dhcp_opts", extraDHCPOpts)

	optsValueTrue := "true"
	optsValueFalse := "false"
	expectedDHCPOptions := []extradhcpopts.UpdateExtraDHCPOpt{
		{
			OptName:   "A",
			OptValue:  &optsValueTrue,
			IPVersion: gophercloud.IPVersion(4),
		},
		{
			OptName:   "B",
			OptValue:  &optsValueFalse,
			IPVersion: gophercloud.IPVersion(6),
		},
	}

	actualDHCPOptions := expandNetworkingPortDHCPOptsV2Update(d.Get("extra_dhcp_opts").(*schema.Set))

	if !reflect.DeepEqual(expectedDHCPOptions, actualDHCPOptions) {
		t.Fatalf("DHCP options differs, want: %+v, but got: %+v",
			expectedDHCPOptions, actualDHCPOptions)
	}
}

func TestExpandNetworkingPortDHCPOptsEmptyV2Update(t *testing.T) {
	r := resourceNetworkingPortV2()
	d := r.TestResourceData()
	d.SetId("1")

	expectedDHCPOptions := []extradhcpopts.UpdateExtraDHCPOpt{}

	actualDHCPOptions := expandNetworkingPortDHCPOptsV2Update(d.Get("extra_dhcp_opts").(*schema.Set))

	if !reflect.DeepEqual(expectedDHCPOptions, actualDHCPOptions) {
		t.Fatalf("DHCP options differs, want: %+v, but got: %+v",
			expectedDHCPOptions, actualDHCPOptions)
	}
}

func TestExpandNetworkingPortDHCPOptsV2Delete(t *testing.T) {
	r := resourceNetworkingPortV2()
	d := r.TestResourceData()
	d.SetId("1")
	dhcpOpts1 := map[string]interface{}{
		"ip_version": 4,
		"opt_name":   "A",
		"opt_value":  "true",
	}
	dhcpOpts2 := map[string]interface{}{
		"ip_version": 6,
		"opt_name":   "B",
		"opt_value":  "false",
	}
	extraDHCPOpts := []map[string]interface{}{dhcpOpts1, dhcpOpts2}
	d.Set("extra_dhcp_opts", extraDHCPOpts)

	expectedDHCPOptions := []extradhcpopts.UpdateExtraDHCPOpt{
		{
			OptName: "A",
		},
		{
			OptName: "B",
		},
	}

	actualDHCPOptions := expandNetworkingPortDHCPOptsV2Delete(d.Get("extra_dhcp_opts").(*schema.Set))

	if !reflect.DeepEqual(expectedDHCPOptions, actualDHCPOptions) {
		t.Fatalf("DHCP options differs, want: %+v, but got: %+v",
			expectedDHCPOptions, actualDHCPOptions)
	}
}

func TestFlattenNetworkingPort2DHCPOptionsV2(t *testing.T) {
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

	expectedDHCPOptions := []map[string]interface{}{
		map[string]interface{}{
			"ip_version": 4,
			"opt_name":   "A",
			"opt_value":  "true",
		},
		map[string]interface{}{
			"ip_version": 6,
			"opt_name":   "B",
			"opt_value":  "false",
		},
	}

	actualDHCPOptions := flattenNetworkingPortDHCPOptsV2(dhcpOptions)

	if !reflect.DeepEqual(actualDHCPOptions, expectedDHCPOptions) {
		t.Fatalf("DHCP options set differs, want: %+v, but got: %+v",
			expectedDHCPOptions, actualDHCPOptions)
	}
}
