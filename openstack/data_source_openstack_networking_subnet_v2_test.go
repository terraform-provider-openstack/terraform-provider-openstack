package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOpenStackNetworkingSubnetV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOpenStackNetworkingSubnetV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetV2DataSourceID("data.openstack_networking_subnet_v2.subnet"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnet_v2.subnet", "name", "tf_test_subnet"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnet_v2.subnet", "cidr", "192.168.199.0/24"),
				),
			},
		},
	})
}

func TestAccOpenStackNetworkingSubnetV2DataSource_name(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOpenStackNetworkingSubnetV2DataSource_name,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetV2DataSourceID("data.openstack_networking_subnet_v2.subnet"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnet_v2.subnet", "name", "tf_test_subnet"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnet_v2.subnet", "cidr", "192.168.199.0/24"),
				),
			},
		},
	})
}

func TestAccOpenStackNetworkingSubnetV2DataSource_cidr(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOpenStackNetworkingSubnetV2DataSource_cidr,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetV2DataSourceID("data.openstack_networking_subnet_v2.subnet"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnet_v2.subnet", "name", "tf_test_subnet"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnet_v2.subnet", "cidr", "192.168.199.0/24"),
				),
			},
		},
	})
}

func testAccCheckNetworkingSubnetV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find network data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Subnet data source ID not set")
		}

		return nil
	}
}

const testAccOpenStackNetworkingSubnetV2DataSource_network = `
resource "openstack_networking_network_v2" "net" {
        name = "tf_test_network"
        admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet" {
  name = "tf_test_subnet"
  cidr = "192.168.199.0/24"
  no_gateway = true
  network_id = "${openstack_networking_network_v2.net.id}"
}
`

var testAccOpenStackNetworkingSubnetV2DataSource_basic = fmt.Sprintf(`
%s

data "openstack_networking_subnet_v2" "subnet" {
	subnet_id = "${openstack_networking_subnet_v2.subnet.id}"
}
`, testAccOpenStackNetworkingSubnetV2DataSource_network)

var testAccOpenStackNetworkingSubnetV2DataSource_name = fmt.Sprintf(`
%s

data "openstack_networking_subnet_v2" "subnet" {
	name = "${openstack_networking_subnet_v2.subnet.name}"
}
`, testAccOpenStackNetworkingSubnetV2DataSource_network)

var testAccOpenStackNetworkingSubnetV2DataSource_cidr = fmt.Sprintf(`
%s

data "openstack_networking_subnet_v2" "subnet" {
	cidr = "${openstack_networking_subnet_v2.subnet.cidr}"
}
`, testAccOpenStackNetworkingSubnetV2DataSource_network)
