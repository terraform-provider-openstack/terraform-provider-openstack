package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/schedulerhints"
)

func blockStorageExtensionsSchedulerHints() schedulerhints.SchedulerHints {
	return schedulerhints.SchedulerHints{
		SameHost:             []string{"83ec2e3b-4321-422b-8706-a84185f52a0a"},
		DifferentHost:        []string{"83ec2e3b-4321-422b-8706-a84185f52a0a"},
		LocalToInstance:      "83ec2e3b-4321-422b-8706-a84185f52a0a",
		Query:                "[“=”, “$backend_id”, “rbd:vol@ceph#cloud”]",
		AdditionalProperties: map[string]interface{}{},
	}
}

func TestFlattenBlockStorageExtensionsSchedulerHints(t *testing.T) {
	expectedSchedulerHints := map[string]interface{}{
		"same_host":             []interface{}{"83ec2e3b-4321-422b-8706-a84185f52a0a"},
		"different_host":        []interface{}{"83ec2e3b-4321-422b-8706-a84185f52a0a"},
		"local_to_instance":     "83ec2e3b-4321-422b-8706-a84185f52a0a",
		"query":                 "[“=”, “$backend_id”, “rbd:vol@ceph#cloud”]",
		"additional_properties": map[string]interface{}{},
	}

	actualSchedulerHints := expandBlockStorageExtensionsSchedulerHints(blockStorageExtensionsSchedulerHints())
	assert.Equal(t, expectedSchedulerHints, actualSchedulerHints)
}

func TestBlockStorageExtensionsSchedulerHintsHash(t *testing.T) {
	s := expandBlockStorageExtensionsSchedulerHints(blockStorageExtensionsSchedulerHints())

	expectedHashcode := 1530836638
	actualHashcode := blockStorageExtensionsSchedulerHintsHash(s)

	assert.Equal(t, expectedHashcode, actualHashcode)
}
