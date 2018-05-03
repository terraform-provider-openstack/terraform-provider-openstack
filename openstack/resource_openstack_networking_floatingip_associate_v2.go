package openstack

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
)

func resourceNetworkingFloatingIPAssociateV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceNetworkingFloatingIPAssociateV2Create,
		Read:   resourceNetworkingFloatingIPAssociateV2Read,
		Delete: resourceNetworkingFloatingIPAssociateV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"floatingip_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"port_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNetworkingFloatingIPAssociateV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack network client: %s", err)
	}

	floatingIPID := d.Get("floatingip_id").(string)
	portID := d.Get("port_id").(string)

	updateOpts := floatingips.UpdateOpts{
		PortID: &portID,
	}

	log.Printf("[DEBUG] Floating IP Associate Create Options: %#v", updateOpts)

	floatingIP, err := floatingips.Update(networkingClient, floatingIPID, updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error associating floating IP %s to port %s: %s",
			floatingIPID, portID, err)
	}

	d.SetId(floatingIP.ID)

	return resourceNetworkFloatingIPV2Read(d, meta)
}

func resourceNetworkingFloatingIPAssociateV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack network client: %s", err)
	}

	floatingIP, err := floatingips.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "floating IP")
	}

	d.Set("floatingip_id", floatingIP.ID)
	d.Set("port_id", floatingIP.PortID)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceNetworkingFloatingIPAssociateV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	networkingClient, err := config.networkingV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack network client: %s", err)
	}

	floatingIPID := d.Get("floatingip_id").(string)
	portID := d.Get("port_id").(string)

	updateOpts := floatingips.UpdateOpts{
		PortID: nil,
	}

	log.Printf("[DEBUG] Floating IP Delete Options: %#v", updateOpts)

	_, err = floatingips.Update(networkingClient, floatingIPID, updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error disassociating floating IP %s from port %s: %s",
			floatingIPID, portID, err)
	}

	return nil
}
