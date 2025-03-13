package openstack

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
)

func TestAccBlockStorageV3VolumeDataSource_basic(t *testing.T) {
	resourceName := "data.openstack_blockstorage_volume_v3.volume_1"
	volumeName := acctest.RandomWithPrefix("tf-acc-volume")

	var volumeID string
	if os.Getenv("TF_ACC") != "" {
		var err error
		volumeID, err = testAccBlockStorageV3CreateVolume(volumeName)
		if err != nil {
			t.Fatal(err)
		}
		defer testAccBlockStorageV3DeleteVolume(t, volumeID)
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV3VolumeDataSourceBasic(volumeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV3VolumeDataSourceID(resourceName, volumeID),
					resource.TestCheckResourceAttr(resourceName, "name", volumeName),
					resource.TestCheckResourceAttr(resourceName, "size", "1"),
				),
			},
		},
	})
}

func testAccBlockStorageV3CreateVolume(volumeName string) (string, error) {
	config, err := testAccAuthFromEnv()
	if err != nil {
		return "", err
	}

	bsClient, err := config.BlockStorageV3Client(context.TODO(), osRegionName)
	if err != nil {
		return "", err
	}

	volCreateOpts := volumes.CreateOpts{
		Size: 1,
		Name: volumeName,
	}

	volume, err := volumes.Create(context.TODO(), bsClient, volCreateOpts, nil).Extract()
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()
	err = volumes.WaitForStatus(ctx, bsClient, volume.ID, "available")
	if err != nil {
		return "", err
	}

	return volume.ID, nil
}

func testAccBlockStorageV3DeleteVolume(t *testing.T, volumeID string) {
	config, err := testAccAuthFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	bsClient, err := config.BlockStorageV3Client(context.TODO(), osRegionName)
	if err != nil {
		t.Fatal(err)
	}

	err = volumes.Delete(context.TODO(), bsClient, volumeID, nil).ExtractErr()
	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 60*time.Second)
	defer cancel()
	err = volumes.WaitForStatus(ctx, bsClient, volumeID, "DELETED")
	if err != nil {
		if !gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
			t.Fatal(err)
		}
	}
}

func testAccCheckBlockStorageV3VolumeDataSourceID(n, id string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find volume data source: %s", n)
		}

		if rs.Primary.ID != id {
			return fmt.Errorf("Volume data source ID not set")
		}

		return nil
	}
}

func testAccBlockStorageV3VolumeDataSourceBasic(snapshotName string) string {
	return fmt.Sprintf(`
    data "openstack_blockstorage_volume_v3" "volume_1" {
      name = "%s"
    }
  `, snapshotName)
}

func TestAccBlockStorageV3VolumeDataSource_attachment(t *testing.T) {
	var dataVolume volumes.Volume

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockStorageV3VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV3VolumeDataSourceAttachment(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV3VolumeExists("data.openstack_blockstorage_volume_v3.volume_1", &dataVolume),
					resource.TestCheckResourceAttrPair("data.openstack_blockstorage_volume_v3.volume_1", "attachment", "openstack_blockstorage_volume_v3.volume_1", "attachment"),
					testAccCheckBlockStorageV3VolumeAttachment(&dataVolume, *regexp.MustCompile(`\/dev\/.dc`)),
				),
			},
		},
	})
}

func testAccBlockStorageV3VolumeDataSourceAttachment() string {
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

data "openstack_blockstorage_volume_v3" "volume_1" {
  name = "volume_1"
  depends_on = [ openstack_compute_volume_attach_v2.va_1 ]
}
`, osNetworkID)
}
