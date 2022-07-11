package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccFirewallPolicyV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallPolicyV1DataSourceGroup,
			},
			{
				Config: testAccFirewallPolicyV1DataSourceBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallPolicyV1DataSourceID("data.openstack_fw_policy_v1.policy_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_policy_v1.policy_1", "name", "policy_1"),
				),
			},
		},
	})
}
func TestAccFirewallPolicyV1DataSource_FWPolicyID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccFirewallPolicyV1DataSourceGroup,
			},
			{
				Config: testAccFirewallPolicyV1DataSourcePolicyID(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFirewallPolicyV1DataSourceID("data.openstack_fw_policy_v1.policy_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_policy_v1.policy_1", "name", "policy_1"),
				),
			},
		},
	})
}

func testAccCheckFirewallPolicyV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find firewall policy data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("firewall policy data source ID not set")
		}

		return nil
	}
}

const testAccFirewallPolicyV1DataSourceGroup = `
resource "openstack_fw_policy_v1" "policy_1" {
    name        = "policy_1"
	description = "My firewall policy"
}
`

func testAccFirewallPolicyV1DataSourceBasic() string {
	return fmt.Sprintf(`
%s

data "openstack_fw_policy_v1" "policy_1" {
	name = "${openstack_fw_policy_v1.policy_1.name}"
}
`, testAccFirewallPolicyV1DataSourceGroup)
}

func testAccFirewallPolicyV1DataSourcePolicyID() string {
	return fmt.Sprintf(`
%s

data "openstack_fw_policy_v1" "policy_1" {
	policy_id = "${openstack_fw_policy_v1.policy_1.id}"
}
`, testAccFirewallPolicyV1DataSourceGroup)
}
