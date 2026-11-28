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

// TODO: implement Metadata in gophercloud's shares.GrantAccessOpts.
// One or more access rule metadata key and value pairs can be provided since
// Shared Filesystem API microversion 2.45.
// https://docs.openstack.org/api-ref/shared-file-system/#grant-access
type shareAccessV2GrantAccessOpts struct {
	// The access rule type that can be "ip", "cert", "user" or "cephx".
	AccessType string `json:"access_type"`
	// The value that defines the access that can be a valid format of IP, cert or user.
	AccessTo string `json:"access_to"`
	// The access level to the share is either "rw" or "ro".
	AccessLevel string `json:"access_level"`
	// One or more access rule metadata key and value pairs as a dictionary of strings.
	Metadata map[string]string `json:"metadata,omitempty"`
}

func (opts shareAccessV2GrantAccessOpts) ToGrantAccessMap() (map[string]any, error) {
	return gophercloud.BuildRequestBody(opts, "allow_access")
}

func shareAccessV2MetadataURL(c *gophercloud.ServiceClient, accessID string) string {
	return c.ServiceURL("share-access-rules", accessID, "metadata")
}

func shareAccessV2MetadatumURL(c *gophercloud.ServiceClient, accessID, key string) string {
	return c.ServiceURL("share-access-rules", accessID, "metadata", key)
}

// TODO: implement this in gophercloud's shareaccessrules package.
// https://docs.openstack.org/api-ref/shared-file-system/#share-access-rule-metadata-since-api-v2-45
func updateShareAccessV2Metadata(ctx context.Context, client *gophercloud.ServiceClient, accessID string, metadata map[string]string) (map[string]string, error) {
	b, err := gophercloud.BuildRequestBody(struct {
		Metadata map[string]string `json:"metadata"`
	}{Metadata: metadata}, "")
	if err != nil {
		return nil, err
	}

	var res struct {
		Metadata map[string]string `json:"metadata"`
	}

	resp, err := client.Put(ctx, shareAccessV2MetadataURL(client, accessID), b, &res, &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return nil, err
	}

	return res.Metadata, nil
}

// TODO: implement this in gophercloud's shareaccessrules package.
// https://docs.openstack.org/api-ref/shared-file-system/#share-access-rule-metadata-since-api-v2-45
func deleteShareAccessV2Metadatum(ctx context.Context, client *gophercloud.ServiceClient, accessID, key string) error {
	resp, err := client.Delete(ctx, shareAccessV2MetadatumURL(client, accessID, key), &gophercloud.RequestOpts{
		OkCodes: []int{200},
	})
	if resp != nil {
		defer resp.Body.Close()
	}

	return err
}

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
