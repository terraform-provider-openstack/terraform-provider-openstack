package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIdentityV3InheritRoleAssignment_importBasic(t *testing.T) {
	resourceName := "openstack_identity_inherit_role_assignment_v3.role_assignment_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIdentityV3InheritRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3InheritRoleAssignmentBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
