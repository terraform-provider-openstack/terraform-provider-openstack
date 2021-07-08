package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNetworkingV2QoSBandwidthLimitRuleDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2QoSBandwidthLimitRuleDataSource,
			},
			{
				Config: testAccOpenStackNetworkingQoSBandwidthLimitRuleV2DataSourceMaxKbps(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingQoSBandwidthLimitRuleV2DataSourceID("data.openstack_networking_qos_bandwidth_limit_rule_v2.bw_limit_rule_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_qos_bandwidth_limit_rule_v2.bw_limit_rule_1", "max_kbps", "3000"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_qos_bandwidth_limit_rule_v2.bw_limit_rule_1", "max_burst_kbps", "300"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_qos_bandwidth_limit_rule_v2.bw_limit_rule_1", "direction", "egress"),
				),
			},
			{
				Config: testAccOpenStackNetworkingQoSBandwidthLimitRuleV2DataSourceMaxBurstKbps(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingQoSBandwidthLimitRuleV2DataSourceID("data.openstack_networking_qos_bandwidth_limit_rule_v2.bw_limit_rule_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_qos_bandwidth_limit_rule_v2.bw_limit_rule_1", "max_kbps", "3000"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_qos_bandwidth_limit_rule_v2.bw_limit_rule_1", "max_burst_kbps", "300"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_qos_bandwidth_limit_rule_v2.bw_limit_rule_1", "direction", "egress"),
				),
			},
		},
	})
}

func testAccCheckNetworkingQoSBandwidthLimitRuleV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find QoS bw limit data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("QoS bw limit data source ID not set")
		}

		return nil
	}
}

const testAccNetworkingV2QoSBandwidthLimitRuleDataSource = `
resource "openstack_networking_qos_policy_v2" "qos_policy_1" {
  name = "qos_policy_1"
}

resource "openstack_networking_qos_bandwidth_limit_rule_v2" "bw_limit_rule_1" {
  qos_policy_id  = "${openstack_networking_qos_policy_v2.qos_policy_1.id}"
  max_kbps       = 3000
  max_burst_kbps = 300
  direction      = "egress"
}
`

func testAccOpenStackNetworkingQoSBandwidthLimitRuleV2DataSourceMaxKbps() string {
	return fmt.Sprintf(`
%s

data "openstack_networking_qos_bandwidth_limit_rule_v2" "bw_limit_rule_1" {
  qos_policy_id = "${openstack_networking_qos_policy_v2.qos_policy_1.id}"
  max_kbps      = "${openstack_networking_qos_bandwidth_limit_rule_v2.bw_limit_rule_1.max_kbps}"
}
`, testAccNetworkingV2QoSBandwidthLimitRuleDataSource)
}

func testAccOpenStackNetworkingQoSBandwidthLimitRuleV2DataSourceMaxBurstKbps() string {
	return fmt.Sprintf(`
%s

data "openstack_networking_qos_bandwidth_limit_rule_v2" "bw_limit_rule_1" {
  qos_policy_id  = "${openstack_networking_qos_policy_v2.qos_policy_1.id}"
  max_burst_kbps = "${openstack_networking_qos_bandwidth_limit_rule_v2.bw_limit_rule_1.max_burst_kbps}"
}
`, testAccNetworkingV2QoSBandwidthLimitRuleDataSource)
}
