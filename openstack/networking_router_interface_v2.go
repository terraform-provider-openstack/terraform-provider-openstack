package openstack

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/hashicorp/terraform/helper/resource"
)

func resourceNetworkingRouterInterfaceV2StateRefreshFunc(networkingClient *gophercloud.ServiceClient, portID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := ports.Get(networkingClient, portID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return r, "DELETED", nil
			}

			return r, "", err
		}

		return r, r.Status, nil
	}
}
