package openstack

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/qos/rules"
)

func resourceNetworkingQoSRuleV2BuildID(qosPolicyID, qosRuleID string) string {
	return fmt.Sprintf("%s/%s", qosPolicyID, qosRuleID)
}

func networkingQoSBandwidthLimitRuleV2StateRefreshFunc(client *gophercloud.ServiceClient, policyID, ruleID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		policy, err := rules.GetBandwidthLimitRule(client, policyID, ruleID).ExtractBandwidthLimitRule()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return policy, "DELETED", nil
			}
			if _, ok := err.(gophercloud.ErrDefault409); ok {
				return policy, "ACTIVE", nil
			}

			return nil, "", err
		}

		return policy, "ACTIVE", nil
	}
}

func networkingQoSDSCPMarkingRuleV2StateRefreshFunc(client *gophercloud.ServiceClient, policyID, ruleID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		policy, err := rules.GetDSCPMarkingRule(client, policyID, ruleID).ExtractDSCPMarkingRule()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return policy, "DELETED", nil
			}
			if _, ok := err.(gophercloud.ErrDefault409); ok {
				return policy, "ACTIVE", nil
			}

			return nil, "", err
		}

		return policy, "ACTIVE", nil
	}
}

func networkingQoSMinimumBandwidthRuleV2StateRefreshFunc(client *gophercloud.ServiceClient, policyID, ruleID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		policy, err := rules.GetMinimumBandwidthRule(client, policyID, ruleID).ExtractMinimumBandwidthRule()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return policy, "DELETED", nil
			}
			if _, ok := err.(gophercloud.ErrDefault409); ok {
				return policy, "ACTIVE", nil
			}

			return nil, "", err
		}

		return policy, "ACTIVE", nil
	}
}
