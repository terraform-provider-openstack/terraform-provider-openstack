package openstack

import (
	"fmt"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/groups"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// GroupCreateOpts represents the attributes used when creating a new firewall group.
type GroupCreateOpts struct {
	groups.CreateOpts
	ValueSpecs map[string]string `json:"value_specs,omitempty"`
}

func fwGroupV2RefreshFunc(networkingClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		var group groups.Group

		err := groups.Get(networkingClient, id).ExtractIntoStructPtr(&group, "firewall_group")
		if err != nil {
			return nil, "", err
		}

		return group, group.Status, nil
	}
}

func fwGroupV2DeleteFunc(networkingClient *gophercloud.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		group, err := groups.Get(networkingClient, id).Extract()

		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return "", "DELETED", nil
			}
			return nil, "", fmt.Errorf("Unexpected error: %s", err)
		}

		return group, "DELETING", nil
	}
}
