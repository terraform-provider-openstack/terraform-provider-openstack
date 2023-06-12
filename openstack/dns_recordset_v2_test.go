package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitDNSRecordSetV2ParseID(t *testing.T) {
	id := "foo/bar"
	expectedZoneID := "foo"
	expectedRecordSetID := "bar"

	actualZoneID, actualRecordSetID, err := dnsRecordSetV2ParseID(id)
	assert.Equal(t, err, nil)
	assert.Equal(t, expectedZoneID, actualZoneID)
	assert.Equal(t, expectedRecordSetID, actualRecordSetID)
}
