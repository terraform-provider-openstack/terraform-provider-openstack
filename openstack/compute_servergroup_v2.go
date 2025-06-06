package openstack

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servergroups"
)

const (
	antiAffinityPolicy = "anti-affinity"
	affinityPolicy     = "affinity"
)

// ServerGroupCreateOpts is a custom ServerGroup struct to include the
// ValueSpecs field.
type ComputeServerGroupV2CreateOpts struct {
	servergroups.CreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}

// ToServerGroupCreateMap casts a CreateOpts struct to a map.
// It overrides routers.ToServerGroupCreateMap to add the ValueSpecs field.
func (opts ComputeServerGroupV2CreateOpts) ToServerGroupCreateMap() (map[string]any, error) {
	return BuildRequest(opts, "server_group")
}

func expandComputeServerGroupV2Policies(client *gophercloud.ServiceClient, raw []any) []string {
	policies := make([]string, len(raw))

	for i, v := range raw {
		client.Microversion = "2.15"

		policy := v.(string)
		policies[i] = policy

		// Set microversion for legacy policies to empty to not change behaviour for those policies
		if policy == antiAffinityPolicy || policy == affinityPolicy {
			client.Microversion = ""
		}
	}

	return policies
}

func expandComputeServerGroupV2RulesMaxServerPerHost(raw []any) int {
	for _, raw := range raw {
		raw, ok := raw.(map[string]any)
		if !ok {
			return 0
		}

		v, ok := raw["max_server_per_host"].(int)
		if !ok {
			return 0
		}
		//nolint:staticcheck // we need the first element
		return v
	}

	return 0
}
