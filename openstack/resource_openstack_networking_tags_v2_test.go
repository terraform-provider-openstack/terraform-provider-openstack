package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccNetworkingV2_tags(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2_config_create,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2Tags(
						"openstack_networking_network_v2.network_1",
						[]string{"a", "b", "c"}),
					testAccCheckNetworkingV2Tags(
						"openstack_networking_subnet_v2.subnet_1",
						[]string{"a", "b", "c"}),
					testAccCheckNetworkingV2Tags(
						"openstack_networking_subnetpool_v2.subnetpool_1",
						[]string{"a", "b", "c"}),
					testAccCheckNetworkingV2Tags(
						"openstack_networking_port_v2.port_1",
						[]string{"a", "b", "c"}),
					testAccCheckNetworkingV2Tags(
						"openstack_networking_secgroup_v2.secgroup_1",
						[]string{"a", "b", "c"}),
					testAccCheckNetworkingV2Tags(
						"openstack_networking_floatingip_v2.fip_1",
						[]string{"a", "b", "c"}),
					testAccCheckNetworkingV2Tags(
						"openstack_networking_router_v2.router_1",
						[]string{"a", "b", "c"}),
				),
			},
			{
				Config: testAccNetworkingV2_config_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2Tags(
						"openstack_networking_network_v2.network_1",
						[]string{"a", "b", "c", "d"}),
					testAccCheckNetworkingV2Tags(
						"openstack_networking_subnet_v2.subnet_1",
						[]string{"a", "b", "c", "d"}),
					testAccCheckNetworkingV2Tags(
						"openstack_networking_subnetpool_v2.subnetpool_1",
						[]string{"a", "b", "c", "d"}),
					testAccCheckNetworkingV2Tags(
						"openstack_networking_port_v2.port_1",
						[]string{"a", "b", "c", "d"}),
					testAccCheckNetworkingV2Tags(
						"openstack_networking_secgroup_v2.secgroup_1",
						[]string{"a", "b", "c", "d"}),
					testAccCheckNetworkingV2Tags(
						"openstack_networking_floatingip_v2.fip_1",
						[]string{"a", "b", "c", "d"}),
					testAccCheckNetworkingV2Tags(
						"openstack_networking_router_v2.router_1",
						[]string{"a", "b", "c", "d"}),
				),
			},
		},
	})
}

const testAccNetworkingV2_config = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
  tags = %[1]s
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  cidr = "192.168.199.0/24"
  network_id = "${openstack_networking_network_v2.network_1.id}"

  dns_nameservers = ["10.0.16.4", "213.186.33.99"]

  allocation_pools {
    start = "192.168.199.100"
    end = "192.168.199.200"
  }

  tags = %[1]s
}

resource "openstack_networking_subnetpool_v2" "subnetpool_1" {
    name = "subnetpool_1"
    description = "terraform subnetpool acceptance test"

    prefixes = ["10.10.0.0/16", "10.11.11.0/24"]

    default_quota = 4

    default_prefixlen = 25
    min_prefixlen = 24
    max_prefixlen = 30

    tags = %[1]s
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }

  tags = %[1]s
}

resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name = "security_group"
  description = "terraform security group acceptance test"
  tags = %[1]s
}

resource "openstack_networking_floatingip_v2" "fip_1" {
	    tags = %[1]s
}

resource "openstack_networking_router_v2" "router_1" {
    name = "router_1"
    admin_state_up = "true"
    tags = %[1]s
}
`

const testAccNetworkingV2_tags_create = `["a", "b", "c"]`

const testAccNetworkingV2_tags_update = `["a", "b", "c", "d"]`

var testAccNetworkingV2_config_create = fmt.Sprintf(
	testAccNetworkingV2_config, testAccNetworkingV2_tags_create)

var testAccNetworkingV2_config_update = fmt.Sprintf(
	testAccNetworkingV2_config, testAccNetworkingV2_tags_update)
