package openstack

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNetworkingV2_tags(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2ConfigCreate(),
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
				Config: testAccNetworkingV2ConfigUpdate(),
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

// Shared acceptance test for network tags.
func testAccCheckNetworkingV2Tags(name string, tags []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]

		if !ok {
			return fmt.Errorf("resource not found: %s", name)
		}

		if _, ok := rs.Primary.Attributes["tags.#"]; !ok {
			return fmt.Errorf("resource tags not found: %s.tags", name)
		}

		var rtags []string
		for key, val := range rs.Primary.Attributes {
			if !strings.HasPrefix(key, "tags.") {
				continue
			}

			if key == "tags.#" {
				continue
			}

			rtags = append(rtags, val)
		}

		sort.Strings(rtags)
		sort.Strings(tags)
		if !reflect.DeepEqual(rtags, tags) {
			return fmt.Errorf(
				"%s.tags: expected: %#v, got %#v", name, tags, rtags)
		}
		return nil
	}
}

const testAccNetworkingV2Config = `
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

const testAccNetworkingV2TagsCreate = `["a", "b", "c"]`

const testAccNetworkingV2TagsUpdate = `["a", "b", "c", "d"]`

func testAccNetworkingV2ConfigCreate() string {
	return fmt.Sprintf(testAccNetworkingV2Config, testAccNetworkingV2TagsCreate)
}

func testAccNetworkingV2ConfigUpdate() string {
	return fmt.Sprintf(testAccNetworkingV2Config, testAccNetworkingV2TagsUpdate)
}
