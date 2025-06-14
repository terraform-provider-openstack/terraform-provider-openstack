package openstack

import (
	"context"
	"log"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/groups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// networkingSecgroupV2StateRefreshFuncDelete returns a special case retry.StateRefreshFunc to try to delete a secgroup.
func networkingSecgroupV2StateRefreshFuncDelete(ctx context.Context, networkingClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		log.Printf("[DEBUG] Attempting to delete openstack_networking_secgroup_v2 %s", id)

		err := groups.Delete(ctx, networkingClient, id).ExtractErr()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				log.Printf("[DEBUG] Successfully deleted openstack_networking_secgroup_v2 %s", id)

				return "", "DELETED", nil
			}

			if gophercloud.ResponseCodeIs(err, http.StatusConflict) {
				return "", "ACTIVE", nil
			}

			return "", "ACTIVE", err
		}

		_, err = groups.Get(ctx, networkingClient, id).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				log.Printf("[DEBUG] Successfully deleted openstack_networking_secgroup_v2 %s", id)

				return "", "DELETED", nil
			}

			return "", "ACTIVE", err
		}

		log.Printf("[DEBUG] openstack_networking_secgroup_v2 %s is still active", id)

		return "", "ACTIVE", nil
	}
}
