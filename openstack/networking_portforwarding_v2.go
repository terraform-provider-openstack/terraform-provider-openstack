package openstack

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/portforwarding"
)

func networkingPortForwardingV2StateRefreshFunc(ctx context.Context, client *gophercloud.ServiceClient, fipID, pfID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		pf, err := portforwarding.Get(ctx, client, fipID, pfID).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return pf, "DELETED", nil
			}

			return nil, "", err
		}

		return pf, "ACTIVE", nil
	}
}
