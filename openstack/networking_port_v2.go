package openstack

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/dns"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/extradhcpopts"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/portsbinding"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/portsecurity"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/qos/policies"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/terraform-provider-openstack/utils/v2/hashcode"
)

type portExtended struct {
	ports.Port
	extradhcpopts.ExtraDHCPOptsExt
	portsecurity.PortSecurityExt
	portsbinding.PortsBindingExt
	dns.PortDNSExt
	policies.QoSPolicyExt
}

func resourceNetworkingPortV2StateRefreshFunc(ctx context.Context, client *gophercloud.ServiceClient, portID string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		n, err := ports.Get(ctx, client, portID).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
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
			rawMap := raw.(map[string]any)

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
			rawMap := raw.(map[string]any)

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
			rawMap := raw.(map[string]any)

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

func flattenNetworkingPortDHCPOptsV2(dhcpOpts extradhcpopts.ExtraDHCPOptsExt) []map[string]any {
	dhcpOptsSet := make([]map[string]any, len(dhcpOpts.ExtraDHCPOpts))

	for i, dhcpOpt := range dhcpOpts.ExtraDHCPOpts {
		dhcpOptsSet[i] = map[string]any{
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
		rawMap := raw.(map[string]any)
		pairs[i] = ports.AddressPair{
			IPAddress:  rawMap["ip_address"].(string),
			MACAddress: rawMap["mac_address"].(string),
		}
	}

	return pairs
}

func flattenNetworkingPortAllowedAddressPairsV2(mac string, allowedAddressPairs []ports.AddressPair) []map[string]any {
	pairs := make([]map[string]any, len(allowedAddressPairs))

	for i, pair := range allowedAddressPairs {
		pairs[i] = map[string]any{
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

func expandNetworkingPortFixedIPV2(d *schema.ResourceData) any {
	// If no_fixed_ip was specified, then just return an empty array.
	// Since no_fixed_ip is mutually exclusive to fixed_ip,
	// we can safely do this.
	//
	// Since we're only concerned about no_fixed_ip being set to "true",
	// GetOk is used.
	if _, ok := d.GetOk("no_fixed_ip"); ok {
		return []any{}
	}

	rawIP := d.Get("fixed_ip").([]any)

	if len(rawIP) == 0 {
		return nil
	}

	ip := make([]ports.IP, 0, len(rawIP))

	for _, raw := range rawIP {
		if raw == nil {
			continue
		}

		rawMap := raw.(map[string]any)
		subnetID := rawMap["subnet_id"].(string)
		ipAddress := rawMap["ip_address"].(string)

		if subnetID == "" && ipAddress == "" {
			continue
		}

		ip = append(ip, ports.IP{
			SubnetID:  rawMap["subnet_id"].(string),
			IPAddress: rawMap["ip_address"].(string),
		})
	}

	return ip
}

func resourceNetworkingPortV2AllowedAddressPairsHash(v any) int {
	var buf bytes.Buffer

	m := v.(map[string]any)
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

func flattenNetworkingPortBindingV2(port portExtended) any {
	var portBinding []map[string]any

	var profile any

	if port.Profile != nil {
		// "TypeMap" with "ValidateFunc", "DiffSuppressFunc" and "StateFunc" combination
		// is not supported by Terraform. Therefore a regular JSON string is used for the
		// port resource.
		tmp, err := json.Marshal(port.Profile)
		if err != nil {
			log.Printf("[DEBUG] flattenNetworkingPortBindingV2: Cannot marshal port.Profile: %s", err)
		}

		profile = string(tmp)
	}

	vifDetails := make(map[string]string)

	for k, v := range port.VIFDetails {
		// don't marshal, if it is a regular string
		if s, ok := v.(string); ok {
			vifDetails[k] = s

			continue
		}

		p, err := json.Marshal(v)
		if err != nil {
			log.Printf("[DEBUG] flattenNetworkingPortBindingV2: Cannot marshal %s key value: %s", k, err)
		}

		vifDetails[k] = string(p)
	}

	portBinding = append(portBinding, map[string]any{
		"profile":     profile,
		"vif_type":    port.VIFType,
		"vif_details": vifDetails,
		"vnic_type":   port.VNICType,
		"host_id":     port.HostID,
	})

	return portBinding
}
