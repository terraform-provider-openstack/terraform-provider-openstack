package openstack

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/flavors"
	"github.com/gophercloud/gophercloud/pagination"
)

func dataSourceLBFlavorV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLBFlavorV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"flavor_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ExactlyOneOf: []string{"name"},
			},

			"name": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Computed:     true,
				ExactlyOneOf: []string{"flavor_id"},
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"flavor_profile_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceLBFlavorV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	lbClient, err := config.LoadBalancerV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack loadbalancer client: %s", err)
	}

	var allFlavors []flavors.Flavor
	if v := d.Get("flavor_id").(string); v != "" {
		var flavor *flavors.Flavor
		flavor, err = flavors.Get(lbClient, v).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return diag.Errorf("No flavor found")
			}
			return diag.Errorf("Unable to retrieve Openstack %s loadbalancer flavor: %s", v, err)
		}

		allFlavors = append(allFlavors, *flavor)

		return diag.FromErr(dataSourceLBFlavorV2Attributes(d, lbClient, &allFlavors[0]))
	} else {
		var allPages pagination.Page
		allPages, err = flavors.List(lbClient, flavors.ListOpts{}).AllPages()
		if err != nil {
			return diag.Errorf("Unable to query OpenStack flavors: %s", err)
		}

		allFlavors, err = flavors.ExtractFlavors(allPages)
		if err != nil {
			return diag.Errorf("Unable to retrieve Openstack loadbalancer flavors: %s", err)
		}
	}

	// Loop through all flavors to find a more specific one
	if len(allFlavors) > 0 {
		var filteredFlavors []flavors.Flavor
		for _, flavor := range allFlavors {
			if v := d.Get("name").(string); v != "" {
				if flavor.Name != v {
					continue
				}

				filteredFlavors = append(filteredFlavors, flavor)
			}
		}

		allFlavors = filteredFlavors
	}

	if len(allFlavors) < 1 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allFlavors) > 1 {
		log.Printf("[DEBUG] Multiple results found: %#v", allFlavors)
		return diag.Errorf("Your query returned more than one result. " +
			"Please try a more specific search criteria")
	}

	return diag.FromErr(dataSourceLBFlavorV2Attributes(d, lbClient, &allFlavors[0]))
}

func dataSourceLBFlavorV2Attributes(d *schema.ResourceData, computeClient *gophercloud.ServiceClient, flavor *flavors.Flavor) error {
	log.Printf("[DEBUG] Retrieved openstack_loadbalancer_flavor_v2 %s: %#v", flavor.ID, flavor)

	d.SetId(flavor.ID)
	d.Set("name", flavor.Name)
	d.Set("description", flavor.Description)
	d.Set("flavor_id", flavor.ID)
	d.Set("flavor_profile_id", flavor.FlavorProfileId)
	d.Set("enabled", flavor.Enabled)

	return nil
}
