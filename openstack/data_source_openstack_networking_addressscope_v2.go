package openstack

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/addressscopes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNetworkingAddressScopeV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingAddressScopeV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"ip_version": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"shared": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceNetworkingAddressScopeV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	listOpts := addressscopes.ListOpts{}

	if v, ok := d.GetOk("name"); ok {
		listOpts.Name = v.(string)
	}

	if v, ok := d.GetOk("ip_version"); ok {
		listOpts.IPVersion = v.(int)
	}

	if v, ok := getOkExists(d, "shared"); ok {
		shared := v.(bool)
		listOpts.Shared = &shared
	}

	if v, ok := d.GetOk("project_id"); ok {
		listOpts.ProjectID = v.(string)
	}

	pages, err := addressscopes.List(networkingClient, listOpts).AllPages(ctx)
	if err != nil {
		return diag.Errorf("Unable to list openstack_networking_addressscope_v2: %s", err)
	}

	allAddressScopes, err := addressscopes.ExtractAddressScopes(pages)
	if err != nil {
		return diag.Errorf("Unable to retrieve openstack_networking_addressscope_v2: %s", err)
	}

	if len(allAddressScopes) < 1 {
		return diag.Errorf("No openstack_networking_addressscope_v2 found")
	}

	if len(allAddressScopes) > 1 {
		return diag.Errorf("More than one openstack_networking_addressscope_v2 found")
	}

	a := allAddressScopes[0]

	log.Printf("[DEBUG] Retrieved openstack_networking_addressscope_v2 %s: %+v", a.ID, a)
	d.SetId(a.ID)

	d.Set("region", GetRegion(d, config))
	d.Set("name", a.Name)
	d.Set("ip_version", a.IPVersion)
	d.Set("shared", a.Shared)
	d.Set("project_id", a.ProjectID)

	return nil
}
