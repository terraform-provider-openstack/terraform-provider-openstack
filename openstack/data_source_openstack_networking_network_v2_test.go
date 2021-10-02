package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOpenStackNetworkingNetworkV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingNetworkV2DataSourceNetwork,
			},
			{
				Config: testAccOpenStackNetworkingNetworkV2DataSourceBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.openstack_networking_network_v2.network_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.network_1", "name", "tf_test_network"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.network_1", "description", "my network description"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.network_1", "admin_state_up", "true"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.network_1", "all_tags.#", "2"),
				),
			},
		},
	})
}

func TestAccOpenStackNetworkingNetworkV2DataSource_subnet(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingNetworkV2DataSourceNetwork,
			},
			{
				Config: testAccOpenStackNetworkingNetworkV2DataSourceSubnet(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.openstack_networking_network_v2.network_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.network_1", "name", "tf_test_network"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.network_1", "admin_state_up", "true"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.network_1", "tags.#", "2"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.network_1", "all_tags.#", "2"),
				),
			},
		},
	})
}

func TestAccOpenStackNetworkingNetworkV2DataSource_networkID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingNetworkV2DataSourceNetwork,
			},
			{
				Config: testAccOpenStackNetworkingNetworkV2DataSourceNetworkID(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.openstack_networking_network_v2.network_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.network_1", "name", "tf_test_network"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.network_1", "admin_state_up", "true"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.network_1", "all_tags.#", "2"),
				),
			},
		},
	})
}

func TestAccOpenStackNetworkingNetworkV2DataSource_externalExplicit(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingNetworkV2DataSourceExternalExplicit(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.openstack_networking_network_v2.network_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.network_1", "name", osPoolName),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.network_1", "admin_state_up", "true"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.network_1", "external", "true"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.network_1", "all_tags.#", "0"),
				),
			},
		},
	})
}

func TestAccOpenStackNetworkingNetworkV2DataSource_externalImplicit(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingNetworkV2DataSourceExternalImplicit(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.openstack_networking_network_v2.network_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.network_1", "name", osPoolName),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.network_1", "admin_state_up", "true"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.network_1", "external", "true"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.network_1", "all_tags.#", "0"),
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
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2NetworkTransparentVlan,
			},
			{
				Config: testAccOpenStackNetworkingNetworkV2DataSourceTransparentVlan(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.openstack_networking_network_v2.network_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.network_1", "name", "network_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.network_1", "admin_state_up", "true"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.network_1", "transparent_vlan", "true"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.network_1", "tags.#", "1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_network_v2.network_1", "all_tags.#", "2"),
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

const testAccOpenStackNetworkingNetworkV2DataSourceNetwork = `
resource "openstack_networking_network_v2" "network_1" {
  name = "tf_test_network"
  description = "my network description"
  admin_state_up = "true"
  tags = [
    "foo",
    "bar",
  ]
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "tf_test_subnet"
  cidr = "192.168.199.0/24"
  no_gateway = true
  network_id = "${openstack_networking_network_v2.network_1.id}"
}
`

func testAccOpenStackNetworkingNetworkV2DataSourceBasic() string {
	return fmt.Sprintf(`
%s

data "openstack_networking_network_v2" "network_1" {
  name = "${openstack_networking_network_v2.network_1.name}"
  description = "${openstack_networking_network_v2.network_1.description}"
}
`, testAccOpenStackNetworkingNetworkV2DataSourceNetwork)
}

func testAccOpenStackNetworkingNetworkV2DataSourceSubnet() string {
	return fmt.Sprintf(`
%s

data "openstack_networking_network_v2" "network_1" {
  matching_subnet_cidr = "${openstack_networking_subnet_v2.subnet_1.cidr}"
  tags = [
    "foo",
    "bar",
  ]
}
`, testAccOpenStackNetworkingNetworkV2DataSourceNetwork)
}

func testAccOpenStackNetworkingNetworkV2DataSourceNetworkID() string {
	return fmt.Sprintf(`
%s

data "openstack_networking_network_v2" "network_1" {
  network_id = "${openstack_networking_network_v2.network_1.id}"
}
`, testAccOpenStackNetworkingNetworkV2DataSourceNetwork)
}

func testAccOpenStackNetworkingNetworkV2DataSourceExternalExplicit() string {
	return fmt.Sprintf(`
data "openstack_networking_network_v2" "network_1" {
  name = "%s"
  external = "true"
}
`, osPoolName)
}

func testAccOpenStackNetworkingNetworkV2DataSourceExternalImplicit() string {
	return fmt.Sprintf(`
data "openstack_networking_network_v2" "network_1" {
  name = "%s"
}
`, osPoolName)
}

func testAccOpenStackNetworkingNetworkV2DataSourceTransparentVlan() string {
	return fmt.Sprintf(`
%s

data "openstack_networking_network_v2" "network_1" {
  transparent_vlan = "${openstack_networking_network_v2.network_1.transparent_vlan}"
  tags = [
    "bar",
  ]
}
`, testAccNetworkingV2NetworkTransparentVlan)
}
