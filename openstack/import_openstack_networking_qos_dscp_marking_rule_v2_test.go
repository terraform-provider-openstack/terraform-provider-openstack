package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNetworkingV2QoSDSCPMarkingRule_importBasic(t *testing.T) {
	resourceName := "openstack_networking_qos_dscp_marking_rule_v2.dscp_marking_rule_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2QoSDSCPMarkingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2QoSDSCPMarkingRule_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
