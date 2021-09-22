package openstack

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas/policies"
)

func TestAccFWPolicyV1_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckFW(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckFWPolicyV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWPolicyV1Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWPolicyV1Exists(
						"openstack_fw_policy_v1.policy_1", "", "", 0),
				),
			},
		},
	})
}

func TestAccFWPolicyV1_addRules(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckFW(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckFWPolicyV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWPolicyV1AddRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWPolicyV1Exists(
						"openstack_fw_policy_v1.policy_1", "policy_1", "terraform acceptance test", 2),
				),
			},
		},
	})
}

func TestAccFWPolicyV1_deleteRules(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckFW(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckFWPolicyV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWPolicyV1DeleteRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWPolicyV1Exists(
						"openstack_fw_policy_v1.policy_1", "policy_1", "terraform acceptance test", 1),
				),
			},
		},
	})
}

func testAccCheckFWPolicyV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.NetworkingV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_fw_policy_v1" {
			continue
		}
		_, err = policies.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Firewall policy (%s) still exists", rs.Primary.ID)
		}
		if _, ok := err.(gophercloud.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckFWPolicyV1Exists(n, name, description string, ruleCount int) resource.TestCheckFunc {
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

		var found *policies.Policy
		for i := 0; i < 5; i++ {
			// Firewall policy creation is asynchronous. Retry some times
			// if we get a 404 error. Fail on any other error.
			found, err = policies.Get(networkingClient, rs.Primary.ID).Extract()
			if err != nil {
				if _, ok := err.(gophercloud.ErrDefault404); ok {
					time.Sleep(time.Second)
					continue
				}
				return err
			}
			break
		}

		switch {
		case name != found.Name:
			err = fmt.Errorf("Expected name <%s>, but found <%s>", name, found.Name)
		case description != found.Description:
			err = fmt.Errorf("Expected description <%s>, but found <%s>", description, found.Description)
		case ruleCount != len(found.Rules):
			err = fmt.Errorf("Expected rule count <%d>, but found <%d>", ruleCount, len(found.Rules))
		}

		if err != nil {
			return err
		}

		return nil
	}
}

const testAccFWPolicyV1Basic = `
resource "openstack_fw_policy_v1" "policy_1" {
}
`

const testAccFWPolicyV1AddRules = `
resource "openstack_fw_policy_v1" "policy_1" {
  name = "policy_1"
  description =  "terraform acceptance test"
  rules = [
    "${openstack_fw_rule_v1.udp_deny.id}",
    "${openstack_fw_rule_v1.tcp_allow.id}"
  ]
}

resource "openstack_fw_rule_v1" "tcp_allow" {
  protocol = "tcp"
  action = "allow"
}

resource "openstack_fw_rule_v1" "udp_deny" {
  protocol = "udp"
  action = "deny"
}
`

const testAccFWPolicyV1DeleteRules = `
resource "openstack_fw_policy_v1" "policy_1" {
  name = "policy_1"
  description =  "terraform acceptance test"
  rules = [
    "${openstack_fw_rule_v1.udp_deny.id}"
  ]
}

resource "openstack_fw_rule_v1" "udp_deny" {
  protocol = "udp"
  action = "deny"
}
`
