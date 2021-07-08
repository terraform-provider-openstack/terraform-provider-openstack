package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkingV2PortIDsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortIDsDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.openstack_networking_port_ids_v2.ports", "ids.#", "2"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_networking_port_ids_v2.ports", "ids.0",
						"openstack_networking_port_v2.port_1", "id"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_networking_port_ids_v2.ports", "ids.1",
						"openstack_networking_port_v2.port_2", "id"),
					resource.TestCheckResourceAttr("data.openstack_networking_port_ids_v2.port_1", "ids.#", "1"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_networking_port_ids_v2.port_1", "ids.0",
						"openstack_networking_port_v2.port_1", "id"),
					resource.TestCheckResourceAttr("data.openstack_networking_port_ids_v2.port_2", "ids.#", "1"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_networking_port_ids_v2.port_2", "ids.0",
						"openstack_networking_port_v2.port_2", "id"),
				),
			},
		},
	})
}

const testAccNetworkingV2PortIDsDataSourceBasic = `
resource "openstack_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

data "openstack_networking_secgroup_v2" "default" {
  name = "default"
}

resource "openstack_networking_port_v2" "port_1" {
  name           = "port_1"
  description    = "test port ids"
  network_id     = "${openstack_networking_network_v2.network_1.id}"
  admin_state_up = "true"

  security_group_ids = [
    "${data.openstack_networking_secgroup_v2.default.id}",
  ]

  tags = [
    "foo",
    "bar",
    "baz",
  ]
}

resource "openstack_networking_port_v2" "port_2" {
  name           = "port_2"
  description    = "test port ids"
  network_id     = "${openstack_networking_network_v2.network_1.id}"
  admin_state_up = "true"

  security_group_ids = [
    "${data.openstack_networking_secgroup_v2.default.id}",
  ]

  tags = [
    "foo",
    "bar",
    "qux",
  ]
}

data "openstack_networking_port_ids_v2" "ports" {
  admin_state_up = "${openstack_networking_port_v2.port_1.admin_state_up == openstack_networking_port_v2.port_2.admin_state_up ? "true" : "true"}"
  description    = "test port ids"
  sort_direction = "asc"
  sort_key       = "name"

  security_group_ids = [
    "${data.openstack_networking_secgroup_v2.default.id}",
  ]

  tags = [
    "foo",
    "bar",
  ]
}

data "openstack_networking_port_ids_v2" "port_1" {
  admin_state_up = "${openstack_networking_port_v2.port_1.admin_state_up == openstack_networking_port_v2.port_2.admin_state_up ? "true" : "true"}"
  description    = "test port ids"
  sort_direction = "asc"
  sort_key       = "name"

  security_group_ids = [
    "${data.openstack_networking_secgroup_v2.default.id}",
  ]

  tags = [
    "foo",
    "bar",
    "baz",
  ]
}

data "openstack_networking_port_ids_v2" "port_2" {
  admin_state_up = "${openstack_networking_port_v2.port_1.admin_state_up == openstack_networking_port_v2.port_2.admin_state_up ? "true" : "true"}"
  description    = "test port ids"
  sort_direction = "asc"
  sort_key       = "name"

  security_group_ids = [
    "${data.openstack_networking_secgroup_v2.default.id}",
  ]

  tags = [
    "foo",
    "bar",
    "qux",
  ]
}
`
