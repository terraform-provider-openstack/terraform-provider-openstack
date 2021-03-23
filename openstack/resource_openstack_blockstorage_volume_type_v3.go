package openstack

import (
	"fmt"
	"log"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumetypes"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBlockStorageVolumeTypeV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceBlockStorageVolumeTypeV3Create,
		Read:   resourceBlockStorageVolumeTypeV3Read,
		Update: resourceBlockStorageVolumeTypeV3Update,
		Delete: resourceBlockStorageVolumeTypeV3Delete,
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

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"is_public": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"extra_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceBlockStorageVolumeTypeV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	blockStorageClient, err := config.BlockStorageV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	isPublic := d.Get("is_public").(bool)
	extraSpecs := d.Get("extra_specs").(map[string]interface{})
	createOpts := volumetypes.CreateOpts{
		Name:        name,
		Description: description,
		IsPublic:    &isPublic,
		ExtraSpecs:  expandToMapStringString(extraSpecs),
	}

	log.Printf("[DEBUG] openstack_blockstorage_volume_type_v3 create options: %#v", createOpts)
	vt, err := volumetypes.Create(blockStorageClient, &createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating openstack_blockstorage_volume_type_v3 %s: %s", name, err)
	}

	d.SetId(vt.ID)

	return resourceBlockStorageVolumeTypeV3Read(d, meta)
}

func resourceBlockStorageVolumeTypeV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	blockStorageClient, err := config.BlockStorageV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	vt, err := volumetypes.Get(blockStorageClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Error retrieving openstack_blockstorage_volume_type_v3")
	}

	log.Printf("[DEBUG] Retrieved openstack_blockstorage_volume_type_v3 %s: %#v", d.Id(), vt)

	d.Set("name", vt.Name)
	d.Set("description", vt.Description)
	d.Set("is_public", vt.IsPublic)
	d.Set("region", GetRegion(d, config))

	es, err := volumetypes.ListExtraSpecs(blockStorageClient, d.Id()).Extract()
	if err != nil {
		return fmt.Errorf("Error reading extra_specs for openstack_blockstorage_volume_type_v3 %s: %s", d.Id(), err)
	}

	if err := d.Set("extra_specs", es); err != nil {
		log.Printf("[WARN] Unable to set extra_specs for openstack_blockstorage_volume_type_v3 %s: %s", d.Id(), err)
	}

	return nil
}

func resourceBlockStorageVolumeTypeV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	blockStorageClient, err := config.BlockStorageV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	hasChange := false
	var updateOpts volumetypes.UpdateOpts

	if d.HasChange("name") {
		hasChange = true
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}

	if d.HasChange("description") {
		hasChange = true
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}

	if d.HasChange("is_public") {
		hasChange = true
		isPublic := d.Get("is_public").(bool)
		updateOpts.IsPublic = &isPublic
	}

	if hasChange {
		_, err = volumetypes.Update(blockStorageClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating openstack_blockstorage_volume_type_v3 %s: %s", d.Id(), err)
		}
	}

	if d.HasChange("extra_specs") {
		oldES, newES := d.GetChange("extra_specs")

		// Delete all old extra specs.
		for oldKey := range oldES.(map[string]interface{}) {
			if err := volumetypes.DeleteExtraSpec(blockStorageClient, d.Id(), oldKey).ExtractErr(); err != nil {
				return fmt.Errorf("Error deleting extra_spec %s from openstack_blockstorage_volume_type_v3 %s: %s", oldKey, d.Id(), err)
			}
		}

		// Add new extra specs.
		newESRaw := newES.(map[string]interface{})
		if len(newESRaw) > 0 {
			extraSpecs := expandBlockStorageVolumeTypeV3ExtraSpecs(newESRaw)

			_, err := volumetypes.CreateExtraSpecs(blockStorageClient, d.Id(), extraSpecs).Extract()
			if err != nil {
				return fmt.Errorf("Error creating extra_specs for openstack_blockstorage_volume_type_v3 %s: %s", d.Id(), err)
			}
		}
	}

	return resourceBlockStorageVolumeTypeV3Read(d, meta)
}

func resourceBlockStorageVolumeTypeV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	blockStorageClient, err := config.BlockStorageV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	err = volumetypes.Delete(blockStorageClient, d.Id()).ExtractErr()
	if err != nil {
		return CheckDeleted(d, err, "Error deleting openstack_blockstorage_volume_type_v3")
	}

	return nil
}
