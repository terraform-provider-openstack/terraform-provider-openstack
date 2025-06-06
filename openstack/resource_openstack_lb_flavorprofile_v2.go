package openstack

import (
	"context"
	"log"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/flavorprofiles"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
)

func resourceLoadBalancerFlavorProfileV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoadBalancerFlavorProfileV2Create,
		ReadContext:   resourceLoadBalancerFlavorProfileV2Read,
		UpdateContext: resourceLoadBalancerFlavorProfileV2Update,
		DeleteContext: resourceLoadBalancerFlavorProfileV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"provider_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			// flavor_data depends on which provider is being used.
			// Therefore we stay close to the API and make it type
			// String. The user can use jsonencode to pass it properly
			"flavor_data": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateFunc:     validateJSONObject,
				DiffSuppressFunc: diffSuppressJSONObject,
				StateFunc: func(v any) string {
					json, _ := structure.NormalizeJsonString(v)

					return json
				},
			},
		},
	}
}

func resourceLoadBalancerFlavorProfileV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack loadbalancing client: %s", err)
	}

	name := d.Get("name").(string)
	providerName := d.Get("provider_name").(string)
	flavorData := d.Get("flavor_data").(string)

	createOpts := flavorprofiles.CreateOpts{
		Name:         name,
		ProviderName: providerName,
		FlavorData:   flavorData,
	}

	q, err := flavorprofiles.Create(ctx, lbClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_lb_flavorprofile_v2: %s", err)
	}

	d.SetId(q.ID)

	log.Printf("[DEBUG] Created openstack_lb_flavorprofile_v2 %#v", q)

	return resourceLoadBalancerFlavorProfileV2Read(ctx, d, meta)
}

func resourceLoadBalancerFlavorProfileV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack loadbalancing client: %s", err)
	}

	q, err := flavorprofiles.Get(ctx, lbClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_lb_flavorprofile_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_lb_flavorprofile_v2 %s: %#v", d.Id(), q)

	d.Set("name", q.Name)
	d.Set("provider_name", q.ProviderName)
	d.Set("flavor_data", q.FlavorData)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceLoadBalancerFlavorProfileV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack loadbalancing client: %s", err)
	}

	var (
		hasChange  bool
		updateOpts flavorprofiles.UpdateOpts
	)

	if d.HasChange("name") {
		hasChange = true
		name := d.Get("name").(string)
		updateOpts.Name = name
	}

	if d.HasChange("provider_name") {
		hasChange = true
		providerName := d.Get("provider_name").(string)
		updateOpts.ProviderName = providerName
	}

	if d.HasChange("flavor_data") {
		hasChange = true
		flavorData := d.Get("flavor_data").(string)
		updateOpts.FlavorData = flavorData
	}

	if hasChange {
		log.Printf("[DEBUG] openstack_lb_flavorprofile_v2 %s update options: %#v", d.Id(), updateOpts)

		_, err := flavorprofiles.Update(ctx, lbClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_lb_flavorprofile_v2: %s", err)
		}
	}

	return resourceLoadBalancerFlavorProfileV2Read(ctx, d, meta)
}

func resourceLoadBalancerFlavorProfileV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack loadbalancing client: %s", err)
	}

	log.Printf("[DEBUG] Deleting openstack_lb_flavorprofile_v2: %s", d.Id())

	if err := flavorprofiles.Delete(ctx, lbClient, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_lb_flavorprofile_v2"))
	}

	d.SetId("")

	return nil
}
