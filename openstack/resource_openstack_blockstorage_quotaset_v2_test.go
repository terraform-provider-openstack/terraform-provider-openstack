package openstack

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/quotasets"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccBlockStorageQuotasetV2_basic(t *testing.T) {
	var (
		project  projects.Project
		quotaset quotasets.QuotaSet
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityV3ProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageQuotasetV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectExists("openstack_identity_project_v3.project_1", &project),
					testAccCheckBlockStorageQuotasetV2Exists("openstack_blockstorage_quotaset_v2.quotaset_1", &quotaset),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v2.quotaset_1", "volumes", "2"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v2.quotaset_1", "snapshots", "2"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v2.quotaset_1", "gigabytes", "2"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v2.quotaset_1", "per_volume_gigabytes", "1"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v2.quotaset_1", "backups", "2"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v2.quotaset_1", "backup_gigabytes", "1"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v2.quotaset_1", "groups", "1"),
				),
			},
			{
				Config: testAccBlockStorageQuotasetV2_update_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectExists("openstack_identity_project_v3.project_1", &project),
					testAccCheckBlockStorageQuotasetV2Exists("openstack_blockstorage_quotaset_v2.quotaset_1", &quotaset),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v2.quotaset_1", "volumes", "3"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v2.quotaset_1", "snapshots", "3"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v2.quotaset_1", "gigabytes", "4"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v2.quotaset_1", "per_volume_gigabytes", "1"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v2.quotaset_1", "backups", "2"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v2.quotaset_1", "backup_gigabytes", "1"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v2.quotaset_1", "groups", "1"),
				),
			},
			{
				Config: testAccBlockStorageQuotasetV2_update_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectExists("openstack_identity_project_v3.project_1", &project),
					testAccCheckBlockStorageQuotasetV2Exists("openstack_blockstorage_quotaset_v2.quotaset_1", &quotaset),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v2.quotaset_1", "volumes", "3"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v2.quotaset_1", "snapshots", "3"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v2.quotaset_1", "gigabytes", "4"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v2.quotaset_1", "per_volume_gigabytes", "2"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v2.quotaset_1", "backups", "4"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v2.quotaset_1", "backup_gigabytes", "4"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v2.quotaset_1", "groups", "4"),
				),
			},
		},
	})
}

func testAccCheckBlockStorageQuotasetV2Exists(n string, quotaset *quotasets.QuotaSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		blockStorageClient, err := config.BlockStorageV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
		}

		found, err := quotasets.Get(blockStorageClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Quotaset not found")
		}

		*quotaset = *found

		return nil
	}
}

const testAccBlockStorageQuotasetV2_basic = `
resource "openstack_identity_project_v3" "project_1" {
  name = "project_1"
}

resource "openstack_blockstorage_quotaset_v2" "quotaset_1" {
  project_id            = "${openstack_identity_project_v3.project_1.id}"
  volumes               = 2
  snapshots             = 2
  gigabytes             = 2
  per_volume_gigabytes = 1
  backups               = 2
  backup_gigabytes      = 1
  groups                = 1
}
`

const testAccBlockStorageQuotasetV2_update_1 = `
resource "openstack_identity_project_v3" "project_1" {
  name = "project_1"
}

resource "openstack_blockstorage_quotaset_v2" "quotaset_1" {
  project_id           = "${openstack_identity_project_v3.project_1.id}"
  volumes              = 3
  snapshots            = 3
  gigabytes            = 4
  per_volume_gigabytes = 1
  backups              = 2
  backup_gigabytes     = 1
  groups               = 1
}
`

const testAccBlockStorageQuotasetV2_update_2 = `
resource "openstack_identity_project_v3" "project_1" {
  name = "project_1"
}

resource "openstack_blockstorage_quotaset_v2" "quotaset_1" {
  project_id           = "${openstack_identity_project_v3.project_1.id}"
  volumes              = 3
  snapshots            = 3
  gigabytes            = 4
  per_volume_gigabytes = 2
  backups              = 4
  backup_gigabytes     = 4
  groups               = 4
}
`
