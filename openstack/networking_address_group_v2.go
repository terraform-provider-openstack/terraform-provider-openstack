package openstack

import (
	"context"
	"log"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/addressgroups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// networkingAddressGroupV2StateRefreshFuncDelete returns a special case retry.StateRefreshFunc to try to delete a secgroup.
func networkingAddressGroupV2StateRefreshFuncDelete(ctx context.Context, networkingClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		log.Printf("[DEBUG] Attempting to delete openstack_networking_address_group_v2 %s", id)

		err := addressgroups.Delete(ctx, networkingClient, id).ExtractErr()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				log.Printf("[DEBUG] Successfully deleted openstack_networking_address_group_v2 %s", id)

				return "", "DELETED", nil
			}

			if gophercloud.ResponseCodeIs(err, http.StatusConflict) {
				return "", "ACTIVE", nil
			}

			return "", "ACTIVE", err
		}

		log.Printf("[DEBUG] openstack_networking_address_group_v2 %s is deleted", id)

		return "", "DELETED", nil
	}
}
