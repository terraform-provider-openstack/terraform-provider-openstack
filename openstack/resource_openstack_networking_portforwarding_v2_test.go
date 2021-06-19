package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/portforwarding"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
)

func TestAccNetworkingV2Portforwarding_basic(t *testing.T) {
	var (
		network networks.Network
		subnet  subnets.Subnet
		router  routers.Router
		port    ports.Port
		pf      portforwarding.PortForwarding
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckPortForwarding(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2PortForwardingDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortForwardingBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2RouterExists("openstack_networking_router_v2.router_1", &router),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2RouterInterfaceExists("openstack_networking_router_interface_v2.int_1"),
					testAccCheckNetworkingV2PortForwardingExists("openstack_networking_portforwarding_v2.pf_1", "openstack_networking_floatingip_v2.fip_1", &pf),
					resource.TestCheckResourceAttr("openstack_networking_portforwarding_v2.pf_1", "internal_port", "25"),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2PortForwardingDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkClient, err := config.NetworkingV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack portforwarding: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_networking_portforwarding_v2" {
			continue
		}
		fipID := rs.Primary.Attributes["floatingip_id"]
		primID := rs.Primary.ID
		_, err := portforwarding.Get(networkClient, fipID, primID).Extract()
		if err == nil {
			return fmt.Errorf("Port Forwarding still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2PortForwardingExists(n string, fipID string, kp *portforwarding.PortForwarding) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		fip, ok := s.RootModule().Resources[fipID]
		if !ok {
			return fmt.Errorf("Floating IP not found: %s", fipID)
		}

		config := testAccProvider.Meta().(*Config)
		networkClient, err := config.NetworkingV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %s", err)
		}

		found, err := portforwarding.Get(networkClient, fip.Primary.ID, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("openstack_networking_portforwarding_v2 not found")
		}

		*kp = *found

		return nil
	}
}

var testAccNetworkingV2PortForwardingBasic = fmt.Sprintf(`
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  description = "Network"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  gateway_ip = "192.168.199.1"
  enable_dhcp = "false"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_router_v2" "router_1" {
  name = "router_1"
  external_network_id = "%s"
  admin_state_up = "true"
}

resource "openstack_networking_port_v2" "port_1" {
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.3"
  }
}

resource "openstack_networking_router_interface_v2" "router_interface_1" {
  router_id = "${openstack_networking_router_v2.router_1.id}"
  port_id = "${openstack_networking_port_v2.port_1.id}"
}

resource "openstack_networking_floatingip_v2" "fip_1" {
  description = "test"
  port_id = ""
  pool = "%s"
  depends_on = [openstack_networking_router_interface_v2.router_interface_1]
}

resource "openstack_networking_portforwarding_v2" "pf_1" {
  protocol = "tcp"
  internal_ip_address = "${openstack_networking_port_v2.port_1.fixed_ip[0].ip_address}"
  internal_port = 25
  internal_port_id = "${openstack_networking_port_v2.port_1.id}"
  external_port = 2230
  floatingip_id = "${openstack_networking_floatingip_v2.fip_1.id}"
  depends_on = [openstack_networking_port_v2.port_1, openstack_networking_floatingip_v2.fip_1]
}
`, osExtGwID, osPoolName)
