package openstack

import (
	"log"

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

		if err != nil {
			switch err.(type) {
			case gophercloud.ErrDefault404:
				// This error usually means that the policy was deleted manually
				log.Printf("[DEBUG] Unable to find openstack_fw_policy_v2 %s: %s", id, err)
				return "", "DELETED", nil
			case gophercloud.ErrDefault409:
				// This error usually means that the policy is attached to a firewall.
				log.Printf("[DEBUG] Error to delete openstack_fw_policy_v2 %s: %s", id, err)
				return nil, "ACTIVE", nil
			}
		}

		return nil, "ACTIVE", err
	}
}
