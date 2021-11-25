package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccBlockStorageV3QuotasetDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV3QuotasetDataSourceBasic,
			},
			{
				Config: testAccBlockStorageV3QuotasetDataSourceSource(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageQuotasetV3DataSourceID("data.openstack_blockstorage_quotaset_v3.source"),
					resource.TestCheckResourceAttrSet("data.openstack_blockstorage_quotaset_v3.source", "volumes"),
					resource.TestCheckResourceAttrSet("data.openstack_blockstorage_quotaset_v3.source", "snapshots"),
					resource.TestCheckResourceAttrSet("data.openstack_blockstorage_quotaset_v3.source", "gigabytes"),
					resource.TestCheckResourceAttrSet("data.openstack_blockstorage_quotaset_v3.source", "per_volume_gigabytes"),
					resource.TestCheckResourceAttrSet("data.openstack_blockstorage_quotaset_v3.source", "backups"),
					resource.TestCheckResourceAttrSet("data.openstack_blockstorage_quotaset_v3.source", "backup_gigabytes"),
					resource.TestCheckResourceAttrSet("data.openstack_blockstorage_quotaset_v3.source", "groups"),
					resource.TestCheckResourceAttrSet("data.openstack_blockstorage_quotaset_v3.source", "volume_type_quota.gigabytes___DEFAULT__"),
				),
			},
		},
	})
}

func testAccCheckBlockStorageQuotasetV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find blockstorage quotaset data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Blockstorage quotaset data source ID not set")
		}

		return nil
	}
}

const testAccBlockStorageV3QuotasetDataSourceBasic = `
resource "openstack_identity_project_v3" "project" {
  name = "test-quotaset-datasource"
}
`

func testAccBlockStorageV3QuotasetDataSourceSource() string {
	return fmt.Sprintf(`
%s

data "openstack_blockstorage_quotaset_v3" "source" {
  project_id = "${openstack_identity_project_v3.project.id}"
}
`, testAccBlockStorageV3QuotasetDataSourceBasic)
}
