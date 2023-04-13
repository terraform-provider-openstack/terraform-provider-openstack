package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitResourceNetworkingRouterRouteV2BuildID(t *testing.T) {
	expected := "d190e837-090a-44f2-adcf-07f9fd392931-route-10.11.12.0/24-192.168.0.111"

	actual := resourceNetworkingRouterRouteV2BuildID(
		"d190e837-090a-44f2-adcf-07f9fd392931",
		"10.11.12.0/24",
		"192.168.0.111",
	)

	assert.Equal(t, expected, actual)
}

func TestUnitPesourceNetworkingRouterRouteV2ParseValidID(t *testing.T) {
	routeID := "40412709-86e2-411a-a66f-16053188ed46-route-192.168.0.0/24-10.11.12.13"

	expectedRouterID := "40412709-86e2-411a-a66f-16053188ed46"
	expectedDstCIDR := "192.168.0.0/24"
	expectedNextHop := "10.11.12.13"

	actualRouterID, actualDstCIDR, actualNextHop, err := resourceNetworkingRouterRouteV2ParseID(routeID)

	assert.NoError(t, err)
	assert.Equal(t, expectedRouterID, actualRouterID)
	assert.Equal(t, expectedDstCIDR, actualDstCIDR)
	assert.Equal(t, expectedNextHop, actualNextHop)
}

func TestUnitPesourceNetworkingRouterRouteV2ParseIDInvalidFirstPart(t *testing.T) {
	routeID := "123-router"

	expectedRouterID := ""
	expectedDstCIDR := ""
	expectedNextHop := ""

	actualRouterID, actualDstCIDR, actualNextHop, err := resourceNetworkingRouterRouteV2ParseID(routeID)

	assert.Error(t, err, "invalid ID format: 123-router")
	assert.Equal(t, expectedRouterID, actualRouterID)
	assert.Equal(t, expectedDstCIDR, actualDstCIDR)
	assert.Equal(t, expectedNextHop, actualNextHop)
}

func TestUnitPesourceNetworkingRouterRouteV2ParseIDInvalidLastPart(t *testing.T) {
	routeID := "123-router-bad"

	expectedRouterID := ""
	expectedDstCIDR := ""
	expectedNextHop := ""

	actualRouterID, actualDstCIDR, actualNextHop, err := resourceNetworkingRouterRouteV2ParseID(routeID)

	assert.Error(t, err, "invalid last part format for 123-router-bad: bad")
	assert.Equal(t, expectedRouterID, actualRouterID)
	assert.Equal(t, expectedDstCIDR, actualDstCIDR)
	assert.Equal(t, expectedNextHop, actualNextHop)
}
