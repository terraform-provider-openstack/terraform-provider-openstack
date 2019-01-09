package openstack

import (
	"fmt"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/hashicorp/terraform/helper/schema"
)

// networkingNetworkV2ID retrieves network ID by the provided name.
func networkingNetworkV2ID(d *schema.ResourceData, meta interface{}, networkName string) (string, error) {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return "", fmt.Errorf("Error creating OpenStack network client: %s", err)
	}

	opts := networks.ListOpts{Name: networkName}
	pager := networks.List(networkingClient, opts)
	networkID := ""

	err = pager.EachPage(func(page pagination.Page) (bool, error) {
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
func networkingNetworkV2Name(d *schema.ResourceData, meta interface{}, networkID string) (string, error) {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return "", fmt.Errorf("Error creating OpenStack network client: %s", err)
	}

	opts := networks.ListOpts{ID: networkID}
	pager := networks.List(networkingClient, opts)
	networkName := ""

	err = pager.EachPage(func(page pagination.Page) (bool, error) {
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
