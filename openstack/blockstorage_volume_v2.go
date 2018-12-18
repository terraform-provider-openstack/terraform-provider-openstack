package openstack

import (
	"bytes"
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v2/volumes"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
)

func flattenBlockStorageVolumeV2Attachments(v []volumes.Attachment) []map[string]interface{} {
	attachments := make([]map[string]interface{}, len(v))
	for i, attachment := range v {
		attachments[i] = make(map[string]interface{})
		attachments[i]["id"] = attachment.ID
		attachments[i]["instance_id"] = attachment.ServerID
		attachments[i]["device"] = attachment.Device
	}

	return attachments
}

func expandBlockStorageVolumeV2Metadata(v map[string]interface{}) map[string]string {
	metadata := make(map[string]string)
	for key, val := range v {
		metadata[key] = val.(string)
	}

	return metadata
}

func blockStorageVolumeV2StateRefreshFunc(client *gophercloud.ServiceClient, volumeID string) resource.StateRefreshFunc {
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

func blockStorageVolumeV2AttachmentHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	if m["instance_id"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["instance_id"].(string)))
	}
	return hashcode.String(buf.String())
}
