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
						"data.openstack_blockstorage_volume_type_v3.vt_ds_1", "name", "volume_type_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_blockstorage_volume_type_v3.vt_ds_1", "description", "volume_type_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_blockstorage_volume_type_v3.vt_ds_1", "is_public", "false"),
					resource.TestCheckResourceAttr(
						"data.openstack_blockstorage_volume_type_v3.vt_ds_1", "extra_specs.spec1", "foo"),
					resource.TestCheckResourceAttr(
						"data.openstack_blockstorage_volume_type_v3.vt_ds_2", "name", "volume_type_2"),
					resource.TestCheckResourceAttr(
						"data.openstack_blockstorage_volume_type_v3.vt_ds_2", "description", "volume_type_2"),
					resource.TestCheckResourceAttr(
						"data.openstack_blockstorage_volume_type_v3.vt_ds_2", "is_public", "true"),
					resource.TestCheckResourceAttr(
						"data.openstack_blockstorage_volume_type_v3.vt_ds_2", "extra_specs.spec3", "foo"),
				),
			},
		},
	})
}

const testAccBlockStorageV3VolumeTypeDataSource = `
resource "openstack_blockstorage_volume_type_v3" "volume_type_1" {
  name = "volume_type_1"
  description = "volume_type_1"
  is_public = false
  extra_specs = {
    spec1 = "foo"
    spec2 = "bar"
  }

}

resource "openstack_blockstorage_volume_type_v3" "volume_type_2" {
  name = "volume_type_2"
  description = "volume_type_2"
  is_public = true
  extra_specs = {
    spec3 = "foo"
    spec4 = "bar"
  }

}

data "openstack_blockstorage_volume_type_v3" "vt_ds_1" {
  name = openstack_blockstorage_volume_type_v3.volume_type_1.name
  is_public = false
  extra_specs = {
    spec1 = "foo"
    spec2 = "bar"
  }
}

data "openstack_blockstorage_volume_type_v3" "vt_ds_2" {
  is_public = true
  extra_specs = {
    spec4 = "bar"
  }
}
`
