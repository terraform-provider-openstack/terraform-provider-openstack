package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccNetworkingV2RBACPolicy_importBasic(t *testing.T) {
	resourceName := "openstack_networking_rbac_policy_v2.rbac_policy_1"
	var projectName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2RBACPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RBACPolicy_basic(projectName),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
