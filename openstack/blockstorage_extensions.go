package openstack

import (
	"bytes"
	"fmt"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/schedulerhints"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
)

func flattenBlockStorageExtensionsSchedulerHints(v schedulerhints.SchedulerHints) map[string]interface{} {
	schedulerHints := make(map[string]interface{})
	schedulerHints["same_host"] = v.SameHost
	schedulerHints["different_host"] = v.DifferentHost
	schedulerHints["local_to_instance"] = v.LocalToInstance
	schedulerHints["query"] = v.Query
	schedulerHints["additional_properties"] = v.AdditionalProperties
	return schedulerHints
}

func blockStorageExtensionsSchedulerHintsHash(v interface{}) int {
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

	buf.WriteString(fmt.Sprintf("%s-", m["different_host"].([]string)))
	buf.WriteString(fmt.Sprintf("%s-", m["same_host"].([]string)))

	return hashcode.String(buf.String())
}
