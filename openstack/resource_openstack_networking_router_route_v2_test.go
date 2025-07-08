package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/extraroutes"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/subnets"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccNetworkingV2RouterRoute_basic(t *testing.T) {
	var router routers.Router

	var network [2]networks.Network

	var subnet [2]subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterRouteCreate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2RouterExists(t.Context(), "openstack_networking_router_v2.router_1", &router),
					testAccCheckNetworkingV2NetworkExists(t.Context(), "openstack_networking_network_v2.network_1", &network[0]),
					testAccCheckNetworkingV2SubnetExists(t.Context(), "openstack_networking_subnet_v2.subnet_1", &subnet[0]),
					testAccCheckNetworkingV2NetworkExists(t.Context(), "openstack_networking_network_v2.network_1", &network[1]),
					testAccCheckNetworkingV2SubnetExists(t.Context(), "openstack_networking_subnet_v2.subnet_1", &subnet[1]),
					testAccCheckNetworkingV2RouterInterfaceExists(t.Context(), "openstack_networking_router_interface_v2.int_1"),
					testAccCheckNetworkingV2RouterInterfaceExists(t.Context(), "openstack_networking_router_interface_v2.int_2"),
					testAccCheckNetworkingV2RouterRouteExists(t.Context(), "openstack_networking_router_route_v2.router_route_1"),
				),
			},
			{
				Config: testAccNetworkingV2RouterRouteUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2RouterRouteExists(t.Context(), "openstack_networking_router_route_v2.router_route_1"),
					testAccCheckNetworkingV2RouterRouteExists(t.Context(), "openstack_networking_router_route_v2.router_route_2"),
				),
			},
			{
				Config: testAccNetworkingV2RouterRouteDestroy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2RouterRouteEmpty(t.Context(), "openstack_networking_router_v2.router_1"),
				),
			},
		},
	})
}

func TestAccNetworkingV2RouterRoute_deleted(t *testing.T) {
	var router routers.Router

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterRouteCreate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2RouterExists(t.Context(), "openstack_networking_router_v2.router_1", &router),
				),
			},
			{
				PreConfig:          testAccNetworkingV2RouterRouteManualDelete(t, &router),
				Config:             testAccNetworkingV2RouterRouteCreate,
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccNetworkingV2RouterRouteManualDelete(t *testing.T, router *routers.Router) func() {
	return func() {
		config := testAccProvider.Meta().(*Config)
		ctx := t.Context()

		networkingClient, err := config.NetworkingV2Client(ctx, osRegionName)
		if err != nil {
			t.Errorf("Error creating OpenStack networking client: %v", err)
		}

		opts := extraroutes.Opts{
			Routes: &router.Routes,
		}
		_, err = extraroutes.Remove(ctx, networkingClient, router.ID, opts).Extract()
		if err != nil {
			t.Errorf("Error removing route from router %s: %v", router.ID, err)
		}
	}
}

func testAccCheckNetworkingV2RouterRouteEmpty(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		networkingClient, err := config.NetworkingV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %w", err)
		}

		router, err := routers.Get(ctx, networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if router.ID != rs.Primary.ID {
			return errors.New("Router not found")
		}

		if len(router.Routes) != 0 {
			return fmt.Errorf("Invalid number of route entries: %d", len(router.Routes))
		}

		return nil
	}
}

func testAccCheckNetworkingV2RouterRouteExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		networkingClient, err := config.NetworkingV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %w", err)
		}

		router, err := routers.Get(ctx, networkingClient, rs.Primary.Attributes["router_id"]).Extract()
		if err != nil {
			return err
		}

		if router.ID != rs.Primary.Attributes["router_id"] {
			return errors.New("Router for route not found")
		}

		found := false

		for _, r := range router.Routes {
			if r.DestinationCIDR == rs.Primary.Attributes["destination_cidr"] && r.NextHop == rs.Primary.Attributes["next_hop"] {
				found = true
			}
		}

		if !found {
			return fmt.Errorf("Could not find route for destination CIDR: %s, next hop: %s", rs.Primary.Attributes["destination_cidr"], rs.Primary.Attributes["next_hop"])
		}

		return nil
	}
}

func testAccCheckNetworkingV2RouterRouteDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		networkingClient, err := config.NetworkingV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_networking_router_route_v2" {
				continue
			}

			routeExists := false

			router, err := routers.Get(ctx, networkingClient, rs.Primary.Attributes["router_id"]).Extract()
			if err == nil {
				rts := router.Routes
				for _, r := range rts {
					if r.DestinationCIDR == rs.Primary.Attributes["destination_cidr"] && r.NextHop == rs.Primary.Attributes["next_hop"] {
						routeExists = true

						break
					}
				}
			}

			if routeExists {
				return errors.New("Route still exists")
			}
		}

		return nil
	}
}

const testAccNetworkingV2RouterRouteCreate = `
resource "openstack_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
}

resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = openstack_networking_network_v2.network_1.id
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = openstack_networking_network_v2.network_1.id

  fixed_ip {
    subnet_id = openstack_networking_subnet_v2.subnet_1.id
    ip_address = "192.168.199.1"
  }
}

resource "openstack_networking_router_interface_v2" "int_1" {
  router_id = openstack_networking_router_v2.router_1.id
  port_id = openstack_networking_port_v2.port_1.id
}

resource "openstack_networking_network_v2" "network_2" {
  name = "network_2"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_2" {
  cidr = "192.168.200.0/24"
  ip_version = 4
  network_id = openstack_networking_network_v2.network_2.id
}

resource "openstack_networking_port_v2" "port_2" {
  name = "port_2"
  admin_state_up = "true"
  network_id = openstack_networking_network_v2.network_2.id

  fixed_ip {
    subnet_id = openstack_networking_subnet_v2.subnet_2.id
    ip_address = "192.168.200.1"
  }
}

resource "openstack_networking_router_interface_v2" "int_2" {
  router_id = openstack_networking_router_v2.router_1.id
  port_id = openstack_networking_port_v2.port_2.id
}

resource "openstack_networking_router_route_v2" "router_route_1" {
  destination_cidr = "10.0.1.0/24"
  next_hop = "192.168.199.254"

  depends_on = ["openstack_networking_router_interface_v2.int_1"]
  router_id = openstack_networking_router_v2.router_1.id
}
`

const testAccNetworkingV2RouterRouteUpdate = `
resource "openstack_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
}

resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = openstack_networking_network_v2.network_1.id
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = openstack_networking_network_v2.network_1.id

  fixed_ip {
    subnet_id = openstack_networking_subnet_v2.subnet_1.id
    ip_address = "192.168.199.1"
  }
}

resource "openstack_networking_router_interface_v2" "int_1" {
  router_id = openstack_networking_router_v2.router_1.id
  port_id = openstack_networking_port_v2.port_1.id
}

resource "openstack_networking_network_v2" "network_2" {
  name = "network_2"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_2" {
  cidr = "192.168.200.0/24"
  ip_version = 4
  network_id = openstack_networking_network_v2.network_2.id
}

resource "openstack_networking_port_v2" "port_2" {
  name = "port_2"
  admin_state_up = "true"
  network_id = openstack_networking_network_v2.network_2.id

  fixed_ip {
    subnet_id = openstack_networking_subnet_v2.subnet_2.id
    ip_address = "192.168.200.1"
  }
}

resource "openstack_networking_router_interface_v2" "int_2" {
  router_id = openstack_networking_router_v2.router_1.id
  port_id = openstack_networking_port_v2.port_2.id
}

resource "openstack_networking_router_route_v2" "router_route_1" {
  destination_cidr = "10.0.1.0/24"
  next_hop = "192.168.199.254"

  depends_on = ["openstack_networking_router_interface_v2.int_1"]
  router_id = openstack_networking_router_v2.router_1.id
}

resource "openstack_networking_router_route_v2" "router_route_2" {
  destination_cidr = "10.0.2.0/24"
  next_hop = "192.168.200.254"

  depends_on = ["openstack_networking_router_interface_v2.int_2"]
  router_id = openstack_networking_router_v2.router_1.id
}
`

const testAccNetworkingV2RouterRouteDestroy = `
resource "openstack_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
}

resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = openstack_networking_network_v2.network_1.id
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = openstack_networking_network_v2.network_1.id

  fixed_ip {
    subnet_id = openstack_networking_subnet_v2.subnet_1.id
    ip_address = "192.168.199.1"
  }
}

resource "openstack_networking_router_interface_v2" "int_1" {
  router_id = openstack_networking_router_v2.router_1.id
  port_id = openstack_networking_port_v2.port_1.id
}

resource "openstack_networking_network_v2" "network_2" {
  name = "network_2"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_2" {
  ip_version = 4
  cidr = "192.168.200.0/24"
  network_id = openstack_networking_network_v2.network_2.id
}

resource "openstack_networking_port_v2" "port_2" {
  name = "port_2"
  admin_state_up = "true"
  network_id = openstack_networking_network_v2.network_2.id

  fixed_ip {
    subnet_id = openstack_networking_subnet_v2.subnet_2.id
    ip_address = "192.168.200.1"
  }
}

resource "openstack_networking_router_interface_v2" "int_2" {
  router_id = openstack_networking_router_v2.router_1.id
  port_id = openstack_networking_port_v2.port_2.id
}
`
