package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccLBV2Pool_importBasic(t *testing.T) {
	resourceName := "openstack_lb_pool_v2.pool_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2PoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccLbV2PoolConfigBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
