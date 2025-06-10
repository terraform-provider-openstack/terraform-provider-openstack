package openstack

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/qos/rules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func resourceNetworkingQoSRuleV2BuildID(qosPolicyID, qosRuleID string) string {
	return fmt.Sprintf("%s/%s", qosPolicyID, qosRuleID)
}

func networkingQoSBandwidthLimitRuleV2StateRefreshFunc(ctx context.Context, client *gophercloud.ServiceClient, policyID, ruleID string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		policy, err := rules.GetBandwidthLimitRule(ctx, client, policyID, ruleID).ExtractBandwidthLimitRule()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return policy, "DELETED", nil
			}

			if gophercloud.ResponseCodeIs(err, http.StatusConflict) {
				return policy, "ACTIVE", nil
			}

			return nil, "", err
		}

		return policy, "ACTIVE", nil
	}
}

func networkingQoSDSCPMarkingRuleV2StateRefreshFunc(ctx context.Context, client *gophercloud.ServiceClient, policyID, ruleID string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		policy, err := rules.GetDSCPMarkingRule(ctx, client, policyID, ruleID).ExtractDSCPMarkingRule()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return policy, "DELETED", nil
			}

			if gophercloud.ResponseCodeIs(err, http.StatusConflict) {
				return policy, "ACTIVE", nil
			}

			return nil, "", err
		}

		return policy, "ACTIVE", nil
	}
}

func networkingQoSMinimumBandwidthRuleV2StateRefreshFunc(ctx context.Context, client *gophercloud.ServiceClient, policyID, ruleID string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		policy, err := rules.GetMinimumBandwidthRule(ctx, client, policyID, ruleID).ExtractMinimumBandwidthRule()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return policy, "DELETED", nil
			}

			if gophercloud.ResponseCodeIs(err, http.StatusConflict) {
				return policy, "ACTIVE", nil
			}

			return nil, "", err
		}

		return policy, "ACTIVE", nil
	}
}
