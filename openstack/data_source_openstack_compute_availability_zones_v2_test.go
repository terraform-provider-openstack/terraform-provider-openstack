package openstack

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccComputeV2AvailabilityZonesV2_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackAvailabilityZonesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.openstack_compute_availability_zones_v2.zones", "names.#", regexp.MustCompile(`[1-9]\d*`)),
				),
			},
		},
	})
}

const testAccOpenStackAvailabilityZonesConfig = `
data "openstack_compute_availability_zones_v2" "zones" {}
`
