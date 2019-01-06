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

func TestExpandDNSRecordSetV2Records(t *testing.T) {
	data := []interface{}{"foo", "[bar]", "baz"}
	expected := []string{"foo", "bar", "baz"}

	actual := expandDNSRecordSetV2Records(data)
	assert.Equal(t, expected, actual)
}
