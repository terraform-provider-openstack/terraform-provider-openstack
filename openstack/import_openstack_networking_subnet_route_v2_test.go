package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkingV2SubnetRoute_importBasic(t *testing.T) {
	resourceName := "openstack_networking_subnet_route_v2.subnet_route_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2SubnetRouteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SubnetRouteCreate,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
