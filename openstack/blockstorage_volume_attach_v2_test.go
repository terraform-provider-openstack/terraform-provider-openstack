package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/volumeactions"
)

func TestExpandBlockStorageV2AttachMode(t *testing.T) {
	expected := volumeactions.ReadWrite

	actual, err := expandBlockStorageV2AttachMode("rw")
	assert.Equal(t, err, nil)
	assert.Equal(t, expected, actual)
}

func TestBlockStorageVolumeAttachV2ParseID(t *testing.T) {
	id := "foo/bar"

	expectedVolumeID := "foo"
	expectedAttachmentID := "bar"

	actualVolumeID, actualAttachmentID, err := blockStorageVolumeAttachV2ParseID(id)

	assert.Equal(t, err, nil)
	assert.Equal(t, expectedVolumeID, actualVolumeID)
	assert.Equal(t, expectedAttachmentID, actualAttachmentID)
}
