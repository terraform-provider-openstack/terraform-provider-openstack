package openstack

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/apiversions"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/shareaccessrules"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/shares"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func sharedFilesystemShareAccessV2StateRefreshFunc(ctx context.Context, client *gophercloud.ServiceClient, shareID string, accessID string) retry.StateRefreshFunc {
	// Set the client to the minimum supported microversion.
	client.Microversion = sharedFilesystemV2MinMicroversion

	// Obtain supported Manila microversions.
	apiInfo, err := apiversions.Get(ctx, client, "v2").Extract()
	if err != nil {
		return func() (any, string, error) {
			return nil, "", fmt.Errorf("Unable to query API endpoint for openstack_sharedfilesystem_share_access_v2: %w", err)
		}
	}

	// Check for newer microversion 2.45 API to get access rules using GET method.
	if ok, err := compatibleMicroversion("min", sharedFilesystemV2ShareAccessRulesMicroversion, apiInfo.Version); err != nil {
		return func() (any, string, error) {
			return nil, "", fmt.Errorf("Error comparing microversions for openstack_sharedfilesystem_share_access_v2 %s: %w", accessID, err)
		}
	} else if ok {
		client.Microversion = sharedFilesystemV2ShareAccessRulesMicroversion

		return sharedFilesystemShareAccessV2StateRefreshStateNew(ctx, client, accessID)
	}

	// Now check and see if the OpenStack environment supports microversion 2.21.
	// If so, use that for the API request for access_key support.
	if ok, err := compatibleMicroversion("min", sharedFilesystemV2SharedAccessMinMicroversion, apiInfo.Version); err != nil {
		return func() (any, string, error) {
			return nil, "", fmt.Errorf("Error comparing microversions for openstack_sharedfilesystem_share_access_v2 %s: %w", accessID, err)
		}
	} else if ok {
		client.Microversion = sharedFilesystemV2SharedAccessMinMicroversion
	}

	return sharedFilesystemShareAccessV2StateRefreshStateOld(ctx, client, shareID, accessID)
}

func sharedFilesystemShareAccessV2StateRefreshStateOld(ctx context.Context, client *gophercloud.ServiceClient, shareID string, accessID string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		access, err := shares.ListAccessRights(ctx, client, shareID).Extract()
		if err != nil {
			return nil, "", err
		}

		for _, v := range access {
			if v.ID == accessID {
				return v, v.State, nil
			}
		}

		return nil, "", gophercloud.ErrUnexpectedResponseCode{Actual: http.StatusNotFound}
	}
}

func sharedFilesystemShareAccessV2StateRefreshStateNew(ctx context.Context, client *gophercloud.ServiceClient, accessID string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		access, err := shareaccessrules.Get(ctx, client, accessID).Extract()
		if err != nil {
			return nil, "", err
		}

		return *access, access.State, nil
	}
}
