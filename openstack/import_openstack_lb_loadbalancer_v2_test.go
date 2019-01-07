package openstack

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccLBV2LoadBalancer_importBasic(t *testing.T) {
	resourceName := "openstack_lb_loadbalancer_v2.loadbalancer_1"

	lbProvider := "haproxy"
	if OS_USE_OCTAVIA != "" {
		lbProvider = "octavia"
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckLB(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2LoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2LoadBalancerConfig_basic(lbProvider),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
