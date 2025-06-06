package openstack

import (
	"context"
	"log"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/flavors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLoadBalancerFlavorV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLoadBalancerFlavorV2Create,
		ReadContext:   resourceLoadBalancerFlavorV2Read,
		UpdateContext: resourceLoadBalancerFlavorV2Update,
		DeleteContext: resourceLoadBalancerFlavorV2Delete,
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

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"flavor_profile_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceLoadBalancerFlavorV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack loadbalancing client: %s", err)
	}

	var enabled *bool

	if v, ok := getOkExists(d, "enabled"); ok {
		v := v.(bool)
		enabled = &v
	}

	// TODO: remove custom struct when gophercloud support pointers,
	// e.g. gophercloud/v3 is released, see https://github.com/gophercloud/gophercloud/pull/3190
	createOpts := flavorsCreateOpts{
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		FlavorProfileID: d.Get("flavor_profile_id").(string),
		Enabled:         enabled,
	}

	flavor, err := flavors.Create(ctx, lbClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_lb_flavor_v2: %s", err)
	}

	d.SetId(flavor.ID)

	log.Printf("[DEBUG] Created openstack_lb_flavor_v2 %#v", flavor)

	return resourceLoadBalancerFlavorV2Read(ctx, d, meta)
}

func resourceLoadBalancerFlavorV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack loadbalancing client: %s", err)
	}

	flavor, err := flavors.Get(ctx, lbClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_lb_flavor_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_lb_flavor_v2 %s: %#v", d.Id(), flavor)

	d.Set("name", flavor.Name)
	d.Set("description", flavor.Description)
	d.Set("flavor_profile_id", flavor.FlavorProfileId)
	d.Set("enabled", flavor.Enabled)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceLoadBalancerFlavorV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack loadbalancing client: %s", err)
	}

	var (
		hasChange bool
		// TODO: remove custom struct when gophercloud support pointers,
		// e.g. gophercloud/v3 is released, see https://github.com/gophercloud/gophercloud/pull/3190
		updateOpts flavorsUpdateOpts
	)

	if d.HasChange("name") {
		hasChange = true
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}

	if d.HasChange("description") {
		hasChange = true
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}

	if d.HasChange("enabled") {
		hasChange = true
		enabled := d.Get("enabled").(bool)
		updateOpts.Enabled = &enabled
	}

	if hasChange {
		log.Printf("[DEBUG] openstack_lb_flavor_v2 %s update options: %#v", d.Id(), updateOpts)

		_, err := flavors.Update(ctx, lbClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_lb_flavor_v2: %s", err)
		}
	}

	return resourceLoadBalancerFlavorV2Read(ctx, d, meta)
}

func resourceLoadBalancerFlavorV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	lbClient, err := config.LoadBalancerV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack loadbalancing client: %s", err)
	}

	log.Printf("[DEBUG] Deleting openstack_lb_flavor_v2: %s", d.Id())

	if err := flavors.Delete(ctx, lbClient, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_lb_flavor_v2"))
	}

	d.SetId("")

	return nil
}
