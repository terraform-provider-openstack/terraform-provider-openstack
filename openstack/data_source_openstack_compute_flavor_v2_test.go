package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOpenStackComputeV2FlavorDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOpenStackComputeV2FlavorDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2FlavorDataSourceID("data.openstack_compute_flavor_v2.flavor_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "name", "m1.acctest"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "ram", "512"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "disk", "5"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "vcpus", "1"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "rx_tx_factor", "1"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "is_public", "true"),
				),
			},
		},
	})
}

func TestAccOpenStackComputeV2FlavorDataSource_testQueries(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOpenStackComputeV2FlavorDataSource_queryDisk,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2FlavorDataSourceID("data.openstack_compute_flavor_v2.flavor_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "name", "m1.resize"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "ram", "512"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "disk", "6"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "vcpus", "1"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "rx_tx_factor", "1"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "is_public", "true"),
				),
			},
			resource.TestStep{
				Config: testAccOpenStackComputeV2FlavorDataSource_queryMinDisk,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2FlavorDataSourceID("data.openstack_compute_flavor_v2.flavor_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "name", "m1.acctest"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "ram", "512"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "disk", "5"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "vcpus", "1"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "rx_tx_factor", "1"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "is_public", "true"),
				),
			},
			resource.TestStep{
				Config: testAccOpenStackComputeV2FlavorDataSource_queryMinRAM,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2FlavorDataSourceID("data.openstack_compute_flavor_v2.flavor_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "name", "m1.acctest"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "ram", "512"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "disk", "5"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "vcpus", "1"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "rx_tx_factor", "1"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "is_public", "true"),
				),
			},
			resource.TestStep{
				Config: testAccOpenStackComputeV2FlavorDataSource_queryVCPUs,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2FlavorDataSourceID("data.openstack_compute_flavor_v2.flavor_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "name", "m1.acctest"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "ram", "512"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "disk", "5"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "vcpus", "1"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "rx_tx_factor", "1"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "is_public", "true"),
				),
			},
		},
	})
}

func testAccCheckComputeV2FlavorDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find flavor data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Flavor data source ID not set")
		}

		return nil
	}
}

const testAccOpenStackComputeV2FlavorDataSource_basic = `
data "openstack_compute_flavor_v2" "flavor_1" {
  name = "m1.acctest"
}
`

const testAccOpenStackComputeV2FlavorDataSource_queryDisk = `
data "openstack_compute_flavor_v2" "flavor_1" {
  disk = 6
}
`

const testAccOpenStackComputeV2FlavorDataSource_queryMinDisk = `
data "openstack_compute_flavor_v2" "flavor_1" {
  name = "m1.acctest"
  min_disk = 5
}
`

const testAccOpenStackComputeV2FlavorDataSource_queryMinRAM = `
data "openstack_compute_flavor_v2" "flavor_1" {
  name = "m1.acctest"
  min_ram = 512
}
`

const testAccOpenStackComputeV2FlavorDataSource_queryVCPUs = `
data "openstack_compute_flavor_v2" "flavor_1" {
  name = "m1.acctest"
  vcpus = 1
}
`
