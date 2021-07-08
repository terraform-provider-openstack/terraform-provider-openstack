package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccLBV1Pool_importBasic(t *testing.T) {
	resourceName := "openstack_lb_pool_v1.pool_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckDeprecated(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV1PoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbV1PoolBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
