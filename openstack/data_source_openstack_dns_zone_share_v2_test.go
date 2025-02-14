package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDNSZoneShareV2_basic(t *testing.T) {
	zoneID := "dummy-zone-id"
	projectID := "dummy-target-project-id"
	sudoProjectID := "dummy-project-id"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: convertTestAccProviders(testAccProviders),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDNSZoneShareV2Config(zoneID, projectID, sudoProjectID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.openstack_dns_zone_share_v2.test", "zone_id", zoneID),
					resource.TestCheckResourceAttr("data.openstack_dns_zone_share_v2.test", "target_project_id", projectID),
					resource.TestCheckResourceAttr("data.openstack_dns_zone_share_v2.test", "project_id", sudoProjectID),
				),
			},
		},
	})
}

func testAccDataSourceDNSZoneShareV2Config(zoneID, projectID, sudoProjectID string) string {
	return fmt.Sprintf(`
data "openstack_dns_zone_share_v2" "test" {
  zone_id            = "%s"
  target_project_id  = "%s"
  project_id         = "%s"
}
`, zoneID, projectID, sudoProjectID)
}
