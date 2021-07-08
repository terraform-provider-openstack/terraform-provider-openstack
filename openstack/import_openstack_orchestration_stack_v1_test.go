package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccOrchestrationStackV1_importBasic(t *testing.T) {
	resourceName := "openstack_orchestration_stack_v1.stack_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckOrchestrationV1StackDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOrchestrationV1StackBasic,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"environment_opts",
					"template_opts",
				},
			},
		},
	})
}
