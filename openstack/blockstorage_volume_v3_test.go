package openstack

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/schedulerhints"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
)

func blockStorageVolumeV3VolumeFixture() volumes.Volume {
	return volumes.Volume{
		ID:   "289da7f8-6440-407c-9fb4-7db01ec49164",
		Name: "vol-001",
		Attachments: []volumes.Attachment{{
			ServerID:     "83ec2e3b-4321-422b-8706-a84185f52a0a",
			AttachmentID: "05551600-a936-4d4a-ba42-79a037c1-c91a",
			AttachedAt:   time.Date(2016, 8, 6, 14, 48, 20, 0, time.UTC),
			HostName:     "foobar",
			VolumeID:     "d6cacb1a-8b59-4c88-ad90-d70ebb82bb75",
			Device:       "/dev/vdc",
			ID:           "d6cacb1a-8b59-4c88-ad90-d70ebb82bb75",
		}},
		AvailabilityZone:   "nova",
		Bootable:           "false",
		ConsistencyGroupID: "",
		CreatedAt:          time.Date(2015, 9, 17, 3, 35, 3, 0, time.UTC),
		Description:        "",
		Encrypted:          false,
		Metadata:           map[string]string{"foo": "bar"},
		ReplicationStatus:  "disabled",
		Size:               75,
		SnapshotID:         "",
		SourceVolID:        "",
		Status:             "available",
		UserID:             "ff1ce52c03ab433aaba9108c2e3ef541",
		VolumeType:         "lvmdriver-1",
	}
}

func TestUnitFlattenBlockStorageVolumeV3Attachments(t *testing.T) {
	expectedAttachments := []map[string]interface{}{
		{
			"id":          "d6cacb1a-8b59-4c88-ad90-d70ebb82bb75",
			"instance_id": "83ec2e3b-4321-422b-8706-a84185f52a0a",
			"device":      "/dev/vdc",
		},
	}

	actualAttachments := flattenBlockStorageVolumeV3Attachments(blockStorageVolumeV3VolumeFixture().Attachments)
	assert.Equal(t, expectedAttachments, actualAttachments)
}

func TestUnitBlockStorageVolumeV3AttachmentHash(t *testing.T) {
	attachments := flattenBlockStorageVolumeV3Attachments(blockStorageVolumeV3VolumeFixture().Attachments)

	expectedHashcode := 236219624
	actualHashcode := blockStorageVolumeV3AttachmentHash(attachments[0])

	assert.Equal(t, expectedHashcode, actualHashcode)
}

func blockStorageVolumeV3SchedulerHints() schedulerhints.SchedulerHints {
	return schedulerhints.SchedulerHints{
		SameHost:             []string{"83ec2e3b-4321-422b-8706-a84185f52a0a"},
		DifferentHost:        []string{"83ec2e3b-4321-422b-8706-a84185f52a0a"},
		LocalToInstance:      "83ec2e3b-4321-422b-8706-a84185f52a0a",
		Query:                "[“=”, “$backend_id”, “rbd:vol@ceph#cloud”]",
		AdditionalProperties: map[string]interface{}{},
	}
}

func TestUnitFlattenBlockStorageVolumeV3SchedulerHints(t *testing.T) {
	expectedSchedulerHints := map[string]interface{}{
		"same_host":             []interface{}{"83ec2e3b-4321-422b-8706-a84185f52a0a"},
		"different_host":        []interface{}{"83ec2e3b-4321-422b-8706-a84185f52a0a"},
		"local_to_instance":     "83ec2e3b-4321-422b-8706-a84185f52a0a",
		"query":                 "[“=”, “$backend_id”, “rbd:vol@ceph#cloud”]",
		"additional_properties": map[string]interface{}{},
	}

	actualSchedulerHints := expandBlockStorageVolumeV3SchedulerHints(blockStorageVolumeV3SchedulerHints())
	assert.Equal(t, expectedSchedulerHints, actualSchedulerHints)
}

func TestUnitBlockStorageVolumeV3SchedulerHintsHash(t *testing.T) {
	s := expandBlockStorageVolumeV3SchedulerHints(blockStorageVolumeV3SchedulerHints())

	expectedHashcode := 1530836638
	actualHashcode := blockStorageVolumeV3SchedulerHintsHash(s)

	assert.Equal(t, expectedHashcode, actualHashcode)
}
