package openstack

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/extraroutes"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/routers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNetworkingRouterRoutesV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingRouterRoutesV2Create,
		ReadContext:   resourceNetworkingRouterRoutesV2Read,
		UpdateContext: resourceNetworkingRouterRoutesV2Update,
		DeleteContext: resourceNetworkingRouterRoutesV2Delete,
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

			"routes": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"destination_cidr": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsCIDR,
						},
						"next_hop": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validation.IsIPAddress,
						},
					},
				},
			},
		},
	}
}

func resourceNetworkingRouterRoutesV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	routerID := d.Get("router_id").(string)
	routes := flattenNetworkingRouterRoutesV2(d.Get("routes").(*schema.Set).List())
	opts := extraroutes.Opts{
		Routes: &routes,
	}
	log.Printf("[DEBUG] openstack_networking_router_routes_v2 %s update options: %#v", routerID, opts)

	_, err = extraroutes.Add(ctx, networkingClient, routerID, opts).Extract()
	if err != nil {
		return diag.Errorf("Error updating openstack_networking_router_routes_v2: %s", err)
	}

	d.SetId(routerID)

	return resourceNetworkingRouterRoutesV2Read(ctx, d, meta)
}

func resourceNetworkingRouterRoutesV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	router, err := routers.Get(ctx, networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting openstack_networking_router_routes_v2 router"))
	}

	d.Set("router_id", router.ID)
	d.Set("routes", expandNetworkingRouterRoutesV2(router.Routes))
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceNetworkingRouterRoutesV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	if d.HasChange("routes") {
		o, n := d.GetChange("routes")
		oldRoutes, newRoutes := o.(*schema.Set), n.(*schema.Set)
		routesToDel := oldRoutes.Difference(newRoutes)
		routesToAdd := newRoutes.Difference(oldRoutes)

		if v := routesToDel.List(); len(v) > 0 {
			log.Printf("[DEBUG] Removing routes '%s' from openstack_networking_router_routes_v2 '%s'", v, d.Get("name"))
			v := flattenNetworkingRouterRoutesV2(v)
			opts := extraroutes.Opts{
				Routes: &v,
			}

			_, err = extraroutes.Remove(ctx, networkingClient, d.Id(), opts).Extract()
			if err != nil {
				return diag.Errorf("Error deleting routes '%s' from openstack_networking_router_routes_v2: %s", v, err)
			}
		}

		if v := routesToAdd.List(); len(v) > 0 {
			log.Printf("[DEBUG] Adding routes '%s' to openstack_networking_router_routes_v2 '%s'", v, d.Get("name"))
			v := flattenNetworkingRouterRoutesV2(v)
			opts := extraroutes.Opts{
				Routes: &v,
			}

			_, err = extraroutes.Add(ctx, networkingClient, d.Id(), opts).Extract()
			if err != nil {
				return diag.Errorf("Error adding routes '%s' to openstack_networking_router_routes_v2: %s", v, err)
			}
		}
	}

	return resourceNetworkingRouterRoutesV2Read(ctx, d, meta)
}

func resourceNetworkingRouterRoutesV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	routes := flattenNetworkingRouterRoutesV2(d.Get("routes").(*schema.Set).List())
	if len(routes) == 0 {
		log.Printf("[DEBUG] No routes to delete for openstack_networking_router_routes_v2 %s", d.Id())
		d.SetId("")

		return nil
	}

	opts := extraroutes.Opts{
		Routes: &routes,
	}
	_, err = extraroutes.Remove(ctx, networkingClient, d.Id(), opts).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_networking_router_routes_v2"))
	}

	d.SetId("")

	return nil
}

func flattenNetworkingRouterRoutesV2(raw []any) []routers.Route {
	if len(raw) == 0 {
		return nil
	}

	routes := make([]routers.Route, 0, len(raw))

	for _, r := range raw {
		routeMap := r.(map[string]any)
		route := routers.Route{
			DestinationCIDR: routeMap["destination_cidr"].(string),
			NextHop:         routeMap["next_hop"].(string),
		}
		routes = append(routes, route)
	}

	return routes
}

func expandNetworkingRouterRoutesV2(routes []routers.Route) []map[string]any {
	if len(routes) == 0 {
		return nil
	}

	expanded := make([]map[string]any, 0, len(routes))

	for _, r := range routes {
		routeMap := map[string]any{
			"destination_cidr": r.DestinationCIDR,
			"next_hop":         r.NextHop,
		}
		expanded = append(expanded, routeMap)
	}

	return expanded
}
