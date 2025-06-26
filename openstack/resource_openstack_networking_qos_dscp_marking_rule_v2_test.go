package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/qos/policies"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/qos/rules"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
		CheckDestroy:      testAccCheckNetworkingV2QoSDSCPMarkingRuleDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2QoSDSCPMarkingRuleBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2QoSPolicyExists(t.Context(),
						"openstack_networking_qos_policy_v2.qos_policy_1", &policy),
					resource.TestCheckResourceAttr(
						"openstack_networking_qos_policy_v2.qos_policy_1", "name", "qos_policy_1"),
					testAccCheckNetworkingV2QoSDSCPMarkingRuleExists(t.Context(),
						"openstack_networking_qos_dscp_marking_rule_v2.dscp_marking_rule_1", &rule),
					resource.TestCheckResourceAttr(
						"openstack_networking_qos_dscp_marking_rule_v2.dscp_marking_rule_1", "dscp_mark", "26"),
				),
			},
			{
				Config: testAccNetworkingV2QoSDSCPMarkingRuleUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2QoSPolicyExists(t.Context(),
						"openstack_networking_qos_policy_v2.qos_policy_1", &policy),
					resource.TestCheckResourceAttr(
						"openstack_networking_qos_policy_v2.qos_policy_1", "name", "qos_policy_1"),
					testAccCheckNetworkingV2QoSDSCPMarkingRuleExists(t.Context(),
						"openstack_networking_qos_dscp_marking_rule_v2.dscp_marking_rule_1", &rule),
					resource.TestCheckResourceAttr(
						"openstack_networking_qos_dscp_marking_rule_v2.dscp_marking_rule_1", "dscp_mark", "20"),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2QoSDSCPMarkingRuleExists(ctx context.Context, n string, rule *rules.DSCPMarkingRule) resource.TestCheckFunc {
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

		qosPolicyID, qosRuleID, err := parsePairedIDs(rs.Primary.ID, "openstack_networking_qos_dscp_marking_rule_v2")
		if err != nil {
			return err
		}

		found, err := rules.GetDSCPMarkingRule(ctx, networkingClient, qosPolicyID, qosRuleID).ExtractDSCPMarkingRule()
		if err != nil {
			return err
		}

		foundID := resourceNetworkingQoSRuleV2BuildID(qosPolicyID, found.ID)

		if foundID != rs.Primary.ID {
			return errors.New("QoS dscp marking rule not found")
		}

		*rule = *found

		return nil
	}
}

func testAccCheckNetworkingV2QoSDSCPMarkingRuleDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		networkingClient, err := config.NetworkingV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_networking_qos_dscp_marking_rule_v2" {
				continue
			}

			qosPolicyID, qosRuleID, err := parsePairedIDs(rs.Primary.ID, "openstack_networking_qos_dscp_marking_rule_v2")
			if err != nil {
				return err
			}

			_, err = rules.GetDSCPMarkingRule(ctx, networkingClient, qosPolicyID, qosRuleID).ExtractDSCPMarkingRule()
			if err == nil {
				return errors.New("QoS rule still exists")
			}
		}

		return nil
	}
}

const testAccNetworkingV2QoSDSCPMarkingRuleBasic = `
resource "openstack_networking_qos_policy_v2" "qos_policy_1" {
  name = "qos_policy_1"
}

resource "openstack_networking_qos_dscp_marking_rule_v2" "dscp_marking_rule_1" {
  qos_policy_id  = openstack_networking_qos_policy_v2.qos_policy_1.id
  dscp_mark      = 26
}
`

const testAccNetworkingV2QoSDSCPMarkingRuleUpdate = `
resource "openstack_networking_qos_policy_v2" "qos_policy_1" {
  name = "qos_policy_1"
}

resource "openstack_networking_qos_dscp_marking_rule_v2" "dscp_marking_rule_1" {
  qos_policy_id  = openstack_networking_qos_policy_v2.qos_policy_1.id
  dscp_mark      = 20
}
`
