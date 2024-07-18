package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/volumeactions"
)

func TestUnitExpandBlockStorageV2AttachMode(t *testing.T) {
	expected := volumeactions.ReadWrite

	actual, err := expandBlockStorageV2AttachMode("rw")
	assert.Equal(t, err, nil)
	assert.Equal(t, expected, actual)
}
