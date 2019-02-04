package openstack

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOpenStackAvailabilityZonesV2_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackAvailabilityZonesConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2AvailabilityZones("data.openstack_compute_availability_zones_v2.zones"),
				),
			},
		},
	})
}

func testAccCheckComputeV2AvailabilityZones(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find AZ resource: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("AZ resource ID not set")
		}
		v, ok := rs.Primary.Attributes["names.#"]
		if !ok {
			return fmt.Errorf("AZ name list missing")
		}
		qty, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		if qty < 1 {
			return fmt.Errorf("No AZs found, something is broken")
		}
		return nil
	}
}

const testAccOpenStackAvailabilityZonesConfig = `
data "openstack_compute_availability_zones_v2" "zones" {}
`
