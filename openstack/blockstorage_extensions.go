package openstack

import (
	"bytes"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/schedulerhints"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
)

func expandBlockStorageExtensionsSchedulerHints(v schedulerhints.SchedulerHints) map[string]interface{} {
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

func resourceBlockStorageSchedulerHints(d *schema.ResourceData, schedulerHintsRaw map[string]interface{}) schedulerhints.SchedulerHints {
	var differentHost []string
	if v, ok := schedulerHintsRaw["different_host"].([]interface{}); ok {
		for _, dh := range v {
			differentHost = append(differentHost, dh.(string))
		}
	}

	var sameHost []string
	if v, ok := schedulerHintsRaw["same_host"].([]interface{}); ok {
		for _, sh := range v {
			sameHost = append(sameHost, sh.(string))
		}
	}

	schedulerHints := schedulerhints.SchedulerHints{
		DifferentHost:        differentHost,
		SameHost:             sameHost,
		Query:                schedulerHintsRaw["query"].(string),
		LocalToInstance:      schedulerHintsRaw["local_to_instance"].(string),
		AdditionalProperties: schedulerHintsRaw["additional_properties"].(map[string]interface{}),
	}

	return schedulerHints
}
