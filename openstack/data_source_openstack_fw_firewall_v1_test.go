package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOpenStackNetworkingFWFirewallV1DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingFWFirewallV1DataSource_group,
			},
			{
				Config: testAccOpenStackNetworkingFWFirewallV1DataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingFWFirewallV1DataSourceID("data.openstack_fw_firewall_v1.firewall_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_firewall_v1.firewall_1", "name", "firewall_1"),
				),
			},
		},
	})
}
func TestAccOpenStackNetworkingFWFirewallV1DataSource_FWFirewallID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingFWFirewallV1DataSource_group,
			},
			{
				Config: testAccOpenStackNetworkingFWFirewallV1DataSource_ID,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingFWFirewallV1DataSourceID("data.openstack_fw_firewall_v1.firewall_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_fw_firewall_v1.firewall_1", "name", "firewall_1"),
				),
			},
		},
	})
}

func testAccCheckNetworkingFWFirewallV1DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find firewall firewall_v1 data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("firewall firewall_v1 data source ID not set")
		}

		return nil
	}
}

const testAccOpenStackNetworkingFWFirewallV1DataSource_group = `
resource "openstack_fw_firewall_v1" "firewall_1" {
    name        = "firewall_1"
}
`

var testAccOpenStackNetworkingFWFirewallV1DataSource_basic = fmt.Sprintf(`
%s

data "openstack_fw_firewall_v1" "firewall_1" {
	name = "${openstack_fw_firewall_v1.firewall_1.name}"
}
`, testAccOpenStackNetworkingFWFirewallV1DataSource_group)

var testAccOpenStackNetworkingFWFirewallV1DataSource_ID = fmt.Sprintf(`
%s

data "openstack_fw_firewall_v1" "firewall_1" {
	id = "${openstack_fw_firewall_v1.firewall_1.id}"
}
`, testAccOpenStackNetworkingFWFirewallV1DataSource_group)
