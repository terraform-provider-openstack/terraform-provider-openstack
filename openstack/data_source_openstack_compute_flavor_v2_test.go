package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccComputeV2FlavorDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2FlavorDataSource_basic,
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

func TestAccComputeV2FlavorDataSource_testQueries(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2FlavorDataSource_queryDisk,
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
			{
				Config: testAccComputeV2FlavorDataSource_queryMinDisk,
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
			{
				Config: testAccComputeV2FlavorDataSource_queryMinRAM,
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
			{
				Config: testAccComputeV2FlavorDataSource_queryVCPUs,
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

func TestAccComputeV2FlavorDataSource_extraSpecs(t *testing.T) {
	var flavorName = acctest.RandomWithPrefix("tf-acc-flavor")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Flavor_extraSpecs_1(flavorName),
			},
			{
				Config: testAccComputeV2FlavorDataSource_extraSpecs(flavorName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2FlavorDataSourceID("data.openstack_compute_flavor_v2.flavor_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "name", flavorName),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "extra_specs.%", "2"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "extra_specs.hw:cpu_policy", "CPU-POLICY"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "extra_specs.hw:cpu_thread_policy", "CPU-THREAD-POLICY"),
				),
			},
		},
	})
}

func TestAccComputeV2FlavorDataSource_flavorID(t *testing.T) {
	var flavorName = acctest.RandomWithPrefix("tf-acc-flavor")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2Flavor_extraSpecs_1(flavorName),
			},
			{
				Config: testAccComputeV2FlavorDataSource_flavorID(flavorName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2FlavorDataSourceID("data.openstack_compute_flavor_v2.flavor_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "name", flavorName),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "extra_specs.%", "2"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "extra_specs.hw:cpu_policy", "CPU-POLICY"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_flavor_v2.flavor_1", "extra_specs.hw:cpu_thread_policy", "CPU-THREAD-POLICY"),
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

const testAccComputeV2FlavorDataSource_basic = `
data "openstack_compute_flavor_v2" "flavor_1" {
  name = "m1.acctest"
}
`

const testAccComputeV2FlavorDataSource_queryDisk = `
data "openstack_compute_flavor_v2" "flavor_1" {
  disk = 6
}
`

const testAccComputeV2FlavorDataSource_queryMinDisk = `
data "openstack_compute_flavor_v2" "flavor_1" {
  name = "m1.acctest"
  min_disk = 5
}
`

const testAccComputeV2FlavorDataSource_queryMinRAM = `
data "openstack_compute_flavor_v2" "flavor_1" {
  name = "m1.acctest"
  min_ram = 512
}
`

const testAccComputeV2FlavorDataSource_queryVCPUs = `
data "openstack_compute_flavor_v2" "flavor_1" {
  name = "m1.acctest"
  vcpus = 1
}
`

func testAccComputeV2FlavorDataSource_extraSpecs(flavorName string) string {
	flavorResource := testAccComputeV2Flavor_extraSpecs_1(flavorName)

	return fmt.Sprintf(`
          %s

          data "openstack_compute_flavor_v2" "flavor_1" {
            name = "${openstack_compute_flavor_v2.flavor_1.name}"
          }
          `, flavorResource)
}

func testAccComputeV2FlavorDataSource_flavorID(flavorName string) string {
	flavorResource := testAccComputeV2Flavor_extraSpecs_1(flavorName)

	return fmt.Sprintf(`
          %s

          data "openstack_compute_flavor_v2" "flavor_1" {
            flavor_id = "${openstack_compute_flavor_v2.flavor_1.id}"
          }
          `, flavorResource)
}
