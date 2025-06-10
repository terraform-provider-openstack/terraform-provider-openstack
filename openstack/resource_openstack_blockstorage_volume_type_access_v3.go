package openstack

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBlockstorageVolumeTypeAccessV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBlockstorageVolumeTypeAccessV3Create,
		ReadContext:   resourceBlockstorageVolumeTypeAccessV3Read,
		DeleteContext: resourceBlockstorageVolumeTypeAccessV3Delete,
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
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"volume_type_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceBlockstorageVolumeTypeAccessV3Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	blockStorageClient, err := config.BlockStorageV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	projectID := d.Get("project_id").(string)
	vtID := d.Get("volume_type_id").(string)

	accessOpts := volumetypes.AddAccessOpts{
		Project: projectID,
	}

	if err := volumetypes.AddAccess(ctx, blockStorageClient, vtID, accessOpts).ExtractErr(); err != nil {
		return diag.Errorf("Error creating openstack_blockstorage_volume_type_access_v3: %s", err)
	}

	id := fmt.Sprintf("%s/%s", vtID, projectID)
	d.SetId(id)

	return resourceBlockstorageVolumeTypeAccessV3Read(ctx, d, meta)
}

func resourceBlockstorageVolumeTypeAccessV3Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	blockStorageClient, err := config.BlockStorageV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	vtID, projectID, err := parsePairedIDs(d.Id(), "openstack_blockstorage_volume_type_access_v3")
	if err != nil {
		return diag.FromErr(err)
	}

	allPages, err := volumetypes.ListAccesses(blockStorageClient, vtID).AllPages(ctx)
	if err != nil {
		return diag.Errorf("Error retrieving accesses openstack_blockstorage_volume_type_access_v3 for vt: %s", vtID)
	}

	allAccesses, err := volumetypes.ExtractAccesses(allPages)
	if err != nil {
		return diag.Errorf("Error extracting accesses openstack_blockstorage_volume_type_access_v3 for vt: %s", vtID)
	}

	found := false

	for _, access := range allAccesses {
		if access.VolumeTypeID == vtID && access.ProjectID == projectID {
			found = true

			break
		}
	}

	if !found {
		return diag.Errorf("Error getting volume type access openstack_blockstorage_volume_type_access_v3 for vt: %s", vtID)
	}

	d.Set("region", GetRegion(d, config))
	d.Set("project_id", projectID)
	d.Set("volume_type_id", vtID)

	return nil
}

func resourceBlockstorageVolumeTypeAccessV3Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	blockStorageClient, err := config.BlockStorageV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	vtID, projectID, err := parsePairedIDs(d.Id(), "openstack_blockstorage_volume_type_access_v3")
	if err != nil {
		return diag.FromErr(err)
	}

	removeOpts := volumetypes.RemoveAccessOpts{
		Project: projectID,
	}

	if err := volumetypes.RemoveAccess(ctx, blockStorageClient, vtID, removeOpts).ExtractErr(); err != nil {
		return diag.Errorf("Error removing openstack_blockstorage_volume_type_access_v3 %s: %s", d.Id(), err)
	}

	return nil
}
