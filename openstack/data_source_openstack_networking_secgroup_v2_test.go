package openstack

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOpenStackNetworkingSecGroupV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingSecGroupV2DataSourceGroup,
			},
			{
				Config: testAccOpenStackNetworkingSecGroupV2DataSourceBasic(),
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
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingSecGroupV2DataSourceGroup,
			},
			{
				Config: testAccOpenStackNetworkingSecGroupV2DataSourceSecGroupID(),
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

func TestAccOpenStackNetworkingSecGroupV2DataSource_stateful(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingSecGroupV2DataSourceStatefulNotSet,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupV2DataSourceID("data.openstack_networking_secgroup_v2.secgroup_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_secgroup_v2.secgroup_1", "stateful", "true"),
				),
			},
			{
				Config: testAccOpenStackNetworkingSecGroupV2DataSourceStatefulSetFalse,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupV2DataSourceID("data.openstack_networking_secgroup_v2.secgroup_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_secgroup_v2.secgroup_1", "stateful", "false"),
				),
			},
			{
				Config: testAccOpenStackNetworkingSecGroupV2DataSourceStatefulSetTrue,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingSecGroupV2DataSourceID("data.openstack_networking_secgroup_v2.secgroup_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_secgroup_v2.secgroup_1", "stateful", "false"),
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
			return errors.New("Security group data source ID not set")
		}

		return nil
	}
}

const testAccOpenStackNetworkingSecGroupV2DataSourceGroup = `
resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name        = "secgroup_1"
  description = "My neutron security group"
  tags = [
    "foo",
    "bar",
  ]
}
`

func testAccOpenStackNetworkingSecGroupV2DataSourceBasic() string {
	return fmt.Sprintf(`
%s

data "openstack_networking_secgroup_v2" "secgroup_1" {
  name = openstack_networking_secgroup_v2.secgroup_1.name
  tags = [
    "bar",
  ]
}
`, testAccOpenStackNetworkingSecGroupV2DataSourceGroup)
}

func testAccOpenStackNetworkingSecGroupV2DataSourceSecGroupID() string {
	return fmt.Sprintf(`
%s

data "openstack_networking_secgroup_v2" "secgroup_1" {
  secgroup_id = openstack_networking_secgroup_v2.secgroup_1.id
  tags = [
    "foo",
  ]
}
`, testAccOpenStackNetworkingSecGroupV2DataSourceGroup)
}

const testAccOpenStackNetworkingSecGroupV2DataSourceStatefulNotSet = `
resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name        = "secgroup_not_stateful_1"
  description = "My neutron security group"
}

data "openstack_networking_secgroup_v2" "secgroup_1" {
  name = openstack_networking_secgroup_v2.secgroup_1.name
}
`

const testAccOpenStackNetworkingSecGroupV2DataSourceStatefulSetFalse = `
resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name        = "secgroup_not_stateful_1"
  description = "My neutron security group"
  stateful    = false
}

data "openstack_networking_secgroup_v2" "secgroup_1" {
  name     = openstack_networking_secgroup_v2.secgroup_1.name
  stateful = false
}
`

const testAccOpenStackNetworkingSecGroupV2DataSourceStatefulSetTrue = `
resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name        = "secgroup_stateful_1"
  description = "My neutron security group"
  stateful    = true
}

resource "openstack_networking_secgroup_v2" "secgroup_2" {
  name        = "secgroup_stateful_1"
  description = "My neutron security group"
  stateful    = false
}

data "openstack_networking_secgroup_v2" "secgroup_1" {
  name     = openstack_networking_secgroup_v2.secgroup_2.name
  description = openstack_networking_secgroup_v2.secgroup_1.description
  stateful = false
}
`
