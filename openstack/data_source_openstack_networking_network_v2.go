package openstack

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/external"
	mtuext "github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/mtu"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/vlantransparent"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/subnets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNetworkingNetworkV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingNetworkV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"matching_subnet_cidr": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"admin_state_up": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"shared": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"external": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"availability_zone_hints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"segments": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"physical_network": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"segmentation_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},

			"transparent_vlan": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"mtu": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"dns_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"subnets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"all_tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceNetworkingNetworkV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	// Prepare basic listOpts.
	var listOpts networks.ListOptsBuilder

	var status string
	if v, ok := d.GetOk("status"); ok {
		status = v.(string)
	}

	listOpts = networks.ListOpts{
		ID:          d.Get("network_id").(string),
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		TenantID:    d.Get("tenant_id").(string),
		Status:      status,
	}

	// Add the external attribute if specified.
	if v, ok := getOkExists(d, "external"); ok {
		isExternal := v.(bool)
		listOpts = external.ListOptsExt{
			ListOptsBuilder: listOpts,
			External:        &isExternal,
		}
	}

	// Add the transparent VLAN attribute if specified.
	if v, ok := getOkExists(d, "transparent_vlan"); ok {
		isVLANTransparent := v.(bool)
		listOpts = vlantransparent.ListOptsExt{
			ListOptsBuilder: listOpts,
			VLANTransparent: &isVLANTransparent,
		}
	}

	// Add the MTU attribute if specified.
	if v, ok := getOkExists(d, "mtu"); ok {
		listOpts = mtuext.ListOptsExt{
			ListOptsBuilder: listOpts,
			MTU:             v.(int),
		}
	}

	tags := networkingV2AttributesTags(d)
	if len(tags) > 0 {
		listOpts = networks.ListOpts{Tags: strings.Join(tags, ",")}
	}

	pages, err := networks.List(networkingClient, listOpts).AllPages(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	// First extract into a normal networks.Network in order to see if
	// there were any results at all.
	tmpAllNetworks, err := networks.ExtractNetworks(pages)
	if err != nil {
		return diag.FromErr(err)
	}

	if len(tmpAllNetworks) < 1 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	var allNetworks []networkExtended

	err = networks.ExtractNetworksInto(pages, &allNetworks)
	if err != nil {
		return diag.Errorf("Unable to retrieve openstack_networking_networks_v2: %s", err)
	}

	var refinedNetworks []networkExtended

	if cidr := d.Get("matching_subnet_cidr").(string); cidr != "" {
		for _, n := range allNetworks {
			for _, s := range n.Subnets {
				subnet, err := subnets.Get(ctx, networkingClient, s).Extract()
				if err != nil {
					if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
						continue
					}

					return diag.Errorf("Unable to retrieve openstack_networking_network_v2 subnet: %s", err)
				}

				if cidr == subnet.CIDR {
					refinedNetworks = append(refinedNetworks, n)
				}
			}
		}
	} else {
		refinedNetworks = allNetworks
	}

	if len(refinedNetworks) < 1 {
		return diag.Errorf("Your query returned no results. " +
			"Please change your search criteria and try again.")
	}

	if len(refinedNetworks) > 1 {
		return diag.Errorf("Your query returned more than one result." +
			" Please try a more specific search criteria")
	}

	network := refinedNetworks[0]

	if err = d.Set("availability_zone_hints", network.AvailabilityZoneHints); err != nil {
		log.Printf("[DEBUG] Unable to set availability_zone_hints for openstack_networking_network_v2 %s: %s", network.ID, err)
	}

	log.Printf("[DEBUG] Retrieved openstack_networking_network_v2 %s: %+v", network.ID, network)
	d.SetId(network.ID)

	d.Set("name", network.Name)
	d.Set("description", network.Description)
	d.Set("admin_state_up", strconv.FormatBool(network.AdminStateUp))
	d.Set("shared", strconv.FormatBool(network.Shared))
	d.Set("external", network.External)
	d.Set("tenant_id", network.TenantID)
	d.Set("segments", flattenNetworkingNetworkSegmentsV2(network))
	d.Set("transparent_vlan", network.VLANTransparent)
	d.Set("subnets", network.Subnets)
	d.Set("all_tags", network.Tags)
	d.Set("mtu", network.MTU)
	d.Set("dns_domain", network.DNSDomain)
	d.Set("region", GetRegion(d, config))

	return nil
}
