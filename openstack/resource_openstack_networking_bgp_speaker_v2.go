package openstack

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/bgp/speakers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNetworkingBGPSpeakerV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingBGPSpeakerV2Create,
		ReadContext:   resourceNetworkingBGPSpeakerV2Read,
		UpdateContext: resourceNetworkingBGPSpeakerV2Update,
		DeleteContext: resourceNetworkingBGPSpeakerV2Delete,
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

			"ip_version": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      4,
				ValidateFunc: validation.IntInSlice([]int{4, 6}),
			},

			"advertise_floating_ip_host_routes": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"advertise_tenant_networks": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"local_as": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},

			"networks": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"peers": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"advertised_routes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"next_hop": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceNetworkingBGPSpeakerV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	opts := speakersCreateOpts{
		TenantID:  d.Get("tenant_id").(string),
		Name:      d.Get("name").(string),
		IPVersion: d.Get("ip_version").(int),
		LocalAS:   d.Get("local_as").(int),
	}

	if v, ok := getOkExists(d, "advertise_floating_ip_host_routes"); ok {
		v := v.(bool)
		opts.AdvertiseFloatingIPHostRoutes = &v
	}

	if v, ok := getOkExists(d, "advertise_tenant_networks"); ok {
		v := v.(bool)
		opts.AdvertiseTenantNetworks = &v
	}

	log.Printf("[DEBUG] openstack_networking_bgp_speaker_v2 create options: %#v", opts)

	bgpSpeaker, err := speakers.Create(ctx, networkingClient, opts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_networking_bgp_speaker_v2: %s", err)
	}

	log.Printf("[DEBUG] Created openstack_networking_bgp_speaker_v2: %#v", bgpSpeaker)

	d.SetId(bgpSpeaker.ID)

	if d.Get("networks") != nil {
		networks := expandToStringSlice(d.Get("networks").(*schema.Set).List())
		for _, network := range networks {
			log.Printf("[DEBUG] Adding network '%s' to openstack_networking_bgp_speaker_v2 '%s'", network, bgpSpeaker.ID)
			opts := speakers.AddGatewayNetworkOpts{
				NetworkID: network,
			}
			_, err = speakers.AddGatewayNetwork(ctx, networkingClient, bgpSpeaker.ID, opts).Extract()
			if err != nil {
				return diag.Errorf("Error adding network '%s' to openstack_networking_bgp_speaker_v2: %s", network, err)
			}

			log.Printf("[DEBUG] Successfully added network '%s' to openstack_networking_bgp_speaker_v2 '%s'", network, bgpSpeaker.ID)
		}
	}

	if d.Get("peers") != nil {
		peers := expandToStringSlice(d.Get("peers").(*schema.Set).List())
		for _, peer := range peers {
			log.Printf("[DEBUG] Adding peer '%s' to openstack_networking_bgp_speaker_v2 '%s'", peer, bgpSpeaker.ID)
			opts := speakers.AddBGPPeerOpts{
				BGPPeerID: peer,
			}
			_, err = speakers.AddBGPPeer(ctx, networkingClient, bgpSpeaker.ID, opts).Extract()
			if err != nil {
				return diag.Errorf("Error adding peer '%s' to openstack_networking_bgp_speaker_v2: %s", peer, err)
			}

			log.Printf("[DEBUG] Successfully added peer '%s' to openstack_networking_bgp_speaker_v2 '%s'", peer, bgpSpeaker.ID)
		}
	}

	return resourceNetworkingBGPSpeakerV2Read(ctx, d, meta)
}

func resourceNetworkingBGPSpeakerV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	bgpSpeaker, err := speakers.Get(ctx, networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_networking_bgp_speaker_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_networking_bgp_speaker_v2: %#v", bgpSpeaker)

	allPages, err := speakers.GetAdvertisedRoutes(networkingClient, d.Id()).AllPages(ctx)
	if err != nil {
		return diag.Errorf("Error retrieving advertised routes for openstack_networking_bgp_speaker_v2: %s", err)
	}

	advertisedRoutes, err := speakers.ExtractAdvertisedRoutes(allPages)
	if err != nil {
		return diag.Errorf("Error extracting advertised routes for openstack_networking_bgp_speaker_v2: %s", err)
	}

	log.Printf("[DEBUG] Retrieved advertised routes for openstack_networking_bgp_speaker_v2: %#v", advertisedRoutes)

	d.Set("tenant_id", bgpSpeaker.TenantID)
	d.Set("name", bgpSpeaker.Name)
	d.Set("ip_version", bgpSpeaker.IPVersion)
	d.Set("advertise_floating_ip_host_routes", bgpSpeaker.AdvertiseFloatingIPHostRoutes)
	d.Set("advertise_tenant_networks", bgpSpeaker.AdvertiseTenantNetworks)
	d.Set("local_as", bgpSpeaker.LocalAS)
	d.Set("networks", bgpSpeaker.Networks)
	d.Set("peers", bgpSpeaker.Peers)
	d.Set("advertised_routes", flattenBGPSpeakerAdvertisedRoutes(advertisedRoutes))
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceNetworkingBGPSpeakerV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var (
		updateOpts speakersUpdateOpts
		changed    bool
	)

	if d.HasChange("name") {
		changed = true
		v := d.Get("name").(string)
		updateOpts.Name = &v
	}

	if d.HasChange("advertise_floating_ip_host_routes") {
		changed = true
		v := d.Get("advertise_floating_ip_host_routes").(bool)
		updateOpts.AdvertiseFloatingIPHostRoutes = &v
	}

	if d.HasChange("advertise_tenant_networks") {
		changed = true
		v := d.Get("advertise_tenant_networks").(bool)
		updateOpts.AdvertiseTenantNetworks = &v
	}

	if changed {
		log.Printf("[DEBUG] Updating openstack_networking_bgp_speaker_v2 %s with options: %#v", d.Id(), updateOpts)

		_, err = speakers.Update(ctx, networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_networking_bgp_speaker_v2: %s", err)
		}
	}

	if d.HasChange("networks") {
		o, n := d.GetChange("networks")
		oldNet, newNet := o.(*schema.Set), n.(*schema.Set)
		netToDel := oldNet.Difference(newNet)
		netToAdd := newNet.Difference(oldNet)

		for _, v := range netToDel.List() {
			log.Printf("[DEBUG] Removing network '%s' from openstack_networking_bgp_speaker_v2 '%s'", v, d.Id())
			opts := speakers.RemoveGatewayNetworkOpts{
				NetworkID: v.(string),
			}
			err = speakers.RemoveGatewayNetwork(ctx, networkingClient, d.Id(), opts).ExtractErr()
			if !gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return diag.Errorf("Error removing network '%s' from openstack_networking_bgp_speaker_v2: %s", v, err)
			}

			log.Printf("[DEBUG] Successfully removed network '%s' from openstack_networking_bgp_speaker_v2 '%s'", v, d.Id())
		}

		for _, v := range netToAdd.List() {
			log.Printf("[DEBUG] Adding network '%s' to openstack_networking_bgp_speaker_v2 '%s'", v, d.Id())
			opts := speakers.AddGatewayNetworkOpts{
				NetworkID: v.(string),
			}
			_, err = speakers.AddGatewayNetwork(ctx, networkingClient, d.Id(), opts).Extract()
			if err != nil {
				return diag.Errorf("Error adding network '%s' to openstack_networking_bgp_speaker_v2: %s", v, err)
			}

			log.Printf("[DEBUG] Successfully added network '%s' to openstack_networking_bgp_speaker_v2 '%s'", v, d.Id())
		}
	}

	if d.HasChange("peers") {
		o, n := d.GetChange("peers")
		oldPeers, newPeers := o.(*schema.Set), n.(*schema.Set)
		peersToDel := oldPeers.Difference(newPeers)
		peersToAdd := newPeers.Difference(oldPeers)

		for _, v := range peersToDel.List() {
			log.Printf("[DEBUG] Removing peer '%s' from openstack_networking_bgp_speaker_v2 '%s'", v, d.Id())
			peerOpts := speakers.RemoveBGPPeerOpts{
				BGPPeerID: v.(string),
			}
			err = speakers.RemoveBGPPeer(ctx, networkingClient, d.Id(), peerOpts).ExtractErr()
			if !gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return diag.Errorf("Error removing peer '%s' from openstack_networking_bgp_speaker_v2: %s", v, err)
			}

			log.Printf("[DEBUG] Successfully removed peer '%s' from openstack_networking_bgp_speaker_v2 '%s'", v, d.Id())
		}

		for _, v := range peersToAdd.List() {
			log.Printf("[DEBUG] Adding peer '%s' to openstack_networking_bgp_speaker_v2 '%s'", v, d.Id())
			peerOpts := speakers.AddBGPPeerOpts{
				BGPPeerID: v.(string),
			}
			_, err = speakers.AddBGPPeer(ctx, networkingClient, d.Id(), peerOpts).Extract()
			if err != nil {
				return diag.Errorf("Error adding peer '%s' to openstack_networking_bgp_speaker_v2: %s", v, err)
			}

			log.Printf("[DEBUG] Successfully added peer '%s' to openstack_networking_bgp_speaker_v2 '%s'", v, d.Id())
		}
	}

	return resourceNetworkingBGPSpeakerV2Read(ctx, d, meta)
}

func resourceNetworkingBGPSpeakerV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	log.Printf("[DEBUG] Deleting openstack_networking_bgp_speaker_v2 %s", d.Id())

	err = speakers.Delete(ctx, networkingClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_networking_bgp_speaker_v2"))
	}

	d.SetId("")

	return nil
}

func flattenBGPSpeakerAdvertisedRoutes(routes []speakers.AdvertisedRoute) []map[string]any {
	flattened := make([]map[string]any, 0, len(routes))

	for _, route := range routes {
		flattenedRoute := map[string]any{
			"destination": route.Destination,
			"next_hop":    route.NextHop,
		}
		flattened = append(flattened, flattenedRoute)
	}

	return flattened
}
