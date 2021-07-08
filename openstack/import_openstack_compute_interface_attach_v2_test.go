package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeV2InterfaceAttachImport_basic(t *testing.T) {
	resourceName := "openstack_compute_interface_attach_v2.ai_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2InterfaceAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InterfaceAttachBasic(),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"admin_pass",
				},
			},
		},
	})
}
