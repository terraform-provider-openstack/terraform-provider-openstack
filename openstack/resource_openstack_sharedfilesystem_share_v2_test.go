package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/shares"
)

func TestAccSFSV2Share_basic(t *testing.T) {
	var share shares.Share

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckSFS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSFSV2ShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSV2ShareConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSV2ShareExists("openstack_sharedfilesystem_share_v2.share_1", &share),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_v2.share_1", "name", "nfs_share"),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_v2.share_1", "description", "test share description"),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_v2.share_1", "share_proto", "NFS"),
				),
			},
			{
				Config: testAccSFSV2ShareConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSV2ShareExists("openstack_sharedfilesystem_share_v2.share_1", &share),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_v2.share_1", "name", "nfs_share_updated"),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_v2.share_1", "is_public", "false"),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_v2.share_1", "description", ""),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_v2.share_1", "share_proto", "NFS"),
				),
			},
			{
				Config: testAccSFSV2ShareConfigExtend,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSV2ShareExists("openstack_sharedfilesystem_share_v2.share_1", &share),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_v2.share_1", "name", "nfs_share_extended"),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_v2.share_1", "is_public", "false"),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_v2.share_1", "share_proto", "NFS"),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_v2.share_1", "size", "2"),
				),
			},
			/*			resource.TestStep{
						Config: testAccSFSV2ShareConfig_shrink,
						Check: resource.ComposeTestCheckFunc(
							testAccCheckSFSV2ShareExists("openstack_sharedfilesystem_share_v2.share_1", &share),
							resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_v2.share_1", "name", "nfs_share_shrunk"),
							resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_v2.share_1", "is_public", "false"),
							resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_v2.share_1", "share_proto", "NFS"),
							resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_v2.share_1", "size", "1"),
						),
					},*/
		},
	})
}

func TestAccSFSV2Share_update(t *testing.T) {
	var share shares.Share

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckSFS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSFSV2ShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSV2ShareConfigMetadataUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSV2ShareExists("openstack_sharedfilesystem_share_v2.share_1", &share),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_v2.share_1", "name", "nfs_share"),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_v2.share_1", "description", "test share description"),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_v2.share_1", "share_proto", "NFS"),
					testAccCheckSFSV2ShareMetadataEquals("key", "value", &share),
				),
			},
			{
				Config: testAccSFSV2ShareConfigMetadataUpdate1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSV2ShareExists("openstack_sharedfilesystem_share_v2.share_1", &share),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_v2.share_1", "name", "nfs_share"),
					testAccCheckSFSV2ShareMetadataEquals("key", "value", &share),
					testAccCheckSFSV2ShareMetadataEquals("new_key", "new_value", &share),
				),
			},
			{
				Config: testAccSFSV2ShareConfigMetadataUpdate2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSV2ShareExists("openstack_sharedfilesystem_share_v2.share_1", &share),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_v2.share_1", "name", "nfs_share"),
					testAccCheckSFSV2ShareMetadataAbsent("key", &share),
					testAccCheckSFSV2ShareMetadataEquals("new_key", "new_value", &share),
				),
			},
			{
				Config: testAccSFSV2ShareConfigMetadataUpdate3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSV2ShareExists("openstack_sharedfilesystem_share_v2.share_1", &share),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_v2.share_1", "name", "nfs_share"),
					testAccCheckSFSV2ShareMetadataAbsent("key", &share),
					testAccCheckSFSV2ShareMetadataAbsent("new_key", &share),
				),
			},
		},
	})
}

func TestAccSFSV2Share_admin(t *testing.T) {
	var share shares.Share

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckSFS(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSFSV2ShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSV2ShareAdminConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSV2ShareExists("openstack_sharedfilesystem_share_v2.share_1", &share),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_v2.share_1", "name", "nfs_share_admin"),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_v2.share_1", "description", "test share description"),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_v2.share_1", "share_proto", "NFS"),
				),
			},
			{
				Config: testAccSFSV2ShareAdminConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSV2ShareExists("openstack_sharedfilesystem_share_v2.share_1", &share),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_v2.share_1", "name", "nfs_share_admin_updated"),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_v2.share_1", "is_public", "true"),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_v2.share_1", "description", ""),
					resource.TestCheckResourceAttr("openstack_sharedfilesystem_share_v2.share_1", "share_proto", "NFS"),
				),
			},
		},
	})
}

func testAccCheckSFSV2ShareDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	sfsClient, err := config.SharedfilesystemV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_sharedfilesystem_securityservice_v2" {
			continue
		}

		_, err := shares.Get(sfsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Manila share still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckSFSV2ShareExists(n string, share *shares.Share) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		sfsClient, err := config.SharedfilesystemV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
		}

		found, err := shares.Get(sfsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Share not found")
		}

		*share = *found

		return nil
	}
}

func testAccCheckSFSV2ShareMetadataEquals(key string, value string, share *shares.Share) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		sfsClient, err := config.SharedfilesystemV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
		}

		metadatum, err := shares.GetMetadatum(sfsClient, share.ID, key).Extract()
		if err != nil {
			return err
		}

		if metadatum[key] != value {
			return fmt.Errorf("Metadata does not match. Expected %v but got %v", metadatum, value)
		}

		return nil
	}
}

func testAccCheckSFSV2ShareMetadataAbsent(key string, share *shares.Share) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		sfsClient, err := config.SharedfilesystemV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
		}

		_, err = shares.GetMetadatum(sfsClient, share.ID, key).Extract()
		if err == nil {
			return fmt.Errorf("Metadata %s key must not exist", key)
		}

		return nil
	}
}

const testAccSFSV2ShareConfigBasic = `
resource "openstack_sharedfilesystem_share_v2" "share_1" {
  name             = "nfs_share"
  description      = "test share description"
  share_proto      = "NFS"
  share_type       = "dhss_false"
  size             = 1
}
`

const testAccSFSV2ShareConfigUpdate = `
resource "openstack_sharedfilesystem_share_v2" "share_1" {
  name             = "nfs_share_updated"
  is_public        = false
  share_proto      = "NFS"
  share_type       = "dhss_false"
  size             = 1
}
`

const testAccSFSV2ShareConfigExtend = `
resource "openstack_sharedfilesystem_share_v2" "share_1" {
  name             = "nfs_share_extended"
  share_proto      = "NFS"
  share_type       = "dhss_false"
  size             = 2
}
`

//const testAccSFSV2ShareConfigShrink = `
//resource "openstack_sharedfilesystem_share_v2" "share_1" {
//  name             = "nfs_share_shrunk"
//  share_proto      = "NFS"
//  share_type       = "dhss_false"
//  size             = 1
//}
//`

const testAccSFSV2ShareConfigMetadataUpdate = `
resource "openstack_sharedfilesystem_share_v2" "share_1" {
  name             = "nfs_share"
  description      = "test share description"
  share_proto      = "NFS"
  share_type       = "dhss_false"
  size             = 1

  metadata = {
    key = "value"
  }
}
`

const testAccSFSV2ShareConfigMetadataUpdate1 = `
resource "openstack_sharedfilesystem_share_v2" "share_1" {
  name             = "nfs_share"
  description      = "test share description"
  share_proto      = "NFS"
  share_type       = "dhss_false"
  size             = 1

  metadata = {
    key = "value"
    new_key = "new_value"
  }
}
`

const testAccSFSV2ShareConfigMetadataUpdate2 = `
resource "openstack_sharedfilesystem_share_v2" "share_1" {
  name             = "nfs_share"
  description      = "test share description"
  share_proto      = "NFS"
  share_type       = "dhss_false"
  size             = 1

  metadata = {
    new_key = "new_value"
  }
}
`

const testAccSFSV2ShareConfigMetadataUpdate3 = `
resource "openstack_sharedfilesystem_share_v2" "share_1" {
  name             = "nfs_share"
  description      = "test share description"
  share_proto      = "NFS"
  share_type       = "dhss_false"
  size             = 1
}
`

const testAccSFSV2ShareAdminConfigBasic = `
resource "openstack_sharedfilesystem_share_v2" "share_1" {
  name             = "nfs_share_admin"
  description      = "test share description"
  share_proto      = "NFS"
  share_type       = "dhss_false"
  size             = 1
}
`

const testAccSFSV2ShareAdminConfigUpdate = `
resource "openstack_sharedfilesystem_share_v2" "share_1" {
  name             = "nfs_share_admin_updated"
  is_public        = true
  share_proto      = "NFS"
  share_type       = "dhss_false"
  size             = 1
}
`
