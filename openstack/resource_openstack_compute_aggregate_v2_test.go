package openstack

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/aggregates"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"strconv"
)

var testAccAggregateConfig = `
resource "openstack_compute_aggregate_v2" "test" {
  name  = "test-aggregate"
  zone  = "nova"
}
`

func TestAccComputeV2Aggregate(t *testing.T) {
	var aggregate aggregates.Aggregate

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAggregateConfig,
				Check:  testAccCheckAggregateExists("openstack_compute_aggregate_v2.test", &aggregate),
			},
			{
				ResourceName: "openstack_compute_aggregate_v2.test",
				ImportState:  true,
			},
		},
	})
}

func testAccCheckAggregateExists(n string, aggregate *aggregates.Aggregate) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		computeClient, err := config.ComputeV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack compute client: %s", err)
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Can't convert ID to integer: %s", err)
		}

		found, err := aggregates.Get(computeClient, id).Extract()
		if err != nil {
			return err
		}

		if found.ID != id {
			return fmt.Errorf("Aggregate not found")
		}

		*aggregate = *found

		return nil
	}
}
