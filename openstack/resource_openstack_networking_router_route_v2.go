package openstack

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/routers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNetworkingRouterRouteV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingRouterRouteV2Create,
		ReadContext:   resourceNetworkingRouterRouteV2Read,
		DeleteContext: resourceNetworkingRouterRouteV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

			"destination_cidr": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"next_hop": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceNetworkingRouterRouteV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	routerID := d.Get("router_id").(string)
	config.Lock(routerID)
	defer config.Unlock(routerID)

	r, err := routers.Get(ctx, networkingClient, routerID).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting openstack_networking_router_route_v2 router"))
	}

	log.Printf("[DEBUG] Retrieved openstack_networking_router_route_v2 router %s: %#v", routerID, r)

	routes := r.Routes
	dstCIDR := d.Get("destination_cidr").(string)
	nextHop := d.Get("next_hop").(string)
	exists := false

	for _, route := range routes {
		if route.DestinationCIDR == dstCIDR && route.NextHop == nextHop {
			exists = true

			break
		}
	}

	if exists {
		log.Printf("[DEBUG] openstack_networking_router_route_v2 %s already has route to %s via %s", routerID, dstCIDR, nextHop)

		return resourceNetworkingRouterRouteV2Read(ctx, d, meta)
	}

	routes = append(routes, routers.Route{
		DestinationCIDR: dstCIDR,
		NextHop:         nextHop,
	})
	updateOpts := routers.UpdateOpts{
		Routes: &routes,
	}
	log.Printf("[DEBUG] openstack_networking_router_route_v2 %s update options: %#v", routerID, updateOpts)

	_, err = routers.Update(ctx, networkingClient, routerID, updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error updating openstack_networking_router_route_v2: %s", err)
	}

	d.SetId(resourceNetworkingRouterRouteV2BuildID(routerID, dstCIDR, nextHop))

	return resourceNetworkingRouterRouteV2Read(ctx, d, meta)
}

func resourceNetworkingRouterRouteV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	idFromResource, dstCIDR, nextHop, err := resourceNetworkingRouterRouteV2ParseID(d.Id())
	if err != nil {
		return diag.Errorf("Error reading openstack_networking_router_route_v2 ID %s: %s", d.Id(), err)
	}

	routerID := d.Get("router_id").(string)
	if routerID == "" {
		routerID = idFromResource
	}

	d.Set("router_id", routerID)

	r, err := routers.Get(ctx, networkingClient, routerID).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting openstack_networking_router_route_v2 router"))
	}

	log.Printf("[DEBUG] Retrieved openstack_networking_router_route_v2 %s: %#v", routerID, r)

	exists := false

	for _, route := range r.Routes {
		if route.DestinationCIDR == dstCIDR && route.NextHop == nextHop {
			d.Set("destination_cidr", dstCIDR)
			d.Set("next_hop", nextHop)

			exists = true

			break
		}
	}

	if !exists {
		log.Printf("[DEBUG] openstack_networking_router_route_v2 %s does not have route to %s via %s", routerID, dstCIDR, nextHop)
		d.SetId("")

		return nil
	}

	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceNetworkingRouterRouteV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	routerID := d.Get("router_id").(string)
	config.Lock(routerID)
	defer config.Unlock(routerID)

	r, err := routers.Get(ctx, networkingClient, routerID).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting openstack_networking_router_route_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_networking_router_route_v2 %s: %#v", routerID, r)

	dstCIDR := d.Get("destination_cidr").(string)
	nextHop := d.Get("next_hop").(string)

	oldRoutes := r.Routes
	newRoute := make([]routers.Route, 0, len(oldRoutes))

	for _, route := range oldRoutes {
		if route.DestinationCIDR != dstCIDR || route.NextHop != nextHop {
			newRoute = append(newRoute, route)
		}
	}

	if len(oldRoutes) == len(newRoute) {
		log.Printf("[DEBUG] openstack_networking_router_route_v2 %s does not have route to %s via %s", routerID, dstCIDR, nextHop)
		d.SetId("")

		return nil
	}

	log.Printf("[DEBUG] Deleting openstack_networking_router_route_v2 %s route to %s via %s", routerID, dstCIDR, nextHop)

	updateOpts := routers.UpdateOpts{
		Routes: &newRoute,
	}

	_, err = routers.Update(ctx, networkingClient, routerID, updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error updating openstack_networking_router_route_v2: %s", err)
	}

	return nil
}
