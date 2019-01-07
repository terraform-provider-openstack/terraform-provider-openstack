package openstack

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/extradhcpopts"
	p "github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
)

func dataSourceNetworkingPortV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNetworkingPortV2Read,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"port_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"admin_state_up": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},

			"network_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"tenant_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"device_owner": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"mac_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"device_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"fixed_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"security_group_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"tags": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"allowed_address_pairs": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Set:      allowedAddressPairsHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ip_address": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"mac_address": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"all_fixed_ips": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"all_security_group_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"all_tags": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"extra_dhcp_option": &schema.Schema{
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"value": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_version": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceNetworkingPortV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))

	listOpts := p.ListOpts{}

	if v, ok := d.GetOk("port_id"); ok {
		listOpts.ID = v.(string)
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

	pages, err := p.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return fmt.Errorf("Unable to list Ports: %s", err)
	}

	allPorts, err := p.ExtractPorts(pages)
	if err != nil {
		return fmt.Errorf("Unable to retrieve Ports: %s", err)
	}

	if len(allPorts) == 0 {
		return fmt.Errorf("No Port found")
	}

	var port struct {
		p.Port
		extradhcpopts.ExtraDHCPOptsExt
	}
	var ports []p.Port

	// Create a slice of all returned Fixed IPs.
	// This will be in the order returned by the API,
	// which is usually alpha-numeric.
	if v, ok := d.GetOk("fixed_ip"); ok {
		for _, p := range allPorts {
			var ips = []string{}
			for _, ipObject := range p.FixedIPs {
				ips = append(ips, ipObject.IPAddress)
				if v == ipObject.IPAddress {
					ports = append(ports, p)
				}
			}
			if len(ports) > 0 && len(ips) > 0 {
				d.Set("all_fixed_ips", ips)
				break
			}
		}
		if len(ports) == 0 {
			return fmt.Errorf("No Port found after the 'fixed_ip' filter")
		}
	} else {
		ports = allPorts
	}

	v := d.Get("security_group_ids").(*schema.Set)
	securityGroups := resourcePortSecurityGroupsV2(v)
	if len(securityGroups) > 0 {
		var sgPorts []p.Port
		for _, p := range ports {
			for _, sg := range p.SecurityGroups {
				if strSliceContains(securityGroups, sg) {
					sgPorts = append(sgPorts, p)
				}
			}
		}
		if len(sgPorts) == 0 {
			return fmt.Errorf("No Port found after the 'security_group_ids' filter")
		}
		ports = sgPorts
	}

	if len(ports) > 1 {
		return fmt.Errorf("More than one Port found (%d). Try to use in a combination with the 'openstack_networking_port_ids_v2' data source", len(ports))
	}

	err = p.Get(networkingClient, ports[0].ID).ExtractInto(&port)
	if err != nil {
		return fmt.Errorf("No Port found: %s", err)
	}

	log.Printf("[DEBUG] Retrieved Port %s: %+v", port.ID, port)
	d.SetId(port.ID)

	d.Set("port_id", port.ID)
	d.Set("name", port.Name)
	d.Set("description", port.Description)
	d.Set("admin_state_up", port.AdminStateUp)
	d.Set("network_id", port.NetworkID)
	d.Set("tenant_id", port.TenantID)
	d.Set("project_id", port.ProjectID)
	d.Set("device_owner", port.DeviceOwner)
	d.Set("mac_address", port.MACAddress)
	d.Set("device_id", port.DeviceID)
	d.Set("region", GetRegion(d, config))
	d.Set("all_tags", port.Tags)
	d.Set("all_security_group_ids", port.SecurityGroups)
	d.Set("allowed_address_pairs", flattenNetworkingPortAllowedAddressPairsV2(port.MACAddress, port.AllowedAddressPairs))
	d.Set("extra_dhcp_option", flattenNetworkingPortDHCPOptsV2(port.ExtraDHCPOptsExt))

	return nil
}
