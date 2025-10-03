package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBlockStorageV3VolumeTypeDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV3VolumeTypeDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.openstack_blockstorage_volume_type_v3.vt_ds", "name", "bar-baz"),
					resource.TestCheckResourceAttr(
						"data.openstack_blockstorage_volume_type_v3.vt_ds", "description", "bar-baz"),
					resource.TestCheckResourceAttr(
						"data.openstack_blockstorage_volume_type_v3.vt_ds", "is_public", "false"),
					resource.TestCheckResourceAttr(
						"data.openstack_blockstorage_volume_type_v3.vt_ds", "extra_specs.#", "2"),
				),
			},
		},
	})
}

const testAccBlockStorageV3VolumeTypeDataSource = `
resource "openstack_blockstorage_volume_type_v3" "volume_type_1" {
  name = "bar-baz"
  description = "bar-baz"
  is_public = false
  extra_specs = {
    bar = "bar"
    baz = "baz"
  }

}

data "openstack_blockstorage_volume_type_v3" "vt_ds" {
  name = openstack_blockstorage_volume_type_v3.volume_type_1.name
  is_public = false
}
`
