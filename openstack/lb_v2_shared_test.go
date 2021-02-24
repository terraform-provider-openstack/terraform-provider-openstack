package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandLBV2ListenerHeadersMap(t *testing.T) {
	raw := map[string]interface{}{
		"header0": "val0",
		"header1": "val1",
	}

	expected := map[string]string{
		"header0": "val0",
		"header1": "val1",
	}

	actual, err := expandLBV2ListenerHeadersMap(raw)

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestExpandLBV2ListenerHeadersMap_err(t *testing.T) {
	raw := map[string]interface{}{
		"header0": "val0",
		"header1": 1,
	}

	actual, err := expandLBV2ListenerHeadersMap(raw)

	assert.Error(t, err)
	assert.Empty(t, actual)
}

func TestParseLBQuotaID(t *testing.T) {
	lbQuotaID := "ed498e81f0cc448bae0ad4f8f21bf67f/region/1"
	expectedProjectID := "ed498e81f0cc448bae0ad4f8f21bf67f"
	expectedRegion := "region/1"

	actualProjectID, actualRegion, err := parseLBQuotaID(lbQuotaID)

	assert.NoError(t, err)
	assert.Equal(t, expectedProjectID, actualProjectID)
	assert.Equal(t, expectedRegion, actualRegion)
}

func TestParseLBQuotaID_err(t *testing.T) {
	wrongLBQuotaID := "foo42"

	actualProjectID, actualRegion, err := parseLBQuotaID(wrongLBQuotaID)

	assert.Error(t, err)
	assert.Empty(t, actualProjectID)
	assert.Empty(t, actualRegion)
}
