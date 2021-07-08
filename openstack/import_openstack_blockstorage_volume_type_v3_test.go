package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBlockStorageV3VolumeType_importBasic(t *testing.T) {
	resourceName := "openstack_blockstorage_volume_type_v3.volume_type_1"

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
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
