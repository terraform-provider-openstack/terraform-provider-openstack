package openstack

import (
	"fmt"
	"log"

	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceSharedFilesystemShareV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSharedFilesystemShareV2Read,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"status": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"is_public": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"share_proto": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"size": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},

			"export_locations": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"preferred": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"metadata": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},

			"availability_zone": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceSharedFilesystemShareV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	sfsClient, err := config.sharedfilesystemV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack sharedfilesystem sfsClient: %s", err)
	}

	sfsClient.Microversion = minManilaShareMicroversion

	isPublic := d.Get("is_public").(bool)
	listOpts := shares.ListOpts{
		Name:               d.Get("name").(string),
		DisplayDescription: d.Get("description").(string),
		ProjectID:          d.Get("project_id").(string),
		Status:             d.Get("status").(string),
		IsPublic:           &isPublic,
	}

	allPages, err := shares.ListDetail(sfsClient, listOpts).AllPages()
	if err != nil {
		return fmt.Errorf("Unable to query shares: %s", err)
	}

	allShares, err := shares.ExtractShares(allPages)
	if err != nil {
		return fmt.Errorf("Unable to retrieve shares: %s", err)
	}

	if len(allShares) < 1 {
		return fmt.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	var share shares.Share
	if len(allShares) > 1 {
		log.Printf("[DEBUG] Multiple results found: %#v", allShares)
		return fmt.Errorf("Your query returned more than one result. Please try a more " +
			"specific search criteria.")
	} else {
		share = allShares[0]
	}

	exportLocationsRaw, err := shares.GetExportLocations(sfsClient, share.ID).Extract()
	if err != nil {
		return fmt.Errorf("Failed to retrieve share's export_locations %s: %s", share.ID, err)
	}

	log.Printf("[DEBUG] Retrieved share's export_locations %s: %#v", share.ID, exportLocationsRaw)

	var exportLocations []map[string]string
	for _, v := range exportLocationsRaw {
		exportLocations = append(exportLocations, map[string]string{
			"path":      v.Path,
			"preferred": fmt.Sprint(v.Preferred),
		})
	}

	return dataSourceSharedFilesystemShareV2Attributes(d, &share, exportLocations, GetRegion(d, config))
}

func dataSourceSharedFilesystemShareV2Attributes(d *schema.ResourceData, share *shares.Share, exportLocations []map[string]string, region string) error {
	d.SetId(share.ID)
	d.Set("name", share.Name)
	d.Set("region", region)
	d.Set("project_id", share.ProjectID)
	d.Set("availability_zone", share.AvailabilityZone)
	d.Set("description", share.Description)
	d.Set("size", share.Size)
	d.Set("status", share.Status)
	d.Set("is_public", share.IsPublic)
	d.Set("share_proto", share.ShareProto)

	if err := d.Set("metadata", share.Metadata); err != nil {
		log.Printf("[DEBUG] Unable to set metadata for share %s: %s", share.ID, err)
	}

	if err := d.Set("export_locations", exportLocations); err != nil {
		log.Printf("[DEBUG] Unable to set export_locations for share %s: %s", share.ID, err)
	}

	return nil
}
