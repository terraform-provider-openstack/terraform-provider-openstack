package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccComputeV2ServerGroup_importBasic(t *testing.T) {
	resourceName := "openstack_compute_servergroup_v2.sg_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2ServerGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2ServerGroupBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
