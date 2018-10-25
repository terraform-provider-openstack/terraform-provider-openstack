package openstack

import (
	"fmt"
	"log"
	"time"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/extradhcpopts"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNetworkingPortExtraDHCPOptionsV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkingPortExtraDHCPOptionsV2Create,
		Read:   resourceNetworkingPortExtraDHCPOptionsV2Read,
		Delete: resourceNetworkingPortExtraDHCPOptionsV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"port_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"extra_dhcp_opts": &schema.Schema{
				Type:     schema.TypeSet,
				Set:      hashDHCPOptionsV2,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"opt_name": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"opt_value": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"ip_version": &schema.Schema{
							Type:     schema.TypeInt,
							Default:  4,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func resourceNetworkingPortExtraDHCPOptionsV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	portID := d.Get("port_id").(string)
	dhcpOpts := d.Get("extra_dhcp_opts").(*schema.Set)
	extraDHCPOpts := expandDHCPOptionsV2Add(dhcpOpts)

	log.Printf("[DEBUG] Adding DHCP options %+v for port %s", extraDHCPOpts, portID)
	if _, err = ports.Update(networkingClient, portID, extradhcpopts.UpdateOptsExt{
		ports.UpdateOpts{},
		extraDHCPOpts,
	}).Extract(); err != nil {
		return fmt.Errorf("Error updating DHCP options for port: %s", err)
	}

	d.SetId(portID)

	return resourceNetworkingPortExtraDHCPOptionsV2Read(d, meta)
}

func resourceNetworkingPortExtraDHCPOptionsV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var port struct {
		ports.Port
		extradhcpopts.ExtraDHCPOptsExt
	}
	err = ports.Get(networkingClient, d.Id()).ExtractInto(&port)
	if err != nil {
		return CheckDeleted(d, err, "port")
	}

	log.Printf("[DEBUG] Retrieved Port %s: %+v", d.Id(), port)

	extraDHCPOpts := flattenDHCPOptionsV2(port.ExtraDHCPOptsExt)
	d.Set("extra_dhcp_opts", extraDHCPOpts)
	d.Set("port_id", port.ID)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceNetworkingPortExtraDHCPOptionsV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	dhcpOpts := d.Get("extra_dhcp_opts").(*schema.Set)
	extraDHCPOpts := expandDHCPOptionsV2Delete(dhcpOpts)

	log.Printf("[DEBUG] Deleting DHCP options from port %s", d.Id())
	if _, err = ports.Update(networkingClient, d.Id(), extradhcpopts.UpdateOptsExt{
		ports.UpdateOpts{},
		extraDHCPOpts,
	}).Extract(); err != nil {
		return fmt.Errorf("Error updating DHCP options for port: %s", err)
	}

	d.SetId("")
	return nil
}
