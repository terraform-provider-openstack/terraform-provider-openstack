package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSFSV2ShareType_importBasic(t *testing.T) {
	resourceName := "openstack_sharedfilesystem_sharetype_v2.sharetype_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
			testAccPreCheckSFS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSFSV2ShareTypeDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccSFSV2ShareTypeConfigBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
