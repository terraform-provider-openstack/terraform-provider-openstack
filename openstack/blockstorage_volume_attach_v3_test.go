package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/volumeactions"
)

func TestExpandBlockStorageV3AttachMode(t *testing.T) {
	expected := volumeactions.ReadWrite

	actual, err := expandBlockStorageV3AttachMode("rw")
	assert.Equal(t, err, nil)
	assert.Equal(t, expected, actual)
}

func TestBlockStorageVolumeAttachV3ParseID(t *testing.T) {
	id := "foo/bar"

	expectedVolumeID := "foo"
	expectedAttachmentID := "bar"

	actualVolumeID, actualAttachmentID, err := blockStorageVolumeAttachV3ParseID(id)

	assert.Equal(t, err, nil)
	assert.Equal(t, expectedVolumeID, actualVolumeID)
	assert.Equal(t, expectedAttachmentID, actualAttachmentID)
}
