package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIdentityV3UserMembership_importBasic(t *testing.T) {
	resourceName := "openstack_identity_user_membership_v3.user_membership_1"

	groupName := "ACCPTTEST-" + acctest.RandString(5)
	userName := "ACCPTTEST-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIdentityV3UserMembershipDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3UserMembershipBasic(groupName, userName),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
