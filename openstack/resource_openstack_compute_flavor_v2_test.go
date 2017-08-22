package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/flavors"
)

func TestAccComputeV2Flavor_basic(t *testing.T) {
	var flavor flavors.Flavor
	var flavorName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2FlavorDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeV2Flavor_basic(flavorName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2FlavorExists("openstack_compute_flavor_v2.flavor_1", &flavor),
				),
			},
		},
	})
}

func testAccCheckComputeV2FlavorDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	computeClient, err := config.computeV2Client(OS_REGION_NAME)
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
		computeClient, err := config.computeV2Client(OS_REGION_NAME)
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

func testAccComputeV2Flavor_basic(flavorName string) string {
	return fmt.Sprintf(`
    resource "openstack_compute_flavor_v2" "flavor_1" {
      name = "%s"
      ram = 2048
      vcpus = 2
      disk = 5
    }
    `, flavorName)
}
