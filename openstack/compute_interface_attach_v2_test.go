package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputeInterfaceAttachV2ParseID(t *testing.T) {
	id := "foo/bar"

	expectedInstanceID := "foo"
	expectedAttachmentID := "bar"

	actualInstanceID, actualAttachmentID, err := computeInterfaceAttachV2ParseID(id)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectedInstanceID, actualInstanceID)
	assert.Equal(t, expectedAttachmentID, actualAttachmentID)
}
