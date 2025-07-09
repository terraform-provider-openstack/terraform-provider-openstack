package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNetworkingV2RouterRoutes_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterRoutesCreate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_networking_router_routes_v2.router_routes_1", "routes.#", "2"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_router_v2.router_1", "routes.#", "2"),
				),
			},
			{
				Config: testAccNetworkingV2RouterRoutesUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_networking_router_routes_v2.router_routes_1", "routes.#", "1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_router_v2.router_1", "routes.#", "1"),
				),
			},
			{
				Config: testAccNetworkingV2RouterRoutesEmpty,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2RouterRouteEmpty(t.Context(), "openstack_networking_router_v2.router_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_router_v2.router_1", "routes.#", "0"),
				),
			},
			{
				Config: testAccNetworkingV2RouterRoutesInit,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2RouterRouteEmpty(t.Context(), "openstack_networking_router_v2.router_1"),
				),
			},
		},
	})
}

const testAccNetworkingV2RouterRoutesInit = `
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
`

const testAccNetworkingV2RouterRoutesCreate = testAccNetworkingV2RouterRoutesInit + `
resource "openstack_networking_router_routes_v2" "router_routes_1" {
  router_id = openstack_networking_router_interface_v2.int_1.router_id
  depends_on = ["openstack_networking_router_interface_v2.int_2"]

  routes {
    destination_cidr = "10.0.1.0/24"
    next_hop = "192.168.199.254"
  }

  routes {
    destination_cidr = "10.0.2.0/24"
    next_hop = "192.168.200.254"
  }
}

data "openstack_networking_router_v2" "router_1" {
  router_id = openstack_networking_router_routes_v2.router_routes_1.router_id
}
`

const testAccNetworkingV2RouterRoutesUpdate = testAccNetworkingV2RouterRoutesInit + `
resource "openstack_networking_router_routes_v2" "router_routes_1" {
  router_id = openstack_networking_router_interface_v2.int_1.router_id
  depends_on = ["openstack_networking_router_interface_v2.int_2"]

  routes {
    destination_cidr = "10.0.1.0/24"
    next_hop = "192.168.199.254"
  }
}

data "openstack_networking_router_v2" "router_1" {
  router_id = openstack_networking_router_routes_v2.router_routes_1.router_id
}
`

const testAccNetworkingV2RouterRoutesEmpty = testAccNetworkingV2RouterRoutesInit + `
resource "openstack_networking_router_routes_v2" "router_routes_1" {
  router_id = openstack_networking_router_interface_v2.int_1.router_id
  depends_on = ["openstack_networking_router_interface_v2.int_2"]
}

data "openstack_networking_router_v2" "router_1" {
  router_id = openstack_networking_router_routes_v2.router_routes_1.router_id
}
`
