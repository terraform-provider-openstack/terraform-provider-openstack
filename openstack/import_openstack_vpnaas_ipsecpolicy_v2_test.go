package openstack

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccIPSecPolicy_importBasic(t *testing.T) {
	resourceName := "openstack_vpnaas_ipsec_policy_v2.policy_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckVPN(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIPSecPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIPSecPolicyV2_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
