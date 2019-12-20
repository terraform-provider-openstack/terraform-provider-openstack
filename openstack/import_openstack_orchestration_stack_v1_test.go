package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccOrchestrationStackV1_importBasic(t *testing.T) {
	resourceName := "openstack_orchestration_stack_v1.stack_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckOrchestrationV1StackDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOrchestrationV1Stack_basic,
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
