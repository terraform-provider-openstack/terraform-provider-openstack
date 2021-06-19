package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/qos/policies"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/qos/rules"
)

func TestAccNetworkingV2QoSDSCPMarkingRule_basic(t *testing.T) {
	var (
		policy policies.Policy
		rule   rules.DSCPMarkingRule
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2QoSDSCPMarkingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2QoSDSCPMarkingRuleBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2QoSPolicyExists(
						"openstack_networking_qos_policy_v2.qos_policy_1", &policy),
					resource.TestCheckResourceAttr(
						"openstack_networking_qos_policy_v2.qos_policy_1", "name", "qos_policy_1"),
					testAccCheckNetworkingV2QoSDSCPMarkingRuleExists(
						"openstack_networking_qos_dscp_marking_rule_v2.dscp_marking_rule_1", &rule),
					resource.TestCheckResourceAttr(
						"openstack_networking_qos_dscp_marking_rule_v2.dscp_marking_rule_1", "dscp_mark", "26"),
				),
			},
			{
				Config: testAccNetworkingV2QoSDSCPMarkingRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2QoSPolicyExists(
						"openstack_networking_qos_policy_v2.qos_policy_1", &policy),
					resource.TestCheckResourceAttr(
						"openstack_networking_qos_policy_v2.qos_policy_1", "name", "qos_policy_1"),
					testAccCheckNetworkingV2QoSDSCPMarkingRuleExists(
						"openstack_networking_qos_dscp_marking_rule_v2.dscp_marking_rule_1", &rule),
					resource.TestCheckResourceAttr(
						"openstack_networking_qos_dscp_marking_rule_v2.dscp_marking_rule_1", "dscp_mark", "20"),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2QoSDSCPMarkingRuleExists(n string, rule *rules.DSCPMarkingRule) resource.TestCheckFunc {
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
			return fmt.Errorf("Error reading openstack_networking_qos_dscp_marking_rule_v2 ID %s: %s", rs.Primary.ID, err)
		}

		found, err := rules.GetDSCPMarkingRule(networkingClient, qosPolicyID, qosRuleID).ExtractDSCPMarkingRule()
		if err != nil {
			return err
		}

		foundID := resourceNetworkingQoSRuleV2BuildID(qosPolicyID, found.ID)

		if foundID != rs.Primary.ID {
			return fmt.Errorf("QoS dscp marking rule not found")
		}

		*rule = *found

		return nil
	}
}

func testAccCheckNetworkingV2QoSDSCPMarkingRuleDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.NetworkingV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_networking_qos_dscp_marking_rule_v2" {
			continue
		}

		qosPolicyID, qosRuleID, err := resourceNetworkingQoSRuleV2ParseID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error reading openstack_networking_qos_dscp_marking_rule_v2 ID %s: %s", rs.Primary.ID, err)
		}

		_, err = rules.GetDSCPMarkingRule(networkingClient, qosPolicyID, qosRuleID).ExtractDSCPMarkingRule()
		if err == nil {
			return fmt.Errorf("QoS rule still exists")
		}
	}

	return nil
}

const testAccNetworkingV2QoSDSCPMarkingRuleBasic = `
resource "openstack_networking_qos_policy_v2" "qos_policy_1" {
  name = "qos_policy_1"
}

resource "openstack_networking_qos_dscp_marking_rule_v2" "dscp_marking_rule_1" {
  qos_policy_id  = "${openstack_networking_qos_policy_v2.qos_policy_1.id}"
  dscp_mark      = 26
}
`

const testAccNetworkingV2QoSDSCPMarkingRuleUpdate = `
resource "openstack_networking_qos_policy_v2" "qos_policy_1" {
  name = "qos_policy_1"
}

resource "openstack_networking_qos_dscp_marking_rule_v2" "dscp_marking_rule_1" {
  qos_policy_id  = "${openstack_networking_qos_policy_v2.qos_policy_1.id}"
  dscp_mark      = 20
}
`
