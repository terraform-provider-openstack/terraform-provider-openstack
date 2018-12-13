package openstack

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccLBV2L7policy_importBasic(t *testing.T) {
	resourceName := "openstack_lb_l7policy_v2.l7policy_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckLB(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2L7policyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLBV2L7policyConfig + testAccCheckLBV2L7policyConfig_basic,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
