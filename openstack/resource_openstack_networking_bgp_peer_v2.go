package openstack

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/bgp/peers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNetworkingBGPPeerV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingBGPPeerV2Create,
		ReadContext:   resourceNetworkingBGPPeerV2Read,
		UpdateContext: resourceNetworkingBGPPeerV2Update,
		DeleteContext: resourceNetworkingBGPPeerV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"auth_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "none",
				ValidateFunc: validation.StringInSlice([]string{
					"none",
					"md5",
				}, false),
			},

			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"remote_as": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},

			"peer_ip": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPAddress,
			},
		},

		CustomizeDiff: func(_ context.Context, d *schema.ResourceDiff, _ any) error {
			authType := d.Get("auth_type").(string)
			password := d.Get("password").(string)

			if authType == "none" && password != "" {
				return errors.New("password can only be set when auth_type is not \"none\"")
			}

			if authType != "none" && password == "" {
				return errors.New("password must be set when auth_type is not \"none\"")
			}

			return nil
		},
	}
}

func resourceNetworkingBGPPeerV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	opts := peersCreateOpts{
		Name:     d.Get("name").(string),
		AuthType: d.Get("auth_type").(string),
		PeerIP:   d.Get("peer_ip").(string),
		RemoteAS: d.Get("remote_as").(int),
		TenantID: d.Get("tenant_id").(string),
	}

	log.Printf("[DEBUG] openstack_networking_bgp_peer_v2 create options: %#v", opts)

	opts.Password = d.Get("password").(string)

	bgpPeer, err := peers.Create(ctx, networkingClient, opts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_networking_bgp_peer_v2: %s", err)
	}

	log.Printf("[DEBUG] Created openstack_networking_bgp_peer_v2: %#v", bgpPeer)

	d.SetId(bgpPeer.ID)

	return resourceNetworkingBGPPeerV2Read(ctx, d, meta)
}

func resourceNetworkingBGPPeerV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	bgpPeer, err := peers.Get(ctx, networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_networking_bgp_peer_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_networking_bgp_peer_v2: %#v", bgpPeer)

	d.Set("tenant_id", bgpPeer.TenantID)
	d.Set("name", bgpPeer.Name)
	d.Set("auth_type", bgpPeer.AuthType)
	d.Set("peer_ip", bgpPeer.PeerIP)
	d.Set("remote_as", bgpPeer.RemoteAS)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceNetworkingBGPPeerV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var (
		updateOpts peersUpdateOpts
		changed    bool
	)

	if d.HasChange("name") {
		changed = true
		v := d.Get("name").(string)
		updateOpts.Name = &v
	}

	log.Printf("[DEBUG] Updating openstack_networking_bgp_peer_v2 %s with options: %#v", d.Id(), updateOpts)

	if d.HasChange("password") {
		changed = true
		v := d.Get("password").(string)
		updateOpts.Password = &v
	}

	if changed {
		_, err = peers.Update(ctx, networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_networking_bgp_peer_v2: %s", err)
		}
	}

	return resourceNetworkingBGPPeerV2Read(ctx, d, meta)
}

func resourceNetworkingBGPPeerV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	log.Printf("[DEBUG] Deleting openstack_networking_bgp_peer_v2 %s", d.Id())

	err = peers.Delete(ctx, networkingClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_networking_bgp_peer_v2"))
	}

	d.SetId("")

	return nil
}
