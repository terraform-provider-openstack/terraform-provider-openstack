package openstack

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccNetworkingV2PortDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.openstack_networking_port_v2.port_1", "id",
						"openstack_networking_port_v2.port_1", "id"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_networking_port_v2.port_2", "id",
						"openstack_networking_port_v2.port_2", "id"),
				),
			},
		},
	})
}

const testAccNetworkingV2PortDataSource_basic = `
resource "openstack_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

data "openstack_networking_secgroup_v2" "default" {
  name = "default"
}

resource "openstack_networking_port_v2" "port_1" {
  name           = "port"
  description    = "test port"
  network_id     = "${openstack_networking_network_v2.network_1.id}"
  admin_state_up = "true"

  security_group_ids = [
    "${data.openstack_networking_secgroup_v2.default.id}",
  ]
}

resource "openstack_networking_port_v2" "port_2" {
  name               = "port"
  description        = "test port"
  network_id         = "${openstack_networking_network_v2.network_1.id}"
  admin_state_up     = "true"
  no_security_groups = "true"

  tags = [
    "foo",
    "bar",
  ]
}

data "openstack_networking_port_v2" "port_1" {
  name           = "${openstack_networking_port_v2.port_1.name}"
  admin_state_up = "${openstack_networking_port_v2.port_2.admin_state_up}"

  security_group_ids = [
    "${data.openstack_networking_secgroup_v2.default.id}",
  ]
}

data "openstack_networking_port_v2" "port_2" {
  name           = "${openstack_networking_port_v2.port_1.name}"
  admin_state_up = "${openstack_networking_port_v2.port_2.admin_state_up}"

  tags = [
    "foo",
    "bar",
  ]
}
`
