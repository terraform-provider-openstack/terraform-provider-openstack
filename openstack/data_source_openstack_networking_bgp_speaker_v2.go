package openstack

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/agents"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/bgp/speakers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNetworkingBGPSpeakerV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingBGPSpeakerV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"speaker_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ip_version": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"advertise_floating_ip_host_routes": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"advertise_tenant_networks": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"local_as": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"networks": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"peers": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"dragents": {
				Type:     schema.TypeSet,
				Computed: true,
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

func dataSourceNetworkingBGPSpeakerV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var speaker *speakers.BGPSpeaker

	if v, ok := d.GetOk("speaker_id"); ok {
		speaker, err = speakers.Get(ctx, networkingClient, v.(string)).Extract()
		if err != nil {
			return diag.Errorf("Unable to retrieve openstack_networking_bgp_speaker_v2 with speaker_id %v: %s", v, err)
		}
	} else {
		pages, err := speakers.List(networkingClient).AllPages(ctx)
		if err != nil {
			return diag.Errorf("Unable to list openstack_networking_bgp_speaker_v2: %s", err)
		}

		allSpeakers, err := speakers.ExtractBGPSpeakers(pages)
		if err != nil {
			return diag.Errorf("Unable to retrieve openstack_networking_bgp_speaker_v2: %s", err)
		}

		if len(allSpeakers) < 1 {
			return diag.Errorf("No openstack_networking_bgp_speaker_v2 found")
		}

		if v, ok := d.GetOk("name"); ok {
			name := v.(string)
			for _, s := range allSpeakers {
				if s.Name == name {
					speaker = &s

					break
				}
			}

			if speaker == nil {
				return diag.Errorf("No openstack_networking_bgp_speaker_v2 found with name %v", v)
			}
		} else {
			if len(allSpeakers) > 1 {
				return diag.Errorf("More than one openstack_networking_bgp_speaker_v2 found")
			}

			speaker = &allSpeakers[0]
		}
	}

	if v, ok := d.GetOk("speaker_id"); ok {
		if speaker.ID != v.(string) {
			return diag.Errorf("No openstack_networking_bgp_speaker_v2 found with speaker_id %v", v)
		}
	}

	if v, ok := d.GetOk("name"); ok {
		if speaker.Name != v.(string) {
			return diag.Errorf("No openstack_networking_bgp_speaker_v2 found with name %v", v)
		}
	}

	log.Printf("[DEBUG] Retrieved openstack_networking_bgp_speaker_v2 %s: %+v", speaker.ID, speaker)

	allPages, err := speakers.GetAdvertisedRoutes(networkingClient, speaker.ID).AllPages(ctx)
	if err != nil {
		return diag.Errorf("Error retrieving advertised routes for openstack_networking_bgp_speaker_v2: %s", err)
	}

	advertisedRoutes, err := speakers.ExtractAdvertisedRoutes(allPages)
	if err != nil {
		return diag.Errorf("Error extracting advertised routes for openstack_networking_bgp_speaker_v2: %s", err)
	}

	log.Printf("[DEBUG] Retrieved advertised routes for openstack_networking_bgp_speaker_v2: %#v", advertisedRoutes)

	allPages, err = agents.ListDRAgentHostingBGPSpeakers(networkingClient, speaker.ID).AllPages(ctx)
	if err != nil {
		return diag.Errorf("Error retrieving dragents for openstack_networking_bgp_speaker_v2: %s", err)
	}

	dragents, err := agents.ExtractAgents(allPages)
	if err != nil {
		return diag.Errorf("Error extracting dragents for openstack_networking_bgp_speaker_v2: %s", err)
	}

	log.Printf("[DEBUG] Retrieved dragents for openstack_networking_bgp_speaker_v2: %#v", advertisedRoutes)

	d.SetId(speaker.ID)
	d.Set("region", GetRegion(d, config))
	d.Set("tenant_id", speaker.TenantID)
	d.Set("name", speaker.Name)
	d.Set("ip_version", speaker.IPVersion)
	d.Set("advertise_floating_ip_host_routes", speaker.AdvertiseFloatingIPHostRoutes)
	d.Set("advertise_tenant_networks", speaker.AdvertiseTenantNetworks)
	d.Set("local_as", speaker.LocalAS)
	d.Set("networks", speaker.Networks)
	d.Set("peers", speaker.Peers)
	d.Set("advertised_routes", flattenBGPSpeakerAdvertisedRoutes(advertisedRoutes))
	d.Set("dragents", flattenNetworkingAgents(dragents))

	return nil
}
