package openstack

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNetworkingRouterInterfaceV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingRouterInterfaceV2Create,
		ReadContext:   resourceNetworkingRouterInterfaceV2Read,
		UpdateContext: resourceNetworkingRouterInterfaceV2Update,
		DeleteContext: resourceNetworkingRouterInterfaceV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

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

			"router_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"subnet_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"port_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"force_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceNetworkingRouterInterfaceV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	createOpts := routers.AddInterfaceOpts{
		SubnetID: d.Get("subnet_id").(string),
		PortID:   d.Get("port_id").(string),
	}

	routerID := d.Get("router_id").(string)
	// the lock is necessary, when multiple interfaces are added to the same router in parallel
	config.Lock(routerID)
	defer config.Unlock(routerID)

	log.Printf("[DEBUG] openstack_networking_router_interface_v2 create options: %#v", createOpts)

	r, err := routers.AddInterface(ctx, networkingClient, routerID, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_networking_router_interface_v2: %s", err)
	}

	log.Printf("[DEBUG] Waiting for openstack_networking_router_interface_v2 %s to become available", r.PortID)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"BUILD", "PENDING_CREATE", "PENDING_UPDATE"},
		Target:     []string{"ACTIVE", "DOWN"},
		Refresh:    resourceNetworkingRouterInterfaceV2StateRefreshFunc(ctx, networkingClient, r.PortID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_networking_router_interface_v2 %s to become available: %s", r.ID, err)
	}

	d.SetId(r.PortID)

	log.Printf("[DEBUG] Created openstack_networking_router_interface_v2 %s: %#v", r.ID, r)

	return resourceNetworkingRouterInterfaceV2Read(ctx, d, meta)
}

func resourceNetworkingRouterInterfaceV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	r, err := ports.Get(ctx, networkingClient, d.Id()).Extract()
	if err != nil {
		if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
			d.SetId("")

			return nil
		}

		return diag.Errorf("Error retrieving openstack_networking_router_interface_v2: %s", err)
	}

	log.Printf("[DEBUG] Retrieved openstack_networking_router_interface_v2 %s: %#v", d.Id(), r)

	d.Set("router_id", r.DeviceID)
	d.Set("port_id", r.ID)
	d.Set("region", GetRegion(d, config))

	// Set the subnet ID by looking at the port's FixedIPs.
	// If there's more than one FixedIP, do not set the subnet
	// as it's not possible to confidently determine which subnet
	// belongs to this interface. However, that situation should
	// not happen.
	if len(r.FixedIPs) != 1 {
		log.Printf("[DEBUG] Unable to set openstack_networking_router_interface_v2 %s subnet_id", d.Id())
	} else {
		d.Set("subnet_id", r.FixedIPs[0].SubnetID)
	}

	return nil
}

func resourceNetworkingRouterInterfaceV2Update(_ context.Context, _ *schema.ResourceData, _ any) diag.Diagnostics {
	return nil
}

func resourceNetworkingRouterInterfaceV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    resourceNetworkingRouterInterfaceV2DeleteRefreshFunc(ctx, networkingClient, d),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      0,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_networking_router_interface_v2 %s to Delete:  %s", d.Id(), err)
	}

	d.SetId("")

	return nil
}
