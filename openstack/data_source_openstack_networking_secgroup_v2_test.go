package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOpenStackNetworkingSecGroupV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingSecGroupV2DataSource_group,
			},
			{
				Config: testAccOpenStackNetworkingSecGroupV2DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupV2DataSourceID("data.openstack_networking_secgroup_v2.secgroup_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_secgroup_v2.secgroup_1", "name", "secgroup_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_secgroup_v2.secgroup_1", "description", "My neutron security group"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_secgroup_v2.secgroup_1", "tags.#", "1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_secgroup_v2.secgroup_1", "all_tags.#", "2"),
				),
			},
		},
	})
}

func TestAccOpenStackNetworkingSecGroupV2DataSource_secGroupID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingSecGroupV2DataSource_group,
			},
			{
				Config: testAccOpenStackNetworkingSecGroupV2DataSource_secGroupID,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupV2DataSourceID("data.openstack_networking_secgroup_v2.secgroup_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_secgroup_v2.secgroup_1", "name", "secgroup_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_secgroup_v2.secgroup_1", "tags.#", "1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_secgroup_v2.secgroup_1", "all_tags.#", "2"),
				),
			},
		},
	})
}

func testAccCheckNetworkingSecGroupV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find security group data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Security group data source ID not set")
		}

		return nil
	}
}

const testAccOpenStackNetworkingSecGroupV2DataSource_group = `
resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name        = "secgroup_1"
  description = "My neutron security group"
  tags = [
    "foo",
    "bar",
  ]
}
`

var testAccOpenStackNetworkingSecGroupV2DataSource_basic = fmt.Sprintf(`
%s

data "openstack_networking_secgroup_v2" "secgroup_1" {
  name = "${openstack_networking_secgroup_v2.secgroup_1.name}"
  tags = [
    "bar",
  ]
}
`, testAccOpenStackNetworkingSecGroupV2DataSource_group)

var testAccOpenStackNetworkingSecGroupV2DataSource_secGroupID = fmt.Sprintf(`
%s

data "openstack_networking_secgroup_v2" "secgroup_1" {
  secgroup_id = "${openstack_networking_secgroup_v2.secgroup_1.id}"
  tags = [
    "foo",
  ]
}
`, testAccOpenStackNetworkingSecGroupV2DataSource_group)
