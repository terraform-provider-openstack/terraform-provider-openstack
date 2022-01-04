package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

const testAccAggregateResource = `
resource "openstack_compute_aggregate_v2" "test2" {
  name = "test2"
  zone = "nova"
  metadata = {
    "test" = "test123"
  }
}
`

const testAccAggregateDataSource = `
data "openstack_compute_aggregate_v2" "test2" {
  name = openstack_compute_aggregate_v2.test2.name
}

resource "openstack_compute_aggregate_v2" "test2" {
  name = "test2"
  zone = "nova"
  metadata = {
    "test" = "test123"
  }
}
`

func TestAccAggregateDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckAdminOnly(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAggregateResource,
			},
			{
				Config: testAccAggregateDataSource,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeAggregateV2DataSourceID("data.openstack_compute_aggregate_v2.test2"),
					resource.TestCheckResourceAttr("data.openstack_compute_aggregate_v2.test2", "name", "test2"),
					resource.TestCheckResourceAttr("data.openstack_compute_aggregate_v2.test2", "zone", "nova"),
					resource.TestCheckResourceAttr("data.openstack_compute_aggregate_v2.test2", "metadata.test", "test123"),
				),
			},
		},
	})
}

func testAccAggregateResourceWithHypervisor() string {
	return fmt.Sprintf(`
resource "openstack_compute_aggregate_v2" "test3" {
  name = "test3"
  zone = "nova"
  hosts = ["%s"]
}
    `, osHypervisorEnvironment)
}

func testAccAggregateDataSourceWithHypervisor() string {
	return fmt.Sprintf(`
data "openstack_compute_aggregate_v2" "test3" {
  name = openstack_compute_aggregate_v2.test3.name
}

resource "openstack_compute_aggregate_v2" "test3" {
  name = "test3"
  zone = "nova"
  hosts = ["%s"]
}
    `, osHypervisorEnvironment)
}

func TestAccAggregateDataSourceWithHypervisor(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckAdminOnly(t)
			testAccPreCheckHypervisor(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccAggregateResourceWithHypervisor(),
			},
			{
				Config: testAccAggregateDataSourceWithHypervisor(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeAggregateV2DataSourceID("data.openstack_compute_aggregate_v2.test3"),
					resource.TestCheckResourceAttr("data.openstack_compute_aggregate_v2.test3", "name", "test3"),
					resource.TestCheckResourceAttr("data.openstack_compute_aggregate_v2.test3", "zone", "nova"),
					resource.TestCheckResourceAttr("data.openstack_compute_aggregate_v2.test3", "hosts.#", "1"),
				),
			},
		},
	})
}

func testAccCheckComputeAggregateV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Data source ID not set")
		}

		return nil
	}
}
