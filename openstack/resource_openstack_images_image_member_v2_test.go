package openstack

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/members"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccImagesImageMemberV2_basic(t *testing.T) {
	var member members.Member

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckImagesImageMemberV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageMemberV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageMemberV2Exists("openstack_images_image_member_v2.image_member_1", &member),
					resource.TestCheckResourceAttrPtr(
						"openstack_images_image_member_v2.image_member_1", "status", &member.Status),
				),
			},
			{
				Config: testAccImagesImageMemberV2_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageMemberV2Exists("openstack_images_image_member_v2.image_member_1", &member),
					resource.TestCheckResourceAttrPtr(
						"openstack_images_image_member_v2.image_member_1", "status", &member.Status),
				),
			},
		},
	})
}

func testAccCheckImagesImageMemberV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	imageClient, err := config.imageV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack Image: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_images_image_member_v2" {
			continue
		}

		imageID, memberID, err := resourceImagesShareV2ParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = members.Get(imageClient, imageID, memberID).Extract()
		if err == nil {
			return fmt.Errorf("Image still exists")
		}
	}

	return nil
}

func testAccCheckImagesImageMemberV2Exists(n string, member *members.Member) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		imageClient, err := config.imageV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack Image: %s", err)
		}

		imageID, memberID, err := resourceImagesShareV2ParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		found, err := members.Get(imageClient, imageID, memberID).Extract()
		if err != nil {
			return err
		}

		id := fmt.Sprintf("%s/%s", found.ImageID, found.MemberID)
		if id != rs.Primary.ID {
			return fmt.Errorf("Image member not found")
		}

		*member = *found

		return nil
	}
}

const testAccImagesImageMemberV2 = `
data "openstack_identity_auth_scope_v3" "scope" {
  name = "scope"
}

resource "openstack_images_image_v2" "image_1" {
  name   = "CirrOS-tf_1"
  image_source_url = "http://download.cirros-cloud.net/0.3.5/cirros-0.3.5-x86_64-disk.img"
  container_format = "bare"
  disk_format = "qcow2"
  visibility = "shared"

  timeouts {
    create = "10m"
  }
}`

var testAccImagesImageMemberV2_basic = fmt.Sprintf(`
%s

resource "openstack_images_image_member_v2" "image_member_1" {
  image_id  = "${openstack_images_image_v2.image_1.id}"
  member_id = "${data.openstack_identity_auth_scope_v3.scope.project_id}"
}
`, testAccImagesImageMemberV2)

var testAccImagesImageMemberV2_update = fmt.Sprintf(`
%s

resource "openstack_images_image_member_v2" "image_member_1" {
  image_id  = "${openstack_images_image_v2.image_1.id}"
  member_id = "${data.openstack_identity_auth_scope_v3.scope.project_id}"
  status    = "accepted"
}
`, testAccImagesImageMemberV2)
