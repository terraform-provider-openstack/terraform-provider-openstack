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
