package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkingV2TrunkDataSource_nosubports(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2TrunkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2TrunkDataSourceNoSubports(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.openstack_networking_trunk_v2.trunk_1", "id",
						"openstack_networking_trunk_v2.trunk_1", "id"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_networking_trunk_v2.trunk_1", "name",
						"openstack_networking_trunk_v2.trunk_1", "name"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_networking_trunk_v2.trunk_1", "port_id",
						"openstack_networking_trunk_v2.trunk_1", "port_id"),
				),
			},
		},
	})
}

func TestAccNetworkingV2TrunkDataSource_subports(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2TrunkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2TrunkDataSourceSubports(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.openstack_networking_trunk_v2.trunk_1", "id",
						"openstack_networking_trunk_v2.trunk_1", "id"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_networking_trunk_v2.trunk_1", "name",
						"openstack_networking_trunk_v2.trunk_1", "name"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_networking_trunk_v2.trunk_1", "port_id",
						"openstack_networking_trunk_v2.trunk_1", "port_id"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_trunk_v2.trunk_1", "sub_port.#", "2"),
				),
			},
		},
	})
}

func TestAccNetworkingV2TrunkDataSource_tags(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2TrunkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2TrunkDataSourceTags(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.openstack_networking_trunk_v2.trunk_1", "id",
						"openstack_networking_trunk_v2.trunk_1", "id"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_networking_trunk_v2.trunk_1", "name",
						"openstack_networking_trunk_v2.trunk_1", "name"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_networking_trunk_v2.trunk_1", "port_id",
						"openstack_networking_trunk_v2.trunk_1", "port_id"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_trunk_v2.trunk_1", "tags.#", "1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_trunk_v2.trunk_1", "all_tags.#", "3"),
				),
			},
		},
	})
}

const testAccNetworkingV2TrunkDataSource = `
resource "openstack_networking_network_v2" "network_1" {
  name = "trunk_network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "parent_port_1" {
  name = "parent_port_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
}
`

func testAccNetworkingV2TrunkDataSourceNoSubports() string {
	return fmt.Sprintf(`
%s

resource "openstack_networking_trunk_v2" "trunk_1" {
  name = "trunk_1"
  port_id = "${openstack_networking_port_v2.parent_port_1.id}"
  admin_state_up = "true"
}

data "openstack_networking_trunk_v2" "trunk_1" {
  name = "${openstack_networking_trunk_v2.trunk_1.name}"
}
`, testAccNetworkingV2TrunkDataSource)
}

func testAccNetworkingV2TrunkDataSourceSubports() string {
	return fmt.Sprintf(`
%s

resource "openstack_networking_port_v2" "subport_1" {
  name = "subport_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "subport_2" {
  name = "subport_2"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_trunk_v2" "trunk_1" {
  name = "trunk_1"
  port_id = "${openstack_networking_port_v2.parent_port_1.id}"
  admin_state_up = "true"

  sub_port {
    port_id = "${openstack_networking_port_v2.subport_1.id}"
    segmentation_id = 1
    segmentation_type = "vlan"
  }

  sub_port {
    port_id = "${openstack_networking_port_v2.subport_2.id}"
    segmentation_id = 2
    segmentation_type = "vlan"
  }
}

data "openstack_networking_trunk_v2" "trunk_1" {
  port_id = "${openstack_networking_trunk_v2.trunk_1.port_id}"
}
`, testAccNetworkingV2TrunkDataSource)
}

func testAccNetworkingV2TrunkDataSourceTags() string {
	return fmt.Sprintf(`
%s

resource "openstack_networking_trunk_v2" "trunk_1" {
  name = "trunk_1"
  port_id = "${openstack_networking_port_v2.parent_port_1.id}"
  admin_state_up = "true"

  tags = [
    "foo",
    "bar",
    "baz"
  ]
}

data "openstack_networking_trunk_v2" "trunk_1" {
  admin_state_up = "${openstack_networking_trunk_v2.trunk_1.admin_state_up}"
  tags = [
    "foo",
  ]
}
`, testAccNetworkingV2TrunkDataSource)
}
