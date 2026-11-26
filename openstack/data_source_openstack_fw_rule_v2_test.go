package openstack

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/fwaas_v2/rules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestUnitFilterFWRulesV2ByPort(t *testing.T) {
	all := []rules.Rule{
		{ID: "a", SourcePort: "9090", DestinationPort: "8080"},
		{ID: "b", SourcePort: "9091", DestinationPort: "8080"},
		{ID: "c", SourcePort: "9090", DestinationPort: "8081"},
		{ID: "d"},
	}

	tests := []struct {
		name     string
		raw      map[string]any
		expected []string
	}{
		{"no port filters", map[string]any{}, []string{"a", "b", "c", "d"}},
		{"source port", map[string]any{"source_port": "9090"}, []string{"a", "c"}},
		{"destination port", map[string]any{"destination_port": "8080"}, []string{"a", "b"}},
		{"both ports", map[string]any{"source_port": "9090", "destination_port": "8081"}, []string{"c"}},
		{"no match", map[string]any{"source_port": "1"}, []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := schema.TestResourceDataRaw(t, dataSourceFWRuleV2().Schema, tt.raw)

			got := filterFWRulesV2ByPort(all, d)

			ids := make([]string, 0, len(got))
			for _, r := range got {
				ids = append(ids, r.ID)
			}

			if fmt.Sprint(ids) != fmt.Sprint(tt.expected) {
				t.Fatalf("expected %v, got %v", tt.expected, ids)
			}
		})
	}
}

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

func TestAccFWRuleV2DataSource_ports(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckFW(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFWRuleV2DataSourcePortsRules,
			},
			{
				Config: testAccFWRuleV2DataSourceSourcePort(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.openstack_fw_rule_v2.by_source_port", "id",
						"openstack_fw_rule_v2.rule_1", "id"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_rule_v2.by_source_port", "source_port", "9090"),
				),
			},
			{
				Config: testAccFWRuleV2DataSourceDestinationPort(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.openstack_fw_rule_v2.by_destination_port", "id",
						"openstack_fw_rule_v2.rule_2", "id"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_rule_v2.by_destination_port", "destination_port", "8081"),
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
			return errors.New("Firewall rule data source ID not set")
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
  tenant_id              = openstack_fw_rule_v2.rule_1.tenant_id
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
  rule_id = openstack_fw_rule_v2.rule_1.id
}
`, testAccFWRuleV2DataSourceBasic)
}

// Two rules that differ only in their ports, so only the port filters can
// tell them apart.
const testAccFWRuleV2DataSourcePortsRules = `
resource "openstack_fw_rule_v2" "rule_1" {
  name             = "rule_ports"
  protocol         = "tcp"
  source_port      = "9090"
  destination_port = "8080"
}

resource "openstack_fw_rule_v2" "rule_2" {
  name             = "rule_ports"
  protocol         = "tcp"
  source_port      = "9091"
  destination_port = "8081"
}
`

func testAccFWRuleV2DataSourceSourcePort() string {
	return fmt.Sprintf(`
%s

data "openstack_fw_rule_v2" "by_source_port" {
  name        = "rule_ports"
  tenant_id   = openstack_fw_rule_v2.rule_1.tenant_id
  source_port = "9090"
}
`, testAccFWRuleV2DataSourcePortsRules)
}

func testAccFWRuleV2DataSourceDestinationPort() string {
	return fmt.Sprintf(`
%s

data "openstack_fw_rule_v2" "by_destination_port" {
  name             = "rule_ports"
  tenant_id        = openstack_fw_rule_v2.rule_2.tenant_id
  destination_port = "8081"
}
`, testAccFWRuleV2DataSourcePortsRules)
}
