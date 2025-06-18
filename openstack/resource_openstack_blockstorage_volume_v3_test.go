package openstack

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/backups"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccBlockStorageV3Volume_basic(t *testing.T) {
	var volume volumes.Volume

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockStorageV3VolumeDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV3VolumeBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV3VolumeExists(t.Context(), "openstack_blockstorage_volume_v3.volume_1", &volume),
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
					testAccCheckBlockStorageV3VolumeExists(t.Context(), "openstack_blockstorage_volume_v3.volume_1", &volume),
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
		CheckDestroy:      testAccCheckBlockStorageV3VolumeDestroy(t.Context()),
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
		CheckDestroy:      testAccCheckBlockStorageV3VolumeDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV3VolumeImage(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV3VolumeExists(t.Context(), "openstack_blockstorage_volume_v3.volume_1", &volume),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_v3.volume_1", "name", "volume_1"),
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
		CheckDestroy:      testAccCheckBlockStorageV3VolumeDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV3VolumeTimeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV3VolumeExists(t.Context(), "openstack_blockstorage_volume_v3.volume_1", &volume),
				),
			},
		},
	})
}

func TestAccBlockStorageV3Volume_attachment(t *testing.T) {
	var volume volumes.Volume

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockStorageV3VolumeDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV3VolumeAttachment(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV3VolumeExists(t.Context(), "openstack_blockstorage_volume_v3.volume_1", &volume),
					testAccCheckBlockStorageV3VolumeAttachment(&volume, *regexp.MustCompile(`\/dev\/.dc`)),
				),
			},
		},
	})
}

func TestAccBlockStorageV3VolumeFromBackup(t *testing.T) {
	var volume volumes.Volume

	volumeName := acctest.RandomWithPrefix("tf-acc-volume")
	backupName := acctest.RandomWithPrefix("tf-acc-backup")

	var volumeID, backupID string

	if os.Getenv("TF_ACC") != "" {
		var err error

		volumeID, backupID, err = testAccBlockStorageV3CreateVolumeAndBackup(t.Context(), volumeName, backupName)
		if err != nil {
			t.Fatal(err)
		}

		defer testAccBlockStorageV3DeleteVolumeAndBackup(t, volumeID, backupID)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			t.Skip("Currently Cinder Backup is not configured properly on GH-A devstack")
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockStorageV3VolumeDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV3VolumeFromBackup(backupID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV3VolumeExists(t.Context(), "openstack_blockstorage_volume_v3.volume_1", &volume),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_v3.volume_1", "name", "volume_1"),
				),
			},
		},
	})
}

func testAccBlockStorageV3CreateVolumeAndBackup(ctx context.Context, volumeName, backupName string) (string, string, error) {
	config, err := testAccAuthFromEnv(ctx)
	if err != nil {
		return "", "", err
	}

	bsClient, err := config.BlockStorageV3Client(ctx, osRegionName)
	if err != nil {
		return "", "", err
	}

	volCreateOpts := volumes.CreateOpts{
		Size: 1,
		Name: volumeName,
	}

	volume, err := volumes.Create(ctx, bsClient, volCreateOpts, nil).Extract()
	if err != nil {
		return "", "", err
	}

	ctx1, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	err = volumes.WaitForStatus(ctx1, bsClient, volume.ID, "available")
	if err != nil {
		return "", "", err
	}

	snapCreateOpts := backups.CreateOpts{
		VolumeID: volume.ID,
		Name:     backupName,
	}

	backup, err := backups.Create(ctx, bsClient, snapCreateOpts).Extract()
	if err != nil {
		return volume.ID, "", err
	}

	ctx2, cancel1 := context.WithTimeout(ctx, 60*time.Second)
	defer cancel1()

	err = testAccBlockStorageV3BackupWaitForStatus(ctx2, bsClient, backup.ID, "available")
	if err != nil {
		return volume.ID, "", err
	}

	return volume.ID, backup.ID, nil
}

func testAccBlockStorageV3DeleteVolumeAndBackup(t *testing.T, volumeID, backupID string) {
	config, err := testAccAuthFromEnv(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	bsClient, err := config.BlockStorageV3Client(t.Context(), osRegionName)
	if err != nil {
		t.Fatal(err)
	}

	err = backups.Delete(t.Context(), bsClient, backupID).ExtractErr()
	if err != nil {
		t.Fatal(err)
	}

	ctx1, cancel := context.WithTimeout(t.Context(), 60*time.Second)
	defer cancel()

	err = testAccBlockStorageV3BackupWaitForStatus(ctx1, bsClient, backupID, "DELETED")
	if err != nil {
		if !gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
			t.Fatal(err)
		}
	}

	err = volumes.Delete(t.Context(), bsClient, volumeID, nil).ExtractErr()
	if err != nil {
		t.Fatal(err)
	}

	ctx2, cancel1 := context.WithTimeout(t.Context(), 60*time.Second)
	defer cancel1()

	err = volumes.WaitForStatus(ctx2, bsClient, volumeID, "DELETED")
	if err != nil {
		if !gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
			t.Fatal(err)
		}
	}
}

func testAccBlockStorageV3BackupWaitForStatus(ctx context.Context, c *gophercloud.ServiceClient, id, status string) error {
	return gophercloud.WaitFor(ctx, func(ctx context.Context) (bool, error) {
		current, err := backups.Get(ctx, c, id).Extract()
		if err != nil {
			return false, err
		}

		if current.Status == status {
			return true, nil
		}

		return false, nil
	})
}

func testAccCheckBlockStorageV3VolumeDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		blockStorageClient, err := config.BlockStorageV3Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack block storage client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_blockstorage_volume_v3" {
				continue
			}

			_, err := volumes.Get(ctx, blockStorageClient, rs.Primary.ID).Extract()
			if err == nil {
				return errors.New("Volume still exists")
			}
		}

		return nil
	}
}

func testAccCheckBlockStorageV3VolumeExists(ctx context.Context, n string, volume *volumes.Volume) resource.TestCheckFunc {
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

		found, err := volumes.Get(ctx, blockStorageClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return errors.New("Volume not found")
		}

		*volume = *found

		return nil
	}
}

func testAccCheckBlockStorageV3VolumeMetadata(
	volume *volumes.Volume, k string, v string,
) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if volume.Metadata == nil {
			return errors.New("No metadata")
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

func testAccCheckBlockStorageV3VolumeAttachment(
	volume *volumes.Volume, r regexp.Regexp,
) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if volume.Attachments == nil {
			return errors.New("No Attachment information")
		}

		if len(volume.Attachments) == 0 {
			return errors.New("Volume shows not being attached to any Instance")
		} else if len(volume.Attachments) > 1 {
			return errors.New("Volume shows being attached to more Instances than expected")
		}

		match := r.MatchString(volume.Attachments[0].Device)
		if match {
			return nil
		}

		return errors.New("Volume shows other mountpoint than expected")
	}
}

func TestAccBlockStorageV3Volume_VolumeTypeUpdate(t *testing.T) {
	var volume volumes.Volume

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockStorageV3VolumeDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV3VolumeRetype(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV3VolumeExists(t.Context(), "openstack_blockstorage_volume_v3.volume_1", &volume),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_v3.volume_1", "volume_type", "initial_type"),
				),
			},
			{
				Config: testAccBlockStorageV3VolumeRetypeUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV3VolumeExists(t.Context(), "openstack_blockstorage_volume_v3.volume_1", &volume),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_v3.volume_1", "volume_type", "new_type"),
				),
			},
		},
	})
}

const testAccBlockStorageV3VolumeBasic = `
resource "openstack_blockstorage_volume_v3" "volume_1" {
  name = "volume_1"
  description = "first test volume"
  metadata = {
    foo = "bar"
  }
  size = 1
  volume_retype_policy = "never"
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

func testAccBlockStorageV3VolumeAttachment() string {
	return fmt.Sprintf(`
resource "openstack_blockstorage_volume_v3" "volume_1" {
  name = "volume_1"
  size = 1
}

resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
}

resource "openstack_compute_volume_attach_v2" "va_1" {
  instance_id = "${openstack_compute_instance_v2.instance_1.id}"
  volume_id = "${openstack_blockstorage_volume_v3.volume_1.id}"
  device = "/dev/vdc"
}
`, osNetworkID)
}

func testAccBlockStorageV3VolumeFromBackup(backupID string) string {
	return fmt.Sprintf(`
resource "openstack_blockstorage_volume_v3" "volume_1" {
  name = "volume_1"
  backup_id = "%s"
  size = 2
}
`, backupID)
}

func testAccBlockStorageV3VolumeRetype() string {
	return `
resource "openstack_blockstorage_volume_type_v3" "initial_type" {
  name        = "initial_type"
  description = "initial_type"
  is_public   = true
}

resource "openstack_blockstorage_volume_type_v3" "new_type" {
  name        = "new_type"
  description = "new_type"
  is_public   = true
}

resource "openstack_blockstorage_volume_v3" "volume_1" {
  name                 = "volume_1"
  size                 = 1
  volume_retype_policy = "on-demand"
  volume_type          = openstack_blockstorage_volume_type_v3.initial_type.name
}`
}

func testAccBlockStorageV3VolumeRetypeUpdate() string {
	return `
resource "openstack_blockstorage_volume_type_v3" "initial_type" {
  name        = "initial_type"
  description = "initial_type"
  is_public   = true
}

resource "openstack_blockstorage_volume_type_v3" "new_type" {
  name        = "new_type"
  description = "new_type"
  is_public   = true
}

resource "openstack_blockstorage_volume_v3" "volume_1" {
  name                 = "volume_1"
  size                 = 1
  volume_retype_policy = "on-demand"
  volume_type          = openstack_blockstorage_volume_type_v3.new_type.name
}`
}
