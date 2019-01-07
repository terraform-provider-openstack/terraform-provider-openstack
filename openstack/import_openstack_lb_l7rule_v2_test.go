package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccLBV2L7Rule_importBasic(t *testing.T) {
	l7ruleResourceName := "openstack_lb_l7rule_v2.l7rule_1"
	l7policyResourceName := "openstack_lb_l7policy_v2.l7policy_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckLB(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2L7RuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLBV2L7RuleConfig_basic,
			},

			{
				ResourceName:      l7ruleResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccLBV2L7RuleImportID(l7policyResourceName, l7ruleResourceName),
			},
		},
	})
}

func testAccLBV2L7RuleImportID(l7policyResource, l7ruleResource string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		l7policy, ok := s.RootModule().Resources[l7policyResource]
		if !ok {
			return "", fmt.Errorf("Pool not found: %s", l7policyResource)
		}

		l7rule, ok := s.RootModule().Resources[l7ruleResource]
		if !ok {
			return "", fmt.Errorf("L7Rule not found: %s", l7ruleResource)
		}

		return fmt.Sprintf("%s/%s", l7policy.Primary.ID, l7rule.Primary.ID), nil
	}
}
