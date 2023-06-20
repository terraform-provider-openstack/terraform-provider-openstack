package openstack

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/apiversions"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shareaccessrules"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
)

func sharedFilesystemShareAccessV2StateRefreshFunc(client *gophercloud.ServiceClient, shareID string, accessID string) resource.StateRefreshFunc {
	// Set the client to the minimum supported microversion.
	client.Microversion = sharedFilesystemV2MinMicroversion

	// Obtain supported Manila microversions.
	apiInfo, err := apiversions.Get(client, "v2").Extract()
	if err != nil {
		return func() (interface{}, string, error) {
			return nil, "", fmt.Errorf("Unable to query API endpoint for openstack_sharedfilesystem_share_access_v2: %s", err)
		}
	}

	// Check for newer microversion 2.45 API to get access rules using GET method.
	if ok, err := compatibleMicroversion("min", sharedFilesystemV2ShareAccessRulesMicroversion, apiInfo.Version); err != nil {
		return func() (interface{}, string, error) {
			return nil, "", fmt.Errorf("Error comparing microversions for openstack_sharedfilesystem_share_access_v2 %s: %s", accessID, err)
		}
	} else if ok {
		client.Microversion = sharedFilesystemV2ShareAccessRulesMicroversion
		return sharedFilesystemShareAccessV2StateRefreshStateNew(client, accessID)
	}

	// Now check and see if the OpenStack environment supports microversion 2.21.
	// If so, use that for the API request for access_key support.
	if ok, err := compatibleMicroversion("min", sharedFilesystemV2SharedAccessMinMicroversion, apiInfo.Version); err != nil {
		return func() (interface{}, string, error) {
			return nil, "", fmt.Errorf("Error comparing microversions for openstack_sharedfilesystem_share_access_v2 %s: %s", accessID, err)
		}
	} else if ok {
		client.Microversion = sharedFilesystemV2SharedAccessMinMicroversion
	}

	return sharedFilesystemShareAccessV2StateRefreshStateOld(client, shareID, accessID)
}

func sharedFilesystemShareAccessV2StateRefreshStateOld(client *gophercloud.ServiceClient, shareID string, accessID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		access, err := shares.ListAccessRights(client, shareID).Extract()
		if err != nil {
			return nil, "", err
		}
		for _, v := range access {
			if v.ID == accessID {
				return v, v.State, nil
			}
		}
		return nil, "", gophercloud.ErrDefault404{}
	}
}

func sharedFilesystemShareAccessV2StateRefreshStateNew(client *gophercloud.ServiceClient, accessID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		access, err := shareaccessrules.Get(client, accessID).Extract()
		if err != nil {
			return nil, "", err
		}
		return *access, access.State, nil
	}
}
