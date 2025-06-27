package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumetypes"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/projects"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccBlockstorageV3VolumeTypeAccess_basic(t *testing.T) {
	var project projects.Project

	projectName := "ACCPTTEST-" + acctest.RandString(5)

	var vt volumetypes.VolumeType

	vtName := "ACCPTTEST-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockstorageV3VolumeTypeAccessDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccBlockstorageV3VolumeTypeAccessBasic(projectName, vtName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectExists(t.Context(), "openstack_identity_project_v3.project_1", &project),
					testAccCheckBlockStorageVolumeTypeV3Exists(t.Context(), "openstack_blockstorage_volume_type_v3.volume_type_1", &vt),
					testAccCheckBlockstorageV3VolumeTypeAccessExists(t.Context(), "openstack_blockstorage_volume_type_access_v3.volume_type_access"),
					resource.TestCheckResourceAttrPtr(
						"openstack_blockstorage_volume_type_access_v3.volume_type_access", "project_id", &project.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_blockstorage_volume_type_access_v3.volume_type_access", "volume_type_id", &vt.ID),
				),
			},
		},
	})
}

func testAccCheckBlockstorageV3VolumeTypeAccessDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		blockStorageClient, err := config.BlockStorageV3Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack block storage client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_blockstorage_volume_type_access_v3" {
				continue
			}

			vtid, pid, err := parsePairedIDs(rs.Primary.ID, "openstack_blockstorage_volume_type_access_v3")
			if err != nil {
				return err
			}

			allPages, err := volumetypes.ListAccesses(blockStorageClient, vtid).AllPages(ctx)
			if err == nil {
				allAccesses, err := volumetypes.ExtractAccesses(allPages)
				if err == nil {
					for _, access := range allAccesses {
						if access.VolumeTypeID == vtid && access.ProjectID == pid {
							return errors.New("VolumeType access still exists")
						}
					}
				}
			}
		}

		return nil
	}
}

func testAccCheckBlockstorageV3VolumeTypeAccessExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		blockStorageClient, err := config.BlockStorageV3Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack block storage client: %w", err)
		}

		vtid, pid, err := parsePairedIDs(rs.Primary.ID, "openstack_blockstorage_volume_type_access_v3")
		if err != nil {
			return err
		}

		allPages, err := volumetypes.ListAccesses(blockStorageClient, vtid).AllPages(ctx)
		if err != nil {
			return fmt.Errorf("Error retrieving accesses for vt: %s", vtid)
		}

		allAccesses, err := volumetypes.ExtractAccesses(allPages)
		if err != nil {
			return fmt.Errorf("Error extracting accesses for vt: %s", vtid)
		}

		found := false

		for _, access := range allAccesses {
			if access.VolumeTypeID == vtid && access.ProjectID == pid {
				found = true

				break
			}
		}

		if !found {
			return fmt.Errorf("Volume Type Access not found for vtid/pid: %s", rs.Primary.ID)
		}

		return nil
	}
}

func testAccBlockstorageV3VolumeTypeAccessBasic(projectName, vtName string) string {
	return fmt.Sprintf(`
resource "openstack_identity_project_v3" "project_1" {
  name = "%s"
}

resource "openstack_blockstorage_volume_type_v3" "volume_type_1" {
  name = "%s"
  is_public = false
}

resource "openstack_blockstorage_volume_type_access_v3" "volume_type_access" {
  project_id = openstack_identity_project_v3.project_1.id
  volume_type_id = openstack_blockstorage_volume_type_v3.volume_type_1.id
}
`, projectName, vtName)
}
