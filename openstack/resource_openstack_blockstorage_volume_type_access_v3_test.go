package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumetypes"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
)

func TestAccBlockstorageV3VolumeTypeAccess_basic(t *testing.T) {
	var project projects.Project
	var projectName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))

	var vt volumetypes.VolumeType
	var vtName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockstorageV3VolumeTypeAccessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockstorageV3VolumeTypeAccessBasic(projectName, vtName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectExists("openstack_identity_project_v3.project_1", &project),
					testAccCheckBlockStorageVolumeTypeV3Exists("openstack_blockstorage_volume_type_v3.volume_type_1", &vt),
					testAccCheckBlockstorageV3VolumeTypeAccessExists("openstack_blockstorage_volume_type_access_v3.volume_type_access"),
					resource.TestCheckResourceAttrPtr(
						"openstack_blockstorage_volume_type_access_v3.volume_type_access", "project_id", &project.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_blockstorage_volume_type_access_v3.volume_type_access", "volume_type_id", &vt.ID),
				),
			},
		},
	})
}

func testAccCheckBlockstorageV3VolumeTypeAccessDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	blockStorageClient, err := config.BlockStorageV3Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_blockstorage_volume_type_access_v3" {
			continue
		}

		vtid, pid, err := parseVolumeTypeAccessID(rs.Primary.ID)
		if err != nil {
			return err
		}

		allPages, err := volumetypes.ListAccesses(blockStorageClient, vtid).AllPages()
		if err == nil {
			allAccesses, err := volumetypes.ExtractAccesses(allPages)
			if err == nil {
				for _, access := range allAccesses {
					if access.VolumeTypeID == vtid && access.ProjectID == pid {
						return fmt.Errorf("VolumeType access still exists")
					}
				}
			}
		}
	}

	return nil
}

func testAccCheckBlockstorageV3VolumeTypeAccessExists(n string) resource.TestCheckFunc {
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

		vtid, pid, err := parseVolumeTypeAccessID(rs.Primary.ID)
		if err != nil {
			return err
		}

		allPages, err := volumetypes.ListAccesses(blockStorageClient, vtid).AllPages()
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
  project_id = "${openstack_identity_project_v3.project_1.id}"
  volume_type_id = "${openstack_blockstorage_volume_type_v3.volume_type_1.id}"
}
`, projectName, vtName)
}
