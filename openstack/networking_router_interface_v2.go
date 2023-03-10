package openstack

import (
	"log"
	"net"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
)

func resourceNetworkingRouterInterfaceV2StateRefreshFunc(networkingClient *gophercloud.ServiceClient, portID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		r, err := ports.Get(networkingClient, portID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return r, "DELETED", nil
			}

			return r, "", err
		}

		return r, r.Status, nil
	}
}

func resourceNetworkingRouterInterfaceV2DeleteRefreshFunc(networkingClient *gophercloud.ServiceClient, d *schema.ResourceData) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
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

		r, err := ports.Get(networkingClient, routerInterfaceID).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted openstack_networking_router_interface_v2 %s", routerInterfaceID)
				return r, "DELETED", nil
			}
			return r, "ACTIVE", err
		}

		_, err = routers.RemoveInterface(networkingClient, routerID, removeOpts).Extract()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				log.Printf("[DEBUG] Successfully deleted openstack_networking_router_interface_v2 %s", routerInterfaceID)
				return r, "DELETED", nil
			}
			if _, ok := err.(gophercloud.ErrDefault409); ok {
				if ok && d.Get("force_destroy").(bool) {
					// The router may have routes preventing the interface to be deleted.
					// Check which routes correspond to a particular router interface.
					var updateRoutes []routers.Route
					if removeOpts.SubnetID != "" {
						// get subnet CIDR
						subnet, err := subnets.Get(networkingClient, removeOpts.SubnetID).Extract()
						if err != nil {
							return r, "ACTIVE", err
						}
						_, cidr, err := net.ParseCIDR(subnet.CIDR)
						if err != nil {
							return r, "ACTIVE", err
						}
						// determine which routes must be removed
						router, err := routers.Get(networkingClient, routerID).Extract()
						if err != nil {
							return r, "ACTIVE", err
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
					_, err := routers.Update(networkingClient, routerID, opts).Extract()
					if err != nil {
						return r, "ACTIVE", err
					}
				}

				log.Printf("[DEBUG] openstack_networking_router_interface_v2 %s is still in use", routerInterfaceID)
				return r, "ACTIVE", nil
			}

			return r, "ACTIVE", err
		}

		log.Printf("[DEBUG] openstack_networking_router_interface_v2 %s is still active", routerInterfaceID)
		return r, "ACTIVE", nil
	}
}
