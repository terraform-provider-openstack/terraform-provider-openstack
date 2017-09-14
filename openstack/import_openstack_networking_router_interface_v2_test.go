package openstack

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccNetworkingV2RouterInterface_importBasic(t *testing.T) {
	resourceName := "openstack_networking_router_interface_v2.int_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2RouterInterfaceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2RouterInterface_basic_port,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
