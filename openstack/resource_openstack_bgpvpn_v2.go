package openstack

import (
	"context"
	"log"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/bgpvpns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceBGPVPNV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBGPVPNV2Create,
		ReadContext:   resourceBGPVPNV2Read,
		UpdateContext: resourceBGPVPNV2Update,
		DeleteContext: resourceBGPVPNV2Delete,
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
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"l2", "l3",
				}, false),
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"vni": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"local_pref": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"route_distinguishers": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"route_targets": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"import_targets": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"export_targets": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"networks": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"routers": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ports": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"shared": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceBGPVPNV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	createOpts := bgpvpns.CreateOpts{
		Name:                d.Get("name").(string),
		RouteDistinguishers: expandToStringSlice(d.Get("route_distinguishers").(*schema.Set).List()),
		RouteTargets:        expandToStringSlice(d.Get("route_targets").(*schema.Set).List()),
		ImportTargets:       expandToStringSlice(d.Get("import_targets").(*schema.Set).List()),
		ExportTargets:       expandToStringSlice(d.Get("export_targets").(*schema.Set).List()),
		Type:                d.Get("type").(string),
		ProjectID:           d.Get("project_id").(string),
		VNI:                 d.Get("vni").(int),
		LocalPref:           d.Get("local_pref").(int),
	}

	log.Printf("[DEBUG] Create openstack_bgpvpn_v2: %#v", createOpts)

	bgpvpn, err := bgpvpns.Create(ctx, networkingClient, createOpts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] openstack_bgpvpn_v2 created: %#v", bgpvpn)

	d.SetId(bgpvpn.ID)

	return resourceBGPVPNV2Read(ctx, d, meta)
}

func resourceBGPVPNV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	log.Printf("[DEBUG] Retrieve information about openstack_bgpvpn_v2: %s", d.Id())

	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	bgpvpn, err := bgpvpns.Get(ctx, networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "openstack_bgpvpn_v2"))
	}

	log.Printf("[DEBUG] Read OpenStack openstack_bgpvpn_v2 %s: %#v", d.Id(), bgpvpn)

	d.Set("name", bgpvpn.Name)
	d.Set("type", bgpvpn.Type)
	d.Set("project_id", bgpvpn.ProjectID)
	d.Set("vni", bgpvpn.VNI)
	d.Set("local_pref", bgpvpn.LocalPref)
	d.Set("route_distinguishers", bgpvpn.RouteDistinguishers)
	d.Set("route_targets", bgpvpn.RouteTargets)
	d.Set("import_targets", bgpvpn.ImportTargets)
	d.Set("export_targets", bgpvpn.ExportTargets)
	d.Set("networks", bgpvpn.Networks)
	d.Set("routers", bgpvpn.Routers)
	d.Set("ports", bgpvpn.Ports)
	d.Set("shared", bgpvpn.Shared)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceBGPVPNV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	opts := bgpvpns.UpdateOpts{}

	var hasChange bool

	if d.HasChange("name") {
		name := d.Get("name").(string)
		opts.Name = &name
		hasChange = true
	}

	if d.HasChange("route_distinguishers") {
		routeDistinguishers := expandToStringSlice(d.Get("route_distinguishers").(*schema.Set).List())
		opts.RouteDistinguishers = &routeDistinguishers
		hasChange = true
	}

	if d.HasChange("route_targets") {
		routeTargets := expandToStringSlice(d.Get("route_targets").(*schema.Set).List())
		opts.RouteTargets = &routeTargets
		hasChange = true
	}

	if d.HasChange("import_targets") {
		importTargets := expandToStringSlice(d.Get("import_targets").(*schema.Set).List())
		opts.ImportTargets = &importTargets
		hasChange = true
	}

	if d.HasChange("export_targets") {
		exportTargets := expandToStringSlice(d.Get("export_targets").(*schema.Set).List())
		opts.ExportTargets = &exportTargets
		hasChange = true
	}

	if d.HasChange("local_pref") {
		localPref := d.Get("local_pref").(int)
		opts.LocalPref = &localPref
		hasChange = true
	}

	log.Printf("[DEBUG] Updating openstack_bgpvpn_v2 with id %s: %#v", d.Id(), opts)

	if hasChange {
		_, err = bgpvpns.Update(ctx, networkingClient, d.Id(), opts).Extract()
		if err != nil {
			return diag.FromErr(err)
		}

		log.Printf("[DEBUG] Updated openstack_bgpvpn_v2 with id %s", d.Id())
	}

	return resourceBGPVPNV2Read(ctx, d, meta)
}

func resourceBGPVPNV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	log.Printf("[DEBUG] Destroy openstack_bgpvpn_v2: %s", d.Id())

	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	err = bgpvpns.Delete(ctx, networkingClient, d.Id()).Err
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_bgpvpn_v2"))
	}

	return nil
}
