package openstack

import (
	"context"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/flavorprofiles"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

func dataSourceLBFlavorProfileV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLBFlavorProfileV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"id": {
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
				ExactlyOneOf: []string{"id"},
			},

			"provider_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"flavor_data": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceLBFlavorProfileV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack loadbalancer client: %s", err)
	}

	var allfps []flavorprofiles.FlavorProfile
	if v := d.Get("id").(string); v != "" {
		var fp *flavorprofiles.FlavorProfile
		fp, err = flavorprofiles.Get(ctx, lbClient, v).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return diag.Errorf("No flavor profile found")
			}
			return diag.Errorf("Unable to retrieve Openstack %s loadbalancer flavor: %s", v, err)
		}

		allfps = append(allfps, *fp)

		return diag.FromErr(dataSourceLBFlavorProfileV2Attributes(d, lbClient, &allfps[0]))
	}

	var allPages pagination.Page
	allPages, err = flavorprofiles.List(lbClient, flavorprofiles.ListOpts{}).AllPages(ctx)
	if err != nil {
		return diag.Errorf("Unable to query OpenStack flavors: %s", err)
	}

	allfps, err = flavorprofiles.ExtractFlavorProfiles(allPages)
	if err != nil {
		return diag.Errorf("Unable to retrieve Openstack loadbalancer flavors: %s", err)
	}

	// Loop through all flavorprofiles to find a more specific one
	if len(allfps) > 0 {
		var filteredfps []flavorprofiles.FlavorProfile
		name := d.Get("name").(string)
		for _, fp := range allfps {
			if fp.Name != name {
				continue
			}

			filteredfps = append(filteredfps, fp)
		}

		allfps = filteredfps
	}

	if len(allfps) < 1 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allfps) > 1 {
		log.Printf("[DEBUG] Multiple results found: %#v", allfps)
		return diag.Errorf("Your query returned more than one result. " +
			"Please try a more specific search criteria")
	}

	return diag.FromErr(dataSourceLBFlavorProfileV2Attributes(d, lbClient, &allfps[0]))
}

func dataSourceLBFlavorProfileV2Attributes(d *schema.ResourceData, computeClient *gophercloud.ServiceClient, fp *flavorprofiles.FlavorProfile) error {
	log.Printf("[DEBUG] Retrieved openstack_loadbalancer_flavorprofile_v2 %s: %#v", fp.ID, fp)

	d.SetId(fp.ID)
	d.Set("name", fp.Name)
	d.Set("provider_name", fp.ProviderName)
	d.Set("flavor_data", fp.FlavorData)

	return nil
}
