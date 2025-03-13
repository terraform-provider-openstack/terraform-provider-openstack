package openstack

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/subnets"
)

func TestAccNetworkingV2RouterInterface_basic_subnet(t *testing.T) {
	var network networks.Network
	var router routers.Router
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2RouterInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterInterfaceBasicSubnet,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2RouterExists("openstack_networking_router_v2.router_1", &router),
					testAccCheckNetworkingV2RouterInterfaceExists("openstack_networking_router_interface_v2.int_1"),
				),
			},
		},
	})
}

func TestAccNetworkingV2RouterInterface_force_destroy(t *testing.T) {
	var network networks.Network
	var router routers.Router
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2RouterInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterInterfaceForceDestroy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2RouterExists("openstack_networking_router_v2.router_1", &router),
					testAccCheckNetworkingV2RouterInterfaceExists("openstack_networking_router_interface_v2.int_1"),
					resource.TestCheckResourceAttr("openstack_networking_router_interface_v2.int_1", "force_destroy", "true"),
					testAccCheckNetworkingV2RouterInterfaceAddExtraRoutes("openstack_networking_router_interface_v2.int_1"),
				),
			},
		},
	})
}

// TestAccNetworkingV2RouterInterface_v6_subnet tests that multiple router interfaces for IPv6 subnets
// which are attached to the same port are handled properly.
func TestAccNetworkingV2RouterInterface_v6_subnet(t *testing.T) {
	var network networks.Network
	var router routers.Router
	var subnet1 subnets.Subnet
	var subnet2 subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2RouterInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterInterfaceV6Subnet + testAccNetworkingV2RouterInterfaceV6SubnetSecondInterface,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet1),
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_2", &subnet2),
					testAccCheckNetworkingV2RouterExists("openstack_networking_router_v2.router_1", &router),
					testAccCheckNetworkingV2RouterInterfaceExists("openstack_networking_router_interface_v2.int_1"),
					testAccCheckNetworkingV2RouterInterfaceExists("openstack_networking_router_interface_v2.int_2"),
				),
			},
			{ // Make sure deleting one of the router interfaces does not remove the other one.
				Config: testAccNetworkingV2RouterInterfaceV6Subnet,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2RouterInterfaceExists("openstack_networking_router_interface_v2.int_1"),
				),
			},
		},
	})
}

func TestAccNetworkingV2RouterInterface_basic_port(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var router routers.Router
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2RouterInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterInterfaceBasicPort,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2RouterExists("openstack_networking_router_v2.router_1", &router),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2RouterInterfaceExists("openstack_networking_router_interface_v2.int_1"),
				),
			},
		},
	})
}

func TestAccNetworkingV2RouterInterface_timeout(t *testing.T) {
	var network networks.Network
	var router routers.Router
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2RouterInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterInterfaceTimeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2RouterExists("openstack_networking_router_v2.router_1", &router),
					testAccCheckNetworkingV2RouterInterfaceExists("openstack_networking_router_interface_v2.int_1"),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2RouterInterfaceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.NetworkingV2Client(context.TODO(), osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_networking_router_interface_v2" {
			continue
		}

		_, err := ports.Get(context.TODO(), networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Router interface still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2RouterInterfaceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.NetworkingV2Client(context.TODO(), osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %s", err)
		}

		found, err := ports.Get(context.TODO(), networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Router interface not found")
		}

		return nil
	}
}

func testAccCheckNetworkingV2RouterInterfaceAddExtraRoutes(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.NetworkingV2Client(context.TODO(), osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %s", err)
		}

		found, err := ports.Get(context.TODO(), networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Router interface not found")
		}

		// add extra routes to the router
		routes := []routers.Route{
			{
				DestinationCIDR: "10.0.0.0/24",
				NextHop:         "192.168.199.200",
			},
		}
		opts := routers.UpdateOpts{
			Routes: &routes,
		}
		_, err = routers.Update(context.TODO(), networkingClient, found.DeviceID, opts).Extract()
		if err != nil {
			return fmt.Errorf("Router cannot be updated: %v", err)
		}

		return nil
	}
}

const testAccNetworkingV2RouterInterfaceBasicSubnet = `
resource "openstack_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
}

resource "openstack_networking_router_interface_v2" "int_1" {
  subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
  router_id = "${openstack_networking_router_v2.router_1.id}"
}

resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}
`

const testAccNetworkingV2RouterInterfaceForceDestroy = `
resource "openstack_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
}

resource "openstack_networking_router_interface_v2" "int_1" {
  subnet_id     = "${openstack_networking_subnet_v2.subnet_1.id}"
  router_id     = "${openstack_networking_router_v2.router_1.id}"
  force_destroy = true
}

resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}
`

const testAccNetworkingV2RouterInterfaceV6Subnet = `
resource "openstack_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
}

resource "openstack_networking_router_interface_v2" "int_1" {
  subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
  router_id = "${openstack_networking_router_v2.router_1.id}"
}

resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  cidr = "fd00:0:0:1::/64"
  ip_version = 6
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_subnet_v2" "subnet_2" {
  cidr = "fd00:0:0:2::/64"
  ip_version = 6
  network_id = "${openstack_networking_network_v2.network_1.id}"
}
`

const testAccNetworkingV2RouterInterfaceV6SubnetSecondInterface = `
resource "openstack_networking_router_interface_v2" "int_2" {
  subnet_id = "${openstack_networking_subnet_v2.subnet_2.id}"
  router_id = "${openstack_networking_router_v2.router_1.id}"
}
`

const testAccNetworkingV2RouterInterfaceBasicPort = `
resource "openstack_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
}

resource "openstack_networking_router_interface_v2" "int_1" {
  router_id = "${openstack_networking_router_v2.router_1.id}"
  port_id = "${openstack_networking_port_v2.port_1.id}"
}

resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.1"
  }
}
`

const testAccNetworkingV2RouterInterfaceTimeout = `
resource "openstack_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
}

resource "openstack_networking_router_interface_v2" "int_1" {
  subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
  router_id = "${openstack_networking_router_v2.router_1.id}"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}

resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}
`
