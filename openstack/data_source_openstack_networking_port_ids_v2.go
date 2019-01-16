package openstack

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
)

func dataSourceNetworkingPortIDsV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetworkingPortIDsV2Read,

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

func dataSourceNetworkingPortIDsV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))

	listOpts := ports.ListOpts{}

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

	if v, ok := d.GetOkExists("admin_state_up"); ok {
		asu := v.(bool)
		listOpts.AdminStateUp = &asu
	}

	if v, ok := d.GetOk("network_id"); ok {
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

	tags := networkV2AttributesTags(d)
	if len(tags) > 0 {
		listOpts.Tags = strings.Join(tags, ",")
	}

	allPages, err := ports.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return fmt.Errorf("Unable to list openstack_networking_ports_v2: %s", err)
	}

	allPorts, err := ports.ExtractPorts(allPages)
	if err != nil {
		return fmt.Errorf("Unable to retrieve openstack_networking_ports_v2: %s", err)
	}

	if len(allPorts) == 0 {
		return fmt.Errorf("No openstack_networking_port_v2 found")
	}

	var portsList []ports.Port
	var portIDs []string

	// Create a slice of all returned Fixed IPs.
	// This will be in the order returned by the API,
	// which is usually alpha-numeric.
	if v, ok := d.GetOk("fixed_ip"); ok {
		for _, p := range allPorts {
			var ips = []string{}
			for _, ipObject := range p.FixedIPs {
				ips = append(ips, ipObject.IPAddress)
				if v == ipObject.IPAddress {
					portsList = append(portsList, p)
				}
			}
		}
		if len(portsList) == 0 {
			log.Printf("No openstack_networking_port_v2 found after the 'fixed_ip' filter")
			return fmt.Errorf("No openstack_networking_port_v2 found")
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
			log.Printf("No openstack_networking_port_v2 found after the 'security_group_ids' filter")
			return fmt.Errorf("No openstack_networking_port_v2 found")
		}
		portsList = sgPorts
	}

	for _, p := range portsList {
		portIDs = append(portIDs, p.ID)
	}

	log.Printf("[DEBUG] Retrieved %d openstack_networking_port_v2: %+v", len(portsList), portsList)

	d.SetId(fmt.Sprintf("%d", hashcode.String(strings.Join(portIDs, ""))))
	d.Set("ids", portIDs)

	return nil
}
