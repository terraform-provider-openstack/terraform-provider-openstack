package openstack

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/segments"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceNetworkingSegmentV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingSegmentV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"segment_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"network_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
				Computed: true,
			},

			"segmentation_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"revision_number": {
				Type:     schema.TypeInt,
				Optional: true,
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

func dataSourceNetworkingSegmentV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	if id, ok := d.Get("segment_id").(string); ok && id != "" {
		// If segment_id is provided, we will try to get the segment by ID
		log.Printf("[DEBUG] Attempting to retrieve openstack_networking_segment_v2 by ID: %s", id)

		seg, err := segments.Get(ctx, networkingClient, id).Extract()
		if err != nil {
			return diag.Errorf("Error retrieving openstack_networking_segment_v2 by ID %s: %s", id, err)
		}

		log.Printf("[DEBUG] Retrieved openstack_networking_segment_v2 %s: %+v", seg.ID, seg)

		d.SetId(seg.ID)
		d.Set("name", seg.Name)
		d.Set("description", seg.Description)
		d.Set("segment_id", seg.ID)
		d.Set("network_id", seg.NetworkID)
		d.Set("network_type", seg.NetworkType)
		d.Set("physical_network", seg.PhysicalNetwork)
		d.Set("segmentation_id", seg.SegmentationID)
		d.Set("revision_number", seg.RevisionNumber)
		d.Set("created_at", seg.CreatedAt.String())
		d.Set("updated_at", seg.UpdatedAt.String())
		d.Set("region", GetRegion(d, config))

		return nil
	}

	listOpts := segments.ListOpts{
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		NetworkID:       d.Get("network_id").(string),
		NetworkType:     d.Get("network_type").(string),
		PhysicalNetwork: d.Get("physical_network").(string),
		SegmentationID:  d.Get("segmentation_id").(int),
		RevisionNumber:  d.Get("revision_number").(int),
	}

	pages, err := segments.List(networkingClient, listOpts).AllPages(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	allSegments, err := segments.ExtractSegments(pages)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(allSegments) < 1 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(allSegments) > 1 {
		return diag.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	seg := allSegments[0]

	log.Printf("[DEBUG] Retrieved openstack_networking_segment_v2 %s: %+v", seg.ID, seg)
	d.SetId(seg.ID)

	d.Set("name", seg.Name)
	d.Set("description", seg.Description)
	d.Set("segment_id", seg.ID)
	d.Set("network_id", seg.NetworkID)
	d.Set("network_type", seg.NetworkType)
	d.Set("physical_network", seg.PhysicalNetwork)
	d.Set("segmentation_id", seg.SegmentationID)
	d.Set("revision_number", seg.RevisionNumber)
	d.Set("created_at", seg.CreatedAt.String())
	d.Set("updated_at", seg.UpdatedAt.String())
	d.Set("region", GetRegion(d, config))

	return nil
}
