package openstack

import (
	"context"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/addressscopes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func resourceNetworkingAddressScopeV2StateRefreshFunc(ctx context.Context, client *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		a, err := addressscopes.Get(ctx, client, id).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return a, "DELETED", nil
			}

			return nil, "", err
		}

		return a, "ACTIVE", nil
	}
}
