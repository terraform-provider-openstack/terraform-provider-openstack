package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDnsZoneShare_basic(t *testing.T) {
	zoneID := "test-zone-id"
	targetProjectID := "test-project-id"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDnsZoneShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDnsZoneShareConfig(zoneID, targetProjectID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openstack_dns_zone_share.test", "zone_id", zoneID),
					resource.TestCheckResourceAttr("openstack_dns_zone_share.test", "target_project_id", targetProjectID),
				),
			},
			{
				ResourceName:      "openstack_dns_zone_share.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccDnsZoneShareConfig(zoneID, targetProjectID string) string {
	return fmt.Sprintf(`
resource "openstack_dns_zone_share" "test" {
	zone_id           = "%s"
	target_project_id = "%s"
}
`, zoneID, targetProjectID)
}

func testAccCheckDnsZoneShareDestroy(s *terraform.State) error {
	// No API exists to verify share deletion.
	return nil
}
