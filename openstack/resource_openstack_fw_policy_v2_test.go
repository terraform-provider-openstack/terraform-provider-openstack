package openstack

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/policies"
)

func TestAccFWPolicyV2_basic(t *testing.T) {
	var policy policies.Policy

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckFW(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckFWPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWPolicyV2Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWPolicyV2Exists(
						"openstack_fw_policy_v2.policy_1", &policy),
				),
			},
		},
	})
}

func TestAccFWPolicyV2_shared(t *testing.T) {
	var policy policies.Policy

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
			testAccPreCheckFW(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckFWPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWPolicyV2Shared,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWPolicyV2Exists(
						"openstack_fw_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttr(
						"openstack_fw_policy_v2.policy_1", "shared", "true"),
				),
			},
		},
	})
}

func TestAccFWPolicyV2_addRules(t *testing.T) {
	var policy policies.Policy

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckFW(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckFWPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWPolicyV2Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWPolicyV2Exists(
						"openstack_fw_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttr(
						"openstack_fw_policy_v2.policy_1", "rules.#", "0"),
				),
			},
			{
				Config: testAccFWPolicyV2AddRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWPolicyV2Exists(
						"openstack_fw_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttr(
						"openstack_fw_policy_v2.policy_1", "rules.#", "2"),
				),
			},
		},
	})
}

func TestAccFWPolicyV2_clearFields(t *testing.T) {
	var policy policies.Policy

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckFW(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckFWPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWPolicyV2Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWPolicyV2Exists(
						"openstack_fw_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttr(
						"openstack_fw_policy_v2.policy_1", "name", ""),
					resource.TestCheckResourceAttr(
						"openstack_fw_policy_v2.policy_1", "description", ""),
				),
			},
			{
				Config: testAccFWPolicyV2FillOutFields,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWPolicyV2Exists(
						"openstack_fw_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttr(
						"openstack_fw_policy_v2.policy_1", "name", "policy_1"),
					resource.TestCheckResourceAttr(
						"openstack_fw_policy_v2.policy_1", "description", "terraform acceptance test"),
				),
			},
			{
				Config: testAccFWPolicyV2ClearFields,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWPolicyV2Exists(
						"openstack_fw_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttr(
						"openstack_fw_policy_v2.policy_1", "name", ""),
					resource.TestCheckResourceAttr(
						"openstack_fw_policy_v2.policy_1", "description", ""),
				),
			},
		},
	})
}

func TestAccFWPolicyV2_deleteRules(t *testing.T) {
	var policy policies.Policy

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckFW(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckFWPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWPolicyV2AddRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWPolicyV2Exists(
						"openstack_fw_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttr(
						"openstack_fw_policy_v2.policy_1", "rules.#", "2"),
				),
			},
			{
				Config: testAccFWPolicyV2DeleteRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWPolicyV2Exists(
						"openstack_fw_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttr(
						"openstack_fw_policy_v2.policy_1", "rules.#", "1"),
				),
			},
			{
				Config: testAccFWPolicyV2Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWPolicyV2Exists(
						"openstack_fw_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttr(
						"openstack_fw_policy_v2.policy_1", "rules.#", "0"),
				),
			},
		},
	})
}

func TestAccFWPolicyV2_rulesOrder(t *testing.T) {
	var policy policies.Policy

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckFW(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckFWPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWPolicyV2RulesOrderBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWPolicyV2Exists("openstack_fw_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttrPair(
						"data.openstack_fw_policy_v2.policy_1", "rules.0",
						"openstack_fw_rule_v2.rule_1", "id"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_fw_policy_v2.policy_1", "rules.1",
						"openstack_fw_rule_v2.rule_2", "id"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_fw_policy_v2.policy_1", "rules.2",
						"openstack_fw_rule_v2.rule_3", "id"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_fw_policy_v2.policy_1", "rules.3",
						"openstack_fw_rule_v2.rule_4", "id"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_policy_v2.policy_1", "rules.#", "4"),
				),
			},

			{
				Config: testAccFWPolicyV2RulesOrderRemove,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWPolicyV2Exists("openstack_fw_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttrPair(
						"data.openstack_fw_policy_v2.policy_1", "rules.0",
						"openstack_fw_rule_v2.rule_4", "id"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_fw_policy_v2.policy_1", "rules.1",
						"openstack_fw_rule_v2.rule_2", "id"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_fw_policy_v2.policy_1", "rules.2",
						"openstack_fw_rule_v2.rule_1", "id"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_policy_v2.policy_1", "rules.#", "3"),
				),
			},

			{
				Config: testAccFWPolicyV2RulesOrderRevert,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWPolicyV2Exists("openstack_fw_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttrPair(
						"data.openstack_fw_policy_v2.policy_1", "rules.0",
						"openstack_fw_rule_v2.rule_4", "id"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_fw_policy_v2.policy_1", "rules.1",
						"openstack_fw_rule_v2.rule_3", "id"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_fw_policy_v2.policy_1", "rules.2",
						"openstack_fw_rule_v2.rule_2", "id"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_fw_policy_v2.policy_1", "rules.3",
						"openstack_fw_rule_v2.rule_1", "id"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_policy_v2.policy_1", "rules.#", "4"),
				),
			},

			{
				Config: testAccFWPolicyV2RulesOrderBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWPolicyV2Exists("openstack_fw_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttrPair(
						"data.openstack_fw_policy_v2.policy_1", "rules.0",
						"openstack_fw_rule_v2.rule_1", "id"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_fw_policy_v2.policy_1", "rules.1",
						"openstack_fw_rule_v2.rule_2", "id"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_fw_policy_v2.policy_1", "rules.2",
						"openstack_fw_rule_v2.rule_3", "id"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_fw_policy_v2.policy_1", "rules.3",
						"openstack_fw_rule_v2.rule_4", "id"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_policy_v2.policy_1", "rules.#", "4"),
				),
			},

			{
				Config: testAccFWPolicyV2RulesOrderShuffle,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWPolicyV2Exists("openstack_fw_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttrPair(
						"data.openstack_fw_policy_v2.policy_1", "rules.0",
						"openstack_fw_rule_v2.rule_1", "id"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_fw_policy_v2.policy_1", "rules.1",
						"openstack_fw_rule_v2.rule_4", "id"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_fw_policy_v2.policy_1", "rules.2",
						"openstack_fw_rule_v2.rule_2", "id"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_fw_policy_v2.policy_1", "rules.3",
						"openstack_fw_rule_v2.rule_3", "id"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_policy_v2.policy_1", "rules.#", "4"),
				),
			},

			{
				Config: testAccFWPolicyV2RulesOrderRemove,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWPolicyV2Exists("openstack_fw_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttrPair(
						"data.openstack_fw_policy_v2.policy_1", "rules.0",
						"openstack_fw_rule_v2.rule_4", "id"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_fw_policy_v2.policy_1", "rules.1",
						"openstack_fw_rule_v2.rule_2", "id"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_fw_policy_v2.policy_1", "rules.2",
						"openstack_fw_rule_v2.rule_1", "id"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_policy_v2.policy_1", "rules.#", "3"),
				),
			},

			{
				Config: testAccFWPolicyV2RulesOrderBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWPolicyV2Exists("openstack_fw_policy_v2.policy_1", &policy),
					resource.TestCheckResourceAttrPair(
						"data.openstack_fw_policy_v2.policy_1", "rules.0",
						"openstack_fw_rule_v2.rule_1", "id"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_fw_policy_v2.policy_1", "rules.1",
						"openstack_fw_rule_v2.rule_2", "id"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_fw_policy_v2.policy_1", "rules.2",
						"openstack_fw_rule_v2.rule_3", "id"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_fw_policy_v2.policy_1", "rules.3",
						"openstack_fw_rule_v2.rule_4", "id"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_policy_v2.policy_1", "rules.#", "4"),
				),
			},
		},
	})
}

func testAccCheckFWPolicyV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.NetworkingV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_fw_policy_v2" {
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

func testAccCheckFWPolicyV2Exists(n string, policy *policies.Policy) resource.TestCheckFunc {
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

		*policy = *found

		return nil
	}
}

const testAccFWPolicyV2Basic = `
resource "openstack_fw_policy_v2" "policy_1" {}
`

const testAccFWPolicyV2Shared = `
resource "openstack_fw_policy_v2" "policy_1" {
  shared = true
}
`

const testAccFWPolicyV2AddRules = `
resource "openstack_fw_rule_v2" "tcp_allow" {
  protocol = "tcp"
  action   = "allow"
}

resource "openstack_fw_rule_v2" "udp_deny" {
  protocol = "udp"
  action   = "deny"
}

resource "openstack_fw_policy_v2" "policy_1" {
  name        = "policy_1"
  description = "terraform acceptance test"
  audited     = true
  rules       = [
    "${openstack_fw_rule_v2.udp_deny.id}",
    "${openstack_fw_rule_v2.tcp_allow.id}"
  ]
}
`

const testAccFWPolicyV2FillOutFields = `
resource "openstack_fw_policy_v2" "policy_1" {
  name        = "policy_1"
  description = "terraform acceptance test"
}
`

const testAccFWPolicyV2ClearFields = `
resource "openstack_fw_policy_v2" "policy_1" {
  name        = ""
  description = ""
}
`

const testAccFWPolicyV2DeleteRules = `
resource "openstack_fw_rule_v2" "udp_deny" {
  protocol = "udp"
  action   = "deny"
}

resource "openstack_fw_policy_v2" "policy_1" {
  name        = "policy_1"
  description = "terraform acceptance test"
  rules       = [
    "${openstack_fw_rule_v2.udp_deny.id}"
  ]
}
`

const testAccFWPolicyV2RulesOrderBasic = `
resource "openstack_fw_rule_v2" "rule_1" {
  protocol = "tcp"
  action   = "deny"
}

resource "openstack_fw_rule_v2" "rule_2" {
  protocol = "tcp"
  action   = "allow"
}

resource "openstack_fw_rule_v2" "rule_3" {
  protocol = "udp"
  action   = "allow"
}

resource "openstack_fw_rule_v2" "rule_4" {
  protocol = "udp"
  action   = "deny"
}

resource "openstack_fw_policy_v2" "policy_1" {
  name        = "policy_1"
  description = "terraform acceptance test"
  rules       = [
    "${openstack_fw_rule_v2.rule_1.id}",
	"${openstack_fw_rule_v2.rule_2.id}",
	"${openstack_fw_rule_v2.rule_3.id}",
	"${openstack_fw_rule_v2.rule_4.id}"
  ]
}

data "openstack_fw_policy_v2" "policy_1" {
  policy_id = "${openstack_fw_policy_v2.policy_1.id}"
}
`

const testAccFWPolicyV2RulesOrderRemove = `
resource "openstack_fw_rule_v2" "rule_1" {
  protocol = "tcp"
  action   = "deny"
}

resource "openstack_fw_rule_v2" "rule_2" {
  protocol = "tcp"
  action   = "allow"
}

resource "openstack_fw_rule_v2" "rule_3" {
  protocol = "udp"
  action   = "allow"
}

resource "openstack_fw_rule_v2" "rule_4" {
  protocol = "udp"
  action   = "deny"
}

resource "openstack_fw_policy_v2" "policy_1" {
  name        = "policy_1"
  description = "terraform acceptance test"
  rules       = [
    "${openstack_fw_rule_v2.rule_4.id}",
	"${openstack_fw_rule_v2.rule_2.id}",
	"${openstack_fw_rule_v2.rule_1.id}"
  ]
}

data "openstack_fw_policy_v2" "policy_1" {
  policy_id = "${openstack_fw_policy_v2.policy_1.id}"
}
`

const testAccFWPolicyV2RulesOrderRevert = `
resource "openstack_fw_rule_v2" "rule_1" {
  protocol = "tcp"
  action   = "deny"
}

resource "openstack_fw_rule_v2" "rule_2" {
  protocol = "tcp"
  action   = "allow"
}

resource "openstack_fw_rule_v2" "rule_3" {
  protocol = "udp"
  action   = "allow"
}

resource "openstack_fw_rule_v2" "rule_4" {
  protocol = "udp"
  action   = "deny"
}

resource "openstack_fw_policy_v2" "policy_1" {
  name        = "policy_1"
  description = "terraform acceptance test"
  rules       = [
    "${openstack_fw_rule_v2.rule_4.id}",
	"${openstack_fw_rule_v2.rule_3.id}",
	"${openstack_fw_rule_v2.rule_2.id}",
	"${openstack_fw_rule_v2.rule_1.id}"
  ]
}

data "openstack_fw_policy_v2" "policy_1" {
  policy_id = "${openstack_fw_policy_v2.policy_1.id}"
}
`

const testAccFWPolicyV2RulesOrderShuffle = `
resource "openstack_fw_rule_v2" "rule_1" {
  protocol = "tcp"
  action   = "deny"
}

resource "openstack_fw_rule_v2" "rule_2" {
  protocol = "tcp"
  action   = "allow"
}

resource "openstack_fw_rule_v2" "rule_3" {
  protocol = "udp"
  action   = "allow"
}

resource "openstack_fw_rule_v2" "rule_4" {
  protocol = "udp"
  action   = "deny"
}

resource "openstack_fw_policy_v2" "policy_1" {
  name        = "policy_1"
  description = "terraform acceptance test"
  rules       = [
    "${openstack_fw_rule_v2.rule_1.id}",
	"${openstack_fw_rule_v2.rule_4.id}",
	"${openstack_fw_rule_v2.rule_2.id}",
	"${openstack_fw_rule_v2.rule_3.id}"
  ]
}

data "openstack_fw_policy_v2" "policy_1" {
  policy_id = "${openstack_fw_policy_v2.policy_1.id}"
}
`
