package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOpenStackNetworkingNetworkV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOpenStackNetworkingNetworkV2DataSource_network,
			},
			resource.TestStep{
				Config: testAccOpenStackNetworkingNetworkV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.openstack_networking_network_v2.net"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.net", "name", "tf_test_network"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.net", "description", "my network description"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.net", "admin_state_up", "true"),
				),
			},
		},
	})
}

func TestAccOpenStackNetworkingNetworkV2DataSource_subnet(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOpenStackNetworkingNetworkV2DataSource_network,
			},
			resource.TestStep{
				Config: testAccOpenStackNetworkingNetworkV2DataSource_subnet,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.openstack_networking_network_v2.net"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.net", "name", "tf_test_network"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.net", "admin_state_up", "true"),
				),
			},
		},
	})
}

func TestAccOpenStackNetworkingNetworkV2DataSource_networkID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOpenStackNetworkingNetworkV2DataSource_network,
			},
			resource.TestStep{
				Config: testAccOpenStackNetworkingNetworkV2DataSource_networkID,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.openstack_networking_network_v2.net"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.net", "name", "tf_test_network"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.net", "admin_state_up", "true"),
				),
			},
		},
	})
}

func TestAccOpenStackNetworkingNetworkV2DataSource_externalExplicit(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOpenStackNetworkingNetworkV2DataSource_externalExplicit,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.openstack_networking_network_v2.net"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.net", "name", OS_POOL_NAME),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.net", "admin_state_up", "true"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.net", "external", "true"),
				),
			},
		},
	})
}

func TestAccOpenStackNetworkingNetworkV2DataSource_externalImplicit(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOpenStackNetworkingNetworkV2DataSource_externalImplicit,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.openstack_networking_network_v2.net"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.net", "name", OS_POOL_NAME),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.net", "admin_state_up", "true"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.net", "external", "true"),
				),
			},
		},
	})
}

func TestAccOpenStackNetworkingNetworkV2DataSource_transparent_vlan(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckTransparentVLAN(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2Network_transparent_vlan,
			},
			resource.TestStep{
				Config: testAccOpenStackNetworkingNetworkV2DataSource_transparent_vlan,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.openstack_networking_network_v2.network_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.network_1", "name", "network_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.network_1", "admin_state_up", "true"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.network_1", "transparent_vlan", "true"),
				),
			},
		},
	})
}

func testAccCheckNetworkingNetworkV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find network data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Network data source ID not set")
		}

		return nil
	}
}

const testAccOpenStackNetworkingNetworkV2DataSource_network = `
resource "openstack_networking_network_v2" "net" {
        name = "tf_test_network"
        description = "my network description"
        admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet" {
  name = "tf_test_subnet"
  cidr = "192.168.199.0/24"
  no_gateway = true
  network_id = "${openstack_networking_network_v2.net.id}"
}
`

var testAccOpenStackNetworkingNetworkV2DataSource_basic = fmt.Sprintf(`
%s

data "openstack_networking_network_v2" "net" {
	name = "${openstack_networking_network_v2.net.name}"
        description = "${openstack_networking_network_v2.net.description}"
}
`, testAccOpenStackNetworkingNetworkV2DataSource_network)

var testAccOpenStackNetworkingNetworkV2DataSource_subnet = fmt.Sprintf(`
%s

data "openstack_networking_network_v2" "net" {
	matching_subnet_cidr = "${openstack_networking_subnet_v2.subnet.cidr}"
}
`, testAccOpenStackNetworkingNetworkV2DataSource_network)

var testAccOpenStackNetworkingNetworkV2DataSource_networkID = fmt.Sprintf(`
%s

data "openstack_networking_network_v2" "net" {
	network_id = "${openstack_networking_network_v2.net.id}"
}
`, testAccOpenStackNetworkingNetworkV2DataSource_network)

var testAccOpenStackNetworkingNetworkV2DataSource_externalExplicit = fmt.Sprintf(`
data "openstack_networking_network_v2" "net" {
	name = "%s"
	external = "true"
}
`, OS_POOL_NAME)

var testAccOpenStackNetworkingNetworkV2DataSource_externalImplicit = fmt.Sprintf(`
data "openstack_networking_network_v2" "net" {
	name = "%s"
}
`, OS_POOL_NAME)

var testAccOpenStackNetworkingNetworkV2DataSource_transparent_vlan = fmt.Sprintf(`
%s

data "openstack_networking_network_v2" "network_1" {
    transparent_vlan = "${openstack_networking_network_v2.network_1.transparent_vlan}"
}
`, testAccNetworkingV2Network_transparent_vlan)
