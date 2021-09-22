package openstack

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
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
					resource.TestCheckResourceAttr(resourceName, "multiattach", "false"),
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

	bsClient, err := config.BlockStorageV3Client(osRegionName)
	if err != nil {
		return "", err
	}

	volCreateOpts := volumes.CreateOpts{
		Size: 1,
		Name: volumeName,
	}

	volume, err := volumes.Create(bsClient, volCreateOpts).Extract()
	if err != nil {
		return "", err
	}

	err = volumes.WaitForStatus(bsClient, volume.ID, "available", 60)
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

	bsClient, err := config.BlockStorageV3Client(osRegionName)
	if err != nil {
		t.Fatal(err)
	}

	err = volumes.Delete(bsClient, volumeID, nil).ExtractErr()
	if err != nil {
		t.Fatal(err)
	}

	err = volumes.WaitForStatus(bsClient, volumeID, "DELETED", 60)
	if err != nil {
		if _, ok := err.(gophercloud.ErrDefault404); !ok {
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
