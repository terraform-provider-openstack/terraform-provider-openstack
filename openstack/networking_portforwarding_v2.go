package openstack

import (
	"context"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/portforwarding"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func networkingPortForwardingV2StateRefreshFunc(ctx context.Context, client *gophercloud.ServiceClient, fipID, pfID string) retry.StateRefreshFunc {
	return func() (any, string, error) {
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
