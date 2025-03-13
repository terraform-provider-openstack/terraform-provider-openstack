package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNetworkingV2QoSBandwidthLimitRule_importBasic(t *testing.T) {
	resourceName := "openstack_networking_qos_bandwidth_limit_rule_v2.bw_limit_rule_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
			testAccSkipReleasesBelow(t, "stable/yoga")
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2QoSBandwidthLimitRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2QoSBandwidthLimitRuleBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
