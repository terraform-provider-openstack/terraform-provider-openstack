package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccLBV2Members_importBasic(t *testing.T) {
	membersResourceName := "openstack_lb_members_v2.members_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckUseOctavia(t)
			testAccPreCheckOctaviaBatchMembersEnv(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV2MembersDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccLbV2MembersConfigBasic,
			},

			{
				ResourceName:      membersResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
