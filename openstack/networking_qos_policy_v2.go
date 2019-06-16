package openstack

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/qos/policies"
	"github.com/hashicorp/terraform/helper/resource"
)

func networkingQoSPolicyV2StateRefreshFunc(client *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		policy, err := policies.Get(client, id).Extract()
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
