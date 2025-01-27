package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceDNSSharedZoneV2_basic(t *testing.T) {
	resourceName := "openstack_dns_shared_zone_v2.test"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDNSSharedZoneV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSSharedZoneV2Exists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "zone_id", "12345"),
					resource.TestCheckResourceAttr(resourceName, "project_id", "67890"),
				),
			},
		},
	})
}

func testAccCheckDNSSharedZoneV2Exists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find DNS Shared Zone resource: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("DNS Shared Zone resource ID not set")
		}

		return nil
	}
}

var testAccResourceDNSSharedZoneV2Config = `
resource "openstack_dns_shared_zone_v2" "test" {
  zone_id    = "12345"
  project_id = "67890"
}
`
