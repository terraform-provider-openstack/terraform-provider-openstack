package openstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/policies"
)

func fwPolicyV2DeleteFunc(networkingClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		err := policies.Delete(networkingClient, id).Err
		if err == nil {
			return "", "DELETED", nil
		}

		if _, ok := err.(gophercloud.ErrDefault409); ok {
			// This error usually means that the policy is attached
			// to a firewall. At this point, the firewall is probably
			// being delete. So, we retry a few times.

			return nil, "ACTIVE", nil
		}

		return nil, "ACTIVE", err
	}
}
