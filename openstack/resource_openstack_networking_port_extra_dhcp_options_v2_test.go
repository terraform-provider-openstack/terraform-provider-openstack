package openstack

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/extradhcpopts"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNetworkingV2PortExtraDHCPOptions_basic(t *testing.T) {
	var network networks.Network
	var subnet subnets.Subnet
	var port ports.Port

	resourceName := "openstack_networking_port_extradhcp_options_v2.opts_1"
	networkName := acctest.RandomWithPrefix("tf-acc-network")
	subnetName := acctest.RandomWithPrefix("tf-acc-subnet")
	portName := acctest.RandomWithPrefix("tf-acc-port")
	optsName := acctest.RandomWithPrefix("tf-acc-extradhcp-opts")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortExtraDHCPOptionsDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2PortExtraDHCPOptions_basic(networkName, subnetName, portName, optsName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(resourceName, "extra_dhcp_opts.#", "2"),
				),
			},
		},
	})
}

func testAccNetworkingV2PortExtraDHCPOptions_basic(networkName, subnetName, portName, optsName string) string {
	return fmt.Sprintf(`
resource "openstack_networking_network_v2" "network_1" {
  name           = "%s"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name       = "%s"
  cidr       = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "port_1" {
  name           = "%s"
  admin_state_up = "true"
  network_id     = "${openstack_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id  =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}

resource "openstack_networking_port_extradhcp_options_v2" "opts_1" {
  name    = "%s"
  port_id = "${openstack_networking_port_v2.port_1.id}"

  extra_dhcp_opts {
    opt_name  = "optionA"
    opt_value = "valueA"
  }

  extra_dhcp_opts {
    opt_name  = "optionB"
    opt_value = "valueB"
  }
}
`, networkName, subnetName, portName, optsName)
}

func testAccCheckNetworkingV2PortExtraDHCPOptionsDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_networking_port_v2" {
			continue
		}

		var port struct {
			ports.Port
			extradhcpopts.ExtraDHCPOptsExt
		}
		err = ports.Get(networkingClient, rs.Primary.ID).ExtractInto(&port)
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return nil
			}
			return fmt.Errorf("Error getting OpenStack Neutron port: %s", err)
		}

		extraDHCPOpts := port.ExtraDHCPOpts
		if len(extraDHCPOpts) != 0 {
			return fmt.Errorf("Port %s still has DHCP options: %+v", port.ID, extraDHCPOpts)
		}
	}

	return nil
}
