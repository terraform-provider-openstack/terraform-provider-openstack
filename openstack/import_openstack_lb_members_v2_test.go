package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccLBV2Members_importBasic(t *testing.T) {
	membersResourceName := "openstack_lb_members_v2.members_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV2MembersDestroy(t.Context()),
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
