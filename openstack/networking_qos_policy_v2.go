package openstack

import (
	"context"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/qos/policies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// QoSPolicyCreateOpts represents the attributes used when creating a new QoS policy.
type QoSPolicyCreateOpts struct {
	policies.CreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}

func networkingQoSPolicyV2StateRefreshFunc(ctx context.Context, client *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		policy, err := policies.Get(ctx, client, id).Extract()
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
