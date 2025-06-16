package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBlockStorageV3Volume_importBasic(t *testing.T) {
	resourceName := "openstack_blockstorage_volume_v3.volume_1"

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
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccBlockStorageV3Volume_importImage(t *testing.T) {
	resourceName := "openstack_blockstorage_volume_v3.volume_1"

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
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
