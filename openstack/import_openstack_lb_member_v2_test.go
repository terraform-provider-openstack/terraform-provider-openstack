package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccLBV2Member_importBasic(t *testing.T) {
	memberResourceName := "openstack_lb_member_v2.member_1"
	poolResourceName := "openstack_lb_pool_v2.pool_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV2MemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccLbV2MemberConfigBasic,
			},

			{
				ResourceName:      memberResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccLBV2MemberImportID(poolResourceName, memberResourceName),
			},
		},
	})
}

func testAccLBV2MemberImportID(poolResource, memberResource string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		pool, ok := s.RootModule().Resources[poolResource]
		if !ok {
			return "", fmt.Errorf("Pool not found: %s", poolResource)
		}

		member, ok := s.RootModule().Resources[memberResource]
		if !ok {
			return "", fmt.Errorf("Member not found: %s", memberResource)
		}

		return fmt.Sprintf("%s/%s", pool.Primary.ID, member.Primary.ID), nil
	}
}
