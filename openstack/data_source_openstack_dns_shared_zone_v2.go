package openstack

import (
	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/zones"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDNSSharedZoneV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDNSSharedZoneV2Read,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter results by zone ID.",
			},
			"shared_zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDNSSharedZoneV2Read(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gophercloud.ServiceClient)

	var listOpts zones.ListOpts

	if v, ok := d.GetOk("zone_id"); ok {
		listOpts.ID = v.(string)
	}

	allPages, err := zones.List(client, listOpts).AllPages()
	if err != nil {
		return err
	}

	allZones, err := zones.ExtractZones(allPages)
	if err != nil {
		return err
	}

	var sharedZones []map[string]interface{}
	for _, zone := range allZones {
		sharedZones = append(sharedZones, map[string]interface{}{
			"zone_id":    zone.ID,
			"project_id": zone.ProjectID,
		})
	}

	if err := d.Set("shared_zones", sharedZones); err != nil {
		return err
	}

	d.SetId("shared-zones-list")
	return nil
}
