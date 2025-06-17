package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/flavorprofiles"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
		CheckDestroy:      testAccCheckLBV2FlavorProfileDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLbV2FlavorProfile,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2FlavorProfileExists(t.Context(), "openstack_lb_flavorprofile_v2.fp_1", &fp),
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
					testAccCheckLBV2FlavorProfileExists(t.Context(), "openstack_lb_flavorprofile_v2.fp_1", &fp),
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

func testAccCheckLBV2FlavorProfileDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		lbClient, err := config.LoadBalancerV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack load balancing client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_lb_flavorprofile_v2" {
				continue
			}

			_, err := flavorprofiles.Get(ctx, lbClient, rs.Primary.ID).Extract()
			if err == nil {
				return fmt.Errorf("Flavor profile still exists: %s", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccCheckLBV2FlavorProfileExists(ctx context.Context, n string, fp *flavorprofiles.FlavorProfile) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		lbClient, err := config.LoadBalancerV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack load balancing client: %w", err)
		}

		found, err := flavorprofiles.Get(ctx, lbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return errors.New("Flavor profile not found")
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
