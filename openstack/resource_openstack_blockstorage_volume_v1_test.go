package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v1/volumes"
)

func TestAccBlockStorageV1Volume_basic(t *testing.T) {
	var volume volumes.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckDeprecated(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockStorageV1VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV1VolumeBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV1VolumeExists("openstack_blockstorage_volume_v1.volume_1", &volume),
					testAccCheckBlockStorageV1VolumeMetadata(&volume, "foo", "bar"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_v1.volume_1", "name", "volume_1"),
				),
			},
			{
				Config: testAccBlockStorageV1VolumeUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV1VolumeExists("openstack_blockstorage_volume_v1.volume_1", &volume),
					testAccCheckBlockStorageV1VolumeMetadata(&volume, "foo", "bar"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_v1.volume_1", "name", "volume_1-updated"),
				),
			},
		},
	})
}

func TestAccBlockStorageV1Volume_image(t *testing.T) {
	var volume volumes.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckDeprecated(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockStorageV1VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV1VolumeImage(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV1VolumeExists("openstack_blockstorage_volume_v1.volume_1", &volume),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_v1.volume_1", "name", "volume_1"),
				),
			},
		},
	})
}

func TestAccBlockStorageV1Volume_timeout(t *testing.T) {
	var volume volumes.Volume

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckDeprecated(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockStorageV1VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV1VolumeTimeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV1VolumeExists("openstack_blockstorage_volume_v1.volume_1", &volume),
				),
			},
		},
	})
}

func testAccCheckBlockStorageV1VolumeDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	blockStorageClient, err := config.BlockStorageV1Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_blockstorage_volume_v1" {
			continue
		}

		_, err := volumes.Get(blockStorageClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Volume still exists")
		}
	}

	return nil
}

func testAccCheckBlockStorageV1VolumeExists(n string, volume *volumes.Volume) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		blockStorageClient, err := config.BlockStorageV1Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
		}

		found, err := volumes.Get(blockStorageClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Volume not found")
		}

		*volume = *found

		return nil
	}
}

func testAccCheckBlockStorageV1VolumeMetadata(
	volume *volumes.Volume, k string, v string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if volume.Metadata == nil {
			return fmt.Errorf("No metadata")
		}

		for key, value := range volume.Metadata {
			if k != key {
				continue
			}

			if v == value {
				return nil
			}

			return fmt.Errorf("Bad value for %s: %s", k, value)
		}

		return fmt.Errorf("Metadata not found: %s", k)
	}
}

const testAccBlockStorageV1VolumeBasic = `
resource "openstack_blockstorage_volume_v1" "volume_1" {
  name = "volume_1"
  description = "first test volume"
  availability_zone = "nova"
  metadata = {
    foo = "bar"
  }
  size = 1
}
`

const testAccBlockStorageV1VolumeUpdate = `
resource "openstack_blockstorage_volume_v1" "volume_1" {
  name = "volume_1-updated"
  description = "first test volume"
  metadata = {
    foo = "bar"
  }
  size = 1
}
`

func testAccBlockStorageV1VolumeImage() string {
	return fmt.Sprintf(`
resource "openstack_blockstorage_volume_v1" "volume_1" {
  name = "volume_1"
  size = 5
  image_id = "%s"
}
`, osImageID)
}

const testAccBlockStorageV1VolumeTimeout = `
resource "openstack_blockstorage_volume_v1" "volume_1" {
  name = "volume_1"
  description = "first test volume"
  size = 1

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`
