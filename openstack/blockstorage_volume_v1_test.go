package openstack

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v1/volumes"
)

func blockStorageVolumeV1VolumeFixture() *volumes.Volume {
	return &volumes.Volume{
		Status: "active",
		Name:   "vol-001",
		Attachments: []map[string]interface{}{
			{
				"attachment_id": "03987cd1-0ad5-40d1-9b2a-7cc48295d4fa",
				"id":            "47e9ecc5-4045-4ee3-9a4b-d859d546a0cf",
				"volume_id":     "6c80f8ac-e3e2-480c-8e6e-f1db92fe4bfe",
				"server_id":     "d1c4788b-9435-42e2-9b81-29f3be1cd01f",
				"host_name":     "mitaka",
				"device":        "/",
			},
		},
		AvailabilityZone: "us-east1",
		Bootable:         "false",
		CreatedAt:        time.Date(2012, 2, 14, 20, 53, 07, 0, time.UTC),
		Description:      "Another volume.",
		VolumeType:       "289da7f8-6440-407c-9fb4-7db01ec49164",
		SnapshotID:       "",
		SourceVolID:      "",
		Metadata: map[string]string{
			"contents": "junk",
		},
		ID:   "521752a6-acf6-4b2d-bc7a-119f9148cd8c",
		Size: 30,
	}
}

func TestFlattenBlockStorageVolumeV1Attachments(t *testing.T) {
	expectedAttachments := []map[string]interface{}{
		{
			"id":          "47e9ecc5-4045-4ee3-9a4b-d859d546a0cf",
			"instance_id": "d1c4788b-9435-42e2-9b81-29f3be1cd01f",
			"device":      "/",
		},
	}

	actualAttachments := flattenBlockStorageVolumeV1Attachments(blockStorageVolumeV1VolumeFixture().Attachments)
	assert.Equal(t, expectedAttachments, actualAttachments)
}

func TestBlockStorageVolumeV1AttachmentHash(t *testing.T) {
	attachments := flattenBlockStorageVolumeV1Attachments(blockStorageVolumeV1VolumeFixture().Attachments)

	expectedHashcode := 258823884
	actualHashcode := blockStorageVolumeV1AttachmentHash(attachments[0])

	assert.Equal(t, expectedHashcode, actualHashcode)
}
