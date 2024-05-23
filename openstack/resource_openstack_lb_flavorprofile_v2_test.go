package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/flavorprofiles"
)

func TestAccLBV2FlavorProfile_basic(t *testing.T) {
	var fp flavorprofiles.FlavorProfile

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV2FlavorProfileDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLbV2FlavorProfile,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2FlavorProfileExists("openstack_lb_flavorprofile_v2.fp_1", &fp),
					resource.TestCheckResourceAttr(
						"openstack_lb_flavorprofile_v2.fp_1", "name", "test"),
					resource.TestCheckResourceAttr(
						"openstack_lb_flavorprofile_v2.fp_1", "provider_name", "amphora"),
					resource.TestCheckResourceAttr(
						"openstack_lb_flavorprofile_v2.fp_1", "flavor_data", "{\"loadbalancer_topology\":\"ACTIVE_STANDBY\"}"),
				),
			},
			{
				Config: testAccCheckLbV2FlavorProfileUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2FlavorProfileExists("openstack_lb_flavorprofile_v2.fp_1", &fp),
					resource.TestCheckResourceAttr(
						"openstack_lb_flavorprofile_v2.fp_1", "name", "test-2"),
					resource.TestCheckResourceAttr(
						"openstack_lb_flavorprofile_v2.fp_1", "provider_name", "amphora"),
					resource.TestCheckResourceAttr(
						"openstack_lb_flavorprofile_v2.fp_1", "flavor_data", "{\"loadbalancer_topology\":\"SINGLE\"}"),
				),
			},
		},
	})
}

func testAccCheckLBV2FlavorProfileDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	lbClient, err := config.LoadBalancerV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack load balancing client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_lb_flavorprofile_v2" {
			continue
		}

		_, err := flavorprofiles.Get(lbClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Flavor profile still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV2FlavorProfileExists(n string, fp *flavorprofiles.FlavorProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		lbClient, err := config.LoadBalancerV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack load balancing client: %s", err)
		}

		found, err := flavorprofiles.Get(lbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Flavor profile not found")
		}

		*fp = *found

		return nil
	}
}

const testAccCheckLbV2FlavorProfile = `
resource "openstack_lb_flavorprofile_v2" "fp_1" {
	name          = "test"
	provider_name = "amphora"
	flavor_data   = jsonencode({
	  "loadbalancer_topology": "ACTIVE_STANDBY",
	})
  }
`

const testAccCheckLbV2FlavorProfileUpdate = `
resource "openstack_lb_flavorprofile_v2" "fp_1" {
	name          = "test-2"
	provider_name = "amphora"
	flavor_data   = jsonencode({
	  "loadbalancer_topology": "SINGLE",
	})
  }
`
