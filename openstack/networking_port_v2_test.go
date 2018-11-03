package openstack

import (
	"reflect"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/extradhcpopts"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/hashicorp/terraform/helper/schema"
)

func TestExpandNetworkingPortDHCPOptsV2Create(t *testing.T) {
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
		F: networkingPortDHCPOptsV2HashSetFunc(),
	}
	dhcpOptsSet.Add(dhcpOptsA)
	dhcpOptsSet.Add(dhcpOptsB)

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

	actualDHCPOptions := expandNetworkingPortDHCPOptsV2Create(dhcpOptsSet)

	if !reflect.DeepEqual(expectedDHCPOptions, actualDHCPOptions) {
		t.Fatalf("DHCP options differs, want: %+v, but got: %+v",
			expectedDHCPOptions, actualDHCPOptions)
	}
}

func TestExpandNetworkingPortDHCPOptsEmptyV2Create(t *testing.T) {
	dhcpOptsSet := &schema.Set{
		F: networkingPortDHCPOptsV2HashSetFunc(),
	}

	expectedDHCPOptions := []extradhcpopts.CreateExtraDHCPOpt{}

	actualDHCPOptions := expandNetworkingPortDHCPOptsV2Create(dhcpOptsSet)

	if !reflect.DeepEqual(expectedDHCPOptions, actualDHCPOptions) {
		t.Fatalf("DHCP options differs, want: %+v, but got: %+v",
			expectedDHCPOptions, actualDHCPOptions)
	}
}

func TestExpandNetworkingPortDHCPOptsV2Update(t *testing.T) {
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
		F: networkingPortDHCPOptsV2HashSetFunc(),
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

	actualDHCPOptions := expandNetworkingPortDHCPOptsV2Update(dhcpOptsSet)

	if !reflect.DeepEqual(expectedDHCPOptions, actualDHCPOptions) {
		t.Fatalf("DHCP options differs, want: %+v, but got: %+v",
			expectedDHCPOptions, actualDHCPOptions)
	}
}

func TestExpandNetworkingPortDHCPOptsEmptyV2Update(t *testing.T) {
	dhcpOptsSet := &schema.Set{
		F: networkingPortDHCPOptsV2HashSetFunc(),
	}

	expectedDHCPOptions := []extradhcpopts.UpdateExtraDHCPOpt{}

	actualDHCPOptions := expandNetworkingPortDHCPOptsV2Update(dhcpOptsSet)

	if !reflect.DeepEqual(expectedDHCPOptions, actualDHCPOptions) {
		t.Fatalf("DHCP options differs, want: %+v, but got: %+v",
			expectedDHCPOptions, actualDHCPOptions)
	}
}

func TestExpandNetworkingPortDHCPOptsV2Delete(t *testing.T) {
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
		F: networkingPortDHCPOptsV2HashSetFunc(),
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

	actualDHCPOptions := expandNetworkingPortDHCPOptsV2Delete(dhcpOptsSet)

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
		F: networkingPortDHCPOptsV2HashSetFunc(),
	}
	expectedDHCPOptionsSet.Add(expectedDHCPOptsA)
	expectedDHCPOptionsSet.Add(expectedDHCPOptsB)

	actualDHCPOptionsSet := flattenNetworkingPortDHCPOptsV2(dhcpOptions)

	if !actualDHCPOptionsSet.Equal(expectedDHCPOptionsSet) {
		t.Fatalf("DHCP options set differs, want: %+v, but got: %+v",
			expectedDHCPOptionsSet, actualDHCPOptionsSet)
	}
}

func TestEnsureEmptyNetworkingPortV2UpdateOpts(t *testing.T) {
	var updateOpts ports.UpdateOpts

	expected := ports.UpdateOpts{}

	actual := ensureNetworkingPortV2UpdateOpts(&updateOpts)

	if actual != expected {
		t.Fatalf("expected empty ports.UpdateOpts{}, but got %v", actual)
	}
}

func TestEnsurePopulatedNetworkingPortV2UpdateOpts(t *testing.T) {
	adminStateUp := true

	expected := ports.UpdateOpts{
		Name:         "eth1",
		AdminStateUp: &adminStateUp,
	}

	actual := ensureNetworkingPortV2UpdateOpts(&expected)

	if actual != expected {
		t.Fatalf("expected empty ports.UpdateOpts{}, but got %v", actual)
	}
}
