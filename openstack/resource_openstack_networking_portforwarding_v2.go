package openstack

import (
	"context"
	"log"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/portforwarding"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNetworkingPortForwardingV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkPortForwardingV2Create,
		ReadContext:   resourceNetworkPortForwardingV2Read,
		UpdateContext: resourceNetworkPortForwardingV2Update,
		DeleteContext: resourceNetworkPortForwardingV2Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"floatingip_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"internal_port_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"internal_ip_address": {
				Type:     schema.TypeString,
				Required: true,
			},

			"internal_port": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"internal_port_range": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"external_port": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"external_port_range": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"protocol": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceNetworkPortForwardingV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	fipID := d.Get("floatingip_id").(string)
	createOpts := portforwarding.CreateOpts{
		InternalIPAddress: d.Get("internal_ip_address").(string),
		InternalPortID:    d.Get("internal_port_id").(string),
		Protocol:          d.Get("protocol").(string),
		Description:       d.Get("description").(string),
	}

	var base_ports = 0
	var range_ports = 0
	if v, ok := d.GetOk("external_port"); ok {
		if v.(int) > 0 {
			createOpts.ExternalPort = v.(int)
			base_ports += 1
		}
	}

	if v, ok := d.GetOk("internal_port"); ok {
		if v.(int) > 0 {
			createOpts.InternalPort = v.(int)
			base_ports += 1
		}
	}

	if v, ok := d.GetOk("external_port_range"); ok {
		if v.(string) != "" {
			createOpts.ExternalPortRange = v.(string)
			range_ports += 1
		}
	}

	if v, ok := d.GetOk("internal_port_range"); ok {
		if v.(string) != "" {
			createOpts.InternalPortRange = v.(string)
			range_ports += 1
		}
	}

	log.Printf("[DEBUG] openstack_networking_portforwarding_v2 create options: %#v", createOpts)

	if ! (base_ports == 2 || range_ports == 2) {
		err := "Either external_port/internal_port or external_port_range/internal_port_range must be specified"
		return diag.Errorf("Error creating openstack_networking_portforwarding_v2: %s", err)
	}

	pf, err := portforwarding.Create(ctx, networkingClient, fipID, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_networking_portforwarding_v2: %s", err)
	}

	log.Printf("[DEBUG] Waiting for openstack_networking_portforwarding_v2 %s to become available.", pf.ID)

	stateConf := &retry.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Refresh:    networkingPortForwardingV2StateRefreshFunc(ctx, networkingClient, fipID, pf.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_networking_portforwarding_v2 %s to become available: %s", pf.ID, err)
	}

	d.SetId(pf.ID)

	log.Printf("[DEBUG] Created openstack_networking_portforwarding_v2 %s: %#v", pf.ID, pf)

	return resourceNetworkPortForwardingV2Read(ctx, d, meta)
}

func resourceNetworkPortForwardingV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	fipID := d.Get("floatingip_id").(string)

	pf, err := portforwarding.Get(ctx, networkingClient, fipID, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting openstack_networking_portforwarding_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_networking_portforwarding_v2 %s: %#v", d.Id(), pf)

	d.Set("id", pf.ID)
	d.Set("description", pf.Description)
	d.Set("internal_port_id", pf.InternalPortID)
	d.Set("internal_ip_address", pf.InternalIPAddress)
	d.Set("protocol", pf.Protocol)
	d.Set("region", GetRegion(d, config))

	if pf.InternalPort > 0 {
		d.Set("internal_port", pf.InternalPort)
	}
	if pf.InternalPortRange != "" {
		d.Set("internal_port_range", pf.InternalPortRange)
	}
	if pf.ExternalPort > 0 {
		d.Set("external_port", pf.ExternalPort)
	}
	if pf.ExternalPortRange != "" {
		d.Set("external_port_range", pf.ExternalPortRange)
	}

	return nil
}

func resourceNetworkPortForwardingV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var hasChange bool

	var updateOpts portforwarding.UpdateOpts

	fipID := d.Get("floatingip_id").(string)

	if d.HasChange("internal_port_id") {
		hasChange = true
		internalPortID := d.Get("internal_port_id").(string)
		updateOpts.InternalPortID = internalPortID
	}

	if d.HasChange("external_port") {
		hasChange = true
		externalPort := d.Get("external_port").(int)
		if externalPort > 0 {
			updateOpts.ExternalPort = externalPort
		}
	}

	if d.HasChange("internal_port") {
		hasChange = true
		internalPort := d.Get("internal_port").(int)
		if internalPort > 0 {
			updateOpts.InternalPort = internalPort
		}
	}

	if d.HasChange("external_port_range") {
		hasChange = true
		externalPortRange := d.Get("external_port_range").(string)
		updateOpts.ExternalPortRange = externalPortRange
	}
 
	if d.HasChange("internal_port_range") {
		hasChange = true
		internalPortRange := d.Get("internal_port_range").(string)
		updateOpts.InternalPortRange = internalPortRange
	}

	if d.HasChange("protocol") {
		hasChange = true
		protocol := d.Get("protocol").(string)
		updateOpts.Protocol = protocol
	}

	if d.HasChange("description") {
		hasChange = true
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}

	if hasChange {
		log.Printf("[DEBUG] openstack_networking_portforwarding_v2 %s update options: %#v", d.Id(), updateOpts)

		_, err = portforwarding.Update(ctx, networkingClient, fipID, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_networking_portforwarding_v2 %s: %s", d.Id(), err)
		}
	}

	return resourceNetworkPortForwardingV2Read(ctx, d, meta)
}

func resourceNetworkPortForwardingV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	fipID := d.Get("floatingip_id").(string)
	if err := portforwarding.Delete(ctx, networkingClient, fipID, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_networking_portforwarding_v2"))
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    networkingPortForwardingV2StateRefreshFunc(ctx, networkingClient, fipID, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      0,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_networking_portforwarding_v2 %s to become deleted: %s", d.Id(), err)
	}

	return nil
}
