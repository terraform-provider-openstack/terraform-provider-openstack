package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFWRuleV1_importBasic(t *testing.T) {
	resourceName := "openstack_fw_rule_v1.rule_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckFW(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckFWRuleV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWRuleV1Basic2,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
