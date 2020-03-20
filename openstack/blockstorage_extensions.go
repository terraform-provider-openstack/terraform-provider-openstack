package openstack

import (
	"bytes"
	"fmt"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/schedulerhints"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
)

func flattenBlockStorageExtensionsSchedulerHints(v schedulerhints.SchedulerHints) map[string]interface{} {
	schedulerHints := make(map[string]interface{})

	var differentHost []interface{}
	for _, dh := range v.DifferentHost {
		differentHost = append(differentHost, dh)
	}

	var sameHost []interface{}
	for _, sh := range v.SameHost {
		sameHost = append(sameHost, sh)
	}

	schedulerHints["different_host"] = differentHost
	schedulerHints["same_host"] = sameHost
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

	buf.WriteString(fmt.Sprintf("%s-", m["different_host"].([]interface{})))
	buf.WriteString(fmt.Sprintf("%s-", m["same_host"].([]interface{})))

	return hashcode.String(buf.String())
}
