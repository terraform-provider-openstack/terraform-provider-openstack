package openstack

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/rules"
)

func TestAccFWRuleV2_basic(t *testing.T) {
	var rule rules.Rule

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckFW(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckFWRuleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWRuleV2Basic1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWRuleV2Exists("openstack_fw_rule_v2.rule_1", &rule),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "name", "rule_1"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "protocol", "udp"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "action", "deny"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "ip_version", "4"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "enabled", "true"),
				),
			},

			{
				Config: testAccFWRuleV2Basic2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWRuleV2Exists("openstack_fw_rule_v2.rule_1", &rule),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "name", "rule_1"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "protocol", "udp"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "action", "deny"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "description", "Terraform accept test"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "ip_version", "4"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "source_ip_address", "1.2.3.4"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "destination_ip_address", "4.3.2.0/24"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "source_port", "444"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "destination_port", "555"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "enabled", "true"),
				),
			},

			{
				Config: testAccFWRuleV2Basic3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWRuleV2Exists("openstack_fw_rule_v2.rule_1", &rule),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "name", "rule_1"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "action", "allow"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "description", "Terraform accept test updated"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "ip_version", "4"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "source_ip_address", "1.2.3.0/24"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "destination_ip_address", "4.3.2.8"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "source_port", "666"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "destination_port", "777"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "enabled", "false"),
				),
			},

			{
				Config: testAccFWRuleV2Basic4,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWRuleV2Exists("openstack_fw_rule_v2.rule_1", &rule),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "name", "rule_1"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "protocol", "udp"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "action", "allow"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "description", "Terraform accept test updated"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "ip_version", "4"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "source_ip_address", "1.2.3.0/24"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "destination_ip_address", "4.3.2.8"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "source_port", "666"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "destination_port", "777"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "enabled", "false"),
				),
			},

			{
				Config: testAccFWRuleV2Basic5,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWRuleV2Exists("openstack_fw_rule_v2.rule_1", &rule),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "name", "rule_1"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "protocol", "udp"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "action", "allow"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "description", "Terraform accept test updated"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "ip_version", "4"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "source_ip_address", "1.2.3.0/24"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "destination_ip_address", "4.3.2.8"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "source_port", "666"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "enabled", "false"),
				),
			},
		},
	})
}

func TestAccFWRuleV2_shared(t *testing.T) {
	var rule rules.Rule

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
			testAccPreCheckFW(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckFWRuleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWRuleV2Shared,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWRuleV2Exists("openstack_fw_rule_v2.rule_1", &rule),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "name", "shared_rule"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "protocol", "any"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "action", "deny"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "ip_version", "4"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "enabled", "true"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "shared", "true"),
				),
			},
		},
	})
}

func TestAccFWRuleV2_anyProtocol(t *testing.T) {
	var rule rules.Rule

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckFW(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckFWRuleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWRuleV2AnyProtocol,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWRuleV2Exists("openstack_fw_rule_v2.rule_1", &rule),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "name", "rule_1"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "action", "allow"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "protocol", "any"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "description", "Allow any protocol"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "ip_version", "4"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "source_ip_address", "192.168.199.0/24"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "enabled", "true"),
				),
			},
		},
	})
}

func TestAccFWRuleV2_updateName(t *testing.T) {
	var rule rules.Rule

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckFW(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckFWRuleV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWRuleV2UpdateName1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWRuleV2Exists("openstack_fw_rule_v2.rule_1", &rule),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "name", "rule_1"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "action", "deny"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "protocol", "udp"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "description", "Terraform accept test"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "ip_version", "4"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "source_ip_address", "1.2.3.4"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "destination_ip_address", "4.3.2.0/24"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "source_port", "444"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "destination_port", "555"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "enabled", "true"),
				),
			},

			{
				Config: testAccFWRuleV2UpdateName2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWRuleV2Exists("openstack_fw_rule_v2.rule_1", &rule),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "name", "updated_rule_1"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "action", "deny"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "protocol", "udp"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "description", "Terraform accept test"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "ip_version", "4"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "source_ip_address", "1.2.3.4"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "destination_ip_address", "4.3.2.0/24"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "source_port", "444"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "destination_port", "555"),
					resource.TestCheckResourceAttr(
						"openstack_fw_rule_v2.rule_1", "enabled", "true"),
				),
			},
		},
	})
}

func testAccCheckFWRuleV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.NetworkingV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_fw_rule_v2" {
			continue
		}
		_, err = rules.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Firewall rule (%s) still exists", rs.Primary.ID)
		}
		if _, ok := err.(gophercloud.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckFWRuleV2Exists(n string, rule *rules.Rule) resource.TestCheckFunc {
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

		var found *rules.Rule
		for i := 0; i < 5; i++ {
			// Firewall rule creation is asynchronous. Retry some times
			// if we get a 404 error. Fail on any other error.
			found, err = rules.Get(networkingClient, rs.Primary.ID).Extract()
			if err != nil {
				if _, ok := err.(gophercloud.ErrDefault404); ok {
					time.Sleep(time.Second)
					continue
				}
				return err
			}
			break
		}

		*rule = *found

		return nil
	}
}

const testAccFWRuleV2Basic1 = `
resource "openstack_fw_rule_v2" "rule_1" {
  name     = "rule_1"
  protocol = "udp"
  action   = "deny"
}
`

const testAccFWRuleV2Shared = `
resource "openstack_fw_rule_v2" "rule_1" {
  name     = "shared_rule"
  shared   = true
}
`

const testAccFWRuleV2Basic2 = `
resource "openstack_fw_rule_v2" "rule_1" {
  name                   = "rule_1"
  description            = "Terraform accept test"
  protocol               = "udp"
  action                 = "deny"
  ip_version             = 4
  source_ip_address      = "1.2.3.4"
  destination_ip_address = "4.3.2.0/24"
  source_port            = "444"
  destination_port       = "555"
  enabled                = true
  shared                 = false
}
`

const testAccFWRuleV2Basic3 = `
resource "openstack_fw_rule_v2" "rule_1" {
  name                   = "rule_1"
  description            = "Terraform accept test updated"
  protocol               = "tcp"
  action                 = "allow"
  ip_version             = 4
  source_ip_address      = "1.2.3.0/24"
  destination_ip_address = "4.3.2.8"
  source_port            = "666"
  destination_port       = "777"
  enabled                = false
}
`

const testAccFWRuleV2Basic4 = `
resource "openstack_fw_rule_v2" "rule_1" {
  name                   = "rule_1"
  description            = "Terraform accept test updated"
  protocol               = "udp"
  action                 = "allow"
  ip_version             = 4
  source_ip_address      = "1.2.3.0/24"
  destination_ip_address = "4.3.2.8"
  source_port            = "666"
  destination_port       = "777"
  enabled                = false
}
`

const testAccFWRuleV2Basic5 = `
resource "openstack_fw_rule_v2" "rule_1" {
  name                   = "rule_1"
  description            = "Terraform accept test updated"
  protocol               = "udp"
  action                 = "allow"
  ip_version             = 4
  source_ip_address      = "1.2.3.0/24"
  destination_ip_address = "4.3.2.8"
  source_port            = "666"
  enabled                = false
}
`

const testAccFWRuleV2AnyProtocol = `
resource "openstack_fw_rule_v2" "rule_1" {
  name              = "rule_1"
  description       = "Allow any protocol"
  protocol          = "any"
  action            = "allow"
  ip_version        = 4
  source_ip_address = "192.168.199.0/24"
  enabled           = true
}
`

const testAccFWRuleV2UpdateName1 = `
resource "openstack_fw_rule_v2" "rule_1" {
  name                   = "rule_1"
  description            = "Terraform accept test"
  protocol               = "udp"
  action                 = "deny"
  ip_version             = 4
  source_ip_address      = "1.2.3.4"
  destination_ip_address = "4.3.2.0/24"
  source_port            = "444"
  destination_port       = "555"
  enabled                = true
}
`

const testAccFWRuleV2UpdateName2 = `
resource "openstack_fw_rule_v2" "rule_1" {
  name                   = "updated_rule_1"
  description            = "Terraform accept test"
  protocol               = "udp"
  action                 = "deny"
  ip_version             = 4
  source_ip_address      = "1.2.3.4"
  destination_ip_address = "4.3.2.0/24"
  source_port            = "444"
  destination_port       = "555"
  enabled                = true
}
`
