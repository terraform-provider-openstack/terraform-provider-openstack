package openstack

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/zones"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDNSSharedZoneV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDNSSharedZoneV2Read,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			// Additional schema fields can be defined here.
		},
	}
}

func dataSourceDNSSharedZoneV2Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gophercloud.ServiceClient)
	zoneID := d.Get("zone_id").(string)

	// Create ListOpts without the undefined field.
	listOpts := zones.ListOpts{}

	pages, err := zones.List(client, listOpts).AllPages(context.Background())
	if err != nil {
		return fmt.Errorf("error listing DNS zones: %s", err)
	}
	allZones, err := zones.ExtractZones(pages)
	if err != nil {
		return fmt.Errorf("error extracting zones: %s", err)
	}
	for _, z := range allZones {
		if z.ID == zoneID {
			d.SetId(z.ID)
			// Set additional fields as needed.
			return nil
		}
	}
	return fmt.Errorf("DNS zone %s not found", zoneID)
}
