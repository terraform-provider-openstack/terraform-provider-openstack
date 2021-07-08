package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSFSV2ShareAccess_importBasic(t *testing.T) {
	shareName := "openstack_sharedfilesystem_share_v2.share_1"
	shareAccessName := "openstack_sharedfilesystem_share_access_v2.share_access_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckSFS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSFSV2ShareAccessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSV2ShareAccessConfigBasic(),
			},

			{
				ResourceName:      shareAccessName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccSFSV2ShareAccessImportID(shareName, shareAccessName),
			},
		},
	})
}

func testAccSFSV2ShareAccessImportID(shareResource, shareAccessResource string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		share, ok := s.RootModule().Resources[shareResource]
		if !ok {
			return "", fmt.Errorf("Share not found: %s", shareResource)
		}

		shareAccess, ok := s.RootModule().Resources[shareAccessResource]
		if !ok {
			return "", fmt.Errorf("Share access not found: %s", shareAccessResource)
		}

		return fmt.Sprintf("%s/%s", share.Primary.ID, shareAccess.Primary.ID), nil
	}
}
