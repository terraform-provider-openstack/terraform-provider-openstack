package openstack

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBlockStorageV3AvailabilityZonesV3_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV3AvailabilityZonesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.openstack_blockstorage_availability_zones_v3.zones", "names.#", regexp.MustCompile("[1-9]\\d*")),
				),
			},
		},
	})
}

const testAccBlockStorageV3AvailabilityZonesConfig = `
data "openstack_blockstorage_availability_zones_v3" "zones" {}
`
