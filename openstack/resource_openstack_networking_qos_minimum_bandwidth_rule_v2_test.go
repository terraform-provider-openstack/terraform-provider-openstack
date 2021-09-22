package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/qos/policies"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/qos/rules"
)

func TestAccNetworkingV2QoSMinimumBandwidthRule_basic(t *testing.T) {
	var (
		policy policies.Policy
		rule   rules.MinimumBandwidthRule
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2QoSMinimumBandwidthRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2QoSMinimumBandwidthRuleBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2QoSPolicyExists(
						"openstack_networking_qos_policy_v2.qos_policy_1", &policy),
					resource.TestCheckResourceAttr(
						"openstack_networking_qos_policy_v2.qos_policy_1", "name", "qos_policy_1"),
					testAccCheckNetworkingV2QoSMinimumBandwidthRuleExists(
						"openstack_networking_qos_minimum_bandwidth_rule_v2.minimum_bandwidth_rule_1", &rule),
					resource.TestCheckResourceAttr(
						"openstack_networking_qos_minimum_bandwidth_rule_v2.minimum_bandwidth_rule_1", "min_kbps", "200"),
				),
			},
			{
				Config: testAccNetworkingV2QoSMinimumBandwidthRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2QoSPolicyExists(
						"openstack_networking_qos_policy_v2.qos_policy_1", &policy),
					resource.TestCheckResourceAttr(
						"openstack_networking_qos_policy_v2.qos_policy_1", "name", "qos_policy_1"),
					testAccCheckNetworkingV2QoSMinimumBandwidthRuleExists(
						"openstack_networking_qos_minimum_bandwidth_rule_v2.minimum_bandwidth_rule_1", &rule),
					resource.TestCheckResourceAttr(
						"openstack_networking_qos_minimum_bandwidth_rule_v2.minimum_bandwidth_rule_1", "min_kbps", "300"),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2QoSMinimumBandwidthRuleExists(n string, rule *rules.MinimumBandwidthRule) resource.TestCheckFunc {
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

		qosPolicyID, qosRuleID, err := resourceNetworkingQoSRuleV2ParseID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error reading openstack_networking_qos_minimum_bandwidth_rule_v2 ID %s: %s", rs.Primary.ID, err)
		}

		found, err := rules.GetMinimumBandwidthRule(networkingClient, qosPolicyID, qosRuleID).ExtractMinimumBandwidthRule()
		if err != nil {
			return err
		}

		foundID := resourceNetworkingQoSRuleV2BuildID(qosPolicyID, found.ID)

		if foundID != rs.Primary.ID {
			return fmt.Errorf("QoS min bw rule not found")
		}

		*rule = *found

		return nil
	}
}

func testAccCheckNetworkingV2QoSMinimumBandwidthRuleDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.NetworkingV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_networking_qos_minimum_bandwidth_rule_v2" {
			continue
		}

		qosPolicyID, qosRuleID, err := resourceNetworkingQoSRuleV2ParseID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error reading openstack_networking_qos_minimum_bandwidth_rule_v2 ID %s: %s", rs.Primary.ID, err)
		}

		_, err = rules.GetMinimumBandwidthRule(networkingClient, qosPolicyID, qosRuleID).ExtractMinimumBandwidthRule()
		if err == nil {
			return fmt.Errorf("QoS rule still exists")
		}
	}

	return nil
}

const testAccNetworkingV2QoSMinimumBandwidthRuleBasic = `
resource "openstack_networking_qos_policy_v2" "qos_policy_1" {
  name = "qos_policy_1"
}

resource "openstack_networking_qos_minimum_bandwidth_rule_v2" "minimum_bandwidth_rule_1" {
  qos_policy_id  = "${openstack_networking_qos_policy_v2.qos_policy_1.id}"
  min_kbps       = 200
}
`

const testAccNetworkingV2QoSMinimumBandwidthRuleUpdate = `
resource "openstack_networking_qos_policy_v2" "qos_policy_1" {
  name = "qos_policy_1"
}

resource "openstack_networking_qos_minimum_bandwidth_rule_v2" "minimum_bandwidth_rule_1" {
  qos_policy_id  = "${openstack_networking_qos_policy_v2.qos_policy_1.id}"
  min_kbps       = 300
}
`
