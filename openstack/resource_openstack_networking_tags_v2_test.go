package openstack

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccNetworkingV2_tags(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2_tags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2Tags(
						"openstack_networking_network_v2.network_1",
						[]string{"a", "b", "c"}),
				),
			},
			resource.TestStep{
				Config: testAccNetworkingV2_tags_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2Tags(
						"openstack_networking_network_v2.network_1",
						[]string{"a", "b", "c", "d"}),
				),
			},
		},
	})
}

const testAccNetworkingV2_tags = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
  tags = ["a", "b", "c"]
}
`

const testAccNetworkingV2_tags_update = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
  tags = ["a", "b", "c", "d"]
}
`
