package openstack

import (
	"context"
	"fmt"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/bgpvpns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceBGPVPNPortAssociateV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBGPVPNPortAssociateV2Create,
		ReadContext:   resourceBGPVPNPortAssociateV2Read,
		UpdateContext: resourceBGPVPNPortAssociateV2Update,
		DeleteContext: resourceBGPVPNPortAssociateV2Delete,
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
			"port_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"advertise_fixed_ips": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"routes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"prefix", "bgpvpn_id",
							}, false),
						},
						"prefix": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsCIDR,
						},
						"bgpvpn_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.IsUUID,
						},
						"local_pref": {
							Type:     schema.TypeInt,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceBGPVPNPortAssociateV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack network client: %s", err)
	}

	bgpvpnID := d.Get("bgpvpn_id").(string)
	portID := d.Get("port_id").(string)
	routes := d.Get("routes").(*schema.Set).List()
	opts := bgpvpns.CreatePortAssociationOpts{
		PortID:    portID,
		ProjectID: d.Get("project_id").(string),
		Routes:    expandBGPVPNPortAssociateRoutesV2(routes),
	}

	if v, ok := getOkExists(d, "advertise_fixed_ips"); ok {
		v := v.(bool)
		opts.AdvertiseFixedIPs = &v
	}

	log.Printf("[DEBUG] openstack_bgpvpn_port_associate_v2 create options: %#v", opts)

	res, err := bgpvpns.CreatePortAssociation(ctx, networkingClient, bgpvpnID, opts).Extract()
	if err != nil {
		return diag.Errorf("Error associating openstack_bgpvpn_port_associate_v2 BGP VPN %s with port %s: %s", bgpvpnID, portID, err)
	}

	id := fmt.Sprintf("%s/%s", bgpvpnID, res.ID)
	d.SetId(id)
	// the project_id is returned only on POST/PUT responses
	d.Set("project_id", res.ProjectID)

	return resourceBGPVPNPortAssociateV2Read(ctx, d, meta)
}

func resourceBGPVPNPortAssociateV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack network client: %s", err)
	}

	bgpvpnID, id, err := parsePairedIDs(d.Id(), "openstack_bgpvpn_port_associate_v2")
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := bgpvpns.GetPortAssociation(ctx, networkingClient, bgpvpnID, id).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting openstack_bgpvpn_port_associate_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_bgpvpn_port_associate_v2 %s: %#v", id, res)

	d.Set("bgpvpn_id", bgpvpnID)
	d.Set("port_id", res.PortID)
	d.Set("advertise_fixed_ips", res.AdvertiseFixedIPs)

	if res.ProjectID != "" {
		// the project_id is returned only on POST/PUT responses
		d.Set("project_id", res.ProjectID)
	}

	d.Set("routes", flattenBGPVPNPortAssociateRoutesV2(res.Routes))

	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceBGPVPNPortAssociateV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack network client: %s", err)
	}

	bgpvpnID, id, err := parsePairedIDs(d.Id(), "openstack_bgpvpn_port_associate_v2")
	if err != nil {
		return diag.FromErr(err)
	}

	opts := bgpvpns.UpdatePortAssociationOpts{}

	if d.HasChange("advertise_fixed_ips") {
		v := d.Get("advertise_fixed_ips").(bool)
		opts.AdvertiseFixedIPs = &v
	}

	if d.HasChange("routes") {
		routes := expandBGPVPNPortAssociateRoutesUpdateV2(d)
		opts.Routes = &routes
	}

	log.Printf("[DEBUG] openstack_bgpvpn_port_associate_v2 %s update options: %#v", id, opts)

	res, err := bgpvpns.UpdatePortAssociation(ctx, networkingClient, bgpvpnID, id, opts).Extract()
	if err != nil {
		return diag.Errorf("Error updating openstack_bgpvpn_port_associate_v2 %s: %s", id, err)
	}

	// the project_id is returned only on POST/PUT responses
	d.Set("project_id", res.ProjectID)

	return resourceBGPVPNPortAssociateV2Read(ctx, d, meta)
}

func resourceBGPVPNPortAssociateV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack network client: %s", err)
	}

	bgpvpnID, id, err := parsePairedIDs(d.Id(), "openstack_bgpvpn_port_associate_v2")
	if err != nil {
		return diag.FromErr(err)
	}

	portID := d.Get("port_id").(string)

	err = bgpvpns.DeletePortAssociation(ctx, networkingClient, bgpvpnID, id).ExtractErr()
	if err != nil && CheckDeleted(d, err, "") == nil {
		return diag.Errorf("Error disassociating openstack_bgpvpn_port_associate_v2 BGP VPN %s with port %s: %s", bgpvpnID, portID, err)
	}

	return nil
}
