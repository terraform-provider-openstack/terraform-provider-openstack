package openstack

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/subnets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNetworkingRouterInterfaceV2StateRefreshFunc(ctx context.Context, networkingClient *gophercloud.ServiceClient, portID string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		r, err := ports.Get(ctx, networkingClient, portID).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return r, "DELETED", nil
			}

			return r, "", err
		}

		return r, r.Status, nil
	}
}

func resourceNetworkingRouterInterfaceV2DeleteRefreshFunc(ctx context.Context, networkingClient *gophercloud.ServiceClient, d *schema.ResourceData) retry.StateRefreshFunc {
	return func() (any, string, error) {
		routerID := d.Get("router_id").(string)
		routerInterfaceID := d.Id()

		log.Printf("[DEBUG] Attempting to delete openstack_networking_router_interface_v2 %s", routerInterfaceID)

		removeOpts := routers.RemoveInterfaceOpts{
			SubnetID: d.Get("subnet_id").(string),
			PortID:   d.Get("port_id").(string),
		}

		if removeOpts.SubnetID != "" {
			// We need to make sure to only send subnet_id, because the port may have multiple
			// openstack_networking_router_interface_v2 attached. Otherwise openstack would delete them too.
			removeOpts.PortID = ""
		}

		_, err := routers.RemoveInterface(ctx, networkingClient, routerID, removeOpts).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				log.Printf("[DEBUG] Successfully deleted openstack_networking_router_interface_v2 %s", routerInterfaceID)

				return "", "DELETED", nil
			}

			if gophercloud.ResponseCodeIs(err, http.StatusConflict) {
				if d.Get("force_destroy").(bool) {
					// The router may have routes preventing the interface to be deleted.
					// Check which routes correspond to a particular router interface.
					var updateRoutes []routers.Route

					if removeOpts.SubnetID != "" {
						// get subnet CIDR
						subnet, err := subnets.Get(ctx, networkingClient, removeOpts.SubnetID).Extract()
						if err != nil {
							return "", "ACTIVE", err
						}

						_, cidr, err := net.ParseCIDR(subnet.CIDR)
						if err != nil {
							return "", "ACTIVE", err
						}
						// determine which routes must be removed
						router, err := routers.Get(ctx, networkingClient, routerID).Extract()
						if err != nil {
							return "", "ACTIVE", err
						}

						for _, route := range router.Routes {
							if ip := net.ParseIP(route.NextHop); ip != nil && !cidr.Contains(ip) {
								updateRoutes = append(updateRoutes, route)
							}
						}
					}

					log.Printf("[DEBUG] Attempting to forceDestroy openstack_networking_router_interface_v2 '%s': %+v", d.Id(), err)

					opts := &routers.UpdateOpts{
						Routes: &updateRoutes,
					}

					_, err := routers.Update(ctx, networkingClient, routerID, opts).Extract()
					if err != nil {
						return "", "ACTIVE", err
					}
				}

				log.Printf("[DEBUG] openstack_networking_router_interface_v2 %s is still in use", routerInterfaceID)

				return "", "ACTIVE", nil
			}

			return "", "ACTIVE", err
		}

		_, err = ports.Get(ctx, networkingClient, routerInterfaceID).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				log.Printf("[DEBUG] Successfully deleted openstack_networking_router_interface_v2 %s", routerInterfaceID)

				return "", "DELETED", nil
			}

			return "", "ACTIVE", err
		}

		log.Printf("[DEBUG] openstack_networking_router_interface_v2 %s is still active", routerInterfaceID)

		return "", "ACTIVE", nil
	}
}
