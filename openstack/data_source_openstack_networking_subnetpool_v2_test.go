package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNetworkingV2SubnetPoolDataSourceBasic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingSubnetPoolV2DataSourceSubnetPool,
			},
			{
				Config: testAccOpenStackNetworkingSubnetPoolV2DataSourceBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetPoolV2DataSourceID("data.openstack_networking_subnetpool_v2.subnetpool_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnetpool_v2.subnetpool_1", "name", "subnetpool_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnetpool_v2.subnetpool_1", "all_tags.#", "2"),
				),
			},
		},
	})
}

func TestAccNetworkingV2SubnetPoolDataSourceTestQueries(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingSubnetPoolV2DataSourceSubnetPool,
			},
			{
				Config: testAccOpenStackNetworkingSubnetPoolV2DataSourceDefaultQuota(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetPoolV2DataSourceID("data.openstack_networking_subnetpool_v2.subnetpool_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnetpool_v2.subnetpool_1", "all_tags.#", "2"),
				),
			},
			{
				Config: testAccOpenStackNetworkingSubnetPoolV2DataSourcePrefixLengths(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetPoolV2DataSourceID("data.openstack_networking_subnetpool_v2.subnetpool_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnetpool_v2.subnetpool_1", "tags.#", "1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnetpool_v2.subnetpool_1", "all_tags.#", "2"),
				),
			},
			{
				Config: testAccOpenStackNetworkingSubnetPoolV2DataSourceDescription(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSubnetPoolV2DataSourceID("data.openstack_networking_subnetpool_v2.subnetpool_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnetpool_v2.subnetpool_1", "tags.#", "1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnetpool_v2.subnetpool_1", "all_tags.#", "2"),
				),
			},
		},
	})
}

func testAccCheckNetworkingSubnetPoolV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find subnetpool data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Subnetpool data source ID not set")
		}

		return nil
	}
}

const testAccOpenStackNetworkingSubnetPoolV2DataSourceSubnetPool = `
resource "openstack_networking_subnetpool_v2" "subnetpool_1" {
  name = "subnetpool_1"
  description = "terraform subnetpool acceptance test"

  prefixes = ["10.10.0.0/16", "10.11.11.0/24"]

  default_quota = 4

  default_prefixlen = 25
  min_prefixlen = 24
  max_prefixlen = 30

  tags = [
    "foo",
    "bar",
  ]
}
`

func testAccOpenStackNetworkingSubnetPoolV2DataSourceBasic() string {
	return fmt.Sprintf(`
%s

data "openstack_networking_subnetpool_v2" "subnetpool_1" {
  name = "${openstack_networking_subnetpool_v2.subnetpool_1.name}"
}
`, testAccOpenStackNetworkingSubnetPoolV2DataSourceSubnetPool)
}

func testAccOpenStackNetworkingSubnetPoolV2DataSourceDefaultQuota() string {
	return fmt.Sprintf(`
%s

data "openstack_networking_subnetpool_v2" "subnetpool_1" {
  default_quota = 4
}
`, testAccOpenStackNetworkingSubnetPoolV2DataSourceSubnetPool)
}

func testAccOpenStackNetworkingSubnetPoolV2DataSourcePrefixLengths() string {
	return fmt.Sprintf(`
%s

data "openstack_networking_subnetpool_v2" "subnetpool_1" {
  default_prefixlen = 25
  min_prefixlen = 24
  max_prefixlen = 30
  tags = [
    "foo",
  ]
}
`, testAccOpenStackNetworkingSubnetPoolV2DataSourceSubnetPool)
}

func testAccOpenStackNetworkingSubnetPoolV2DataSourceDescription() string {
	return fmt.Sprintf(`
%s

data "openstack_networking_subnetpool_v2" "subnetpool_1" {
  description = "${openstack_networking_subnetpool_v2.subnetpool_1.description}"
  tags = [
    "bar",
  ]
}
`, testAccOpenStackNetworkingSubnetPoolV2DataSourceSubnetPool)
}
