package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNetworkingV2QoSPolicyDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2QoSPolicyDataSource,
			},
			{
				Config: testAccOpenStackNetworkingQoSPolicyV2DataSourceBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingQoSPolicyV2DataSourceID("data.openstack_networking_qos_policy_v2.qos_policy_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_qos_policy_v2.qos_policy_1", "name", "qos_policy_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_qos_policy_v2.qos_policy_1", "description", "terraform qos policy acceptance test"),
				),
			},
		},
	})
}

func TestAccNetworkingV2QoSPolicyDataSource_description(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2QoSPolicyDataSource,
			},
			{
				Config: testAccOpenStackNetworkingQoSPolicyV2DataSourceDescription(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingQoSPolicyV2DataSourceID("data.openstack_networking_qos_policy_v2.qos_policy_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_qos_policy_v2.qos_policy_1", "name", "qos_policy_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_qos_policy_v2.qos_policy_1", "description", "terraform qos policy acceptance test"),
				),
			},
		},
	})
}

func testAccCheckNetworkingQoSPolicyV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find QoS policy data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("QoS policy data source ID not set")
		}

		return nil
	}
}

const testAccNetworkingV2QoSPolicyDataSource = `
resource "openstack_networking_qos_policy_v2" "qos_policy_1" {
  name        = "qos_policy_1"
  description = "terraform qos policy acceptance test"
}
`

func testAccOpenStackNetworkingQoSPolicyV2DataSourceBasic() string {
	return fmt.Sprintf(`
%s

data "openstack_networking_qos_policy_v2" "qos_policy_1" {
  name = "${openstack_networking_qos_policy_v2.qos_policy_1.name}"
}
`, testAccNetworkingV2QoSPolicyDataSource)
}

func testAccOpenStackNetworkingQoSPolicyV2DataSourceDescription() string {
	return fmt.Sprintf(`
%s

data "openstack_networking_qos_policy_v2" "qos_policy_1" {
  description = "${openstack_networking_qos_policy_v2.qos_policy_1.description}"
}
`, testAccNetworkingV2QoSPolicyDataSource)
}
