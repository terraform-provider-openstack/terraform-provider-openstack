package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccImagesImageV2_importBasic(t *testing.T) {
	resourceName := "openstack_images_image_v2.image_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckImagesImageV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageV2Basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"region",
					"local_file_path",
					"image_cache_path",
					"image_source_url",
					"verify_checksum",
				},
			},
		},
	})
}
