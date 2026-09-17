package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/sharetypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccSFSV2ShareTypeAccess_basic(t *testing.T) {
	var project projects.Project

	var sharetype sharetypes.ShareType

	projectName := "ACCPTTEST-" + acctest.RandString(5)
	shareTypeName := "ACCPTTEST-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
			testAccPreCheckSFS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSFSV2ShareTypeAccessDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccSFSV2ShareTypeAccessBasic(projectName, shareTypeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectExists(t.Context(), "openstack_identity_project_v3.project_1", &project),
					testAccCheckSFSV2ShareTypeExists(t.Context(), "openstack_sharedfilesystem_sharetype_v2.sharetype_1", &sharetype),
					testAccCheckSFSV2ShareTypeAccessExists(t.Context(), "openstack_sharedfilesystem_sharetype_access_v2.access_1"),
					resource.TestCheckResourceAttrPtr(
						"openstack_sharedfilesystem_sharetype_access_v2.access_1", "project_id", &project.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_sharedfilesystem_sharetype_access_v2.access_1", "share_type_id", &sharetype.ID),
				),
			},
		},
	})
}

func testAccCheckSFSV2ShareTypeAccessDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		sfsClient, err := config.SharedfilesystemV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_sharedfilesystem_sharetype_access_v2" {
				continue
			}

			shareTypeID, projectID, err := parsePairedIDs(rs.Primary.ID, "openstack_sharedfilesystem_sharetype_access_v2")
			if err != nil {
				return err
			}

			if _, err := getSharedFilesystemShareTypeAccess(ctx, sfsClient, shareTypeID, projectID); err == nil {
				return errors.New("Manila share type access still exists")
			}
		}

		return nil
	}
}

func testAccCheckSFSV2ShareTypeAccessExists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		sfsClient, err := config.SharedfilesystemV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %w", err)
		}

		shareTypeID, projectID, err := parsePairedIDs(rs.Primary.ID, "openstack_sharedfilesystem_sharetype_access_v2")
		if err != nil {
			return err
		}

		if _, err := getSharedFilesystemShareTypeAccess(ctx, sfsClient, shareTypeID, projectID); err != nil {
			return fmt.Errorf("Share type access not found for %s: %w", rs.Primary.ID, err)
		}

		return nil
	}
}

func testAccSFSV2ShareTypeAccessBasic(projectName, shareTypeName string) string {
	return fmt.Sprintf(`
resource "openstack_identity_project_v3" "project_1" {
  name = "%s"
}

resource "openstack_sharedfilesystem_sharetype_v2" "sharetype_1" {
  name      = "%s"
  is_public = false

  extra_specs = {
    driver_handles_share_servers = "true"
  }
}

resource "openstack_sharedfilesystem_sharetype_access_v2" "access_1" {
  share_type_id = openstack_sharedfilesystem_sharetype_v2.sharetype_1.id
  project_id    = openstack_identity_project_v3.project_1.id
}
`, projectName, shareTypeName)
}
