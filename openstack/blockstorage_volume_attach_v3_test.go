package openstack

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitExpandBlockStorageV3AttachMode(t *testing.T) {
	expected := volumes.ReadWrite

	actual, err := expandBlockStorageV3AttachMode("rw")
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}
