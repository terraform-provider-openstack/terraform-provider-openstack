package openstack

import (
	"context"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/flavorprofiles"
)

func dataSourceLBFlavorProfileV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceLBFlavorProfileV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"flavorprofile_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"name", "provider_name"},
			},

			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"flavorprofile_id"},
			},

			"provider_name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"flavorprofile_id"},
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

	if id := d.Get("flavorprofile_id").(string); id != "" {
		fp, err := flavorprofiles.Get(ctx, lbClient, id).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return diag.Errorf("No flavor profile found")
			}
			return diag.Errorf("Unable to retrieve OpenStack %s loadbalancer flavor: %s", id, err)
		}

		return dataSourceLBFlavorProfileV2Attributes(d, lbClient, fp)
	}

	opts := flavorprofiles.ListOpts{
		Name:         d.Get("name").(string),
		ProviderName: d.Get("provider_name").(string),
	}
	allPages, err := flavorprofiles.List(lbClient, opts).AllPages(ctx)
	if err != nil {
		return diag.Errorf("Unable to query OpenStack flavors: %s", err)
	}

	allfps, err := flavorprofiles.ExtractFlavorProfiles(allPages)
	if err != nil {
		return diag.Errorf("Unable to retrieve OpenStack loadbalancer flavors: %s", err)
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

	return dataSourceLBFlavorProfileV2Attributes(d, lbClient, &allfps[0])
}

func dataSourceLBFlavorProfileV2Attributes(d *schema.ResourceData, computeClient *gophercloud.ServiceClient, fp *flavorprofiles.FlavorProfile) diag.Diagnostics {
	log.Printf("[DEBUG] Retrieved openstack_lb_flavorprofile_v2 %s: %#v", fp.ID, fp)

	d.SetId(fp.ID)
	d.Set("name", fp.Name)
	d.Set("provider_name", fp.ProviderName)
	d.Set("flavor_data", fp.FlavorData)

	return nil
}
