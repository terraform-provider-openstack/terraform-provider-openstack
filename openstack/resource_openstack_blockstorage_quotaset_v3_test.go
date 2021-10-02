package openstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/extensions/quotasets"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumetypes"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
)

func TestAccBlockStorageQuotasetV3_basic(t *testing.T) {
	var (
		project    projects.Project
		quotaset   quotasets.QuotaSet
		volumeType volumetypes.VolumeType
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockStorageQuotasetV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageQuotasetV3Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectExists("openstack_identity_project_v3.project_1", &project),
					testAccCheckBlockStorageQuotasetV3Exists("openstack_blockstorage_quotaset_v3.quotaset_1", &quotaset),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v3.quotaset_1", "volumes", "2"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v3.quotaset_1", "snapshots", "2"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v3.quotaset_1", "gigabytes", "2"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v3.quotaset_1", "per_volume_gigabytes", "1"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v3.quotaset_1", "backups", "2"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v3.quotaset_1", "backup_gigabytes", "1"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v3.quotaset_1", "groups", "1"),
				),
			},
			{
				Config: testAccBlockStorageQuotasetV3Update1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectExists("openstack_identity_project_v3.project_1", &project),
					testAccCheckBlockStorageVolumeTypeV3Exists("openstack_blockstorage_volume_type_v3.volume_type_1", &volumeType),
					testAccCheckBlockStorageQuotasetV3Exists("openstack_blockstorage_quotaset_v3.quotaset_1", &quotaset),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v3.quotaset_1", "volumes", "3"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v3.quotaset_1", "snapshots", "3"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v3.quotaset_1", "gigabytes", "4"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v3.quotaset_1", "per_volume_gigabytes", "1"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v3.quotaset_1", "backups", "2"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v3.quotaset_1", "backup_gigabytes", "1"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v3.quotaset_1", "groups", "1"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v3.quotaset_1", "volume_type_quota.%", "9"), // Extra:map[string]interface {}{"gigabytes___DEFAULT__":-1, "gigabytes_foo":100, "gigabytes_lvmdriver-1":-1, "snapshots___DEFAULT__":-1, "snapshots_foo":100, "snapshots_lvmdriver-1":-1, "volumes___DEFAULT__":-1, "volumes_foo":100, "volumes_lvmdriver-1":-1}}
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v3.quotaset_1", "volume_type_quota.volumes_foo", "100"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v3.quotaset_1", "volume_type_quota.snapshots_foo", "100"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v3.quotaset_1", "volume_type_quota.gigabytes_foo", "100"),
				),
			},
			{
				Config: testAccBlockStorageQuotasetV3Update2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectExists("openstack_identity_project_v3.project_1", &project),
					testAccCheckBlockStorageVolumeTypeV3Exists("openstack_blockstorage_volume_type_v3.volume_type_1", &volumeType),
					testAccCheckBlockStorageQuotasetV3Exists("openstack_blockstorage_quotaset_v3.quotaset_1", &quotaset),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v3.quotaset_1", "volumes", "3"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v3.quotaset_1", "snapshots", "3"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v3.quotaset_1", "gigabytes", "4"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v3.quotaset_1", "per_volume_gigabytes", "2"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v3.quotaset_1", "backups", "4"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v3.quotaset_1", "backup_gigabytes", "4"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v3.quotaset_1", "groups", "4"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v3.quotaset_1", "volume_type_quota.%", "9"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v3.quotaset_1", "volume_type_quota.volumes_foo", "10"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v3.quotaset_1", "volume_type_quota.snapshots_foo", "-1"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_quotaset_v3.quotaset_1", "volume_type_quota.gigabytes_foo", "-1"),
				),
			},
		},
	})
}

func testAccCheckBlockStorageQuotasetV3Exists(n string, quotaset *quotasets.QuotaSet) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		blockStorageClient, err := config.BlockStorageV3Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
		}

		projectID := strings.Split(rs.Primary.ID, "/")[0]

		found, err := quotasets.Get(blockStorageClient, projectID).Extract()
		if err != nil {
			return err
		}

		if found.ID != projectID {
			return fmt.Errorf("Quotaset not found")
		}

		*quotaset = *found

		return nil
	}
}

func testAccCheckBlockStorageQuotasetV3Destroy(s *terraform.State) error {
	err := testAccCheckIdentityV3ProjectDestroy(s)
	if err != nil {
		return err
	}

	err = testAccCheckBlockStorageVolumeTypeV3Destroy(s)
	if err != nil {
		return err
	}

	return nil
}

const testAccBlockStorageQuotasetV3Basic = `
resource "openstack_identity_project_v3" "project_1" {
  name = "project_1"
}

resource "openstack_blockstorage_quotaset_v3" "quotaset_1" {
  project_id            = "${openstack_identity_project_v3.project_1.id}"
  volumes               = 2
  snapshots             = 2
  gigabytes             = 2
  per_volume_gigabytes  = 1
  backups               = 2
  backup_gigabytes      = 1
  groups                = 1
}
`

// NOTE: __DEFAULT__ and lvmdriver-1 types were added while migrating to SDK V2 since we had errors like
//     TestAccBlockStorageQuotasetV3_basic: resource_openstack_blockstorage_quotaset_v3_test.go:23: Step 2/3 error: After applying this test step, the plan was not empty.
//         stdout:
//
//
//         Terraform used the selected providers to generate the following execution
//         plan. Resource actions are indicated with the following symbols:
//           ~ update in-place
//
//         Terraform will perform the following actions:
//
//           # openstack_blockstorage_quotaset_v3.quotaset_1 will be updated in-place
//           ~ resource "openstack_blockstorage_quotaset_v3" "quotaset_1" {
//                 id                   = "c16b6831bca24f3ea568639501806d80/RegionOne"
//               ~ volume_type_quota    = {
//                   - "gigabytes___DEFAULT__" = "-1" -> null
//                   - "gigabytes_lvmdriver-1" = "-1" -> null
//                   - "snapshots___DEFAULT__" = "-1" -> null
//                   - "snapshots_lvmdriver-1" = "-1" -> null
//                   - "volumes___DEFAULT__"   = "-1" -> null
//                   - "volumes_lvmdriver-1"   = "-1" -> null
//                     # (3 unchanged elements hidden)
//                 }
//                 # (9 unchanged attributes hidden)
//             }
//
//         Plan: 0 to add, 1 to change, 0 to destroy.
const testAccBlockStorageQuotasetV3Update1 = `
resource "openstack_identity_project_v3" "project_1" {
  name = "project_1"
}

resource "openstack_blockstorage_volume_type_v3" "volume_type_1" {
  name = "foo"
  description = "foo"
  is_public = true
}

resource "openstack_blockstorage_quotaset_v3" "quotaset_1" {
  project_id           = "${openstack_identity_project_v3.project_1.id}"
  volumes              = 3
  snapshots            = 3
  gigabytes            = 4
  per_volume_gigabytes = 1
  backups              = 2
  backup_gigabytes     = 1
  groups               = 1
  volume_type_quota     = {
    volumes___DEFAULT__ = -1
    volumes_lvmdriver-1 = -1
    volumes_foo   = 100
    gigabytes___DEFAULT__ = -1
    gigabytes_lvmdriver-1 = -1
    gigabytes_foo = 100
    snapshots___DEFAULT__ = -1
    snapshots_lvmdriver-1 = -1
    snapshots_foo = 100
  }

  depends_on = [openstack_blockstorage_volume_type_v3.volume_type_1]
}
`

const testAccBlockStorageQuotasetV3Update2 = `
resource "openstack_identity_project_v3" "project_1" {
  name = "project_1"
}

resource "openstack_blockstorage_volume_type_v3" "volume_type_1" {
  name = "foo"
  description = "foo"
  is_public = true
}

resource "openstack_blockstorage_quotaset_v3" "quotaset_1" {
  project_id           = "${openstack_identity_project_v3.project_1.id}"
  volumes              = 3
  snapshots            = 3
  gigabytes            = 4
  per_volume_gigabytes = 2
  backups              = 4
  backup_gigabytes     = 4
  groups               = 4
  volume_type_quota     = {
    volumes___DEFAULT__ = -1
    volumes_lvmdriver-1 = -1
    volumes_foo   = 10
    gigabytes___DEFAULT__ = -1
    gigabytes_lvmdriver-1 = -1
    gigabytes_foo = -1
    snapshots___DEFAULT__ = -1
    snapshots_lvmdriver-1 = -1
    snapshots_foo = -1
  }

  depends_on = [openstack_blockstorage_volume_type_v3.volume_type_1]
}
`
