package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccFWRuleV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckFW(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFWRuleV2DataSourceBasic,
			},
			{
				Config: testAccFWRuleV2DataSourceName(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWRuleV2DataSourceID("data.openstack_fw_rule_v2.rule_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_rule_v2.rule_1", "name", "rule_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_rule_v2.rule_1", "protocol", "tcp"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_rule_v2.rule_1", "action", "deny"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_rule_v2.rule_1", "ip_version", "4"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_rule_v2.rule_1", "shared", "false"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_rule_v2.rule_1", "enabled", "true"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_fw_rule_v2.rule_1", "tenant_id"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_fw_rule_v2.rule_1", "source_ip_address"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_fw_rule_v2.rule_1", "source_port"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_fw_rule_v2.rule_1", "destination_ip_address"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_fw_rule_v2.rule_1", "destination_port"),
				),
			},
		},
	})
}

func TestAccFWRuleV2DataSource_FWRuleID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckFW(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFWRuleV2DataSourceBasic,
			},
			{
				Config: testAccFWRuleV2DataSourceRuleID(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWRuleV2DataSourceID("data.openstack_fw_rule_v2.rule_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_rule_v2.rule_1", "name", "rule_1"),
				),
			},
		},
	})
}

func testAccCheckFWRuleV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find firewall rule data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Firewall rule data source ID not set")
		}

		return nil
	}
}

const testAccFWRuleV2DataSourceBasic = `
resource "openstack_fw_rule_v2" "rule_1" {
  name                   = "rule_1"
  description            = "My firewall rule"
  protocol			     = "tcp"
  source_ip_address      = "10.20.30.40/32"
  source_port            = "9090"
  destination_ip_address = "10.11.12.13/32"
  destination_port       = "8080"
}
`

func testAccFWRuleV2DataSourceName() string {
	return fmt.Sprintf(`
%s

data "openstack_fw_rule_v2" "rule_1" {
  name                   = "rule_1"
  description            = "My firewall rule"
  tenant_id              = "${openstack_fw_rule_v2.rule_1.tenant_id}"
  protocol               = "TCP"
  action                 = "deny"
  ip_version             = 4
  source_ip_address      = "10.20.30.40/32"
  source_port            = "9090"
  destination_ip_address = "10.11.12.13/32"
  destination_port       = "8080"
  shared                 = false
  enabled                = false
}
`, testAccFWRuleV2DataSourceBasic)
}

func testAccFWRuleV2DataSourceRuleID() string {
	return fmt.Sprintf(`
%s

data "openstack_fw_rule_v2" "rule_1" {
  rule_id = "${openstack_fw_rule_v2.rule_1.id}"
}
`, testAccFWRuleV2DataSourceBasic)
}
