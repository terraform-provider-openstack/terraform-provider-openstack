package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDNSRecordSetV2ParseID(t *testing.T) {
	id := "foo/bar"
	expectedZoneID := "foo"
	expectedRecordSetID := "bar"

	actualZoneID, actualRecordSetID, err := dnsRecordSetV2ParseID(id)
	assert.Equal(t, err, nil)
	assert.Equal(t, expectedZoneID, actualZoneID)
	assert.Equal(t, expectedRecordSetID, actualRecordSetID)
}

func TestDNSRecordSetV2RecordsStateFunc(t *testing.T) {
	data := []interface{}{"foo", "[bar]", "baz"}
	expected := []string{"foo", "bar", "baz"}

	for i, record := range data {
		actual := dnsRecordSetV2RecordsStateFunc(record)
		assert.Equal(t, expected[i], actual)
	}
}
