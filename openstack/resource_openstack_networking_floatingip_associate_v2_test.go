package openstack

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/floatingips"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccNetworkingV2FloatingIPAssociate_basic(t *testing.T) {
	var fip floatingips.FloatingIP

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2FloatingIPAssociateDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2FloatingIPAssociateBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2FloatingIPExists(t.Context(),
						"openstack_networking_floatingip_associate_v2.fip_1", &fip),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_floatingip_associate_v2.fip_1", "floating_ip", &fip.FloatingIP),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_floatingip_associate_v2.fip_1", "port_id", &fip.PortID),
				),
			},
		},
	})
}

func TestAccNetworkingV2FloatingIPAssociate_twoFixedIPs(t *testing.T) {
	var fip floatingips.FloatingIP

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2FloatingIPAssociateDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2FloatingIPAssociateTwoFixedIPs1(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2FloatingIPExists(t.Context(),
						"openstack_networking_floatingip_associate_v2.fip_1", &fip),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_floatingip_associate_v2.fip_1", "floating_ip", &fip.FloatingIP),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_floatingip_associate_v2.fip_1", "port_id", &fip.PortID),
					testAccCheckNetworkingV2FloatingIPBoundToCorrectIP(&fip, "192.168.199.20"),
					resource.TestCheckResourceAttr("openstack_networking_floatingip_associate_v2.fip_1", "fixed_ip", "192.168.199.20"),
				),
			},
			{
				Config: testAccNetworkingV2FloatingIPAssociateTwoFixedIPs2(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2FloatingIPExists(t.Context(),
						"openstack_networking_floatingip_associate_v2.fip_1", &fip),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_floatingip_associate_v2.fip_1", "floating_ip", &fip.FloatingIP),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_floatingip_associate_v2.fip_1", "port_id", &fip.PortID),
					testAccCheckNetworkingV2FloatingIPBoundToCorrectIP(&fip, "192.168.199.21"),
					resource.TestCheckResourceAttr("openstack_networking_floatingip_associate_v2.fip_1", "fixed_ip", "192.168.199.21"),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2FloatingIPAssociateDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		networkClient, err := config.NetworkingV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack network client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_networking_floatingip_v2" {
				continue
			}

			fip, err := floatingips.Get(ctx, networkClient, rs.Primary.ID).Extract()
			if err != nil {
				if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
					return nil
				}

				return fmt.Errorf("Error retrieving Floating IP: %w", err)
			}

			if fip.PortID != "" {
				return errors.New("Floating IP is still associated")
			}
		}

		return nil
	}
}

func testAccNetworkingV2FloatingIPAssociateBasic() string {
	return fmt.Sprintf(`
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = openstack_networking_network_v2.network_1.id
}

resource "openstack_networking_router_interface_v2" "router_interface_1" {
  router_id = openstack_networking_router_v2.router_1.id
  subnet_id = openstack_networking_subnet_v2.subnet_1.id
}

resource "openstack_networking_router_v2" "router_1" {
  name = "router_1"
  external_network_id = "%s"
}

resource "openstack_networking_port_v2" "port_1" {
  admin_state_up = "true"
  network_id = openstack_networking_subnet_v2.subnet_1.network_id

  fixed_ip {
    subnet_id = openstack_networking_router_interface_v2.router_interface_1.subnet_id
    ip_address = "192.168.199.20"
  }
}

resource "openstack_networking_floatingip_v2" "fip_1" {
  pool = "%s"
}

resource "openstack_networking_floatingip_associate_v2" "fip_1" {
  floating_ip = openstack_networking_floatingip_v2.fip_1.address
  port_id = openstack_networking_port_v2.port_1.id

}
`, osExtGwID, osPoolName)
}

func testAccNetworkingV2FloatingIPAssociateTwoFixedIPs1() string {
	return fmt.Sprintf(`
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = openstack_networking_network_v2.network_1.id
}

resource "openstack_networking_router_interface_v2" "router_interface_1" {
  router_id = openstack_networking_router_v2.router_1.id
  subnet_id = openstack_networking_subnet_v2.subnet_1.id
}

resource "openstack_networking_router_v2" "router_1" {
  name = "router_1"
  external_network_id = "%s"
}

resource "openstack_networking_port_v2" "port_1" {
  admin_state_up = "true"
  network_id = openstack_networking_subnet_v2.subnet_1.network_id

  fixed_ip {
    subnet_id = openstack_networking_router_interface_v2.router_interface_1.subnet_id
    ip_address = "192.168.199.20"
  }

  fixed_ip {
    subnet_id = openstack_networking_router_interface_v2.router_interface_1.subnet_id
    ip_address = "192.168.199.21"
  }
}

resource "openstack_networking_floatingip_v2" "fip_1" {
  pool = "%s"
}

resource "openstack_networking_floatingip_associate_v2" "fip_1" {
  floating_ip = openstack_networking_floatingip_v2.fip_1.address
  port_id = openstack_networking_port_v2.port_1.id
  fixed_ip = "192.168.199.20"
}
`, osExtGwID, osPoolName)
}

func testAccNetworkingV2FloatingIPAssociateTwoFixedIPs2() string {
	return fmt.Sprintf(`
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = openstack_networking_network_v2.network_1.id
}

resource "openstack_networking_router_interface_v2" "router_interface_1" {
  router_id = openstack_networking_router_v2.router_1.id
  subnet_id = openstack_networking_subnet_v2.subnet_1.id
}

resource "openstack_networking_router_v2" "router_1" {
  name = "router_1"
  external_network_id = "%s"
}

resource "openstack_networking_port_v2" "port_1" {
  admin_state_up = "true"
  network_id = openstack_networking_subnet_v2.subnet_1.network_id

  fixed_ip {
    subnet_id = openstack_networking_router_interface_v2.router_interface_1.subnet_id
    ip_address = "192.168.199.20"
  }

  fixed_ip {
    subnet_id = openstack_networking_router_interface_v2.router_interface_1.subnet_id
    ip_address = "192.168.199.21"
  }
}

resource "openstack_networking_floatingip_v2" "fip_1" {
  pool = "%s"
}

resource "openstack_networking_floatingip_associate_v2" "fip_1" {
  floating_ip = openstack_networking_floatingip_v2.fip_1.address
  port_id = openstack_networking_port_v2.port_1.id
  fixed_ip = "192.168.199.21"
}
`, osExtGwID, osPoolName)
}
