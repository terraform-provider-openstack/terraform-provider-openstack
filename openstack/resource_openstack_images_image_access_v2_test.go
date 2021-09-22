package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/members"
)

func TestAccImagesImageAccessV2_basic(t *testing.T) {
	var member members.Member

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckImagesImageAccessV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccImagesImageAccessV2Basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageAccessV2Exists("openstack_images_image_access_v2.image_access_1", &member),
					resource.TestCheckResourceAttrPtr(
						"openstack_images_image_access_v2.image_access_1", "status", &member.Status),
					resource.TestCheckResourceAttr(
						"openstack_images_image_access_v2.image_access_1", "status", "pending"),
				),
			},
			{
				Config: testAccImagesImageAccessV2Update(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImagesImageAccessV2Exists("openstack_images_image_access_v2.image_access_1", &member),
					resource.TestCheckResourceAttrPtr(
						"openstack_images_image_access_v2.image_access_1", "status", &member.Status),
					resource.TestCheckResourceAttr(
						"openstack_images_image_access_v2.image_access_1", "status", "accepted"),
				),
			},
		},
	})
}

func testAccCheckImagesImageAccessV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	imageClient, err := config.ImageV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack Image: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_images_image_access_v2" {
			continue
		}

		imageID, memberID, err := resourceImagesImageAccessV2ParseID(rs.Primary.ID)
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

func testAccCheckImagesImageAccessV2Exists(n string, member *members.Member) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		imageClient, err := config.ImageV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack Image: %s", err)
		}

		imageID, memberID, err := resourceImagesImageAccessV2ParseID(rs.Primary.ID)
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

const testAccImagesImageAccessV2 = `
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

func testAccImagesImageAccessV2Basic() string {
	return fmt.Sprintf(`
%s

resource "openstack_images_image_access_v2" "image_access_1" {
  image_id  = "${openstack_images_image_v2.image_1.id}"
  member_id = "${data.openstack_identity_auth_scope_v3.scope.project_id}"
}
`, testAccImagesImageAccessV2)
}

func testAccImagesImageAccessV2Update() string {
	return fmt.Sprintf(`
%s

resource "openstack_images_image_access_v2" "image_access_1" {
  image_id  = "${openstack_images_image_v2.image_1.id}"
  member_id = "${data.openstack_identity_auth_scope_v3.scope.project_id}"
  status    = "accepted"
}
`, testAccImagesImageAccessV2)
}
