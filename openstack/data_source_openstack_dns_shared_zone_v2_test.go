package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceDNSSharedZoneV2_basic(t *testing.T) {
	dataSourceName := "data.openstack_dns_shared_zone_v2.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDNSSharedZoneV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSSharedZoneV2DataSourceExists(dataSourceName),
					resource.TestCheckResourceAttr(dataSourceName, "shared_zones.#", "1"),
					resource.TestCheckResourceAttr(dataSourceName, "shared_zones.0.zone_id", "12345"),
					resource.TestCheckResourceAttr(dataSourceName, "shared_zones.0.project_id", "67890"),
				),
			},
		},
	})
}

func testAccCheckDNSSharedZoneV2DataSourceExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find DNS Shared Zone data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("DNS Shared Zone data source ID not set")
		}

		return nil
	}
}

var testAccDataSourceDNSSharedZoneV2Config = `
data "openstack_dns_shared_zone_v2" "test" {
  zone_id = "12345"
}
`
