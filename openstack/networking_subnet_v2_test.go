package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
