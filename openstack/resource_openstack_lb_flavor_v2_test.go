package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/flavors"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccLBV2Flavor_basic(t *testing.T) {
	var fp flavors.Flavor

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV2FlavorDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLbV2Flavor,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2FlavorExists(t.Context(), "openstack_lb_flavor_v2.flavor_1", &fp),
					resource.TestCheckResourceAttr(
						"openstack_lb_flavor_v2.flavor_1", "name", "test"),
					resource.TestCheckResourceAttr(
						"openstack_lb_flavor_v2.flavor_1", "description", "test"),
					resource.TestCheckResourceAttr(
						"openstack_lb_flavor_v2.flavor_1", "enabled", "true"),
				),
			},
			{
				Config: testAccCheckLbV2FlavorUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2FlavorExists(t.Context(), "openstack_lb_flavor_v2.flavor_1", &fp),
					resource.TestCheckResourceAttr(
						"openstack_lb_flavor_v2.flavor_1", "name", "test-disabled"),
					resource.TestCheckResourceAttr(
						"openstack_lb_flavor_v2.flavor_1", "description", "test-disabled"),
					resource.TestCheckResourceAttr(
						"openstack_lb_flavor_v2.flavor_1", "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckLBV2FlavorDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		lbClient, err := config.LoadBalancerV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack load balancing client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_lb_flavor_v2" {
				continue
			}

			_, err := flavors.Get(ctx, lbClient, rs.Primary.ID).Extract()
			if err == nil {
				return fmt.Errorf("Flavor still exists: %s", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccCheckLBV2FlavorExists(ctx context.Context, n string, fp *flavors.Flavor) resource.TestCheckFunc {
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

		found, err := flavors.Get(ctx, lbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return errors.New("Flavor not found")
		}

		*fp = *found

		return nil
	}
}

const testAccCheckLbV2Flavor = `
resource "openstack_lb_flavorprofile_v2" "fp_1" {
	name          = "test"
	provider_name = "amphora"
	flavor_data   = jsonencode({
	  "loadbalancer_topology": "ACTIVE_STANDBY",
	})
}

resource "openstack_lb_flavor_v2" "flavor_1" {
	name              = "test"
	description       = "test"
	flavor_profile_id = openstack_lb_flavorprofile_v2.fp_1.id
}
`

const testAccCheckLbV2FlavorUpdate = `
resource "openstack_lb_flavorprofile_v2" "fp_1" {
	name          = "test"
	provider_name = "amphora"
	flavor_data   = jsonencode({
		"loadbalancer_topology": "ACTIVE_STANDBY",
	})
}

resource "openstack_lb_flavor_v2" "flavor_1" {
	name              = "test-disabled"
	description       = "test-disabled"
	enabled           = false
	flavor_profile_id = openstack_lb_flavorprofile_v2.fp_1.id
}
`
