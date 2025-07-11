package openstack

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/dns"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/terraform-provider-openstack/utils/v2/hashcode"
)

func dataSourceNetworkingPortIDsV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingPortIDsV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"device_owner": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"mac_address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"device_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"fixed_ip": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsIPAddress,
			},

			"status": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"security_group_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"dns_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"sort_direction": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"asc", "desc",
				}, true),
			},

			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceNetworkingPortIDsV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	listOpts := ports.ListOpts{}

	var listOptsBuilder ports.ListOptsBuilder

	if v, ok := d.GetOk("sort_key"); ok {
		listOpts.SortKey = v.(string)
	}

	if v, ok := d.GetOk("sort_direction"); ok {
		listOpts.SortDir = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		listOpts.Name = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		listOpts.Description = v.(string)
	}

	if v, ok := getOkExists(d, "admin_state_up"); ok {
		asu := v.(bool)
		listOpts.AdminStateUp = &asu
	}

	if v, ok := d.GetOk("network_id"); ok {
		listOpts.NetworkID = v.(string)
	}

	if v, ok := d.GetOk("status"); ok {
		listOpts.Status = v.(string)
	}

	if v, ok := d.GetOk("tenant_id"); ok {
		listOpts.TenantID = v.(string)
	}

	if v, ok := d.GetOk("project_id"); ok {
		listOpts.ProjectID = v.(string)
	}

	if v, ok := d.GetOk("device_owner"); ok {
		listOpts.DeviceOwner = v.(string)
	}

	if v, ok := d.GetOk("mac_address"); ok {
		listOpts.MACAddress = v.(string)
	}

	if v, ok := d.GetOk("device_id"); ok {
		listOpts.DeviceID = v.(string)
	}

	tags := networkingV2AttributesTags(d)
	if len(tags) > 0 {
		listOpts.Tags = strings.Join(tags, ",")
	}

	listOptsBuilder = listOpts

	if v, ok := d.GetOk("dns_name"); ok {
		listOptsBuilder = dns.PortListOptsExt{
			ListOptsBuilder: listOptsBuilder,
			DNSName:         v.(string),
		}
	}

	allPages, err := ports.List(networkingClient, listOptsBuilder).AllPages(ctx)
	if err != nil {
		return diag.Errorf("Unable to list openstack_networking_port_ids_v2: %s", err)
	}

	allPorts, err := ports.ExtractPorts(allPages)
	if err != nil {
		return diag.Errorf("Unable to retrieve openstack_networking_port_ids_v2: %s", err)
	}

	if len(allPorts) == 0 {
		log.Printf("[DEBUG] No ports in openstack_networking_port_ids_v2 found")
	}

	portsList := make([]ports.Port, 0, len(allPorts))
	portIDs := make([]string, 0, len(allPorts))

	// Filter returned Fixed IPs by a "fixed_ip".
	if v, ok := d.GetOk("fixed_ip"); ok {
		for _, p := range allPorts {
			for _, ipObject := range p.FixedIPs {
				if v.(string) == ipObject.IPAddress {
					portsList = append(portsList, p)
				}
			}
		}

		if len(portsList) == 0 {
			log.Printf("[DEBUG] No ports in openstack_networking_port_ids_v2 found after the 'fixed_ip' filter")
		}
	} else {
		portsList = allPorts
	}

	securityGroups := expandToStringSlice(d.Get("security_group_ids").(*schema.Set).List())
	if len(securityGroups) > 0 {
		var sgPorts []ports.Port

		for _, p := range portsList {
			for _, sg := range p.SecurityGroups {
				if strSliceContains(securityGroups, sg) {
					sgPorts = append(sgPorts, p)
				}
			}
		}

		if len(sgPorts) == 0 {
			log.Printf("[DEBUG] No ports in openstack_networking_port_ids_v2 found after the 'security_group_ids' filter")
		}

		portsList = sgPorts
	}

	for _, p := range portsList {
		portIDs = append(portIDs, p.ID)
	}

	log.Printf("[DEBUG] Retrieved %d ports in openstack_networking_port_ids_v2: %+v", len(portsList), portsList)

	d.SetId(strconv.Itoa(hashcode.String(strings.Join(portIDs, ""))))
	d.Set("ids", portIDs)
	d.Set("region", GetRegion(d, config))

	return nil
}
