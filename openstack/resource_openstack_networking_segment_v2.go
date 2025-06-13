package openstack

import (
	"context"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/segments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNetworkingSegmentV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingSegmentCreate,
		ReadContext:   resourceNetworkingSegmentRead,
		UpdateContext: resourceNetworkingSegmentUpdate,
		DeleteContext: resourceNetworkingSegmentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"network_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"vlan",
					"vxlan",
					"flat",
					"gre",
					"geneve",
					"local",
				}, true),
			},

			"physical_network": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"segmentation_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"revision_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceNetworkingSegmentCreate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	opts := segments.CreateOpts{
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		NetworkID:       d.Get("network_id").(string),
		NetworkType:     d.Get("network_type").(string),
		PhysicalNetwork: d.Get("physical_network").(string),
		SegmentationID:  d.Get("segmentation_id").(int),
	}

	seg, err := segments.Create(ctx, networkingClient, opts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_networking_segment_v2: %s", err)
	}

	d.SetId(seg.ID)

	return resourceNetworkingSegmentRead(ctx, d, meta)
}

func resourceNetworkingSegmentRead(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	seg, err := segments.Get(ctx, networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting openstack_networking_segment_v2"))
	}

	d.Set("name", seg.Name)
	d.Set("description", seg.Description)
	d.Set("network_id", seg.NetworkID)
	d.Set("network_type", seg.NetworkType)
	d.Set("physical_network", seg.PhysicalNetwork)
	d.Set("segmentation_id", seg.SegmentationID)
	d.Set("revision_number", seg.RevisionNumber)
	d.Set("created_at", seg.CreatedAt.String())
	d.Set("updated_at", seg.UpdatedAt.String())

	return nil
}

func resourceNetworkingSegmentUpdate(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var (
		opts    segments.UpdateOpts
		changed bool
	)

	if d.HasChange("name") {
		changed = true
		v := d.Get("name").(string)
		opts.Name = &v
	}

	if d.HasChange("description") {
		changed = true
		v := d.Get("description").(string)
		opts.Description = &v
	}

	// it is possible to update the VLAN network_type segmentation_id
	if d.HasChange("segmentation_id") {
		changed = true
		v := d.Get("segmentation_id").(int)
		opts.SegmentationID = &v
	}

	if changed {
		_, err := segments.Update(ctx, networkingClient, d.Id(), opts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_networking_segment_v2 %s: %s", d.Id(), err)
		}
	}

	return resourceNetworkingSegmentRead(ctx, d, meta)
}

func resourceNetworkingSegmentDelete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	err = segments.Delete(ctx, networkingClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_networking_segment_v2"))
	}

	d.SetId("")

	return nil
}
