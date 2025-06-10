package openstack

import (
	"context"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/subnetpools"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func networkingSubnetpoolV2StateRefreshFunc(ctx context.Context, client *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		subnetpool, err := subnetpools.Get(ctx, client, id).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return subnetpool, "DELETED", nil
			}

			if gophercloud.ResponseCodeIs(err, http.StatusConflict) {
				return subnetpool, "ACTIVE", nil
			}

			return nil, "", err
		}

		return subnetpool, "ACTIVE", nil
	}
}
