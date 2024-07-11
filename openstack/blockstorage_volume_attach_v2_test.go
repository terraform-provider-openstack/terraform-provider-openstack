package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v2/volumes"
)

func TestUnitExpandBlockStorageV2AttachMode(t *testing.T) {
	expected := volumes.ReadWrite

	actual, err := expandBlockStorageV2AttachMode("rw")
	assert.Equal(t, err, nil)
	assert.Equal(t, expected, actual)
}

func TestUnitBlockStorageVolumeAttachV2ParseID(t *testing.T) {
	id := "foo/bar"

	expectedVolumeID := "foo"
	expectedAttachmentID := "bar"

	actualVolumeID, actualAttachmentID, err := blockStorageVolumeAttachV2ParseID(id)

	assert.Equal(t, err, nil)
	assert.Equal(t, expectedVolumeID, actualVolumeID)
	assert.Equal(t, expectedAttachmentID, actualAttachmentID)
}
