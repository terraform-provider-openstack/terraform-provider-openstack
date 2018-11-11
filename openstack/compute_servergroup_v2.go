package openstack

import (
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/servergroups"
)

// ServerGroupCreateOpts is a custom ServerGroup struct to include the
// ValueSpecs field.
type ComputeServerGroupV2CreateOpts struct {
	servergroups.CreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}

// ToServerGroupCreateMap casts a CreateOpts struct to a map.
// It overrides routers.ToServerGroupCreateMap to add the ValueSpecs field.
func (opts ComputeServerGroupV2CreateOpts) ToServerGroupCreateMap() (map[string]interface{}, error) {
	return BuildRequest(opts, "server_group")
}

func expandComputeServerGroupV2Policies(raw []interface{}) []string {
	policies := make([]string, len(raw))
	for i, v := range raw {
		policies[i] = v.(string)
	}

	return policies
}
