package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccLBV2LoadBalancer_importBasic(t *testing.T) {
	resourceName := "openstack_lb_loadbalancer_v2.loadbalancer_1"

	lbProvider := "haproxy"
	if osUseOctavia != "" {
		lbProvider = "octavia"
	}

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckLB(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2LoadBalancerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbV2LoadBalancerConfigBasic(lbProvider),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
