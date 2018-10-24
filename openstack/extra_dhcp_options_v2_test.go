package openstack

import (
	"reflect"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/extradhcpopts"
	"github.com/hashicorp/terraform/helper/schema"
)

func TestExpandDHCPOptionsV2Add(t *testing.T) {
	dhcpOptsA := map[string]interface{}{
		"ip_version": 4,
		"opt_name":   "A",
		"opt_value":  "true",
	}
	dhcpOptsB := map[string]interface{}{
		"ip_version": 6,
		"opt_name":   "B",
		"opt_value":  "false",
	}
	dhcpOptsSet := &schema.Set{
		F: dhcpOptionsV2HashSetFunc(),
	}
	dhcpOptsSet.Add(dhcpOptsA)
	dhcpOptsSet.Add(dhcpOptsB)

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

	actualDHCPOptions := expandDHCPOptionsV2Add(dhcpOptsSet)

	if !reflect.DeepEqual(expectedDHCPOptions, actualDHCPOptions) {
		t.Fatalf("DHCP options differs, want: %+v, but got: %+v",
			expectedDHCPOptions, actualDHCPOptions)
	}
}

func TestExpandDHCPOptionsV2Delete(t *testing.T) {
	dhcpOptsA := map[string]interface{}{
		"ip_version": 4,
		"opt_name":   "A",
		"opt_value":  "true",
	}
	dhcpOptsB := map[string]interface{}{
		"ip_version": 6,
		"opt_name":   "B",
		"opt_value":  "false",
	}
	dhcpOptsSet := &schema.Set{
		F: dhcpOptionsV2HashSetFunc(),
	}
	dhcpOptsSet.Add(dhcpOptsA)
	dhcpOptsSet.Add(dhcpOptsB)

	expectedDHCPOptions := []extradhcpopts.UpdateExtraDHCPOpt{
		{
			OptName: "A",
		},
		{
			OptName: "B",
		},
	}

	actualDHCPOptions := expandDHCPOptionsV2Delete(dhcpOptsSet)

	if !reflect.DeepEqual(expectedDHCPOptions, actualDHCPOptions) {
		t.Fatalf("DHCP options differs, want: %+v, but got: %+v",
			expectedDHCPOptions, actualDHCPOptions)
	}
}

func TestFlattenDHCPOptionsV2(t *testing.T) {
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

	expectedDHCPOptsA := map[string]interface{}{
		"ip_version": 4,
		"opt_name":   "A",
		"opt_value":  "true",
	}
	expectedDHCPOptsB := map[string]interface{}{
		"ip_version": 6,
		"opt_name":   "B",
		"opt_value":  "false",
	}
	expectedDHCPOptionsSet := &schema.Set{
		F: dhcpOptionsV2HashSetFunc(),
	}
	expectedDHCPOptionsSet.Add(expectedDHCPOptsA)
	expectedDHCPOptionsSet.Add(expectedDHCPOptsB)

	actualDHCPOptionsSet := flattenDHCPOptionsV2(dhcpOptions)

	if !actualDHCPOptionsSet.Equal(expectedDHCPOptionsSet) {
		t.Fatalf("DHCP options set differs, want: %+v, but got: %+v",
			expectedDHCPOptionsSet, actualDHCPOptionsSet)
	}
}
