package openstack

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccLBV2Listener_importBasic(t *testing.T) {
	resourceName := "openstack_lb_listener_v2.listener_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckLB(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2ListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccLBV2ListenerConfig_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
