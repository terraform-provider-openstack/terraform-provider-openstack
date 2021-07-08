package openstack

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/aggregates"
)

var testAccAggregateConfig = `
resource "openstack_compute_aggregate_v2" "test" {
  name  = "test-aggregate"
  zone  = "nova"
}
`

func testAccAggregateHypervisorConfig() string {
	return fmt.Sprintf(`
resource "openstack_compute_aggregate_v2" "test" {
  name = "test-aggregate"
  zone = "nova"
  hosts = [ "%s" ]
  metadata = {
    test = "123"
  }
}
    `, osHypervisorEnvironment)
}

func TestAccComputeV2Aggregate(t *testing.T) {
	var aggregate aggregates.Aggregate

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckAdminOnly(t) },
		ProviderFactories: testAccProviders,
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

func TestAccComputeV2AggregateWithHypervisor(t *testing.T) {
	var aggregate aggregates.Aggregate

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckAdminOnly(t)
			testAccPreCheckHypervisor(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAggregateConfig,
				Check:  testAccCheckAggregateExists("openstack_compute_aggregate_v2.test", &aggregate),
			},
			{
				Config: testAccAggregateHypervisorConfig(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAggregateExists("openstack_compute_aggregate_v2.test", &aggregate),
					resource.TestCheckResourceAttr("openstack_compute_aggregate_v2.test", "hosts.#", "1"),
					resource.TestCheckResourceAttr("openstack_compute_aggregate_v2.test", "metadata.test", "123"),
				),
			},
			{
				Config: testAccAggregateConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAggregateExists("openstack_compute_aggregate_v2.test", &aggregate),
					resource.TestCheckResourceAttr("openstack_compute_aggregate_v2.test", "hosts.#", "0"),
					resource.TestCheckNoResourceAttr("openstack_compute_aggregate_v2.test", "metadata.test"),
				),
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
