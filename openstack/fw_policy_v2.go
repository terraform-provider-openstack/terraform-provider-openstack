package openstack

import (
	"context"
	"log"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/fwaas_v2/policies"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func fwPolicyV2DeleteFunc(ctx context.Context, networkingClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		err := policies.Delete(ctx, networkingClient, id).Err
		if err == nil {
			return "", "DELETED", nil
		}

		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				// This error usually means that the policy was deleted manually
				log.Printf("[DEBUG] Unable to find openstack_fw_policy_v2 %s: %s", id, err)

				return "", "DELETED", nil
			}

			if gophercloud.ResponseCodeIs(err, http.StatusConflict) {
				// This error usually means that the policy is attached to a firewall.
				log.Printf("[DEBUG] Error to delete openstack_fw_policy_v2 %s: %s", id, err)

				return nil, "ACTIVE", nil
			}
		}

		return nil, "ACTIVE", err
	}
}
