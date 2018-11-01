package openstack

import (
	"bytes"
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/extradhcpopts"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func expandDHCPOptionsV2Create(dhcpOpts *schema.Set) []extradhcpopts.CreateExtraDHCPOpt {
	rawDHCPOpts := dhcpOpts.List()

	extraDHCPOpts := make([]extradhcpopts.CreateExtraDHCPOpt, dhcpOpts.Len())
	for i, raw := range rawDHCPOpts {
		rawMap := raw.(map[string]interface{})

		ipVersion := rawMap["ip_version"].(int)
		optName := rawMap["opt_name"].(string)
		optValue := rawMap["opt_value"].(string)

		extraDHCPOpts[i] = extradhcpopts.CreateExtraDHCPOpt{
			OptName:   optName,
			OptValue:  optValue,
			IPVersion: gophercloud.IPVersion(ipVersion),
		}
	}

	return extraDHCPOpts
}

func expandDHCPOptionsV2Update(dhcpOpts *schema.Set) []extradhcpopts.UpdateExtraDHCPOpt {
	rawDHCPOpts := dhcpOpts.List()

	extraDHCPOpts := make([]extradhcpopts.UpdateExtraDHCPOpt, dhcpOpts.Len())
	for i, raw := range rawDHCPOpts {
		rawMap := raw.(map[string]interface{})

		ipVersion := rawMap["ip_version"].(int)
		optName := rawMap["opt_name"].(string)
		optValue := rawMap["opt_value"].(string)

		extraDHCPOpts[i] = extradhcpopts.UpdateExtraDHCPOpt{
			OptName:   optName,
			OptValue:  &optValue,
			IPVersion: gophercloud.IPVersion(ipVersion),
		}
	}

	return extraDHCPOpts
}

func expandDHCPOptionsV2Delete(dhcpOpts *schema.Set) []extradhcpopts.UpdateExtraDHCPOpt {
	rawDHCPOpts := dhcpOpts.List()

	extraDHCPOpts := make([]extradhcpopts.UpdateExtraDHCPOpt, dhcpOpts.Len())
	for i, raw := range rawDHCPOpts {
		rawMap := raw.(map[string]interface{})
		extraDHCPOpts[i] = extradhcpopts.UpdateExtraDHCPOpt{
			OptName:  rawMap["opt_name"].(string),
			OptValue: nil,
		}
	}

	return extraDHCPOpts
}

func flattenDHCPOptionsV2(dhcpOpts extradhcpopts.ExtraDHCPOptsExt) *schema.Set {
	dhcpOptsSet := &schema.Set{
		F: dhcpOptionsV2HashSetFunc(),
	}

	for _, dhcpOpt := range dhcpOpts.ExtraDHCPOpts {
		dhcpOptsSet.Add(map[string]interface{}{
			"ip_version": dhcpOpt.IPVersion,
			"opt_name":   dhcpOpt.OptName,
			"opt_value":  dhcpOpt.OptValue,
		})
	}

	return dhcpOptsSet
}

// dhcpOptionsV2Schema returns *schema.Resource from the "extra_dhcp_opts" attribute.
func dhcpOptionsV2Schema() *schema.Resource {
	return resourceNetworkingPortExtraDHCPOptionsV2().Schema["extra_dhcp_opts"].Elem.(*schema.Resource)
}

// dhcpOptionsV2HashSetFunc returns schema.SchemaSetFunc that can be used to
// create a new schema.Set for the "extra_dhcp_opts" attribute.
func dhcpOptionsV2HashSetFunc() schema.SchemaSetFunc {
	return schema.HashResource(dhcpOptionsV2Schema())
}

// hashDHCPOptionsV2 is a hash function to use with the "extra_dhcp_opts" set.
func hashDHCPOptionsV2(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["opt_name"].(string)))
	if m["ip_version"] != "" {
		buf.WriteString(fmt.Sprintf("%d-", m["ip_version"].(int)))
	}
	return hashcode.String(buf.String())
}
