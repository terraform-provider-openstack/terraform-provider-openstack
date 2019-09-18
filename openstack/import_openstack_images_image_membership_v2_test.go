package openstack

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccImagesImageMembershipV2_importBasic(t *testing.T) {
	resourceName := "openstack_images_image_membership_v2.image_membership_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckImagesImageMembershipV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageMembershipV2_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
