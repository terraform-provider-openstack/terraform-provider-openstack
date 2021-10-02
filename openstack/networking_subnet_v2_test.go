package openstack

import (
	"errors"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
)

func TestNetworkingSubnetV2AllocationPools(t *testing.T) {
	allocationPools := []subnets.AllocationPool{
		{
			Start: "192.168.0.2",
			End:   "192.168.0.254",
		},
		{
			Start: "10.0.0.2",
			End:   "10.255.255.254",
		},
	}

	expected := []map[string]interface{}{
		{
			"start": "192.168.0.2",
			"end":   "192.168.0.254",
		},
		{
			"start": "10.0.0.2",
			"end":   "10.255.255.254",
		},
	}

	actual := flattenNetworkingSubnetV2AllocationPools(allocationPools)

	assert.ElementsMatch(t, expected, actual)
}

func TestExpandNetworkingSubnetV2AllocationPools(t *testing.T) {
	r := resourceNetworkingSubnetV2()
	d := r.TestResourceData()
	d.SetId("1")

	allocationPools := []map[string]interface{}{
		{
			"start": "192.168.0.2",
			"end":   "192.168.0.254",
		},
		{
			"start": "10.0.0.2",
			"end":   "10.255.255.254",
		},
	}

	d.Set("allocation_pool", allocationPools)

	expected := []subnets.AllocationPool{
		{
			Start: "192.168.0.2",
			End:   "192.168.0.254",
		},
		{
			Start: "10.0.0.2",
			End:   "10.255.255.254",
		},
	}

	actual := expandNetworkingSubnetV2AllocationPools(d.Get("allocation_pool").(*schema.Set).List())

	assert.ElementsMatch(t, expected, actual)
}

func TestExpandNetworkingSubnetV2HostRoutes(t *testing.T) {
	r := resourceNetworkingSubnetV2()
	d := r.TestResourceData()
	d.SetId("1")

	hostRoutes := []map[string]interface{}{
		{
			"destination_cidr": "192.168.0.0/24",
			"next_hop":         "10.0.0.1",
		},
		{
			"destination_cidr": "10.0.0.0/8",
			"next_hop":         "192.168.0.1",
		},
	}

	d.Set("host_routes", hostRoutes)

	expected := []subnets.HostRoute{
		{
			DestinationCIDR: "192.168.0.0/24",
			NextHop:         "10.0.0.1",
		},
		{
			DestinationCIDR: "10.0.0.0/8",
			NextHop:         "192.168.0.1",
		},
	}

	actual := expandNetworkingSubnetV2HostRoutes(d.Get("host_routes").([]interface{}))

	assert.ElementsMatch(t, expected, actual)
}

func TestNetworkingSubnetV2AllocationPoolsMatch(t *testing.T) {
	oldPools := []interface{}{
		map[string]interface{}{
			"start": "192.168.199.2",
			"end":   "192.168.199.100",
		},

		map[string]interface{}{
			"start": "10.3.0.1",
			"end":   "10.3.0.100",
		},
	}

	newPools := []interface{}{
		map[string]interface{}{
			"start": "192.168.199.2",
			"end":   "192.168.199.100",
		},

		map[string]interface{}{
			"start": "10.3.0.1",
			"end":   "10.3.0.100",
		},
	}

	same := networkingSubnetV2AllocationPoolsMatch(oldPools, newPools)
	assert.Equal(t, same, true)

	oldPools = []interface{}{
		map[string]interface{}{
			"start": "192.168.199.2",
			"end":   "192.168.199.100",
		},
	}

	newPools = []interface{}{
		map[string]interface{}{
			"start": "192.168.199.2",
			"end":   "192.168.199.100",
		},

		map[string]interface{}{
			"start": "10.3.0.1",
			"end":   "10.3.0.100",
		},
	}

	same = networkingSubnetV2AllocationPoolsMatch(oldPools, newPools)
	assert.Equal(t, same, false)

	oldPools = []interface{}{
		map[string]interface{}{
			"start": "192.168.199.2",
			"end":   "192.168.199.100",
		},

		map[string]interface{}{
			"start": "10.3.0.1",
			"end":   "10.3.0.100",
		},
	}

	newPools = []interface{}{
		map[string]interface{}{
			"start": "10.3.0.1",
			"end":   "10.3.0.100",
		},
	}

	same = networkingSubnetV2AllocationPoolsMatch(oldPools, newPools)
	assert.Equal(t, same, false)

	oldPools = []interface{}{
		map[string]interface{}{
			"start": "192.168.199.10",
			"end":   "192.168.199.150",
		},

		map[string]interface{}{
			"start": "10.3.0.1",
			"end":   "10.3.0.100",
		},
	}

	newPools = []interface{}{
		map[string]interface{}{
			"start": "192.168.199.2",
			"end":   "192.168.199.100",
		},

		map[string]interface{}{
			"start": "10.3.0.1",
			"end":   "10.3.0.100",
		},
	}

	same = networkingSubnetV2AllocationPoolsMatch(oldPools, newPools)
	assert.Equal(t, same, false)
}

func TestNetworkingSubnetV2DNSNameserverAreUnique(t *testing.T) {
	tableTest := []struct {
		input []interface{}
		err   error
	}{
		{
			input: []interface{}{"192.168.199.2", "192.168.199.3"},
			err:   nil,
		},
		{
			input: []interface{}{"192.168.199.1", "192.168.199.5", "192.168.199.1"},
			err:   errors.New("got duplicate nameserver 192.168.199.1"),
		},
		{
			input: []interface{}{},
			err:   nil,
		},
	}

	for _, test := range tableTest {
		assert.Equal(t, test.err, networkingSubnetV2DNSNameserverAreUnique(test.input))
	}
}
