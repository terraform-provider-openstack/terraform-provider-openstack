package openstack

import (
	"context"
	"log"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/flavors"
	"github.com/gophercloud/gophercloud/v2/pagination"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLBFlavorV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLBFlavorV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

func dataSourceLBFlavorV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack loadbalancer client: %s", err)
	}

	var allFlavors []flavors.Flavor

	if v := d.Get("flavor_id").(string); v != "" {
		var flavor *flavors.Flavor

		flavor, err = flavors.Get(ctx, lbClient, v).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return diag.Errorf("No flavor found")
			}

			return diag.Errorf("Unable to retrieve Openstack %s loadbalancer flavor: %s", v, err)
		}

		dataSourceLBFlavorV2Attributes(d, flavor)
		d.Set("region", GetRegion(d, config))

		return nil
	}

	var allPages pagination.Page

	allPages, err = flavors.List(lbClient, flavors.ListOpts{}).AllPages(ctx)
	if err != nil {
		return diag.Errorf("Unable to query OpenStack flavors: %s", err)
	}

	allFlavors, err = flavors.ExtractFlavors(allPages)
	if err != nil {
		return diag.Errorf("Unable to retrieve Openstack loadbalancer flavors: %s", err)
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

	dataSourceLBFlavorV2Attributes(d, &allFlavors[0])
	d.Set("region", GetRegion(d, config))

	return nil
}

func dataSourceLBFlavorV2Attributes(d *schema.ResourceData, flavor *flavors.Flavor) {
	log.Printf("[DEBUG] Retrieved openstack_lb_flavor_v2 %s: %#v", flavor.ID, flavor)

	d.SetId(flavor.ID)
	d.Set("name", flavor.Name)
	d.Set("description", flavor.Description)
	d.Set("flavor_id", flavor.ID)
	d.Set("flavor_profile_id", flavor.FlavorProfileId)
	d.Set("enabled", flavor.Enabled)
}
