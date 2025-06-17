package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNetworkingV2QoSPolicyImportBasic(t *testing.T) {
	resourceName := "openstack_networking_qos_policy_v2.qos_policy_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2QoSPolicyDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2QoSPolicyBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
