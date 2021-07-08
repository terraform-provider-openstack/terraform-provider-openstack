package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/qos/policies"
)

func TestAccNetworkingV2QoSPolicyBasic(t *testing.T) {
	var policy policies.Policy

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2QoSPolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2QoSPolicyBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2QoSPolicyExists(
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

func testAccCheckNetworkingV2QoSPolicyExists(n string, policy *policies.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.NetworkingV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %s", err)
		}

		found, err := policies.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("QoS policy not found")
		}

		*policy = *found

		return nil
	}
}

func testAccCheckNetworkingV2QoSPolicyDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.NetworkingV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_networking_qos_policy_v2" {
			continue
		}

		_, err := policies.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("QoS policy still exists")
		}
	}

	return nil
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
