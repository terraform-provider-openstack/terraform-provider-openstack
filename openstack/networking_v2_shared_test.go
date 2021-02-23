package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseNetworkingQuotaID(t *testing.T) {
	networkingQuotaID := "557f2600ca3b40a0b8c8621a1e7ba559/region/1"
	expectedProjectID := "557f2600ca3b40a0b8c8621a1e7ba559"
	expectedRegion := "region/1"

	actualProjectID, actualRegion, err := parseNetworkingQuotaID(networkingQuotaID)

	assert.NoError(t, err)
	assert.Equal(t, expectedProjectID, actualProjectID)
	assert.Equal(t, expectedRegion, actualRegion)
}

func TestParseNetworkingQuotaID_err(t *testing.T) {
	wrongNetworkingQuotaID := "foo42"

	actualProjectID, actualRegion, err := parseNetworkingQuotaID(wrongNetworkingQuotaID)

	assert.Error(t, err)
	assert.Empty(t, actualProjectID)
	assert.Empty(t, actualRegion)
}
