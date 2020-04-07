package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccLBV2Members_importBasic(t *testing.T) {
	membersResourceName := "openstack_lb_members_v2.members_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckLB(t)
			testAccPreCheckUseOctavia(t)
			testAccPreCheckOctaviaBatchMembersEnv(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2MembersDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccLBV2MembersConfig_basic,
			},

			{
				ResourceName:      membersResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
