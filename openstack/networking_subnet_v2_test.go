package openstack

import (
	"errors"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/subnets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestUnitNetworkingSubnetV2AllocationPools(t *testing.T) {
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

	expected := []map[string]any{
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

func TestUnitExpandNetworkingSubnetV2AllocationPools(t *testing.T) {
	r := resourceNetworkingSubnetV2()
	d := r.TestResourceData()
	d.SetId("1")

	allocationPools := []map[string]any{
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

func TestUnitNetworkingSubnetV2AllocationPoolsMatch(t *testing.T) {
	oldPools := []any{
		map[string]any{
			"start": "192.168.199.2",
			"end":   "192.168.199.100",
		},

		map[string]any{
			"start": "10.3.0.1",
			"end":   "10.3.0.100",
		},
	}

	newPools := []any{
		map[string]any{
			"start": "192.168.199.2",
			"end":   "192.168.199.100",
		},

		map[string]any{
			"start": "10.3.0.1",
			"end":   "10.3.0.100",
		},
	}

	same := networkingSubnetV2AllocationPoolsMatch(oldPools, newPools)
	assert.True(t, same)

	oldPools = []any{
		map[string]any{
			"start": "192.168.199.2",
			"end":   "192.168.199.100",
		},
	}

	newPools = []any{
		map[string]any{
			"start": "192.168.199.2",
			"end":   "192.168.199.100",
		},

		map[string]any{
			"start": "10.3.0.1",
			"end":   "10.3.0.100",
		},
	}

	same = networkingSubnetV2AllocationPoolsMatch(oldPools, newPools)
	assert.False(t, same)

	oldPools = []any{
		map[string]any{
			"start": "192.168.199.2",
			"end":   "192.168.199.100",
		},

		map[string]any{
			"start": "10.3.0.1",
			"end":   "10.3.0.100",
		},
	}

	newPools = []any{
		map[string]any{
			"start": "10.3.0.1",
			"end":   "10.3.0.100",
		},
	}

	same = networkingSubnetV2AllocationPoolsMatch(oldPools, newPools)
	assert.False(t, same)

	oldPools = []any{
		map[string]any{
			"start": "192.168.199.10",
			"end":   "192.168.199.150",
		},

		map[string]any{
			"start": "10.3.0.1",
			"end":   "10.3.0.100",
		},
	}

	newPools = []any{
		map[string]any{
			"start": "192.168.199.2",
			"end":   "192.168.199.100",
		},

		map[string]any{
			"start": "10.3.0.1",
			"end":   "10.3.0.100",
		},
	}

	same = networkingSubnetV2AllocationPoolsMatch(oldPools, newPools)
	assert.False(t, same)
}

func TestUnitNetworkingSubnetV2DNSNameserverAreUnique(t *testing.T) {
	tableTest := []struct {
		input []any
		err   error
	}{
		{
			input: []any{"192.168.199.2", "192.168.199.3"},
			err:   nil,
		},
		{
			input: []any{"192.168.199.1", "192.168.199.5", "192.168.199.1"},
			err:   errors.New("got duplicate nameserver 192.168.199.1"),
		},
		{
			input: []any{},
			err:   nil,
		},
	}

	for _, test := range tableTest {
		assert.Equal(t, test.err, networkingSubnetV2DNSNameserverAreUnique(test.input))
	}
}
