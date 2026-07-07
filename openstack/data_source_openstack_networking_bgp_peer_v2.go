package openstack

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/bgp/peers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNetworkingBGPPeerV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingBGPPeerV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"peer_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"auth_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"peer_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"remote_as": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceNetworkingBGPPeerV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var peer *peers.BGPPeer

	if v, ok := d.GetOk("peer_id"); ok {
		peer, err = peers.Get(ctx, networkingClient, v.(string)).Extract()
		if err != nil {
			return diag.Errorf("Unable to retrieve openstack_networking_bgp_peer_v2 with peer_id %v: %s", v, err)
		}
	} else {
		pages, err := peers.List(networkingClient).AllPages(ctx)
		if err != nil {
			return diag.Errorf("Unable to list openstack_networking_bgp_peer_v2: %s", err)
		}

		allPeers, err := peers.ExtractBGPPeers(pages)
		if err != nil {
			return diag.Errorf("Unable to retrieve openstack_networking_bgp_peer_v2: %s", err)
		}

		if len(allPeers) < 1 {
			return diag.Errorf("No openstack_networking_bgp_peer_v2 found")
		}

		if v, ok := d.GetOk("name"); ok {
			name := v.(string)
			for _, p := range allPeers {
				if p.Name == name {
					peer = &p

					break
				}
			}

			if peer == nil {
				return diag.Errorf("No openstack_networking_bgp_peer_v2 found with name %v", v)
			}
		} else {
			if len(allPeers) > 1 {
				return diag.Errorf("More than one openstack_networking_bgp_peer_v2 found")
			}

			peer = &allPeers[0]
		}
	}

	if v, ok := d.GetOk("name"); ok {
		if peer.Name != v.(string) {
			return diag.Errorf("No openstack_networking_bgp_peer_v2 found with name %v", v)
		}
	}

	if v, ok := d.GetOk("peer_id"); ok {
		if peer.ID != v.(string) {
			return diag.Errorf("No openstack_networking_bgp_peer_v2 found with peer_id %v", v)
		}
	}

	log.Printf("[DEBUG] Retrieved openstack_networking_bgp_peer_v2 %s: %+v", peer.ID, peer)

	d.SetId(peer.ID)
	d.Set("region", GetRegion(d, config))
	d.Set("name", peer.Name)
	d.Set("auth_type", peer.AuthType)
	d.Set("peer_ip", peer.PeerIP)
	d.Set("remote_as", peer.RemoteAS)
	d.Set("tenant_id", peer.TenantID)

	return nil
}
