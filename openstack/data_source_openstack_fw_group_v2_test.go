package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccFWGroupV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckFW(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFWGroupV2DataSourceGroup,
			},
			{
				Config: testAccFWGroupV2DataSourceBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWGroupV2DataSourceID("data.openstack_fw_group_v2.group_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_group_v2.group_1", "name", "group_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_group_v2.group_1", "description", "My firewall group"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_group_v2.group_1", "shared", "false"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_fw_group_v2.group_1", "ingress_firewall_policy_id"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_fw_group_v2.group_1", "egress_firewall_policy_id"),
				),
			},
		},
	})
}

func TestAccFWGroupV2DataSource_shared(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
			testAccPreCheckFW(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFWGroupV2DataSourceGroupShared,
			},
			{
				Config: testAccFWGroupV2DataSourceShared(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWGroupV2DataSourceID("data.openstack_fw_group_v2.group_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_group_v2.group_1", "name", "group_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_group_v2.group_1", "description", "My firewall group"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_group_v2.group_1", "shared", "true"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_fw_group_v2.group_1", "ingress_firewall_policy_id"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_fw_group_v2.group_1", "egress_firewall_policy_id"),
				),
			},
		},
	})
}

func TestAccFWGroupV2DataSource_FWGroupID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckFW(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFWGroupV2DataSourceGroup,
			},
			{
				Config: testAccFWGroupV2DataSourceGroupID(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWGroupV2DataSourceID("data.openstack_fw_group_v2.group_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_group_v2.group_1", "name", "group_1"),
				),
			},
		},
	})
}

func testAccCheckFWGroupV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find firewall group data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Firewall group data source ID not set")
		}

		return nil
	}
}

const testAccFWGroupV2DataSourceGroup = `
resource "openstack_fw_policy_v2" "policy_1" {
  name        = "policy_1"
  description = "My firewall ingress policy"
}

resource "openstack_fw_policy_v2" "policy_2" {
  name        = "policy_2"
  description = "My firewall egress policy"
}

resource "openstack_fw_group_v2" "group_1" {
  name                       = "group_1"
  description                = "My firewall group"
  ingress_firewall_policy_id = "${openstack_fw_policy_v2.policy_1.id}"
  egress_firewall_policy_id  = "${openstack_fw_policy_v2.policy_2.id}"
}
`

func testAccFWGroupV2DataSourceBasic() string {
	return fmt.Sprintf(`
%s

data "openstack_fw_group_v2" "group_1" {
  name                       = "group_1"
  description                = "My firewall group"
  project_id                 = "${openstack_fw_group_v2.group_1.project_id}"
  shared                     = false
  ingress_firewall_policy_id = "${openstack_fw_policy_v2.policy_1.id}"
  egress_firewall_policy_id  = "${openstack_fw_policy_v2.policy_2.id}"
}
`, testAccFWGroupV2DataSourceGroup)
}

const testAccFWGroupV2DataSourceGroupShared = `
resource "openstack_fw_policy_v2" "policy_1" {
  name        = "policy_1"
  description = "My firewall ingress policy"
}

resource "openstack_fw_policy_v2" "policy_2" {
  name        = "policy_2"
  description = "My firewall egress policy"
}

resource "openstack_fw_group_v2" "group_1" {
  name                       = "group_1"
  description                = "My firewall group"
  shared                     = true
  ingress_firewall_policy_id = "${openstack_fw_policy_v2.policy_1.id}"
  egress_firewall_policy_id  = "${openstack_fw_policy_v2.policy_2.id}"
}
`

func testAccFWGroupV2DataSourceShared() string {
	return fmt.Sprintf(`
%s

data "openstack_fw_group_v2" "group_1" {
  name                       = "group_1"
  description                = "My firewall group"
  tenant_id                  = "${openstack_fw_group_v2.group_1.tenant_id}"
  shared                     = true
  ingress_firewall_policy_id = "${openstack_fw_policy_v2.policy_1.id}"
  egress_firewall_policy_id  = "${openstack_fw_policy_v2.policy_2.id}"
}
`, testAccFWGroupV2DataSourceGroupShared)
}

func testAccFWGroupV2DataSourceGroupID() string {
	return fmt.Sprintf(`
%s

data "openstack_fw_group_v2" "group_1" {
  group_id = "${openstack_fw_group_v2.group_1.id}"
}
`, testAccFWGroupV2DataSourceGroup)
}
