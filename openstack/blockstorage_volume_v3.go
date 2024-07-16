package openstack

import (
	"bytes"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/schedulerhints"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/utils/terraform/hashcode"
)

const blockstorageV3VolumeFromBackupMicroversion = "3.47"
const blockstorageV3ResizeOnlineInUse = "3.42"

func flattenBlockStorageVolumeV3Attachments(v []volumes.Attachment) []map[string]interface{} {
	attachments := make([]map[string]interface{}, len(v))
	for i, attachment := range v {
		attachments[i] = make(map[string]interface{})
		attachments[i]["id"] = attachment.ID
		attachments[i]["instance_id"] = attachment.ServerID
		attachments[i]["device"] = attachment.Device
	}

	return attachments
}

func blockStorageVolumeV3StateRefreshFunc(client *gophercloud.ServiceClient, volumeID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		v, err := volumes.Get(client, volumeID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return v, "deleted", nil
			}

			return nil, "", err
		}

		if v.Status == "error" {
			return v, v.Status, fmt.Errorf("The volume is in error status. " +
				"Please check with your cloud admin or check the Block Storage " +
				"API logs to see why this error occurred.")
		}

		return v, v.Status, nil
	}
}

func blockStorageVolumeV3AttachmentHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if m["instance_id"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["instance_id"].(string)))
	}
	return hashcode.String(buf.String())
}

func expandBlockStorageVolumeV3SchedulerHints(v schedulerhints.SchedulerHints) map[string]interface{} {
	schedulerHints := make(map[string]interface{})

	differentHost := make([]interface{}, len(v.DifferentHost))
	for i, dh := range v.DifferentHost {
		differentHost[i] = dh
	}

	sameHost := make([]interface{}, len(v.SameHost))
	for i, sh := range v.SameHost {
		sameHost[i] = sh
	}

	schedulerHints["different_host"] = differentHost
	schedulerHints["same_host"] = sameHost
	schedulerHints["local_to_instance"] = v.LocalToInstance
	schedulerHints["query"] = v.Query
	schedulerHints["additional_properties"] = v.AdditionalProperties
	return schedulerHints
}

func blockStorageVolumeV3SchedulerHintsHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if m["query"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["query"].(string)))
	}

	if m["local_to_instance"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["local_to_instance"].(string)))
	}

	if m["additional_properties"] != nil {
		for _, v := range m["additional_properties"].(map[string]interface{}) {
			buf.WriteString(fmt.Sprintf("%s-", v))
		}
	}

	buf.WriteString(fmt.Sprintf("%s-", m["different_host"].([]interface{})))
	buf.WriteString(fmt.Sprintf("%s-", m["same_host"].([]interface{})))

	return hashcode.String(buf.String())
}

func resourceBlockStorageVolumeV3SchedulerHints(schedulerHintsRaw map[string]interface{}) schedulerhints.SchedulerHints {
	schedulerHints := schedulerhints.SchedulerHints{
		Query:                schedulerHintsRaw["query"].(string),
		LocalToInstance:      schedulerHintsRaw["local_to_instance"].(string),
		AdditionalProperties: schedulerHintsRaw["additional_properties"].(map[string]interface{}),
	}

	if v, ok := schedulerHintsRaw["different_host"].([]interface{}); ok {
		differentHost := make([]string, len(v))

		for i, dh := range v {
			differentHost[i] = dh.(string)
		}

		schedulerHints.DifferentHost = differentHost
	}

	if v, ok := schedulerHintsRaw["same_host"].([]interface{}); ok {
		sameHost := make([]string, len(v))

		for i, sh := range v {
			sameHost[i] = sh.(string)
		}

		schedulerHints.SameHost = sameHost
	}

	return schedulerHints
}
