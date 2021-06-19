package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
)

func TestAccBlockStorageV3Volume_basic(t *testing.T) {
	var volume volumes.Volume

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockStorageV3VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV3VolumeBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV3VolumeExists("openstack_blockstorage_volume_v3.volume_1", &volume),
					testAccCheckBlockStorageV3VolumeMetadata(&volume, "foo", "bar"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_v3.volume_1", "name", "volume_1"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_v3.volume_1", "size", "1"),
				),
			},
			{
				Config: testAccBlockStorageV3VolumeUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV3VolumeExists("openstack_blockstorage_volume_v3.volume_1", &volume),
					testAccCheckBlockStorageV3VolumeMetadata(&volume, "foo", "bar"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_v3.volume_1", "name", "volume_1-updated"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_v3.volume_1", "size", "2"),
				),
			},
		},
	})
}

func TestAccBlockStorageV3Volume_online_resize(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreOnlineResize(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockStorageV3VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV3VolumeOnlineResize(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_v3.volume_1", "size", "1"),
				),
			},
			{
				Config: testAccBlockStorageV3VolumeOnlineResizeUpdate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_v3.volume_1", "size", "2"),
				),
			},
		},
	})
}

func TestAccBlockStorageV3Volume_image(t *testing.T) {
	var volume volumes.Volume

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockStorageV3VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV3VolumeImage(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV3VolumeExists("openstack_blockstorage_volume_v3.volume_1", &volume),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_v3.volume_1", "name", "volume_1"),
				),
			},
		},
	})
}

func TestAccBlockStorageV3Volume_image_multiattach(t *testing.T) {
	var volume volumes.Volume

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockStorageV3VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV3VolumeImageMultiattach(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV3VolumeExists("openstack_blockstorage_volume_v3.volume_1", &volume),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_v3.volume_1", "name", "volume_1"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_v3.volume_1", "multiattach", "true"),
				),
			},
		},
	})
}

func TestAccBlockStorageV3Volume_timeout(t *testing.T) {
	var volume volumes.Volume

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockStorageV3VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV3VolumeTimeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV3VolumeExists("openstack_blockstorage_volume_v3.volume_1", &volume),
				),
			},
		},
	})
}

func testAccCheckBlockStorageV3VolumeDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	blockStorageClient, err := config.BlockStorageV3Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_blockstorage_volume_v3" {
			continue
		}

		_, err := volumes.Get(blockStorageClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Volume still exists")
		}
	}

	return nil
}

func testAccCheckBlockStorageV3VolumeExists(n string, volume *volumes.Volume) resource.TestCheckFunc {
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

func testAccCheckBlockStorageV3VolumeMetadata(
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

const testAccBlockStorageV3VolumeBasic = `
resource "openstack_blockstorage_volume_v3" "volume_1" {
  name = "volume_1"
  description = "first test volume"
  metadata = {
    foo = "bar"
  }
  size = 1
}
`

func testAccBlockStorageV3VolumeOnlineResize() string {
	return fmt.Sprintf(`
resource "openstack_compute_instance_v2" "basic" {
  name            = "instance_1"
  flavor_name     = "%s"
  image_name      = "%s"
}

resource "openstack_blockstorage_volume_v3" "volume_1" {
  name = "volume_1"
  description = "test volume"
  size = 1
  enable_online_resize = true
}

resource "openstack_compute_volume_attach_v2" "va_1" {
  instance_id = "${openstack_compute_instance_v2.basic.id}"
  volume_id   = "${openstack_blockstorage_volume_v3.volume_1.id}"
}
`, osFlavorName, osImageID)
}

func testAccBlockStorageV3VolumeOnlineResizeUpdate() string {
	return fmt.Sprintf(`
resource "openstack_compute_instance_v2" "basic" {
  name            = "instance_1"
  flavor_name     = "%s"
  image_name      = "%s"
}

resource "openstack_blockstorage_volume_v3" "volume_1" {
  name = "volume_1"
  description = "test volume"
  size = 2
  enable_online_resize = true
}

resource "openstack_compute_volume_attach_v2" "va_1" {
  instance_id = "${openstack_compute_instance_v2.basic.id}"
  volume_id   = "${openstack_blockstorage_volume_v3.volume_1.id}"
}
`, osFlavorName, osImageID)
}

const testAccBlockStorageV3VolumeUpdate = `
resource "openstack_blockstorage_volume_v3" "volume_1" {
  name = "volume_1-updated"
  description = "first test volume"
  metadata = {
    foo = "bar"
  }
  size = 2
}
`

func testAccBlockStorageV3VolumeImage() string {
	return fmt.Sprintf(`
resource "openstack_blockstorage_volume_v3" "volume_1" {
  name = "volume_1"
  size = 5
  image_id = "%s"
}
`, osImageID)
}

func testAccBlockStorageV3VolumeImageMultiattach() string {
	return fmt.Sprintf(`
resource "openstack_blockstorage_volume_v3" "volume_1" {
  name = "volume_1"
  size = 5
  image_id = "%s"
  multiattach = true
}
`, osImageID)
}

const testAccBlockStorageV3VolumeTimeout = `
resource "openstack_blockstorage_volume_v3" "volume_1" {
  name = "volume_1"
  description = "first test volume"
  size = 1

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`
