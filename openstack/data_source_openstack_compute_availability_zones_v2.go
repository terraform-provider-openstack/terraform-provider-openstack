package openstack

import (
	"context"
	"sort"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/availabilityzones"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/terraform-provider-openstack/utils/v2/hashcode"
)

func dataSourceComputeAvailabilityZonesV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComputeAvailabilityZonesV2Read,
		Schema: map[string]*schema.Schema{
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"region": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"state": {
				Type:         schema.TypeString,
				Default:      "available",
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"available", "unavailable"}, true),
			},
		},
	}
}

func dataSourceComputeAvailabilityZonesV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)
	region := GetRegion(d, config)

	computeClient, err := config.ComputeV2Client(ctx, region)
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	allPages, err := availabilityzones.List(computeClient).AllPages(ctx)
	if err != nil {
		return diag.Errorf("Error retrieving openstack_compute_availability_zones_v2: %s", err)
	}

	zoneInfo, err := availabilityzones.ExtractAvailabilityZones(allPages)
	if err != nil {
		return diag.Errorf("Error extracting openstack_compute_availability_zones_v2 from response: %s", err)
	}

	stateBool := d.Get("state").(string) == "available"
	zones := make([]string, 0, len(zoneInfo))

	for _, z := range zoneInfo {
		if z.ZoneState.Available == stateBool {
			zones = append(zones, z.ZoneName)
		}
	}

	// sort.Strings sorts in place, returns nothing
	sort.Strings(zones)

	d.SetId(hashcode.Strings(zones))
	d.Set("names", zones)
	d.Set("region", region)

	return nil
}
