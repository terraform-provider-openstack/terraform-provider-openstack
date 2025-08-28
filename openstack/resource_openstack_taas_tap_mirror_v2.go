package openstack

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/taas/tapmirrors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceTapMirrorV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTapMirrorV2Create,
		ReadContext:   resourceTapMirrorV2Read,
		UpdateContext: resourceTapMirrorV2Update,
		DeleteContext: resourceTapMirrorV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"mirror_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"erspanv1", "gre",
				}, false),
			},
			"remote_ip": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"directions": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"in": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							AtLeastOneOf: []string{"directions.0.out"},
						},
						"out": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func resourceTapMirrorV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var createOpts tapmirrors.CreateOptsBuilder

	mirrorType := resourceTapMirrorV2MirrorType(d.Get("mirror_type").(string))
	directions := resourceTapMirrorV2Directions(d.Get("directions").([]any))

	createOpts = tapmirrors.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		TenantID:    d.Get("tenant_id").(string),
		PortID:      d.Get("port_id").(string),
		MirrorType:  mirrorType,
		RemoteIP:    d.Get("remote_ip").(string),
		Directions:  directions,
	}

	log.Printf("[DEBUG] Create tapMirror: %#v", createOpts)

	tapMirror, err := tapmirrors.Create(ctx, networkingClient, createOpts).Extract()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] TapMirror created: %#v", tapMirror)

	d.SetId(tapMirror.ID)

	return resourceTapMirrorV2Read(ctx, d, meta)
}

func resourceTapMirrorV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	log.Printf("[DEBUG] Retrieve information about tapMirror: %s", d.Id())

	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	tapMirror, err := tapmirrors.Get(ctx, networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "tapMirror"))
	}

	log.Printf("[DEBUG] Read OpenStack Endpoint TapMirror %s: %#v", d.Id(), tapMirror)

	d.Set("name", tapMirror.Name)
	d.Set("description", tapMirror.Description)
	d.Set("tenant_id", tapMirror.TenantID)
	d.Set("project_id", tapMirror.ProjectID)
	d.Set("port_id", tapMirror.PortID)
	d.Set("mirror_type", tapMirror.MirrorType)
	d.Set("remote_ip", tapMirror.RemoteIP)
	d.Set("directions", resourceTapMirrorV2DirectionsToMap(tapMirror.Directions))

	return nil
}

func resourceTapMirrorV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	opts := tapmirrors.UpdateOpts{}

	var hasChange bool

	if d.HasChange("name") {
		name := d.Get("name").(string)
		opts.Name = &name
		hasChange = true
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		opts.Description = &description
		hasChange = true
	}

	log.Printf("[DEBUG] Updating tapMirror with id %s: %#v", d.Id(), opts)

	if hasChange {
		_, err := tapmirrors.Update(ctx, networkingClient, d.Id(), opts).Extract()
		if err != nil {
			return diag.FromErr(err)
		}

		log.Printf("[DEBUG] Updated tapMirror with id %s", d.Id())
	}

	return resourceTapMirrorV2Read(ctx, d, meta)
}

func resourceTapMirrorV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	log.Printf("[DEBUG] Destroy tapMirror: %s", d.Id())

	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	err = tapmirrors.Delete(ctx, networkingClient, d.Id()).Err
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(err)
}

func resourceTapMirrorV2MirrorType(mirrorType string) tapmirrors.MirrorType {
	var result tapmirrors.MirrorType

	switch mirrorType {
	case "erspanv1":
		result = tapmirrors.MirrorTypeErspanv1
	case "gre":
		result = tapmirrors.MirrorTypeGre
	}

	return result
}

func resourceTapMirrorV2Directions(directions []any) tapmirrors.Directions {
	var result tapmirrors.Directions

	for _, raw := range directions {
		binding := raw.(map[string]any)
		if value, exists := binding["in"]; exists {
			result.In = value.(int)
		}

		if value, exists := binding["out"]; exists {
			result.Out = value.(int)
		}
	}

	return result
}

func resourceTapMirrorV2DirectionsToMap(directions tapmirrors.Directions) []map[string]any {
	result := make(map[string]any, 2)

	result["in"] = directions.In
	result["out"] = directions.Out

	return []map[string]any{result}
}
