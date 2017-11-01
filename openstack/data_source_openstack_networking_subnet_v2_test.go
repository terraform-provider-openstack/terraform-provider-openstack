package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNetworkingV2SubnetDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOpenStackNetworkingSubnetV2DataSource_subnet,
			},
			resource.TestStep{
				Config: testAccOpenStackNetworkingSubnetV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetV2DataSourceID("data.openstack_networking_subnet_v2.subnet_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnet_v2.subnet_1", "name", "subnet_1"),
				),
			},
		},
	})
}

func TestAccNetworkingV2SubnetDataSource_testQueries(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOpenStackNetworkingSubnetV2DataSource_subnet,
			},
			resource.TestStep{
				Config: testAccOpenStackNetworkingSubnetV2DataSource_cidr,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetV2DataSourceID("data.openstack_networking_subnet_v2.subnet_1"),
				),
			},
			resource.TestStep{
				Config: testAccOpenStackNetworkingSubnetV2DataSource_dhcpEnabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetV2DataSourceID("data.openstack_networking_subnet_v2.subnet_1"),
				),
			},
			resource.TestStep{
				Config: testAccOpenStackNetworkingSubnetV2DataSource_ipVersion,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetV2DataSourceID("data.openstack_networking_subnet_v2.subnet_1"),
				),
			},
			resource.TestStep{
				Config: testAccOpenStackNetworkingSubnetV2DataSource_gatewayIP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetV2DataSourceID("data.openstack_networking_subnet_v2.subnet_1"),
				),
			},
		},
	})
}

func testAccCheckNetworkingSubnetV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find subnet data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Subnet data source ID not set")
		}

		return nil
	}
}

const testAccOpenStackNetworkingSubnetV2DataSource_subnet = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  network_id = "${openstack_networking_network_v2.network_1.id}"
}
`

var testAccOpenStackNetworkingSubnetV2DataSource_basic = fmt.Sprintf(`
%s

data "openstack_networking_subnet_v2" "subnet_1" {
	name = "${openstack_networking_subnet_v2.subnet_1.name}"
}
`, testAccOpenStackNetworkingSubnetV2DataSource_subnet)

var testAccOpenStackNetworkingSubnetV2DataSource_cidr = fmt.Sprintf(`
%s

data "openstack_networking_subnet_v2" "subnet_1" {
	cidr = "192.168.199.0/24"
}
`, testAccOpenStackNetworkingSubnetV2DataSource_subnet)

var testAccOpenStackNetworkingSubnetV2DataSource_dhcpEnabled = fmt.Sprintf(`
%s

data "openstack_networking_subnet_v2" "subnet_1" {
  network_id = "${openstack_networking_network_v2.network_1.id}"
	dhcp_enabled = true
}
`, testAccOpenStackNetworkingSubnetV2DataSource_subnet)

var testAccOpenStackNetworkingSubnetV2DataSource_ipVersion = fmt.Sprintf(`
%s

data "openstack_networking_subnet_v2" "subnet_1" {
  network_id = "${openstack_networking_network_v2.network_1.id}"
  ip_version = 4
}
`, testAccOpenStackNetworkingSubnetV2DataSource_subnet)

var testAccOpenStackNetworkingSubnetV2DataSource_gatewayIP = fmt.Sprintf(`
%s

data "openstack_networking_subnet_v2" "subnet_1" {
  gateway_ip = "${openstack_networking_subnet_v2.subnet_1.gateway_ip}"
}
`, testAccOpenStackNetworkingSubnetV2DataSource_subnet)
