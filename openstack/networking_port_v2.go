package openstack

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/extradhcpopts"
	"github.com/hashicorp/terraform/helper/schema"
)

func expandNetworkingPortDHCPOptsV2Create(dhcpOpts *schema.Set) []extradhcpopts.CreateExtraDHCPOpt {
	if dhcpOpts == nil {
		return []extradhcpopts.CreateExtraDHCPOpt{}
	}

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

func expandNetworkingPortDHCPOptsV2Update(dhcpOpts *schema.Set) []extradhcpopts.UpdateExtraDHCPOpt {
	if dhcpOpts == nil {
		return []extradhcpopts.UpdateExtraDHCPOpt{}
	}

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

func expandNetworkingPortDHCPOptsV2Delete(dhcpOpts *schema.Set) []extradhcpopts.UpdateExtraDHCPOpt {
	if dhcpOpts == nil {
		return []extradhcpopts.UpdateExtraDHCPOpt{}
	}

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

func flattenNetworkingPortDHCPOptsV2(dhcpOpts extradhcpopts.ExtraDHCPOptsExt) *schema.Set {
	dhcpOptsSet := &schema.Set{}

	for _, dhcpOpt := range dhcpOpts.ExtraDHCPOpts {
		dhcpOptsSet.Add(map[string]interface{}{
			"ip_version": dhcpOpt.IPVersion,
			"opt_name":   dhcpOpt.OptName,
			"opt_value":  dhcpOpt.OptValue,
		})
	}

	return dhcpOptsSet
}
