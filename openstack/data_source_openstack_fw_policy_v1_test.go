package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOpenStackNetworkingFWPolicyV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingFWPolicyV1DataSourceGroup,
			},
			{
				Config: testAccOpenStackNetworkingFWPolicyV1DataSourceBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingFWPolicyV1DataSourceID("data.openstack_fw_policy_v1.policy_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_policy_v1.policy_1", "name", "policy_1"),
				),
			},
		},
	})
}
func TestAccOpenStackNetworkingFWPolicyV1DataSource_FWPolicyID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingFWPolicyV1DataSourceGroup,
			},
			{
				Config: testAccOpenStackNetworkingFWPolicyV1DataSourcePolicyID(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingFWPolicyV1DataSourceID("data.openstack_fw_policy_v1.policy_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_policy_v1.policy_1", "name", "policy_1"),
				),
			},
		},
	})
}

func testAccCheckNetworkingFWPolicyV1DataSourceID(n string) resource.TestCheckFunc {
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

const testAccOpenStackNetworkingFWPolicyV1DataSourceGroup = `
resource "openstack_fw_policy_v1" "policy_1" {
        name        = "policy_1"
	description = "My firewall policy"
}
`

func testAccOpenStackNetworkingFWPolicyV1DataSourceBasic() string {
	return fmt.Sprintf(`
%s

data "openstack_fw_policy_v1" "policy_1" {
	name = "${openstack_fw_policy_v1.policy_1.name}"
}
`, testAccOpenStackNetworkingFWPolicyV1DataSourceGroup)
}

func testAccOpenStackNetworkingFWPolicyV1DataSourcePolicyID() string {
	return fmt.Sprintf(`
%s

data "openstack_fw_policy_v1" "policy_1" {
	policy_id = "${openstack_fw_policy_v1.policy_1.id}"
}
`, testAccOpenStackNetworkingFWPolicyV1DataSourceGroup)
}
