package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccFWPolicyV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckFW(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFWPolicyV2DataSourceBasic,
			},
			{
				Config: testAccFWPolicyV2DataSourceName(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWPolicyV2DataSourceID("data.openstack_fw_policy_v2.policy_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_policy_v2.policy_1", "name", "policy_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_policy_v2.policy_1", "description", "My firewall policy"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_policy_v2.policy_1", "shared", "false"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_policy_v2.policy_1", "audited", "false"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_fw_policy_v2.policy_1", "tenant_id"),
				),
			},
		},
	})
}

func TestAccFWPolicyV2DataSource_shared(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
			testAccPreCheckFW(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFWPolicyV2DataSourceBasic,
			},
			{
				Config: testAccFWPolicyV2DataSourceShared,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWPolicyV2DataSourceID("data.openstack_fw_policy_v2.policy_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_policy_v2.policy_1", "name", "policy_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_policy_v2.policy_1", "description", "My firewall policy"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_policy_v2.policy_1", "shared", "true"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_policy_v2.policy_1", "audited", "true"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_fw_policy_v2.policy_1", "tenant_id"),
				),
			},
		},
	})
}

func TestAccFWPolicyV2DataSource_FWPolicyID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckFW(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFWPolicyV2DataSourceBasic,
			},
			{
				Config: testAccFWPolicyV2DataSourcePolicyID(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWPolicyV2DataSourceID("data.openstack_fw_policy_v2.policy_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_policy_v2.policy_1", "name", "policy_1"),
				),
			},
		},
	})
}

func testAccCheckFWPolicyV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find firewall policy data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Firewall policy data source ID not set")
		}

		return nil
	}
}

const testAccFWPolicyV2DataSourceBasic = `
resource "openstack_fw_policy_v2" "policy_1" {
  name        = "policy_1"
  description = "My firewall policy"
}
`

const testAccFWPolicyV2DataSourceShared = `
resource "openstack_fw_policy_v2" "policy_1" {
  name        = "policy_1"
  description = "My firewall policy"
  shared      = true
  audited     = true
}

data "openstack_fw_policy_v2" "policy_1" {
  name        = "${openstack_fw_policy_v2.policy_1.name}"
  description = "${openstack_fw_policy_v2.policy_1.description}"
  tenant_id   = "${openstack_fw_policy_v2.policy_1.tenant_id}"
  shared      = true
  audited     = true
}
`

func testAccFWPolicyV2DataSourceName() string {
	return fmt.Sprintf(`
%s

data "openstack_fw_policy_v2" "policy_1" {
  name        = "policy_1"
  description = "My firewall policy"
  tenant_id   = "${openstack_fw_policy_v2.policy_1.tenant_id}"
  shared      = false
  audited     = false
}
`, testAccFWPolicyV2DataSourceBasic)
}

func testAccFWPolicyV2DataSourcePolicyID() string {
	return fmt.Sprintf(`
%s

data "openstack_fw_policy_v2" "policy_1" {
  policy_id = "${openstack_fw_policy_v2.policy_1.id}"
}
`, testAccFWPolicyV2DataSourceBasic)
}
