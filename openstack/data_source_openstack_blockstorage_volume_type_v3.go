package openstack

import (
	"context"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceBlockStorageVolumeTypeV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBlockStorageVolumeTypeV3Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"is_public": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"qos_specs_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"extra_specs": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"public_access": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceBlockStorageVolumeTypeV3Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	client, err := config.BlockStorageV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	listOpts := volumetypes.ListOpts{}

	if v, ok := d.GetOk("is_public"); ok {
		if v.(bool) {
			listOpts.IsPublic = volumetypes.VisibilityPublic
		} else {
			listOpts.IsPublic = volumetypes.VisibilityPrivate
		}
	} else {
		listOpts.IsPublic = volumetypes.VisibilityDefault
	}

	allPages, err := volumetypes.List(client, listOpts).AllPages(ctx)
	if err != nil {
		return diag.Errorf("Unable to query openstack_blockstorage_volume_type_v3: %s", err)
	}

	var allVolumeTypes []volumetypes.VolumeType

	err = volumetypes.ExtractVolumeTypesInto(allPages, &allVolumeTypes)
	if err != nil {
		return diag.Errorf("Unable to retrieve openstack_blockstorage_volume_v3: %s", err)
	}

	if name, ok := d.GetOk("name"); ok {
		filtered := make([]volumetypes.VolumeType, 0, len(allVolumeTypes))

		for _, vt := range allVolumeTypes {
			if vt.Name == name.(string) {
				filtered = append(filtered, vt)
			}
		}

		allVolumeTypes = filtered
	}

	if len(allVolumeTypes) > 1 {
		return diag.Errorf("Your openstack_blockstorage_volume_type_v3 query returned multiple results")
	}

	if len(allVolumeTypes) < 1 {
		return diag.Errorf("Your openstack_blockstorage_volume_type_v3 query returned no results")
	}

	dataSourceBlockStorageVolumeTypeV3Attributes(d, allVolumeTypes[0])
	d.Set("region", GetRegion(d, config))

	return nil
}

func dataSourceBlockStorageVolumeTypeV3Attributes(d *schema.ResourceData, volumetype volumetypes.VolumeType) {
	d.SetId(volumetype.ID)
	d.Set("name", volumetype.Name)
	d.Set("description", volumetype.Description)
	d.Set("is_public", volumetype.IsPublic)
	d.Set("qos_specs_id", volumetype.QosSpecID)
	d.Set("public_access", volumetype.PublicAccess)

	if volumetype.ExtraSpecs != nil {
		d.Set("extra_specs", volumetype.ExtraSpecs)
	} else {
		d.Set("extra_specs", map[string]string{})
	}
}
