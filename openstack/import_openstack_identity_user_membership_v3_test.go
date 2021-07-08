package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIdentityV3UserMembership_importBasic(t *testing.T) {
	resourceName := "openstack_identity_user_membership_v3.user_membership_1"

	groupName := fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))
	userName := fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIdentityV3UserMembershipDestroy,
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
