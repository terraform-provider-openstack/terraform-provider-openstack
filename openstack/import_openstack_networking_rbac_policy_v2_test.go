package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNetworkingV2RBACPolicy_importBasic(t *testing.T) {
	resourceName := "openstack_networking_rbac_policy_v2.rbac_policy_1"

	projectName := "ACCPTTEST-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2RBACPolicyDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RBACPolicyBasic(projectName),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
