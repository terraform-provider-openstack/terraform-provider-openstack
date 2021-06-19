package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccSFSV2ShareDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckSFS(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSFSV2ShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSV2ShareDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSV2ShareDataSourceID("data.openstack_sharedfilesystem_share_v2.share_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_sharedfilesystem_share_v2.share_1", "name", "nfs_share"),
					resource.TestCheckResourceAttr(
						"data.openstack_sharedfilesystem_share_v2.share_1", "description", "test share description"),
					resource.TestCheckResourceAttr(
						"data.openstack_sharedfilesystem_share_v2.share_1", "share_proto", "NFS"),
					resource.TestCheckResourceAttr(
						"data.openstack_sharedfilesystem_share_v2.share_1", "is_public", "false"),
					resource.TestCheckResourceAttr(
						"data.openstack_sharedfilesystem_share_v2.share_1", "size", "1"),
				),
			},
		},
	})
}

func testAccCheckSFSV2ShareDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find share data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Share data source ID not set")
		}

		return nil
	}
}

const testAccSFSV2ShareDataSourceBasic = `
resource "openstack_sharedfilesystem_share_v2" "share_1" {
  name        = "nfs_share"
  description = "test share description"
  share_proto = "NFS"
  share_type  = "dhss_false"
  size        = 1
}

data "openstack_sharedfilesystem_share_v2" "share_1" {
  name        = "${openstack_sharedfilesystem_share_v2.share_1.name}"
  description = "${openstack_sharedfilesystem_share_v2.share_1.description}"
}
`
