package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/qos/policies"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccNetworkingV2QoSPolicyBasic(t *testing.T) {
	var policy policies.Policy

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2QoSPolicyDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2QoSPolicyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2QoSPolicyExists(t.Context(),
						"openstack_networking_qos_policy_v2.qos_policy_1", &policy),
					resource.TestCheckResourceAttr(
						"openstack_networking_qos_policy_v2.qos_policy_1", "name", "qos_policy_1"),
					resource.TestCheckResourceAttr(
						"openstack_networking_qos_policy_v2.qos_policy_1", "description", "terraform qos policy acceptance test"),
				),
			},
			{
				Config: testAccNetworkingV2QoSPolicyUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_networking_qos_policy_v2.qos_policy_1", "name", "qos_policy_1"),
					resource.TestCheckResourceAttr(
						"openstack_networking_qos_policy_v2.qos_policy_1", "description", "terraform qos policy acceptance test updated"),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2QoSPolicyExists(ctx context.Context, n string, policy *policies.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		networkingClient, err := config.NetworkingV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %w", err)
		}

		found, err := policies.Get(ctx, networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return errors.New("QoS policy not found")
		}

		*policy = *found

		return nil
	}
}

func testAccCheckNetworkingV2QoSPolicyDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		networkingClient, err := config.NetworkingV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_networking_qos_policy_v2" {
				continue
			}

			_, err := policies.Get(ctx, networkingClient, rs.Primary.ID).Extract()
			if err == nil {
				return errors.New("QoS policy still exists")
			}
		}

		return nil
	}
}

const testAccNetworkingV2QoSPolicyBasic = `
resource "openstack_networking_qos_policy_v2" "qos_policy_1" {
	name        = "qos_policy_1"
	description = "terraform qos policy acceptance test"
}
`

const testAccNetworkingV2QoSPolicyUpdate = `
resource "openstack_networking_qos_policy_v2" "qos_policy_1" {
	name        = "qos_policy_1"
	description = "terraform qos policy acceptance test updated"
}
`
