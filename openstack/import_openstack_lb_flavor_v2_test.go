package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLBV2Flavor_importBasic(t *testing.T) {
	resourceName := "openstack_lb_flavor_v2.flavor_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV2FlavorDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLbV2Flavor,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
