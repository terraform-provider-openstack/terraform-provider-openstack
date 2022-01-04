package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkingV2PortDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.openstack_networking_port_v2.port_1", "id",
						"openstack_networking_port_v2.port_1", "id"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_networking_port_v2.port_2", "id",
						"openstack_networking_port_v2.port_2", "id"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_port_v2.port_2", "allowed_address_pairs.#", "2"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_networking_port_v2.port_3", "id",
						"openstack_networking_port_v2.port_1", "id"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_port_v2.port_3", "all_fixed_ips.#", "2"),
				),
			},
		},
	})
}

const testAccNetworkingV2PortDataSourceBasic = `
resource "openstack_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name       = "subnet_1"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  cidr       = "10.0.0.0/24"
  ip_version = 4
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

  fixed_ip {
    subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
  }

  fixed_ip {
    subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
  }
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

  allowed_address_pairs {
    ip_address  = "10.0.0.201"
    mac_address = "fa:16:3e:f8:ab:da"
  }

  allowed_address_pairs {
    ip_address  = "10.0.0.202"
    mac_address = "fa:16:3e:ab:4b:58"
  }
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

data "openstack_networking_port_v2" "port_3" {
  fixed_ip = "${openstack_networking_port_v2.port_1.all_fixed_ips.1}"
}
`
