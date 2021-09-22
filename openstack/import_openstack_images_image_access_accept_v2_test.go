package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccImagesImageAccessAcceptV2_importBasic(t *testing.T) {
	resourceName := "openstack_images_image_access_accept_v2.image_access_accept_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckImagesImageAccessAcceptV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageAccessAcceptV2Basic(),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
