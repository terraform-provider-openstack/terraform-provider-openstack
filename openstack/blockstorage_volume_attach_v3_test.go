package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
)

func TestUnitExpandBlockStorageV3AttachMode(t *testing.T) {
	expected := volumes.ReadWrite

	actual, err := expandBlockStorageV3AttachMode("rw")
	assert.Equal(t, err, nil)
	assert.Equal(t, expected, actual)
}

func TestUnitBlockStorageVolumeAttachV3ParseID(t *testing.T) {
	id := "foo/bar"

	expectedVolumeID := "foo"
	expectedAttachmentID := "bar"

	actualVolumeID, actualAttachmentID, err := blockStorageVolumeAttachV3ParseID(id)

	assert.Equal(t, err, nil)
	assert.Equal(t, expectedVolumeID, actualVolumeID)
	assert.Equal(t, expectedAttachmentID, actualAttachmentID)
}
