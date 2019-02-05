package openstack

import (
	"fmt"
	"time"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/availabilityzones"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func dataSourceComputeAvailabilityZonesV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceComputeAvailabilityZonesV2Read,
		Schema: map[string]*schema.Schema{
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"state": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"available", "unavailable"}, true),
			},
		},
	}
}

func dataSourceComputeAvailabilityZonesV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	region := d.Get("region").(string)
	if region == "" {
		region = GetRegion(d, config)
	}
	computeClient, err := config.computeV2Client(region)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err.Error())
	}

	allPages, err := availabilityzones.List(computeClient).AllPages()
	if err != nil {
		return fmt.Errorf("Error retrieving openstack_compute_availability_zones_v2: %s", err.Error())
	}
	zoneInfo, err := availabilityzones.ExtractAvailabilityZones(allPages)
	if err != nil {
		return fmt.Errorf("Error extracting openstack_compute_availability_zones_v2 from response: %s", err.Error())
	}

	state := d.Get("state").(string)
	if state == "" {
		state = "available"
	}
	stateBool := state == "available"
	zones := make([]string, 0, len(zoneInfo))
	for _, z := range zoneInfo {
		if z.ZoneState.Available == stateBool {
			zones = append(zones, z.ZoneName)
		}
	}

	d.SetId(time.Now().UTC().String())
	d.Set("names", zones)

	return nil
}
