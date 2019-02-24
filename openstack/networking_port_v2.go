package openstack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/extradhcpopts"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/portsbinding"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

type portExtended struct {
	ports.Port
	extradhcpopts.ExtraDHCPOptsExt
	portsbinding.PortsBindingExt
}

func resourceNetworkingPortV2StateRefreshFunc(client *gophercloud.ServiceClient, portID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		n, err := ports.Get(client, portID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return n, "DELETED", nil
			}

			return n, "", err
		}

		return n, n.Status, nil
	}
}

func expandNetworkingPortDHCPOptsV2Create(dhcpOpts *schema.Set) []extradhcpopts.CreateExtraDHCPOpt {
	var extraDHCPOpts []extradhcpopts.CreateExtraDHCPOpt

	if dhcpOpts != nil {
		for _, raw := range dhcpOpts.List() {
			rawMap := raw.(map[string]interface{})

			ipVersion := rawMap["ip_version"].(int)
			optName := rawMap["name"].(string)
			optValue := rawMap["value"].(string)

			extraDHCPOpts = append(extraDHCPOpts, extradhcpopts.CreateExtraDHCPOpt{
				OptName:   optName,
				OptValue:  optValue,
				IPVersion: gophercloud.IPVersion(ipVersion),
			})
		}
	}

	return extraDHCPOpts
}

func expandNetworkingPortDHCPOptsV2Update(oldDHCPopts, newDHCPopts *schema.Set) []extradhcpopts.UpdateExtraDHCPOpt {
	var extraDHCPOpts []extradhcpopts.UpdateExtraDHCPOpt
	var newOptNames []string

	if newDHCPopts != nil {
		for _, raw := range newDHCPopts.List() {
			rawMap := raw.(map[string]interface{})

			ipVersion := rawMap["ip_version"].(int)
			optName := rawMap["name"].(string)
			optValue := rawMap["value"].(string)
			// DHCP option name is the primary key, we will check this key below
			newOptNames = append(newOptNames, optName)

			extraDHCPOpts = append(extraDHCPOpts, extradhcpopts.UpdateExtraDHCPOpt{
				OptName:   optName,
				OptValue:  &optValue,
				IPVersion: gophercloud.IPVersion(ipVersion),
			})
		}
	}

	if oldDHCPopts != nil {
		for _, raw := range oldDHCPopts.List() {
			rawMap := raw.(map[string]interface{})

			optName := rawMap["name"].(string)

			// if we already add a new option with the same name, it means that we update it, no need to delete
			if !strSliceContains(newOptNames, optName) {
				extraDHCPOpts = append(extraDHCPOpts, extradhcpopts.UpdateExtraDHCPOpt{
					OptName:  optName,
					OptValue: nil,
				})
			}
		}
	}

	return extraDHCPOpts
}

func flattenNetworkingPortDHCPOptsV2(dhcpOpts extradhcpopts.ExtraDHCPOptsExt) []map[string]interface{} {
	dhcpOptsSet := make([]map[string]interface{}, len(dhcpOpts.ExtraDHCPOpts))

	for i, dhcpOpt := range dhcpOpts.ExtraDHCPOpts {
		dhcpOptsSet[i] = map[string]interface{}{
			"ip_version": dhcpOpt.IPVersion,
			"name":       dhcpOpt.OptName,
			"value":      dhcpOpt.OptValue,
		}
	}

	return dhcpOptsSet
}

func expandNetworkingPortAllowedAddressPairsV2(allowedAddressPairs *schema.Set) []ports.AddressPair {
	rawPairs := allowedAddressPairs.List()

	pairs := make([]ports.AddressPair, len(rawPairs))
	for i, raw := range rawPairs {
		rawMap := raw.(map[string]interface{})
		pairs[i] = ports.AddressPair{
			IPAddress:  rawMap["ip_address"].(string),
			MACAddress: rawMap["mac_address"].(string),
		}
	}

	return pairs
}

func flattenNetworkingPortAllowedAddressPairsV2(mac string, allowedAddressPairs []ports.AddressPair) []map[string]interface{} {
	pairs := make([]map[string]interface{}, len(allowedAddressPairs))

	for i, pair := range allowedAddressPairs {
		pairs[i] = map[string]interface{}{
			"ip_address": pair.IPAddress,
		}
		// Only set the MAC address if it is different than the
		// port's MAC. This means that a specific MAC was set.
		if pair.MACAddress != mac {
			pairs[i]["mac_address"] = pair.MACAddress
		}
	}

	return pairs
}

func expandNetworkingPortFixedIPV2(d *schema.ResourceData) interface{} {
	// If no_fixed_ip was specified, then just return an empty array.
	// Since no_fixed_ip is mutually exclusive to fixed_ip,
	// we can safely do this.
	//
	// Since we're only concerned about no_fixed_ip being set to "true",
	// GetOk is used.
	if _, ok := d.GetOk("no_fixed_ip"); ok {
		return []interface{}{}
	}

	rawIP := d.Get("fixed_ip").([]interface{})

	if len(rawIP) == 0 {
		return nil
	}

	ip := make([]ports.IP, len(rawIP))
	for i, raw := range rawIP {
		rawMap := raw.(map[string]interface{})
		ip[i] = ports.IP{
			SubnetID:  rawMap["subnet_id"].(string),
			IPAddress: rawMap["ip_address"].(string),
		}
	}
	return ip
}

func resourceNetworkingPortV2AllowedAddressPairsHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-%s", m["ip_address"].(string), m["mac_address"].(string)))

	return hashcode.String(buf.String())
}

func expandNetworkingPortFixedIPToStringSlice(fixedIPs []ports.IP) []string {
	s := make([]string, len(fixedIPs))
	for i, fixedIP := range fixedIPs {
		s[i] = fixedIP.IPAddress
	}

	return s
}

func flattenNetworkingPortBindingV2(port portExtended) []map[string]interface{} {
	var portBinding []map[string]interface{}

	profile, err := json.Marshal(port.Profile)
	if err != nil {
		log.Printf("[DEBUG] flattenNetworkingPortBindingV2: Cannot marshal port.Profile: %s", err)
	}

	vifDetails := make(map[string]string)
	for k, v := range port.VIFDetails {
		p, err := json.Marshal(v)
		if err != nil {
			log.Printf("[DEBUG] flattenNetworkingPortBindingV2: Cannot marshal %s key value: %s", k, err)
		}
		vifDetails[k] = string(p)
	}

	portBinding = append(portBinding, map[string]interface{}{
		"profile":     string(profile),
		"vif_type":    port.VIFType,
		"vif_details": vifDetails,
		"vnic_type":   port.VNICType,
		"host_id":     port.HostID,
	})

	return portBinding
}

func expandNetworkingPortBindingProfileCreateV2(portBinding []interface{}, opts *ports.CreateOptsBuilder) error {
	for _, raw := range portBinding {
		binding := raw.(map[string]interface{})
		profile := map[string]interface{}{}

		// Convert raw string into the map
		rawProfile := binding["profile"].(string)
		if len(rawProfile) > 0 {
			err := json.Unmarshal([]byte(rawProfile), &profile)
			if err != nil {
				return fmt.Errorf("Failed to unmarshal the JSON: %s", err)
			}
		}

		*opts = portsbinding.CreateOptsExt{
			CreateOptsBuilder: *opts,
			HostID:            binding["host_id"].(string),
			Profile:           profile,
			VNICType:          binding["vnic_type"].(string),
		}
	}

	return nil
}

func expandNetworkingPortBindingProfileUpdateV2(portBinding []interface{}, opts *ports.UpdateOptsBuilder) error {
	profile := map[string]interface{}{}

	// default options, when unsetting the port bindings
	newOpts := portsbinding.UpdateOptsExt{
		UpdateOptsBuilder: *opts,
		HostID:            new(string),
		Profile:           profile,
		VNICType:          "normal",
	}

	for _, raw := range portBinding {
		binding := raw.(map[string]interface{})

		// Convert raw string into the map
		rawProfile := binding["profile"].(string)
		if len(rawProfile) > 0 {
			err := json.Unmarshal([]byte(rawProfile), &profile)
			if err != nil {
				return fmt.Errorf("Failed to unmarshal the JSON: %s", err)
			}
		}

		hostID := binding["host_id"].(string)
		newOpts = portsbinding.UpdateOptsExt{
			UpdateOptsBuilder: *opts,
			HostID:            &hostID,
			Profile:           profile,
			VNICType:          binding["vnic_type"].(string),
		}
	}

	*opts = newOpts

	return nil
}

func suppressDiffPortBindingProfileV2(k, old, new string, d *schema.ResourceData) bool {
	if old == "{}" && new == "" {
		return true
	}

	if old == "" && new == "{}" {
		return true
	}

	return false
}
