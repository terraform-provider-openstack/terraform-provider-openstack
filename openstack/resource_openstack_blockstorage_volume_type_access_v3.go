package openstack

import (
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumetypes"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBlockstorageVolumeTypeAccessV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceBlockstorageVolumeTypeAccessV3Create,
		Read:   resourceBlockstorageVolumeTypeAccessV3Read,
		Delete: resourceBlockstorageVolumeTypeAccessV3Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

func resourceBlockstorageVolumeTypeAccessV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	blockStorageClient, err := config.BlockStorageV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	projectID := d.Get("project_id").(string)
	vtID := d.Get("volume_type_id").(string)

	accessOpts := volumetypes.AddAccessOpts{
		Project: projectID,
	}

	if err := volumetypes.AddAccess(blockStorageClient, vtID, accessOpts).ExtractErr(); err != nil {
		return fmt.Errorf("Error creating openstack_blockstorage_volume_type_access_v3: %s", err)
	}

	id := fmt.Sprintf("%s/%s", vtID, projectID)
	d.SetId(id)

	return resourceBlockstorageVolumeTypeAccessV3Read(d, meta)
}

func resourceBlockstorageVolumeTypeAccessV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	blockStorageClient, err := config.BlockStorageV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	vtID, projectID, err := parseVolumeTypeAccessID(d.Id())
	if err != nil {
		return CheckDeleted(d, err, "Error parsing ID of openstack_blockstorage_volume_type_access_v3")
	}

	allPages, err := volumetypes.ListAccesses(blockStorageClient, vtID).AllPages()
	if err != nil {
		return fmt.Errorf("Error retrieving accesses openstack_blockstorage_volume_type_access_v3 for vt: %s", vtID)
	}

	allAccesses, err := volumetypes.ExtractAccesses(allPages)
	if err != nil {
		return fmt.Errorf("Error extracting accesses openstack_blockstorage_volume_type_access_v3 for vt: %s", vtID)
	}

	found := false
	for _, access := range allAccesses {
		if access.VolumeTypeID == vtID && access.ProjectID == projectID {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("Error getting volume type access openstack_blockstorage_volume_type_access_v3 for vt: %s", vtID)
	}

	d.Set("region", GetRegion(d, config))
	d.Set("project_id", projectID)
	d.Set("volume_type_id", vtID)

	return nil
}

func resourceBlockstorageVolumeTypeAccessV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	blockStorageClient, err := config.BlockStorageV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	vtID, projectID, err := parseVolumeTypeAccessID(d.Id())
	if err != nil {
		return fmt.Errorf("Error parsing ID of openstack_blockstorage_volume_type_access_v3 %s: %s", d.Id(), err)
	}

	removeOpts := volumetypes.RemoveAccessOpts{
		Project: projectID,
	}

	if err := volumetypes.RemoveAccess(blockStorageClient, vtID, removeOpts).ExtractErr(); err != nil {
		return fmt.Errorf("Error removing openstack_blockstorage_volume_type_access_v3 %s: %s", d.Id(), err)
	}

	return nil
}

func parseVolumeTypeAccessID(id string) (string, string, error) {
	idParts := strings.Split(id, "/")
	if len(idParts) < 2 {
		return "", "", fmt.Errorf("Unable to determine volume type access ID %s", id)
	}

	return idParts[0], idParts[1], nil
}
