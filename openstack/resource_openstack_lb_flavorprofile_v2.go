package openstack

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/flavorprofiles"
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
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if strings.ReplaceAll(old, " ", "") == strings.ReplaceAll(new, " ", "") {
						return true
					}
					return false
				},
			},
		},
	}
}

func resourceLoadBalancerFlavorProfileV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	lbClient, err := config.LoadBalancerV2Client(GetRegion(d, config))
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

	q, err := flavorprofiles.Create(lbClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_lb_flavorprofile_v2: %s", err)
	}

	d.SetId(q.ID)

	log.Printf("[DEBUG] Created openstack_lb_flavorprofile_v2 %#v", q)

	return resourceLoadBalancerFlavorProfileV2Read(ctx, d, meta)
}

func resourceLoadBalancerFlavorProfileV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	lbClient, err := config.LoadBalancerV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack loadbalancing client: %s", err)
	}

	q, err := flavorprofiles.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_lb_flavorprofile_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_lb_flavorprofile_v2 %s: %#v", d.Id(), q)

	d.Set("name", q.Name)
	d.Set("provider_name", q.ProviderName)
	d.Set("flavor_data", q.FlavorData)

	return nil
}

func resourceLoadBalancerFlavorProfileV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	lbClient, err := config.LoadBalancerV2Client(GetRegion(d, config))
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
		_, err := flavorprofiles.Update(lbClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_lb_flavorprofile_v2: %s", err)
		}
	}

	return resourceLoadBalancerFlavorProfileV2Read(ctx, d, meta)
}

func resourceLoadBalancerFlavorProfileV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	lbClient, err := config.LoadBalancerV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack loadbalancing client: %s", err)
	}

	log.Printf("[DEBUG] Deleting openstack_lb_flavorprofile_v2: %s", d.Id())
	if err := flavorprofiles.Delete(lbClient, d.Id()).Err; err != nil {
		return diag.Errorf("Error deleting openstack_lb_flavorprofile_v2: %s", err)
	}

	d.SetId("")
	return nil
}
