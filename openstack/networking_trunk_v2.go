package openstack

import (
	"context"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/trunks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func networkingTrunkV2StateRefreshFunc(ctx context.Context, client *gophercloud.ServiceClient, trunkID string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		trunk, err := trunks.Get(ctx, client, trunkID).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return trunk, "DELETED", nil
			}

			return nil, "", err
		}

		return trunk, trunk.Status, nil
	}
}

func flattenNetworkingTrunkV2Subports(subports []trunks.Subport) []map[string]any {
	trunkSubports := make([]map[string]any, len(subports))

	for i, subport := range subports {
		trunkSubports[i] = map[string]any{
			"port_id":           subport.PortID,
			"segmentation_type": subport.SegmentationType,
			"segmentation_id":   subport.SegmentationID,
		}
	}

	return trunkSubports
}

func expandNetworkingTrunkV2Subports(subports *schema.Set) []trunks.Subport {
	rawSubports := subports.List()

	trunkSubports := make([]trunks.Subport, len(rawSubports))

	for i, raw := range rawSubports {
		rawMap := raw.(map[string]any)

		trunkSubports[i] = trunks.Subport{
			PortID:           rawMap["port_id"].(string),
			SegmentationType: rawMap["segmentation_type"].(string),
			SegmentationID:   rawMap["segmentation_id"].(int),
		}
	}

	return trunkSubports
}

func expandNetworkingTrunkV2SubportsRemove(subports *schema.Set) []trunks.RemoveSubport {
	rawSubports := subports.List()

	subportsToRemove := make([]trunks.RemoveSubport, len(rawSubports))

	for i, raw := range rawSubports {
		rawMap := raw.(map[string]any)

		subportsToRemove[i] = trunks.RemoveSubport{
			PortID: rawMap["port_id"].(string),
		}
	}

	return subportsToRemove
}
