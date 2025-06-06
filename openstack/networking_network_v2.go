package openstack

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/dns"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/external"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/mtu"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/portsecurity"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/provider"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/qos/policies"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/vlantransparent"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/v2/pagination"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type networkExtended struct {
	networks.Network
	external.NetworkExternalExt
	vlantransparent.TransparentExt
	portsecurity.PortSecurityExt
	mtu.NetworkMTUExt
	dns.NetworkDNSExt
	policies.QoSPolicyExt
	provider.NetworkProviderExt
}

// networkingNetworkV2ID retrieves network ID by the provided name.
func networkingNetworkV2ID(ctx context.Context, d *schema.ResourceData, meta any, networkName string) (string, error) {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return "", fmt.Errorf("Error creating OpenStack network client: %w", err)
	}

	opts := networks.ListOpts{Name: networkName}
	pager := networks.List(networkingClient, opts)
	networkID := ""

	err = pager.EachPage(ctx, func(_ context.Context, page pagination.Page) (bool, error) {
		networkList, err := networks.ExtractNetworks(page)
		if err != nil {
			return false, err
		}

		for _, n := range networkList {
			if n.Name == networkName {
				networkID = n.ID

				return false, nil
			}
		}

		return true, nil
	})

	return networkID, err
}

// networkingNetworkV2Name retrieves network name by the provided ID.
func networkingNetworkV2Name(ctx context.Context, d *schema.ResourceData, meta any, networkID string) (string, error) {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return "", fmt.Errorf("Error creating OpenStack network client: %w", err)
	}

	opts := networks.ListOpts{ID: networkID}
	pager := networks.List(networkingClient, opts)
	networkName := ""

	err = pager.EachPage(ctx, func(_ context.Context, page pagination.Page) (bool, error) {
		networkList, err := networks.ExtractNetworks(page)
		if err != nil {
			return false, err
		}

		for _, n := range networkList {
			if n.ID == networkID {
				networkName = n.Name

				return false, nil
			}
		}

		return true, nil
	})

	return networkName, err
}

func resourceNetworkingNetworkV2StateRefreshFunc(ctx context.Context, client *gophercloud.ServiceClient, networkID string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		n, err := networks.Get(ctx, client, networkID).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return n, "DELETED", nil
			}

			if gophercloud.ResponseCodeIs(err, http.StatusConflict) {
				return n, "ACTIVE", nil
			}

			return n, "", err
		}

		return n, n.Status, nil
	}
}

func expandNetworkingNetworkSegmentsV2(segments *schema.Set) []provider.Segment {
	rawSegments := segments.List()

	if len(rawSegments) == 1 {
		// unset segments
		rawMap := rawSegments[0].(map[string]any)
		if rawMap["physical_network"] == "" &&
			rawMap["network_type"] == "" &&
			rawMap["segmentation_id"] == 0 {
			return nil
		}
	}

	providerSegments := make([]provider.Segment, len(rawSegments))

	for i, raw := range rawSegments {
		rawMap := raw.(map[string]any)
		providerSegments[i] = provider.Segment{
			PhysicalNetwork: rawMap["physical_network"].(string),
			NetworkType:     rawMap["network_type"].(string),
			SegmentationID:  rawMap["segmentation_id"].(int),
		}
	}

	return providerSegments
}

func flattenNetworkingNetworkSegmentsV2(network networkExtended) []map[string]any {
	singleSegment := 0
	if network.NetworkType != "" ||
		network.PhysicalNetwork != "" ||
		network.SegmentationID != "" {
		singleSegment = 1
	}

	segmentsSet := make([]map[string]any, len(network.Segments)+singleSegment)

	if singleSegment > 0 {
		segmentationID, err := strconv.Atoi(network.SegmentationID)
		if err != nil {
			log.Printf("[DEBUG] Unable to convert %q segmentation ID to an integer: %s", network.SegmentationID, err)
		}

		segmentsSet[0] = map[string]any{
			"physical_network": network.PhysicalNetwork,
			"network_type":     network.NetworkType,
			"segmentation_id":  segmentationID,
		}
	}

	for i, segment := range network.Segments {
		segmentsSet[i+singleSegment] = map[string]any{
			"physical_network": segment.PhysicalNetwork,
			"network_type":     segment.NetworkType,
			"segmentation_id":  segment.SegmentationID,
		}
	}

	return segmentsSet
}
