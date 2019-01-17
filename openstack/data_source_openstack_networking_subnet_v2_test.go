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
			{
				Config: testAccOpenStackNetworkingSubnetV2DataSource_subnet,
			},
			{
				Config: testAccOpenStackNetworkingSubnetV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetV2DataSourceID("data.openstack_networking_subnet_v2.subnet_1"),
					testAccCheckNetworkingSubnetV2DataSourceGoodNetwork("data.openstack_networking_subnet_v2.subnet_1", "openstack_networking_network_v2.network_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnet_v2.subnet_1", "name", "subnet_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnet_v2.subnet_1", "all_tags.#", "2"),
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
			{
				Config: testAccOpenStackNetworkingSubnetV2DataSource_subnet,
			},
			{
				Config: testAccOpenStackNetworkingSubnetV2DataSource_cidr,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetV2DataSourceID("data.openstack_networking_subnet_v2.subnet_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnet_v2.subnet_1", "description", "my subnet description"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnet_v2.subnet_1", "all_tags.#", "2"),
				),
			},
			{
				Config: testAccOpenStackNetworkingSubnetV2DataSource_dhcpEnabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetV2DataSourceID("data.openstack_networking_subnet_v2.subnet_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnet_v2.subnet_1", "tags.#", "1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnet_v2.subnet_1", "all_tags.#", "2"),
				),
			},
			{
				Config: testAccOpenStackNetworkingSubnetV2DataSource_ipVersion,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetV2DataSourceID("data.openstack_networking_subnet_v2.subnet_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnet_v2.subnet_1", "all_tags.#", "2"),
				),
			},
			{
				Config: testAccOpenStackNetworkingSubnetV2DataSource_gatewayIP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetV2DataSourceID("data.openstack_networking_subnet_v2.subnet_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnet_v2.subnet_1", "all_tags.#", "2"),
				),
			},
		},
	})
}

func TestAccNetworkingV2SubnetDataSource_networkIdAttribute(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingSubnetV2DataSource_networkIdAttribute,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetV2DataSourceID("data.openstack_networking_subnet_v2.subnet_1"),
					testAccCheckNetworkingSubnetV2DataSourceGoodNetwork("data.openstack_networking_subnet_v2.subnet_1", "openstack_networking_network_v2.network_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnet_v2.subnet_1", "tags.#", "1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnet_v2.subnet_1", "all_tags.#", "2"),
					testAccCheckNetworkingPortV2ID("openstack_networking_port_v2.port_1"),
				),
			},
		},
	})
}

func TestAccNetworkingV2SubnetDataSource_subnetPoolIdAttribute(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingSubnetV2DataSource_subnetPoolIdAttribute,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetV2DataSourceID("data.openstack_networking_subnet_v2.subnet_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnet_v2.subnet_1", "tags.#", "2"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnet_v2.subnet_1", "all_tags.#", "2"),
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

func testAccCheckNetworkingPortV2ID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find port resource: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Port resource ID not set")
		}

		return nil
	}
}

func testAccCheckNetworkingSubnetV2DataSourceGoodNetwork(n1, n2 string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds1, ok := s.RootModule().Resources[n1]
		if !ok {
			return fmt.Errorf("Can't find subnet data source: %s", n1)
		}

		if ds1.Primary.ID == "" {
			return fmt.Errorf("Subnet data source ID not set")
		}

		rs2, ok := s.RootModule().Resources[n2]
		if !ok {
			return fmt.Errorf("Can't find network resource: %s", n2)
		}

		if rs2.Primary.ID == "" {
			return fmt.Errorf("Network resource ID not set")
		}

		if rs2.Primary.ID != ds1.Primary.Attributes["network_id"] {
			return fmt.Errorf("Network id and subnet network_id don't match")
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
  description = "my subnet description"
  cidr = "192.168.199.0/24"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  tags = [
    "foo",
    "bar",
  ]
}
`

const testAccOpenStackNetworkingSubnetV2DataSource_subnetWithSubnetPool = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnetpool_v2" "subnetpool_1" {
  name = "my_ipv4_pool"
  prefixes = ["10.11.12.0/24"]
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "10.11.12.0/25"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  subnetpool_id = "${openstack_networking_subnetpool_v2.subnetpool_1.id}"
  tags = [
    "foo",
    "bar",
  ]
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
  tags = []
}
`, testAccOpenStackNetworkingSubnetV2DataSource_subnet)

var testAccOpenStackNetworkingSubnetV2DataSource_dhcpEnabled = fmt.Sprintf(`
%s

data "openstack_networking_subnet_v2" "subnet_1" {
  network_id = "${openstack_networking_network_v2.network_1.id}"
  dhcp_enabled = true
  tags = [
    "bar",
  ]
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

var testAccOpenStackNetworkingSubnetV2DataSource_networkIdAttribute = fmt.Sprintf(`
%s

data "openstack_networking_subnet_v2" "subnet_1" {
  subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
  tags = [
    "foo",
  ]
}

resource "openstack_networking_port_v2" "port_1" {
  name            = "test_port"
  network_id      = "${data.openstack_networking_subnet_v2.subnet_1.network_id}"
  admin_state_up  = "true"
}

`, testAccOpenStackNetworkingSubnetV2DataSource_subnet)

var testAccOpenStackNetworkingSubnetV2DataSource_subnetPoolIdAttribute = fmt.Sprintf(`
%s

data "openstack_networking_subnet_v2" "subnet_1" {
  subnetpool_id = "${openstack_networking_subnet_v2.subnet_1.subnetpool_id}"
  tags = [
    "foo",
    "bar",
  ]
}
`, testAccOpenStackNetworkingSubnetV2DataSource_subnetWithSubnetPool)
