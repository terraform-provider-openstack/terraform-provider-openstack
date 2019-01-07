package openstack

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccNetworkingV2SubnetRoute_importBasic(t *testing.T) {
	resourceName := "openstack_networking_subnet_route_v2.subnet_route_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SubnetRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SubnetRoute_create,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
