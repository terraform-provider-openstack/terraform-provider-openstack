package openstack

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/rules"
	"github.com/gophercloud/utils/terraform/hashcode"
)

// networkingSecgroupV2StateRefreshFuncDelete returns a special case resource.StateRefreshFunc to try to delete a secgroup.
func networkingSecgroupV2StateRefreshFuncDelete(networkingClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete openstack_networking_secgroup_v2 %s", id)

		r, err := groups.Get(networkingClient, id).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted openstack_networking_secgroup_v2 %s", id)
				return r, "DELETED", nil
			}

			return r, "ACTIVE", err
		}

		err = groups.Delete(networkingClient, id).ExtractErr()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted openstack_networking_secgroup_v2 %s", id)
				return r, "DELETED", nil
			}
			if _, ok := err.(gophercloud.ErrDefault409); ok {
				return r, "ACTIVE", nil
			}

			return r, "ACTIVE", err
		}

		log.Printf("[DEBUG] openstack_networking_secgroup_v2 %s is still active", id)

		return r, "ACTIVE", nil
	}
}

func networkingSecgroupV2RuleHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	buf.WriteString(fmt.Sprintf("%s-", m["description"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["direction"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["ethertype"].(string)))
	buf.WriteString(fmt.Sprintf("%s-", m["protocol"].(string)))
	buf.WriteString(fmt.Sprintf("%d-", m["port_range_min"].(int)))
	buf.WriteString(fmt.Sprintf("%d-", m["port_range_max"].(int)))
	buf.WriteString(fmt.Sprintf("%s-", strings.ToLower(m["remote_ip_prefix"].(string))))
	buf.WriteString(fmt.Sprintf("%s-", m["remote_group_id"].(string)))
	buf.WriteString(fmt.Sprintf("%t-", m["self"].(bool)))

	return hashcode.String(buf.String())
}

func networkingSecgroupV2RulesCheckForErrors(d *schema.ResourceData) error {
	rawRules := d.Get("rule").(*schema.Set).List()

	for _, rawRule := range rawRules {
		rawRuleMap := rawRule.(map[string]interface{})

		// only one of cidr, from_group_id, or self can be set
		remoteIPPrefix := rawRuleMap["remote_ip_prefix"].(string)
		remoteGroupID := rawRuleMap["remote_group_id"].(string)
		self := rawRuleMap["self"].(bool)
		errorMessage := fmt.Errorf("Only one of remote_ip_prefix, remote_group_id, or self can be set")

		// if remote_ip_prefix is set, remote_group_id and self cannot be set
		if remoteIPPrefix != "" {
			if remoteGroupID != "" || self {
				return errorMessage
			}
		}

		// if remote_group_id is set, remote_ip_prefix and self cannot be set
		if remoteGroupID != "" {
			if remoteIPPrefix != "" || self {
				return errorMessage
			}
		}

		// if self is set, remote_ip_prefix and remote_group_id cannot be set
		if self {
			if remoteIPPrefix != "" || remoteGroupID != "" {
				return errorMessage
			}
		}
	}

	return nil
}

func expandNetworkingSecgroupV2CreateRules(d *schema.ResourceData) ([]rules.CreateOpts, error) {
	rawRules := d.Get("rule").(*schema.Set).List()
	createRuleOptsList := make([]rules.CreateOpts, len(rawRules))

	for i, rawRule := range rawRules {
		rule, err := expandNetworkingSecgroupV2CreateRule(d, rawRule)
		if err != nil {
			return []rules.CreateOpts{}, err
		}
		createRuleOptsList[i] = rule
	}

	return createRuleOptsList, nil
}

func expandNetworkingSecgroupV2CreateRule(d *schema.ResourceData, rawRule interface{}) (rules.CreateOpts, error) {
	rawRuleMap := rawRule.(map[string]interface{})
	remoteGroupID := rawRuleMap["remote_group_id"].(string)
	if rawRuleMap["self"].(bool) {
		remoteGroupID = d.Id()
	}

	protocol, err := resourceNetworkingSecGroupRuleV2Protocol(rawRuleMap["protocol"].(string))
	if err != nil {
		return rules.CreateOpts{}, err
	}

	direction, err := resourceNetworkingSecGroupRuleV2Direction(rawRuleMap["direction"].(string))
	if err != nil {
		return rules.CreateOpts{}, err
	}

	ethertype, err := resourceNetworkingSecGroupRuleV2EtherType(rawRuleMap["ethertype"].(string))
	if err != nil {
		return rules.CreateOpts{}, err
	}

	return rules.CreateOpts{
		SecGroupID:     d.Id(),
		Description:    rawRuleMap["description"].(string),
		Direction:      direction,
		EtherType:      ethertype,
		Protocol:       protocol,
		PortRangeMin:   rawRuleMap["port_range_min"].(int),
		PortRangeMax:   rawRuleMap["port_range_max"].(int),
		RemoteIPPrefix: rawRuleMap["remote_ip_prefix"].(string),
		RemoteGroupID:  remoteGroupID,
	}, nil
}

func flattenNetworkingSecgroupV2Rules(computeClient *gophercloud.ServiceClient, d *schema.ResourceData, sgrs []rules.SecGroupRule) ([]map[string]interface{}, error) {
	sgrMap := make([]map[string]interface{}, len(sgrs))
	for i, sgr := range sgrs {
		self := false

		if sgr.RemoteGroupID == d.Id() {
			self = true
		}

		sgrMap[i] = map[string]interface{}{
			"id":               sgr.ID,
			"description":      sgr.Description,
			"direction":        sgr.Direction,
			"ethertype":        sgr.EtherType,
			"protocol":         sgr.Protocol,
			"port_range_min":   sgr.PortRangeMin,
			"port_range_max":   sgr.PortRangeMax,
			"remote_ip_prefix": sgr.RemoteIPPrefix,
			"remote_group_id":  sgr.RemoteGroupID,
			"self":             self,
		}
	}
	return sgrMap, nil
}

func expandNetworkingSecgroupV2Rule(d *schema.ResourceData, rawRule interface{}) rules.SecGroupRule {
	rawRuleMap := rawRule.(map[string]interface{})

	return rules.SecGroupRule{
		ID:             rawRuleMap["id"].(string),
		SecGroupID:     d.Id(),
		Description:    rawRuleMap["description"].(string),
		Direction:      rawRuleMap["direction"].(string),
		EtherType:      rawRuleMap["ethertype"].(string),
		PortRangeMin:   rawRuleMap["port_range_min"].(int),
		PortRangeMax:   rawRuleMap["port_range_max"].(int),
		RemoteIPPrefix: rawRuleMap["remote_ip_prefix"].(string),
		RemoteGroupID:  rawRuleMap["remote_group_id"].(string),
	}
}
