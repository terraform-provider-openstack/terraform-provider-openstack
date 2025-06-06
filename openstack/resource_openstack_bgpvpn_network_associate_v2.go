package openstack

import (
	"context"
	"fmt"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/bgpvpns"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBGPVPNNetworkAssociateV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBGPVPNNetworkAssociateV2Create,
		ReadContext:   resourceBGPVPNNetworkAssociateV2Read,
		DeleteContext: resourceBGPVPNNetworkAssociateV2Delete,
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
			"network_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceBGPVPNNetworkAssociateV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack network client: %s", err)
	}

	bgpvpnID := d.Get("bgpvpn_id").(string)
	networkID := d.Get("network_id").(string)
	opts := bgpvpns.CreateNetworkAssociationOpts{
		NetworkID: networkID,
		ProjectID: d.Get("project_id").(string),
	}

	log.Printf("[DEBUG] openstack_bgpvpn_network_associate_v2 create options: %#v", opts)

	res, err := bgpvpns.CreateNetworkAssociation(ctx, networkingClient, bgpvpnID, opts).Extract()
	if err != nil {
		return diag.Errorf("Error associating openstack_bgpvpn_network_associate_v2 BGP VPN %s with network %s: %s", bgpvpnID, networkID, err)
	}

	id := fmt.Sprintf("%s/%s", bgpvpnID, res.ID)
	d.SetId(id)

	return resourceBGPVPNNetworkAssociateV2Read(ctx, d, meta)
}

func resourceBGPVPNNetworkAssociateV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack network client: %s", err)
	}

	bgpvpnID, id, err := parsePairedIDs(d.Id(), "openstack_bgpvpn_network_associate_v2")
	if err != nil {
		return diag.FromErr(err)
	}

	res, err := bgpvpns.GetNetworkAssociation(ctx, networkingClient, bgpvpnID, id).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting openstack_bgpvpn_network_associate_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_bgpvpn_network_associate_v2 %s: %#v", id, res)

	d.Set("bgpvpn_id", bgpvpnID)
	d.Set("network_id", res.NetworkID)
	d.Set("project_id", res.ProjectID)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceBGPVPNNetworkAssociateV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack network client: %s", err)
	}

	bgpvpnID, id, err := parsePairedIDs(d.Id(), "openstack_bgpvpn_network_associate_v2")
	if err != nil {
		return diag.FromErr(err)
	}

	networkID := d.Get("network_id").(string)

	err = bgpvpns.DeleteNetworkAssociation(ctx, networkingClient, bgpvpnID, id).ExtractErr()
	if err != nil && CheckDeleted(d, err, "") != nil {
		return diag.Errorf("Error disassociating openstack_bgpvpn_network_associate_v2 BGP VPN %s with network %s: %s", bgpvpnID, networkID, err)
	}

	return nil
}
