package openstack

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/sharetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSharedFilesystemShareTypeAccessV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSharedFilesystemShareTypeAccessV2Create,
		ReadContext:   resourceSharedFilesystemShareTypeAccessV2Read,
		DeleteContext: resourceSharedFilesystemShareTypeAccessV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"share_type_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSharedFilesystemShareTypeAccessV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	sfsClient, err := config.SharedfilesystemV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	shareTypeID := d.Get("share_type_id").(string)
	projectID := d.Get("project_id").(string)

	accessOpts := sharetypes.AccessOpts{
		Project: projectID,
	}

	log.Printf("[DEBUG] openstack_sharedfilesystem_sharetype_access_v2 create options: share_type_id=%s, %#v", shareTypeID, accessOpts)

	if err := sharetypes.AddAccess(ctx, sfsClient, shareTypeID, accessOpts).ExtractErr(); err != nil {
		return diag.Errorf("Error creating openstack_sharedfilesystem_sharetype_access_v2 for share type %s and project %s: %s", shareTypeID, projectID, err)
	}

	id := fmt.Sprintf("%s/%s", shareTypeID, projectID)
	d.SetId(id)

	return resourceSharedFilesystemShareTypeAccessV2Read(ctx, d, meta)
}

func resourceSharedFilesystemShareTypeAccessV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	sfsClient, err := config.SharedfilesystemV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	shareTypeID, projectID, err := parsePairedIDs(d.Id(), "openstack_sharedfilesystem_sharetype_access_v2")
	if err != nil {
		return diag.FromErr(err)
	}

	access, err := getSharedFilesystemShareTypeAccess(ctx, sfsClient, shareTypeID, projectID)
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting openstack_sharedfilesystem_sharetype_access_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_sharedfilesystem_sharetype_access_v2 %s: %#v", d.Id(), access)

	d.Set("region", GetRegion(d, config))
	d.Set("share_type_id", access.ShareTypeID)
	d.Set("project_id", access.ProjectID)

	return nil
}

func resourceSharedFilesystemShareTypeAccessV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	sfsClient, err := config.SharedfilesystemV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	shareTypeID, projectID, err := parsePairedIDs(d.Id(), "openstack_sharedfilesystem_sharetype_access_v2")
	if err != nil {
		return diag.FromErr(err)
	}

	// Make sure the access still exists before trying to remove it, so
	// that a resource which has already been removed out-of-band (or a
	// share type which no longer exists) does not surface as an error
	// during "terraform destroy".
	if _, err := getSharedFilesystemShareTypeAccess(ctx, sfsClient, shareTypeID, projectID); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_sharedfilesystem_sharetype_access_v2"))
	}

	removeOpts := sharetypes.AccessOpts{
		Project: projectID,
	}

	if err := sharetypes.RemoveAccess(ctx, sfsClient, shareTypeID, removeOpts).ExtractErr(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_sharedfilesystem_sharetype_access_v2"))
	}

	return nil
}

// getSharedFilesystemShareTypeAccess retrieves the ShareTypeAccess entry
// for the given share type / project pair. The Shared File Systems API
// does not expose a way to fetch a single access rule directly, so the
// full access list for the share type is retrieved and filtered locally.
//
// If no matching entry is found, a synthetic 404 error is returned so that
// callers can use CheckDeleted to treat the resource as deleted.
func getSharedFilesystemShareTypeAccess(ctx context.Context, client *gophercloud.ServiceClient, shareTypeID, projectID string) (*sharetypes.ShareTypeAccess, error) {
	accessList, err := sharetypes.ShowAccess(ctx, client, shareTypeID).Extract()
	if err != nil {
		return nil, err
	}

	for _, access := range accessList {
		if access.ProjectID == projectID {
			return &access, nil
		}
	}

	return nil, gophercloud.ErrUnexpectedResponseCode{Actual: http.StatusNotFound}
}
