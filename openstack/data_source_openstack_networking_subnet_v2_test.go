package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNetworkingV2SubnetDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingSubnetV2DataSourceSubnet,
			},
			{
				Config: testAccOpenStackNetworkingSubnetV2DataSourceBasic(),
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
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingSubnetV2DataSourceSubnet,
			},
			{
				Config: testAccOpenStackNetworkingSubnetV2DataSourceCidr(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetV2DataSourceID("data.openstack_networking_subnet_v2.subnet_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnet_v2.subnet_1", "description", "my subnet description"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnet_v2.subnet_1", "all_tags.#", "2"),
				),
			},
			{
				Config: testAccOpenStackNetworkingSubnetV2DataSourceDhcpEnabled(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetV2DataSourceID("data.openstack_networking_subnet_v2.subnet_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnet_v2.subnet_1", "tags.#", "1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnet_v2.subnet_1", "all_tags.#", "2"),
				),
			},
			{
				Config: testAccOpenStackNetworkingSubnetV2DataSourceIPVersion(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetV2DataSourceID("data.openstack_networking_subnet_v2.subnet_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnet_v2.subnet_1", "all_tags.#", "2"),
				),
			},
			{
				Config: testAccOpenStackNetworkingSubnetV2DataSourceGatewayIP(),
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
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingSubnetV2DataSourceNetworkIDAttribute(),
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
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingSubnetV2DataSourceSubnetPoolIDAttribute(),
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

const testAccOpenStackNetworkingSubnetV2DataSourceSubnet = `
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

const testAccOpenStackNetworkingSubnetV2DataSourceSubnetWithSubnetPool = `
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

func testAccOpenStackNetworkingSubnetV2DataSourceBasic() string {
	return fmt.Sprintf(`
%s

data "openstack_networking_subnet_v2" "subnet_1" {
  name = "${openstack_networking_subnet_v2.subnet_1.name}"
}
`, testAccOpenStackNetworkingSubnetV2DataSourceSubnet)
}

func testAccOpenStackNetworkingSubnetV2DataSourceCidr() string {
	return fmt.Sprintf(`
%s

data "openstack_networking_subnet_v2" "subnet_1" {
  cidr = "192.168.199.0/24"
  tags = []
}
`, testAccOpenStackNetworkingSubnetV2DataSourceSubnet)
}

func testAccOpenStackNetworkingSubnetV2DataSourceDhcpEnabled() string {
	return fmt.Sprintf(`
%s

data "openstack_networking_subnet_v2" "subnet_1" {
  network_id = "${openstack_networking_network_v2.network_1.id}"
  dhcp_enabled = true
  tags = [
    "bar",
  ]
}
`, testAccOpenStackNetworkingSubnetV2DataSourceSubnet)
}

func testAccOpenStackNetworkingSubnetV2DataSourceIPVersion() string {
	return fmt.Sprintf(`
%s

data "openstack_networking_subnet_v2" "subnet_1" {
  network_id = "${openstack_networking_network_v2.network_1.id}"
  ip_version = 4
}
`, testAccOpenStackNetworkingSubnetV2DataSourceSubnet)
}

func testAccOpenStackNetworkingSubnetV2DataSourceGatewayIP() string {
	return fmt.Sprintf(`
%s

data "openstack_networking_subnet_v2" "subnet_1" {
  gateway_ip = "${openstack_networking_subnet_v2.subnet_1.gateway_ip}"
}
`, testAccOpenStackNetworkingSubnetV2DataSourceSubnet)
}

func testAccOpenStackNetworkingSubnetV2DataSourceNetworkIDAttribute() string {
	return fmt.Sprintf(`
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

`, testAccOpenStackNetworkingSubnetV2DataSourceSubnet)
}

func testAccOpenStackNetworkingSubnetV2DataSourceSubnetPoolIDAttribute() string {
	return fmt.Sprintf(`
%s

data "openstack_networking_subnet_v2" "subnet_1" {
  subnetpool_id = "${openstack_networking_subnet_v2.subnet_1.subnetpool_id}"
  tags = [
    "foo",
    "bar",
  ]
}
`, testAccOpenStackNetworkingSubnetV2DataSourceSubnetWithSubnetPool)
}
