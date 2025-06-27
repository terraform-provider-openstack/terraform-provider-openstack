package openstack

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccBlockStorageVolumeAttachV3_basic(t *testing.T) {
	var va volumes.Attachment

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockStorageVolumeAttachV3Destroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageVolumeAttachV3Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageVolumeAttachV3Exists(t.Context(), "openstack_blockstorage_volume_attach_v3.va_1", &va),
				),
			},
		},
	})
}

func TestAccBlockStorageVolumeAttachV3_timeout(t *testing.T) {
	var va volumes.Attachment

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockStorageVolumeAttachV3Destroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageVolumeAttachV3Timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageVolumeAttachV3Exists(t.Context(), "openstack_blockstorage_volume_attach_v3.va_1", &va),
				),
			},
		},
	})
}

func testAccCheckBlockStorageVolumeAttachV3Destroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		client, err := config.BlockStorageV3Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack block storage client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_blockstorage_volume_attach_v3" {
				continue
			}

			volumeID, attachmentID, err := parsePairedIDs(rs.Primary.ID, "openstack_blockstorage_volume_attach_v3")
			if err != nil {
				return err
			}

			volume, err := volumes.Get(ctx, client, volumeID).Extract()
			if err != nil {
				if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
					return nil
				}

				return err
			}

			for _, v := range volume.Attachments {
				if attachmentID == v.AttachmentID {
					return errors.New("Volume attachment still exists")
				}
			}
		}

		return nil
	}
}

func testAccCheckBlockStorageVolumeAttachV3Exists(ctx context.Context, n string, va *volumes.Attachment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		client, err := config.BlockStorageV3Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack block storage client: %w", err)
		}

		volumeID, attachmentID, err := parsePairedIDs(rs.Primary.ID, "openstack_blockstorage_volume_attach_v3")
		if err != nil {
			return err
		}

		volume, err := volumes.Get(ctx, client, volumeID).Extract()
		if err != nil {
			return err
		}

		var found bool

		for _, v := range volume.Attachments {
			if attachmentID == v.AttachmentID {
				found = true
				*va = v
			}
		}

		if !found {
			return errors.New("Volume Attachment not found")
		}

		return nil
	}
}

const testAccBlockStorageVolumeAttachV3Basic = `
resource "openstack_blockstorage_volume_v3" "volume_1" {
  name = "volume_1"
  size = 1
}

resource "openstack_blockstorage_volume_attach_v3" "va_1" {
  volume_id = openstack_blockstorage_volume_v3.volume_1.id
  device = "auto"

  host_name = "devstack"
  ip_address = "192.168.255.10"
  initiator = "iqn.1993-08.org.debian:01:e9861fb1859"
  os_type = "linux2"
  platform = "x86_64"
}
`

const testAccBlockStorageVolumeAttachV3Timeout = `
resource "openstack_blockstorage_volume_v3" "volume_1" {
  name = "volume_1"
  size = 1
}

resource "openstack_blockstorage_volume_attach_v3" "va_1" {
  volume_id = openstack_blockstorage_volume_v3.volume_1.id
  device = "auto"

  host_name = "devstack"
  ip_address = "192.168.255.10"
  initiator = "iqn.1993-08.org.debian:01:e9861fb1859"
  os_type = "linux2"
  platform = "x86_64"

  timeouts {
    create = "5m"
  }
}
`
