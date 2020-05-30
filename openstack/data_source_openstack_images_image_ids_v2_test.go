package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOpenStackImagesV2ImageIDsDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackImagesV2ImageIDsDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.openstack_images_image_ids_v2.images_empty", "ids.#", "0"),
					resource.TestCheckResourceAttr(
						"data.openstack_images_image_ids_v2.images_by_name", "ids.#", "1"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_images_image_ids_v2.images_by_name", "ids.0",
                                                "openstack_images_image_v2.image_1", "id"),
					resource.TestCheckResourceAttr(
						"data.openstack_images_image_ids_v2.images_by_name_regex", "ids.#", "1"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_images_image_ids_v2.images_by_name_regex", "ids.0",
                                                "openstack_images_image_v2.image_2", "id"),
				),
				),
			},
		},
	})
}

// Standard CirrOS image
const testAccOpenStackImagesV2ImageIDsDataSource_basic = `
resource "openstack_images_image_v2" "image_1" {
  name = "CirrOS-tf_1"
  container_format = "bare"
  disk_format = "qcow2"
  image_source_url = "http://download.cirros-cloud.net/0.3.5/cirros-0.3.5-x86_64-disk.img"
  tags = ["cirros-tf_1"]
  properties = {
    foo = "bar"
    bar = "foo"
  }
  visibility = "private"
}

resource "openstack_images_image_v2" "image_2" {
  name = "CirrOS-tf_2"
  container_format = "bare"
  disk_format = "qcow2"
  image_source_url = "http://download.cirros-cloud.net/0.5.1/cirros-0.5.1-x86_64-disk.img"
  tags = ["cirros-tf_2"]
  properties = {
    foo = "bar"
  }
  visibility = "private"
}

data "openstack_images_image_ids_v2" "images_empty" {
        name = "non-existed-image"
}

data "openstack_images_image_ids_v2" "images_by_name" {
	name = "${openstack_images_image_v2.image_1.name}"
}

data "openstack_images_image_ids_v2" "images_by_name_regex" {
	name = "CirrOS-.+2"
}

data "openstack_images_image_ids_v2" "images_by_tag" {
	visibility = "private"
	tag = "cirros-tf_1"
}

data "openstack_images_image_ids_v2" "images_by_size_min" {
	visibility = "private"
	size_min = "15000000"
}

data "openstack_images_image_ids_v2" "images_by_size_max" {
	visibility = "private"
	size_max = "15000000"
}

data "openstack_images_image_ids_v2" "images_by_properties" {
  properties = {
    foo = "bar"
    bar = "foo"
  }
}
`
