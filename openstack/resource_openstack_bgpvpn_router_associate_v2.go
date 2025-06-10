package openstack

import (
	"context"
	"fmt"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/bgpvpns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBGPVPNRouterAssociateV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBGPVPNRouterAssociateV2Create,
		ReadContext:   resourceBGPVPNRouterAssociateV2Read,
		UpdateContext: resourceBGPVPNRouterAssociateV2Update,
		DeleteContext: resourceBGPVPNRouterAssociateV2Delete,
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
			"bgpvpn_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"advertise_extra_routes": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceBGPVPNRouterAssociateV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack network client: %s", err)
	}

	bgpvpnID := d.Get("bgpvpn_id").(string)
	routerID := d.Get("router_id").(string)
	opts := bgpvpns.CreateRouterAssociationOpts{
		RouterID:  routerID,
		ProjectID: d.Get("project_id").(string),
	}

	if v, ok := getOkExists(d, "advertise_extra_routes"); ok {
		v := v.(bool)
		opts.AdvertiseExtraRoutes = &v
	}

	log.Printf("[DEBUG] openstack_bgpvpn_router_associate_v2 create options: %#v", opts)

	res, err := bgpvpns.CreateRouterAssociation(ctx, networkingClient, bgpvpnID, opts).Extract()
	if err != nil {
		return diag.Errorf("Error associating openstack_bgpvpn_router_associate_v2 BGP VPN %s with router %s: %s", bgpvpnID, routerID, err)
	}

	id := fmt.Sprintf("%s/%s", bgpvpnID, res.ID)
	d.SetId(id)

	return resourceBGPVPNRouterAssociateV2Read(ctx, d, meta)
}

func resourceBGPVPNRouterAssociateV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack network client: %s", err)
	}

	bgpvpnID, id, err := parsePairedIDs(d.Id(), "openstack_bgpvpn_router_associate_v2")
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := bgpvpns.GetRouterAssociation(ctx, networkingClient, bgpvpnID, id).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting openstack_bgpvpn_router_associate_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_bgpvpn_router_associate_v2 %s: %#v", id, res)

	d.Set("bgpvpn_id", bgpvpnID)
	d.Set("router_id", res.RouterID)
	d.Set("advertise_extra_routes", res.AdvertiseExtraRoutes)
	d.Set("project_id", res.ProjectID)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceBGPVPNRouterAssociateV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack network client: %s", err)
	}

	bgpvpnID, id, err := parsePairedIDs(d.Id(), "openstack_bgpvpn_router_associate_v2")
	if err != nil {
		return diag.FromErr(err)
	}

	opts := bgpvpns.UpdateRouterAssociationOpts{}

	if d.HasChange("advertise_extra_routes") {
		v := d.Get("advertise_extra_routes").(bool)
		opts.AdvertiseExtraRoutes = &v
	}

	log.Printf("[DEBUG] openstack_bgpvpn_router_associate_v2 %s update options: %#v", id, opts)

	_, err = bgpvpns.UpdateRouterAssociation(ctx, networkingClient, bgpvpnID, id, opts).Extract()
	if err != nil {
		return diag.Errorf("Error updating openstack_bgpvpn_router_associate_v2 %s: %s", id, err)
	}

	return resourceBGPVPNRouterAssociateV2Read(ctx, d, meta)
}

func resourceBGPVPNRouterAssociateV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack network client: %s", err)
	}

	bgpvpnID, id, err := parsePairedIDs(d.Id(), "openstack_bgpvpn_router_associate_v2")
	if err != nil {
		return diag.FromErr(err)
	}

	routerID := d.Get("router_id").(string)

	err = bgpvpns.DeleteRouterAssociation(ctx, networkingClient, bgpvpnID, id).ExtractErr()
	if err != nil && CheckDeleted(d, err, "") != nil {
		return diag.Errorf("Error disassociating openstack_bgpvpn_router_associate_v2 BGP VPN %s with router %s: %s", bgpvpnID, routerID, err)
	}

	return nil
}
