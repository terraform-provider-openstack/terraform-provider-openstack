package openstack

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/subnets"
)

// networkingSubnetV2StateRefreshFunc returns a standard retry.StateRefreshFunc to wait for subnet status.
func networkingSubnetV2StateRefreshFunc(ctx context.Context, client *gophercloud.ServiceClient, subnetID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		subnet, err := subnets.Get(ctx, client, subnetID).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return subnet, "DELETED", nil
			}

			return nil, "", err
		}

		return subnet, "ACTIVE", nil
	}
}

// networkingSubnetV2StateRefreshFuncDelete returns a special case retry.StateRefreshFunc to try to delete a subnet.
func networkingSubnetV2StateRefreshFuncDelete(ctx context.Context, networkingClient *gophercloud.ServiceClient, subnetID string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		log.Printf("[DEBUG] Attempting to delete openstack_networking_subnet_v2 %s", subnetID)

		s, err := subnets.Get(ctx, networkingClient, subnetID).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				log.Printf("[DEBUG] Successfully deleted openstack_networking_subnet_v2 %s", subnetID)
				return s, "DELETED", nil
			}

			return s, "ACTIVE", err
		}

		err = subnets.Delete(ctx, networkingClient, subnetID).ExtractErr()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				log.Printf("[DEBUG] Successfully deleted openstack_networking_subnet_v2 %s", subnetID)
				return s, "DELETED", nil
			}
			// Subnet is still in use - we can retry.
			if gophercloud.ResponseCodeIs(err, http.StatusConflict) {
				return s, "ACTIVE", nil
			}

			return s, "ACTIVE", err
		}

		log.Printf("[DEBUG] openstack_networking_subnet_v2 %s is still active", subnetID)

		return s, "ACTIVE", nil
	}
}

// expandNetworkingSubnetV2AllocationPools returns a slice of subnets.AllocationPool structs.
func expandNetworkingSubnetV2AllocationPools(allocationPools []interface{}) []subnets.AllocationPool {
	result := make([]subnets.AllocationPool, len(allocationPools))
	for i, raw := range allocationPools {
		rawMap := raw.(map[string]interface{})

		result[i] = subnets.AllocationPool{
			Start: rawMap["start"].(string),
			End:   rawMap["end"].(string),
		}
	}

	return result
}

// flattenNetworkingSubnetV2AllocationPools allows to flatten slice of subnets.AllocationPool structs into
// a slice of maps.
func flattenNetworkingSubnetV2AllocationPools(allocationPools []subnets.AllocationPool) []map[string]interface{} {
	result := make([]map[string]interface{}, len(allocationPools))
	for i, allocationPool := range allocationPools {
		pool := make(map[string]interface{})
		pool["start"] = allocationPool.Start
		pool["end"] = allocationPool.End

		result[i] = pool
	}

	return result
}

// flattenNetworkingSubnetV2HostRoutes allows to flatten slice of subnets.HostRoute structs into
// a slice of maps.
func flattenNetworkingSubnetV2HostRoutes(hostRoutes []subnets.HostRoute) []map[string]interface{} {
	result := make([]map[string]interface{}, len(hostRoutes))
	for i, hostRoute := range hostRoutes {
		route := make(map[string]interface{})
		route["destination_cidr"] = hostRoute.DestinationCIDR
		route["next_hop"] = hostRoute.NextHop

		result[i] = route
	}

	return result
}

func networkingSubnetV2AllocationPoolsMatch(oldPools, newPools []interface{}) bool {
	if len(oldPools) != len(newPools) {
		return false
	}

	for _, newPool := range newPools {
		var found bool

		newPoolPool := newPool.(map[string]interface{})
		newStart := newPoolPool["start"].(string)
		newEnd := newPoolPool["end"].(string)

		for _, oldPool := range oldPools {
			oldPoolPool := oldPool.(map[string]interface{})
			oldStart := oldPoolPool["start"].(string)
			oldEnd := oldPoolPool["end"].(string)

			if oldStart == newStart && oldEnd == newEnd {
				found = true
			}
		}

		if !found {
			return false
		}
	}

	return true
}

func networkingSubnetV2DNSNameserverAreUnique(raw []interface{}) error {
	set := make(map[string]struct{})
	for _, rawNS := range raw {
		nameserver, ok := rawNS.(string)
		if ok {
			if _, exists := set[nameserver]; exists {
				return fmt.Errorf("got duplicate nameserver %s", nameserver)
			}
			set[nameserver] = struct{}{}
		}
	}

	return nil
}
