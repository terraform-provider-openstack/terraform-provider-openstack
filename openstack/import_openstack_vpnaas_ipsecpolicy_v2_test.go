package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIPSecPolicy_importBasic(t *testing.T) {
	resourceName := "openstack_vpnaas_ipsec_policy_v2.policy_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckVPN(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIPSecPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecPolicyV2Basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
