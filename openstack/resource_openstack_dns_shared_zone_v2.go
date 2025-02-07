package openstack

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/zones"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// resourceDnsZoneShare creates a resource for sharing DNS zones.
func resourceDnsZoneShare() *schema.Resource {
	return &schema.Resource{
		Create: resourceDnsZoneShareCreate,
		Read:   resourceDnsZoneShareRead,
		Delete: resourceDnsZoneShareDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"target_project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceDnsZoneShareCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gophercloud.ServiceClient)
	zoneID := d.Get("zone_id").(string)
	targetProjectID := d.Get("target_project_id").(string)

	shareOpts := zones.ShareZoneOpts{
		TargetProjectID: targetProjectID,
	}

	err := zones.Share(context.Background(), client, zoneID, shareOpts).ExtractErr()
	if err != nil {
		return fmt.Errorf("error sharing DNS zone %s: %s", zoneID, err)
	}

	// Use a composite ID (zoneID:targetProjectID)
	d.SetId(fmt.Sprintf("%s:%s", zoneID, targetProjectID))
	return resourceDnsZoneShareRead(d, meta)
}

func resourceDnsZoneShareRead(d *schema.ResourceData, meta interface{}) error {
	// There is no dedicated API to read share details.
	// We assume that if creation succeeded, the share exists.
	return nil
}

func resourceDnsZoneShareDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gophercloud.ServiceClient)
	zoneID := d.Get("zone_id").(string)
	shareID := d.Get("target_project_id").(string)

	err := zones.Unshare(context.Background(), client, zoneID, shareID).ExtractErr()
	if err != nil {
		return fmt.Errorf("error unsharing DNS zone %s: %s", zoneID, err)
	}
	d.SetId("")
	return nil
}
