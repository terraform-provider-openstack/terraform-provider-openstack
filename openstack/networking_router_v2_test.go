package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
)

func TestExpandNetworkingRouterExternalFixedIPsV2(t *testing.T) {
	r := resourceNetworkingRouterV2()
	d := r.TestResourceData()
	d.SetId("1")
	fixedIPs1 := map[string]string{
		"subnet_id":  "subnet_1",
		"ip_address": "192.168.101.1",
	}
	fixedIPs2 := map[string]string{
		"subnet_id":  "subnet_2",
		"ip_address": "192.168.201.1",
	}
	externalFixedIPs := []map[string]string{fixedIPs1, fixedIPs2}
	d.Set("external_fixed_ip", externalFixedIPs)

	expectedExternalFixedIPs := []routers.ExternalFixedIP{
		{
			SubnetID:  "subnet_1",
			IPAddress: "192.168.101.1",
		},
		{
			SubnetID:  "subnet_2",
			IPAddress: "192.168.201.1",
		},
	}

	actualExternalFixedIPs := expandNetworkingRouterExternalFixedIPsV2(d.Get("external_fixed_ip").([]interface{}))

	assert.ElementsMatch(t, expectedExternalFixedIPs, actualExternalFixedIPs)
}

func TestFlattenNetworkingRouterExternalFixedIPsV2(t *testing.T) {
	externalFixedIPs := []routers.ExternalFixedIP{
		{
			SubnetID:  "subnet_1",
			IPAddress: "192.168.101.1",
		},
		{
			SubnetID:  "subnet_2",
			IPAddress: "192.168.201.1",
		},
	}

	expectedExternalFixedIPs := []map[string]string{
		{
			"subnet_id":  "subnet_1",
			"ip_address": "192.168.101.1",
		},
		{
			"subnet_id":  "subnet_2",
			"ip_address": "192.168.201.1",
		},
	}

	actualExternalFixedIPs := flattenNetworkingRouterExternalFixedIPsV2(externalFixedIPs)

	assert.ElementsMatch(t, expectedExternalFixedIPs, actualExternalFixedIPs)
}
