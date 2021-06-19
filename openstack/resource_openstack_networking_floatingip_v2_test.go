package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
)

func TestAccNetworkingV2FloatingIP_basic(t *testing.T) {
	var fip floatingips.FloatingIP

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2FloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2FloatingIPBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2FloatingIPExists("openstack_networking_floatingip_v2.fip_1", &fip),
					resource.TestCheckResourceAttr("openstack_networking_floatingip_v2.fip_1", "description", "test floating IP"),
				),
			},
		},
	})
}

func TestAccNetworkingV2FloatingIP_fixedip_bind(t *testing.T) {
	var fip floatingips.FloatingIP

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2FloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2FloatingIPFixedIPBind1(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2FloatingIPExists("openstack_networking_floatingip_v2.fip_1", &fip),
					testAccCheckNetworkingV2FloatingIPBoundToCorrectIP(&fip, "192.168.199.20"),
					resource.TestCheckResourceAttr("openstack_networking_floatingip_v2.fip_1", "description", "test"),
					resource.TestCheckResourceAttr("openstack_networking_floatingip_v2.fip_1", "fixed_ip", "192.168.199.20"),
				),
			},
			{
				Config: testAccNetworkingV2FloatingIPFixedipBind2(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2FloatingIPExists("openstack_networking_floatingip_v2.fip_1", &fip),
					testAccCheckNetworkingV2FloatingIPBoundToCorrectIP(&fip, "192.168.199.10"),
					resource.TestCheckResourceAttr("openstack_networking_floatingip_v2.fip_1", "description", ""),
					resource.TestCheckResourceAttr("openstack_networking_floatingip_v2.fip_1", "fixed_ip", "192.168.199.10"),
				),
			},
		},
	})
}

func TestAccNetworkingV2FloatingIP_subnetIDs(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2FloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2FloatingIPSubnetIDs(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openstack_networking_floatingip_v2.fip_1", "description", "test"),
				),
			},
		},
	})
}

func TestAccNetworkingV2FloatingIP_timeout(t *testing.T) {
	var fip floatingips.FloatingIP

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2FloatingIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2FloatingIPTimeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2FloatingIPExists("openstack_networking_floatingip_v2.fip_1", &fip),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2FloatingIPDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkClient, err := config.NetworkingV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack floating IP: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_networking_floatingip_v2" {
			continue
		}

		_, err := floatingips.Get(networkClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Floating IP still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2FloatingIPExists(n string, kp *floatingips.FloatingIP) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkClient, err := config.NetworkingV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %s", err)
		}

		found, err := floatingips.Get(networkClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Floating IP not found")
		}

		*kp = *found

		return nil
	}
}

func testAccCheckNetworkingV2FloatingIPBoundToCorrectIP(fip *floatingips.FloatingIP, fixedIP string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if fip.FixedIP != fixedIP {
			return fmt.Errorf("Floating IP associated with wrong fixed ip")
		}

		return nil
	}
}

const testAccNetworkingV2FloatingIPBasic = `
resource "openstack_networking_floatingip_v2" "fip_1" {
  description = "test floating IP"
}
`

func testAccNetworkingV2FloatingIPFixedIPBind1() string {
	return fmt.Sprintf(`
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_router_interface_v2" "router_interface_1" {
  router_id = "${openstack_networking_router_v2.router_1.id}"
  subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
}

resource "openstack_networking_router_v2" "router_1" {
  name = "router_1"
  external_gateway = "%s"
}

resource "openstack_networking_port_v2" "port_1" {
  admin_state_up = "true"
  network_id = "${openstack_networking_subnet_v2.subnet_1.network_id}"

  fixed_ip {
    subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.10"
  }

  fixed_ip {
    subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.20"
  }
}

resource "openstack_networking_floatingip_v2" "fip_1" {
  pool = "%s"
  description = "test"
  port_id = "${openstack_networking_port_v2.port_1.id}"
  fixed_ip = "${openstack_networking_port_v2.port_1.fixed_ip.1.ip_address}"
}
`, osExtGwID, osPoolName)
}

func testAccNetworkingV2FloatingIPFixedipBind2() string {
	return fmt.Sprintf(`
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_router_interface_v2" "router_interface_1" {
  router_id = "${openstack_networking_router_v2.router_1.id}"
  subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
}

resource "openstack_networking_router_v2" "router_1" {
  name = "router_1"
  external_gateway = "%s"
}

resource "openstack_networking_port_v2" "port_1" {
  admin_state_up = "true"
  network_id = "${openstack_networking_subnet_v2.subnet_1.network_id}"

  fixed_ip {
    subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.10"
  }

  fixed_ip {
    subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.20"
  }
}

resource "openstack_networking_floatingip_v2" "fip_1" {
  pool = "%s"
  port_id = "${openstack_networking_port_v2.port_1.id}"
  fixed_ip = "${openstack_networking_port_v2.port_1.fixed_ip.0.ip_address}"
}
`, osExtGwID, osPoolName)
}

const testAccNetworkingV2FloatingIPTimeout = `
resource "openstack_networking_floatingip_v2" "fip_1" {
  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`

func testAccNetworkingV2FloatingIPSubnetIDs() string {
	return fmt.Sprintf(`
data "openstack_networking_network_v2" "ext_network" {
  name = "%s"
}

resource "openstack_networking_floatingip_v2" "fip_1" {
  pool = data.openstack_networking_network_v2.ext_network.name
  description = "test"
  subnet_ids = flatten([
    data.openstack_networking_network_v2.ext_network.id, # wrong UUID
    data.openstack_networking_network_v2.ext_network.subnets,
    data.openstack_networking_network_v2.ext_network.id, # wrong UUID again
  ])
}
`, osPoolName)
}
