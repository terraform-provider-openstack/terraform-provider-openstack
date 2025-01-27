package openstack

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/zones"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDNSSharedZoneV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceDNSSharedZoneV2Create,
		Read:   resourceDNSSharedZoneV2Read,
		Delete: resourceDNSSharedZoneV2Delete,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the zone to be shared.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the project with which the zone is shared.",
			},
		},
	}
}

func resourceDNSSharedZoneV2Create(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gophercloud.ServiceClient)
	zoneID := d.Get("zone_id").(string)
	projectID := d.Get("project_id").(string)

	shareOpts := zones.ShareZoneOpts{
		TargetProjectID: projectID,
	}

	err := zones.Share(client, zoneID, shareOpts).ExtractErr()
	if err != nil {
		return err
	}

	d.SetId(zoneID)
	return resourceDNSSharedZoneV2Read(d, meta)
}

func resourceDNSSharedZoneV2Read(d *schema.ResourceData, meta interface{}) error {
	// Fetch shared zone details using the OpenStack API.
	return nil
}

func resourceDNSSharedZoneV2Delete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gophercloud.ServiceClient)
	zoneID := d.Get("zone_id").(string)

	err := zones.Unshare(client, zoneID).ExtractErr()
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}
