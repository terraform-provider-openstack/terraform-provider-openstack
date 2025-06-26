package openstack

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/shares"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccSFSV2ShareAccess_basic(t *testing.T) {
	var shareAccess1 shares.AccessRight

	var shareAccess2 shares.AccessRight

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckSFS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSFSV2ShareAccessDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccSFSV2ShareAccessConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSV2ShareAccessExists(t.Context(), "openstack_sharedfilesystem_share_access_v2.share_access_1", &shareAccess1),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_access_v2.share_access_1", "access_type", "ip"),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_access_v2.share_access_1", "access_to", "192.168.199.10"),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_access_v2.share_access_1", "access_level", "rw"),
					resource.TestMatchResourceAttr("openstack_sharedfilesystem_share_access_v2.share_access_1", "share_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
					testAccCheckSFSV2ShareAccessExists(t.Context(), "openstack_sharedfilesystem_share_access_v2.share_access_2", &shareAccess2),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_access_v2.share_access_2", "access_type", "ip"),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_access_v2.share_access_2", "access_to", "192.168.199.11"),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_access_v2.share_access_2", "access_level", "rw"),
					resource.TestMatchResourceAttr("openstack_sharedfilesystem_share_access_v2.share_access_2", "share_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
					testAccCheckSFSV2ShareAccessDiffers(&shareAccess1, &shareAccess2),
				),
			},
			{
				Config: testAccSFSV2ShareAccessConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSV2ShareAccessExists(t.Context(), "openstack_sharedfilesystem_share_access_v2.share_access_1", &shareAccess1),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_access_v2.share_access_1", "access_type", "ip"),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_access_v2.share_access_1", "access_to", "192.168.199.10"),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_access_v2.share_access_1", "access_level", "ro"),
					resource.TestMatchResourceAttr("openstack_sharedfilesystem_share_access_v2.share_access_1", "share_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
					testAccCheckSFSV2ShareAccessExists(t.Context(), "openstack_sharedfilesystem_share_access_v2.share_access_2", &shareAccess2),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_access_v2.share_access_2", "access_type", "ip"),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_access_v2.share_access_2", "access_to", "192.168.199.11"),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_access_v2.share_access_2", "access_level", "ro"),
					resource.TestMatchResourceAttr("openstack_sharedfilesystem_share_access_v2.share_access_2", "share_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
					testAccCheckSFSV2ShareAccessDiffers(&shareAccess1, &shareAccess2),
				),
			},
		},
	})
}

func testAccCheckSFSV2ShareAccessDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		sfsClient, err := config.SharedfilesystemV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_sharedfilesystem_share_access_v2" {
				continue
			}

			var shareID string

			for k, v := range rs.Primary.Attributes {
				if k == "share_id" {
					shareID = v

					break
				}
			}

			access, err := shares.ListAccessRights(ctx, sfsClient, shareID).Extract()
			if err == nil {
				for _, v := range access {
					if v.ID == rs.Primary.ID {
						return fmt.Errorf("Manila share access still exists: %s", rs.Primary.ID)
					}
				}
			}
		}

		return nil
	}
}

func testAccCheckSFSV2ShareAccessExists(ctx context.Context, n string, share *shares.AccessRight) resource.TestCheckFunc {
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

		var shareID string

		for k, v := range rs.Primary.Attributes {
			if k == "share_id" {
				shareID = v

				break
			}
		}

		sfsClient.Microversion = sharedFilesystemV2MinMicroversion

		access, err := shares.ListAccessRights(ctx, sfsClient, shareID).Extract()
		if err != nil {
			return fmt.Errorf("Unable to get %s share: %w", shareID, err)
		}

		var found shares.AccessRight

		for _, v := range access {
			if v.ID == rs.Primary.ID {
				found = v

				break
			}
		}

		if found.ID != rs.Primary.ID {
			return errors.New("ShareAccess not found")
		}

		*share = found

		return nil
	}
}

func testAccCheckSFSV2ShareAccessDiffers(shareAccess1, shareAccess2 *shares.AccessRight) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if shareAccess1.ID != shareAccess2.ID {
			return nil
		}

		return errors.New("Share accesses should differ")
	}
}

const testAccSFSV2ShareAccessConfig = `
resource "openstack_sharedfilesystem_share_v2" "share_1" {
  name             = "nfs_share"
  description      = "test share description"
  share_proto      = "NFS"
  share_type       = "dhss_false"
  size             = 1
}
`

func testAccSFSV2ShareAccessConfigBasic() string {
	return fmt.Sprintf(`
%s

resource "openstack_sharedfilesystem_share_access_v2" "share_access_1" {
  share_id     = openstack_sharedfilesystem_share_v2.share_1.id
  access_type  = "ip"
  access_to    = "192.168.199.10"
  access_level = "rw"
}

resource "openstack_sharedfilesystem_share_access_v2" "share_access_2" {
  share_id     = openstack_sharedfilesystem_share_v2.share_1.id
  access_type  = "ip"
  access_to    = "192.168.199.11"
  access_level = "rw"
}
`, testAccSFSV2ShareAccessConfig)
}

func testAccSFSV2ShareAccessConfigUpdate() string {
	return fmt.Sprintf(`
%s

resource "openstack_sharedfilesystem_share_access_v2" "share_access_1" {
  share_id     = openstack_sharedfilesystem_share_v2.share_1.id
  access_type  = "ip"
  access_to    = "192.168.199.10"
  access_level = "ro"
}

resource "openstack_sharedfilesystem_share_access_v2" "share_access_2" {
  share_id     = openstack_sharedfilesystem_share_v2.share_1.id
  access_type  = "ip"
  access_to    = "192.168.199.11"
  access_level = "ro"
}
`, testAccSFSV2ShareAccessConfig)
}
