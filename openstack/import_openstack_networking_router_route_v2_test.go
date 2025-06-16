package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNetworkingV2RouterRoute_importBasic(t *testing.T) {
	resourceName := "openstack_networking_router_route_v2.router_route_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2RouterRouteDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterRouteCreate,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
