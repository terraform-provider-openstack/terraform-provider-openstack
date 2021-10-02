package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOpenStackSubnetsV2SubnetIDsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackSubnetsV2SubnetIDsDataSourceEmpty(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnet_ids_v2.subnets_empty", "ids.#", "0"),
				),
			},
			{
				Config: testAccOpenStackSubnetsV2SubnetIDsDataSourceName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnet_ids_v2.subnets_by_name", "ids.#", "1"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_networking_subnet_ids_v2.subnets_by_name", "ids.0",
						"openstack_networking_subnet_v2.subnet_1", "id"),
				),
			},
			{
				Config: testAccOpenStackSubnetsV2SubnetIDsDataSourceRegex(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnet_ids_v2.subnets_by_name_regex", "ids.#", "1"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_networking_subnet_ids_v2.subnets_by_name_regex", "ids.0",
						"openstack_networking_subnet_v2.subnet_2", "id"),
				),
			},
			{
				Config: testAccOpenStackSubnetsV2SubnetIDsDataSourceTag(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnet_ids_v2.subnets_by_tag", "ids.#", "2"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_networking_subnet_ids_v2.subnets_by_tag", "ids.0",
						"openstack_networking_subnet_v2.subnet_1", "id"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_networking_subnet_ids_v2.subnets_by_tag", "ids.1",
						"openstack_networking_subnet_v2.subnet_2", "id"),
				),
			},
			{
				Config: testAccOpenStackSubnetsV2SubnetIDsDataSourceProperties(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.openstack_networking_subnet_ids_v2.subnets_by_tags", "ids.#", "1"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_networking_subnet_ids_v2.subnets_by_tags", "ids.0",
						"openstack_networking_subnet_v2.subnet_2", "id"),
				),
			},
		},
	})
}

const testAccOpenStackSubnetsV2SubnetIDsDataSource = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_one"
  description = "my subnet description"
  cidr = "192.168.198.0/24"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  tags = [
    "foo",
  ]
}

resource "openstack_networking_subnet_v2" "subnet_2" {
  name = "subnet_two"
  description = "my subnet description"
  cidr = "192.168.199.0/24"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  tags = [
    "foo",
    "bar",
  ]
}
`

func testAccOpenStackSubnetsV2SubnetIDsDataSourceEmpty() string {
	return fmt.Sprintf(`
%s

data "openstack_networking_subnet_ids_v2" "subnets_empty" {
    name = "non-existed-subnet"
}
`, testAccOpenStackSubnetsV2SubnetIDsDataSource)
}

func testAccOpenStackSubnetsV2SubnetIDsDataSourceName() string {
	return fmt.Sprintf(`
%s

data "openstack_networking_subnet_ids_v2" "subnets_by_name" {
    name = "${openstack_networking_subnet_v2.subnet_1.name}"
    description = "${openstack_networking_subnet_v2.subnet_2.description}" # to avoid race condition for further tests
}
`, testAccOpenStackSubnetsV2SubnetIDsDataSource)
}

func testAccOpenStackSubnetsV2SubnetIDsDataSourceRegex() string {
	return fmt.Sprintf(`
%s

data "openstack_networking_subnet_ids_v2" "subnets_by_name_regex" {
    name_regex = "two$"
}
`, testAccOpenStackSubnetsV2SubnetIDsDataSource)
}

func testAccOpenStackSubnetsV2SubnetIDsDataSourceTag() string {
	return fmt.Sprintf(`
%s

data "openstack_networking_subnet_ids_v2" "subnets_by_tag" {
    sort_key = "name"
    sort_direction = "asc"
    tags = [
      "foo",
    ]
}
`, testAccOpenStackSubnetsV2SubnetIDsDataSource)
}

func testAccOpenStackSubnetsV2SubnetIDsDataSourceProperties() string {
	return fmt.Sprintf(`
%s

data "openstack_networking_subnet_ids_v2" "subnets_by_tags" {
    tags = [
      "bar",
      "foo",
    ]
}
`, testAccOpenStackSubnetsV2SubnetIDsDataSource)
}
