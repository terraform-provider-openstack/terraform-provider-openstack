package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNetworkingV2SecGroupRule_importBasic(t *testing.T) {
	resourceName := "openstack_networking_secgroup_rule_v2.secgroup_rule_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2SecGroupRuleDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SecGroupRuleBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
