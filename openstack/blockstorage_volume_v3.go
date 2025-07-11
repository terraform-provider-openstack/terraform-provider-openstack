package openstack

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/terraform-provider-openstack/utils/v2/hashcode"
)

const (
	blockstorageV3VolumeFromBackupMicroversion = "3.47"
	blockstorageV3ResizeOnlineInUse            = "3.42"
)

func flattenBlockStorageVolumeV3Attachments(v []volumes.Attachment) []map[string]any {
	attachments := make([]map[string]any, len(v))
	for i, attachment := range v {
		attachments[i] = make(map[string]any)
		attachments[i]["id"] = attachment.ID
		attachments[i]["instance_id"] = attachment.ServerID
		attachments[i]["device"] = attachment.Device
	}

	return attachments
}

func blockStorageVolumeV3StateRefreshFunc(ctx context.Context, client *gophercloud.ServiceClient, volumeID string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		v, err := volumes.Get(ctx, client, volumeID).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return v, "deleted", nil
			}

			return nil, "", err
		}

		if v.Status == "error" {
			return v, v.Status, errors.New("The volume is in error status. " +
				"Please check with your cloud admin or check the Block Storage " +
				"API logs to see why this error occurred.")
		}

		return v, v.Status, nil
	}
}

func blockStorageVolumeV3AttachmentHash(v any) int {
	var buf bytes.Buffer

	m := v.(map[string]any)
	if m["instance_id"] != nil {
		buf.WriteString(m["instance_id"].(string) + "-")
	}

	return hashcode.String(buf.String())
}

func expandBlockStorageVolumeV3SchedulerHints(v volumes.SchedulerHintOpts) map[string]any {
	schedulerHints := make(map[string]any)

	differentHost := make([]any, len(v.DifferentHost))
	for i, dh := range v.DifferentHost {
		differentHost[i] = dh
	}

	sameHost := make([]any, len(v.SameHost))
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

func blockStorageVolumeV3SchedulerHintsHash(v any) int {
	var buf bytes.Buffer

	m := v.(map[string]any)

	if m["query"] != nil {
		buf.WriteString(m["query"].(string) + "-")
	}

	if m["local_to_instance"] != nil {
		buf.WriteString(m["local_to_instance"].(string) + "-")
	}

	if m["additional_properties"] != nil {
		for _, v := range m["additional_properties"].(map[string]any) {
			buf.WriteString(fmt.Sprintf("%s-", v))
		}
	}

	buf.WriteString(fmt.Sprintf("%s-", m["different_host"].([]any)))
	buf.WriteString(fmt.Sprintf("%s-", m["same_host"].([]any)))

	return hashcode.String(buf.String())
}

func resourceBlockStorageVolumeV3SchedulerHints(schedulerHintsRaw map[string]any) volumes.SchedulerHintOpts {
	schedulerHints := volumes.SchedulerHintOpts{
		Query:                schedulerHintsRaw["query"].(string),
		LocalToInstance:      schedulerHintsRaw["local_to_instance"].(string),
		AdditionalProperties: schedulerHintsRaw["additional_properties"].(map[string]any),
	}

	if v, ok := schedulerHintsRaw["different_host"].([]any); ok {
		differentHost := make([]string, len(v))

		for i, dh := range v {
			differentHost[i] = dh.(string)
		}

		schedulerHints.DifferentHost = differentHost
	}

	if v, ok := schedulerHintsRaw["same_host"].([]any); ok {
		sameHost := make([]string, len(v))

		for i, sh := range v {
			sameHost[i] = sh.(string)
		}

		schedulerHints.SameHost = sameHost
	}

	return schedulerHints
}
