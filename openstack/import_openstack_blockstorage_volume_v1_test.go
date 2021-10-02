package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBlockStorageV1Volume_importBasic(t *testing.T) {
	resourceName := "openstack_blockstorage_volume_v1.volume_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckDeprecated(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockStorageV1VolumeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV1VolumeBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
