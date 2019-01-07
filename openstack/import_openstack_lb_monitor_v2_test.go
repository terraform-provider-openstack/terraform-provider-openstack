package openstack

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccLBV2Monitor_importBasic(t *testing.T) {
	resourceName := "openstack_lb_monitor_v2.monitor_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckLB(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2MonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccLBV2MonitorConfig_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
