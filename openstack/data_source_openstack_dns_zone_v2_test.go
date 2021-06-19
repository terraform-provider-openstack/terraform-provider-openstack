package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func zoneName() string {
	return fmt.Sprintf("ACPTTEST%s.com.", acctest.RandString(5))
}

func TestAccOpenStackDNSZoneV2DataSource_basic(t *testing.T) {
	zoneName := zoneName()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckDNS(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackDNSZoneV2DataSourceZone(zoneName),
			},
			{
				Config: testAccOpenStackDNSZoneV2DataSourceBasic(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSZoneV2DataSourceID("data.openstack_dns_zone_v2.z1"),
					resource.TestCheckResourceAttr(
						"data.openstack_dns_zone_v2.z1", "name", zoneName),
					resource.TestCheckResourceAttr(
						"data.openstack_dns_zone_v2.z1", "type", "PRIMARY"),
					resource.TestCheckResourceAttr(
						"data.openstack_dns_zone_v2.z1", "ttl", "7200"),
				),
			},
		},
	})
}

func testAccCheckDNSZoneV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find DNS Zone data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("DNS Zone data source ID not set")
		}

		return nil
	}
}

func testAccOpenStackDNSZoneV2DataSourceZone(zoneName string) string {
	return fmt.Sprintf(`
resource "openstack_dns_zone_v2" "z1" {
  name = "%s"
  email = "terraform-dns-zone-v2-test-name@example.com"
  type = "PRIMARY"
  ttl = 7200
}`, zoneName)
}

func testAccOpenStackDNSZoneV2DataSourceBasic(zoneName string) string {
	return fmt.Sprintf(`
%s
data "openstack_dns_zone_v2" "z1" {
	name = "${openstack_dns_zone_v2.z1.name}"
}
`, testAccOpenStackDNSZoneV2DataSourceZone(zoneName))
}
