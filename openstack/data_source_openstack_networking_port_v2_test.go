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
			resource.TestStep{
				Config: testAccNetworkingV2PortDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.openstack_networking_port_v2.port", "id",
						"openstack_networking_port_v2.port_1", "id"),
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

resource "openstack_networking_port_v2" "port_1" {
  name           = "port"
  description    = "test port"
  network_id     = "${openstack_networking_network_v2.network_1.id}"
  admin_state_up = "true"
}

data "openstack_networking_port_v2" "port" {
  name = "${openstack_networking_port_v2.port_1.name}"
}
`
