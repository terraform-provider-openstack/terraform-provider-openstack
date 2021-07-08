package openstack

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSFSV2AvailabilityZonesV2_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSV2AvailabilityZonesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.openstack_sharedfilesystem_availability_zones_v2.zones", "names.#", regexp.MustCompile("[1-9]\\d*")),
				),
			},
		},
	})
}

const testAccSFSV2AvailabilityZonesConfig = `
data "openstack_sharedfilesystem_availability_zones_v2" "zones" {}
`
