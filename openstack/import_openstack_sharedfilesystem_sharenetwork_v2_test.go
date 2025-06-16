package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccSFSV2ShareNetwork_importBasic(t *testing.T) {
	resourceName := "openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckSFS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSFSV2ShareNetworkDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccSFSV2ShareNetworkConfigBasic(),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
