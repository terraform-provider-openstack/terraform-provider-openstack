package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitResourceNetworkingSubnetRouteV2BuildID(t *testing.T) {
	expected := "2f6a2fa9-b2ef-4f74-a5ce-4132d1c35455-route-10.13.14.0/24-192.168.0.112"

	actual := resourceNetworkingSubnetRouteV2BuildID(
		"2f6a2fa9-b2ef-4f74-a5ce-4132d1c35455",
		"10.13.14.0/24",
		"192.168.0.112",
	)

	assert.Equal(t, expected, actual)
}

func TestUnitResourceNetworkingSubnetRouteV2ParseValidID(t *testing.T) {
	routeID := "5d621e5d-aa5a-4dd0-ad38-a73d6a17367f-route-192.168.100.0/24-10.11.12.13"

	expectedRouterID := "5d621e5d-aa5a-4dd0-ad38-a73d6a17367f"
	expectedDstCIDR := "192.168.100.0/24"
	expectedNextHop := "10.11.12.13"

	actualRouterID, actualDstCIDR, actualNextHop, err := resourceNetworkingSubnetRouteV2ParseID(routeID)

	assert.NoError(t, err)
	assert.Equal(t, expectedRouterID, actualRouterID)
	assert.Equal(t, expectedDstCIDR, actualDstCIDR)
	assert.Equal(t, expectedNextHop, actualNextHop)
}

func TestUnitResourceNetworkingSubnetRouteV2ParseIDInvalidFirstPart(t *testing.T) {
	routeID := "123-router"

	expectedRouterID := ""
	expectedDstCIDR := ""
	expectedNextHop := ""

	actualRouterID, actualDstCIDR, actualNextHop, err := resourceNetworkingSubnetRouteV2ParseID(routeID)

	assert.Error(t, err, "invalid ID format: 123-router")
	assert.Equal(t, expectedRouterID, actualRouterID)
	assert.Equal(t, expectedDstCIDR, actualDstCIDR)
	assert.Equal(t, expectedNextHop, actualNextHop)
}

func TestUnitResourceNetworkingSubnetRouteV2ParseIDInvalidLastPart(t *testing.T) {
	routeID := "123-router-bad"

	expectedRouterID := ""
	expectedDstCIDR := ""
	expectedNextHop := ""

	actualRouterID, actualDstCIDR, actualNextHop, err := resourceNetworkingSubnetRouteV2ParseID(routeID)

	assert.Error(t, err, "invalid last part format for 123-router-bad: bad")
	assert.Equal(t, expectedRouterID, actualRouterID)
	assert.Equal(t, expectedDstCIDR, actualDstCIDR)
	assert.Equal(t, expectedNextHop, actualNextHop)
}
