package openstack

import (
	"context"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/routers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func resourceNetworkingRouterV2StateRefreshFunc(ctx context.Context, client *gophercloud.ServiceClient, routerID string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		n, err := routers.Get(ctx, client, routerID).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return n, "DELETED", nil
			}

			return n, "", err
		}

		return n, n.Status, nil
	}
}

func expandNetworkingRouterExternalFixedIPsV2(externalFixedIPs []any) []routersExternalFixedIP {
	fixedIPs := make([]routersExternalFixedIP, len(externalFixedIPs))

	for i, raw := range externalFixedIPs {
		rawMap := raw.(map[string]any)

		fixedIPs[i] = routersExternalFixedIP{
			SubnetID:  rawMap["subnet_id"].(string),
			IPAddress: rawMap["ip_address"].(string),
		}
	}

	return fixedIPs
}

func expandNetworkingRouterExternalSubnetIDsV2(externalSubnetIDs []any) []routersExternalFixedIP {
	subnetIDs := make([]routersExternalFixedIP, len(externalSubnetIDs))

	for i, raw := range externalSubnetIDs {
		subnetIDs[i] = routersExternalFixedIP{
			SubnetID: raw.(string),
		}
	}

	return subnetIDs
}

func flattenNetworkingRouterExternalFixedIPsV2(externalFixedIPs []routers.ExternalFixedIP) []map[string]string {
	fixedIPs := make([]map[string]string, len(externalFixedIPs))

	for i, fixedIP := range externalFixedIPs {
		fixedIPs[i] = map[string]string{
			"subnet_id":  fixedIP.SubnetID,
			"ip_address": fixedIP.IPAddress,
		}
	}

	return fixedIPs
}
