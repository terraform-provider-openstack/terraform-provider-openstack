package openstack

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccIKEPolicy_importBasic(t *testing.T) {
	resourceName := "openstack_vpnaas_ike_policy_v2.policy_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckVPN(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIKEPolicyV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccIKEPolicyV2_basic,
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
