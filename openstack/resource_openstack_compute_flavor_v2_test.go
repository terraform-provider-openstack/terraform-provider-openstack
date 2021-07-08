package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
)

func TestAccComputeV2Flavor_basic(t *testing.T) {
	var flavor flavors.Flavor
	var flavorName = acctest.RandomWithPrefix("tf-acc-flavor")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2FlavorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2FlavorBasic(flavorName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2FlavorExists("openstack_compute_flavor_v2.flavor_1", &flavor),
					resource.TestCheckResourceAttr(
						"openstack_compute_flavor_v2.flavor_1", "ram", "2048"),
					resource.TestCheckResourceAttr(
						"openstack_compute_flavor_v2.flavor_1", "vcpus", "2"),
					resource.TestCheckResourceAttr(
						"openstack_compute_flavor_v2.flavor_1", "disk", "5"),
					resource.TestCheckResourceAttr(
						"openstack_compute_flavor_v2.flavor_1", "ephemeral", "64"),
				),
			},
			{
				Config: testAccComputeV2FlavorBasicWithID(flavorName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2FlavorExists("openstack_compute_flavor_v2.flavor_1", &flavor),
					resource.TestCheckResourceAttr(
						"openstack_compute_flavor_v2.flavor_1", "ram", "2048"),
					resource.TestCheckResourceAttr(
						"openstack_compute_flavor_v2.flavor_1", "vcpus", "2"),
					resource.TestCheckResourceAttr(
						"openstack_compute_flavor_v2.flavor_1", "disk", "5"),
					resource.TestCheckResourceAttr(
						"openstack_compute_flavor_v2.flavor_1", "flavor_id", "b50e603d-29d0-461a-88f7-bd6750d4ce3d"),
				),
			},
		},
	})
}

func TestAccComputeV2Flavor_extraSpecs(t *testing.T) {
	var flavor flavors.Flavor
	var flavorName = acctest.RandomWithPrefix("tf-acc-flavor")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2FlavorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2FlavorExtraSpecs1(flavorName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2FlavorExists("openstack_compute_flavor_v2.flavor_1", &flavor),
					resource.TestCheckResourceAttr(
						"openstack_compute_flavor_v2.flavor_1", "extra_specs.%", "2"),
					resource.TestCheckResourceAttr(
						"openstack_compute_flavor_v2.flavor_1", "extra_specs.hw:cpu_policy", "CPU-POLICY"),
					resource.TestCheckResourceAttr(
						"openstack_compute_flavor_v2.flavor_1", "extra_specs.hw:cpu_thread_policy", "CPU-THREAD-POLICY"),
				),
			},
			{
				Config: testAccComputeV2FlavorExtraSpecs2(flavorName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2FlavorExists("openstack_compute_flavor_v2.flavor_1", &flavor),
					resource.TestCheckResourceAttr(
						"openstack_compute_flavor_v2.flavor_1", "extra_specs.%", "1"),
					resource.TestCheckResourceAttr(
						"openstack_compute_flavor_v2.flavor_1", "extra_specs.hw:cpu_policy", "CPU-POLICY-2"),
				),
			},
		},
	})
}

func testAccCheckComputeV2FlavorDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	computeClient, err := config.ComputeV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_compute_flavor_v2" {
			continue
		}

		_, err := flavors.Get(computeClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Flavor still exists")
		}
	}

	return nil
}

func testAccCheckComputeV2FlavorExists(n string, flavor *flavors.Flavor) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		computeClient, err := config.ComputeV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack compute client: %s", err)
		}

		found, err := flavors.Get(computeClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Flavor not found")
		}

		*flavor = *found

		return nil
	}
}

func testAccComputeV2FlavorBasic(flavorName string) string {
	return fmt.Sprintf(`
    resource "openstack_compute_flavor_v2" "flavor_1" {
      name = "%s"
      ram = 2048
      vcpus = 2
      disk = 5
      ephemeral = 64

      is_public = true
    }
    `, flavorName)
}

func testAccComputeV2FlavorBasicWithID(flavorName string) string {
	return fmt.Sprintf(`
    resource "openstack_compute_flavor_v2" "flavor_1" {
      name = "%s"
      flavor_id = "b50e603d-29d0-461a-88f7-bd6750d4ce3d"
      ram = 2048
      vcpus = 2
      disk = 5

      is_public = true
    }
    `, flavorName)
}
func testAccComputeV2FlavorExtraSpecs1(flavorName string) string {
	return fmt.Sprintf(`
    resource "openstack_compute_flavor_v2" "flavor_1" {
      name = "%s"
      ram = 2048
      vcpus = 2
      disk = 5

      is_public = true

      extra_specs = {
        "hw:cpu_policy" = "CPU-POLICY",
        "hw:cpu_thread_policy" = "CPU-THREAD-POLICY"
      }
    }
    `, flavorName)
}

func testAccComputeV2FlavorExtraSpecs2(flavorName string) string {
	return fmt.Sprintf(`
    resource "openstack_compute_flavor_v2" "flavor_1" {
      name = "%s"
      ram = 2048
      vcpus = 2
      disk = 5

      is_public = true

      extra_specs = {
        "hw:cpu_policy" = "CPU-POLICY-2"
      }
    }
    `, flavorName)
}
