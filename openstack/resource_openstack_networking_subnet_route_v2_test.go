package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/subnets"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccNetworkingV2SubnetRoute_basic(t *testing.T) {
	var (
		router  routers.Router
		network networks.Network
		subnet  subnets.Subnet
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SubnetRouteCreate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2RouterExists(t.Context(), "openstack_networking_router_v2.router_1", &router),
					testAccCheckNetworkingV2NetworkExists(t.Context(), "openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2SubnetExists(t.Context(), "openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2RouterInterfaceExists(t.Context(), "openstack_networking_router_interface_v2.int_1"),
					testAccCheckNetworkingV2SubnetRouteExists(t.Context(), "openstack_networking_subnet_route_v2.subnet_route_1"),
				),
			},
			{
				Config: testAccNetworkingV2SubnetRouteUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetRouteExists(t.Context(), "openstack_networking_subnet_route_v2.subnet_route_1"),
					testAccCheckNetworkingV2SubnetRouteExists(t.Context(), "openstack_networking_subnet_route_v2.subnet_route_2"),
				),
			},
			{
				Config: testAccNetworkingV2SubnetRouteDestroy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetRouteEmpty(t.Context(), "openstack_networking_subnet_v2.subnet_1"),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2SubnetRouteEmpty(ctx context.Context, n string) resource.TestCheckFunc {
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

		subnet, err := subnets.Get(ctx, networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if subnet.ID != rs.Primary.ID {
			return errors.New("Subnet not found")
		}

		if len(subnet.HostRoutes) != 0 {
			return fmt.Errorf("Invalid number of route entries: %d", len(subnet.HostRoutes))
		}

		return nil
	}
}

func testAccCheckNetworkingV2SubnetRouteExists(ctx context.Context, n string) resource.TestCheckFunc {
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

		subnet, err := subnets.Get(ctx, networkingClient, rs.Primary.Attributes["subnet_id"]).Extract()
		if err != nil {
			return err
		}

		if subnet.ID != rs.Primary.Attributes["subnet_id"] {
			return errors.New("Subnet for route not found")
		}

		found := false

		for _, r := range subnet.HostRoutes {
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

func testAccCheckNetworkingV2SubnetRouteDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		networkingClient, err := config.NetworkingV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_networking_subnet_route_v2" {
				continue
			}

			routeExists := false

			subnet, err := subnets.Get(ctx, networkingClient, rs.Primary.Attributes["subnet_id"]).Extract()
			if err == nil {
				rts := subnet.HostRoutes
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

const testAccNetworkingV2SubnetRouteCreate = `
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

resource "openstack_networking_subnet_route_v2" "subnet_route_1" {
  destination_cidr = "10.0.1.0/24"
  next_hop = "192.168.199.254"

  subnet_id = openstack_networking_subnet_v2.subnet_1.id
}
`

const testAccNetworkingV2SubnetRouteUpdate = `
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

resource "openstack_networking_subnet_route_v2" "subnet_route_1" {
  destination_cidr = "10.0.1.0/24"
  next_hop = "192.168.199.254"

  subnet_id = openstack_networking_subnet_v2.subnet_1.id
}

resource "openstack_networking_subnet_route_v2" "subnet_route_2" {
  destination_cidr = "10.0.2.0/24"
  next_hop = "192.168.199.254"

  subnet_id = openstack_networking_subnet_v2.subnet_1.id
}
`

const testAccNetworkingV2SubnetRouteDestroy = `
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
`
