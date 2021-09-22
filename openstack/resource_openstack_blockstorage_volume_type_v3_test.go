package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumetypes"
)

func TestAccBlockStorageVolumeTypeV3_basic(t *testing.T) {
	var volumetype volumetypes.VolumeType

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockStorageVolumeTypeV3Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageVolumeTypeV3Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageVolumeTypeV3Exists("openstack_blockstorage_volume_type_v3.volume_type_1", &volumetype),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "name", "foo"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "description", "foo"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "is_public", "true"),
				),
			},
			{
				Config: testAccBlockStorageVolumeTypeV3Update1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageVolumeTypeV3Exists("openstack_blockstorage_volume_type_v3.volume_type_1", &volumetype),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "name", "bar-baz"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "description", "bar-baz"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "is_public", "false"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "extra_specs.%", "2"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "extra_specs.bar", "bar"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "extra_specs.baz", "baz"),
				),
			},
			{
				Config: testAccBlockStorageVolumeTypeV3Update2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageVolumeTypeV3Exists("openstack_blockstorage_volume_type_v3.volume_type_1", &volumetype),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "name", "foo-foo"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "description", "bar-bar"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "is_public", "false"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "extra_specs.%", "2"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "extra_specs.bar", "baz"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "extra_specs.foo", "foo"),
				),
			},
		},
	})
}

func testAccCheckBlockStorageVolumeTypeV3Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	blockStorageClient, err := config.BlockStorageV3Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_blockstorage_volume_type_v3" {
			continue
		}

		_, err := volumetypes.Get(blockStorageClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("VolumeType still exists")
		}
	}

	return nil
}

func testAccCheckBlockStorageVolumeTypeV3Exists(n string, volumetype *volumetypes.VolumeType) resource.TestCheckFunc {
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

		found, err := volumetypes.Get(blockStorageClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("VolumeType not found")
		}

		*volumetype = *found

		return nil
	}
}

const testAccBlockStorageVolumeTypeV3Basic = `
resource "openstack_blockstorage_volume_type_v3" "volume_type_1" {
	name = "foo"
	description = "foo"
	is_public = true

}
`

const testAccBlockStorageVolumeTypeV3Update1 = `
resource "openstack_blockstorage_volume_type_v3" "volume_type_1" {
	name = "bar-baz"
	description = "bar-baz"
	is_public = false
	extra_specs = {
	  bar = "bar"
	  baz = "baz"
	}

}
`

const testAccBlockStorageVolumeTypeV3Update2 = `
resource "openstack_blockstorage_volume_type_v3" "volume_type_1" {
	name = "foo-foo"
	description = "bar-bar"
	is_public = false
	extra_specs = {
      bar = "baz"
	  foo = "foo"
	}
}
`
