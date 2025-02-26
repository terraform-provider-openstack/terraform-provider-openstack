package openstack

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceDNSZoneShareV2_basic(t *testing.T) {
	zoneName := fmt.Sprintf("ACPTTEST%s.com.", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckDNS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDNSZoneShareV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDNSZoneShareV2Config(zoneName),
				Check: resource.ComposeTestCheckFunc(
					// Check data source with only zone_id provided.
					testAccCheckDNSZoneShareV2DataSourceHasShares("data.openstack_dns_zone_share_v2.zone_only"),
					resource.TestMatchResourceAttr("data.openstack_dns_zone_share_v2.zone_only", "zone_id", regexp.MustCompile(`^ACPTTEST.+\.com\.$`)),

					// Check data source with zone_id and project_id.
					testAccCheckDNSZoneShareV2DataSourceHasShares("data.openstack_dns_zone_share_v2.zone_and_project_id"),
					resource.TestMatchResourceAttr("data.openstack_dns_zone_share_v2.zone_and_project_id", "zone_id", regexp.MustCompile(`^ACPTTEST.+\.com\.$`)),
					resource.TestCheckResourceAttrPair("data.openstack_dns_zone_share_v2.zone_and_project_id", "project_id", "openstack_dns_zone_v2.zone", "project_id"),

					// Check data source with zone_id, project_id and target_project_id.
					testAccCheckDNSZoneShareV2DataSourceHasShares("data.openstack_dns_zone_share_v2.zone_project_id_and_target"),
					resource.TestMatchResourceAttr("data.openstack_dns_zone_share_v2.zone_project_id_and_target", "zone_id", regexp.MustCompile(`^ACPTTEST.+\.com\.$`)),
					resource.TestCheckResourceAttrPair("data.openstack_dns_zone_share_v2.zone_project_id_and_target", "project_id", "openstack_dns_zone_v2.zone", "project_id"),
					testAccCheckTargetProjectNotAdmin("data.openstack_dns_zone_share_v2.zone_project_id_and_target"),
				),
			},
		},
	})
}

// testAccDataSourceDNSZoneShareV2Config returns a configuration that:
//   - Retrieves all project IDs,
//   - Creates a DNS zone,
//   - Uses a local to select the first project that is not the zone’s project,
//   - Shares the DNS zone with that alternate project,
//   - And defines three data sources to read the share.
func testAccDataSourceDNSZoneShareV2Config(zoneName string) string {
	return fmt.Sprintf(`
data "openstack_identity_project_ids_v3" "projects" {
  name_regex = ".*"
}

resource "openstack_dns_zone_v2" "zone" {
  name  = "%s"
  email = "admin@example.com"
  ttl   = 3000
  type  = "PRIMARY"
}

locals {
  # Exclude the admin (zone's) project ID from the list.
  non_admin_projects = [
    for id in data.openstack_identity_project_ids_v3.projects.ids : id
    if id != openstack_dns_zone_v2.zone.project_id
  ]
  alternate_project = length(local.non_admin_projects) > 0 ? local.non_admin_projects[0] : ""
}

resource "openstack_dns_zone_share_v2" "share" {
  zone_id           = openstack_dns_zone_v2.zone.id
  target_project_id = local.alternate_project
}

data "openstack_dns_zone_share_v2" "zone_only" {
  zone_id = openstack_dns_zone_v2.zone.id

  depends_on = [ openstack_dns_zone_share_v2.share ]
}

data "openstack_dns_zone_share_v2" "zone_and_project_id" {
  zone_id    = openstack_dns_zone_v2.zone.id
  project_id = openstack_dns_zone_v2.zone.project_id

  depends_on = [ openstack_dns_zone_share_v2.share ]
}

data "openstack_dns_zone_share_v2" "zone_project_id_and_target" {
  zone_id           = openstack_dns_zone_v2.zone.id
  project_id        = openstack_dns_zone_v2.zone.project_id
  target_project_id = local.alternate_project
  
  depends_on = [ openstack_dns_zone_share_v2.share ]
}
`, zoneName)
}

// testAccCheckDNSZoneShareV2DataSourceHasShares checks that the data source returns at least one share.
func testAccCheckDNSZoneShareV2DataSourceHasShares(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("DNS zone share data source not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("DNS zone share data source has no ID set")
		}
		sharesCount, ok := rs.Primary.Attributes["shares.#"]
		if !ok {
			return fmt.Errorf("no shares attribute found in data source %s", n)
		}
		if sharesCount == "0" {
			return fmt.Errorf("no shares returned in data source %s", n)
		}
		return nil
	}
}

// testAccCheckTargetProjectNotAdmin checks that the target_project_id in the data source
// is not empty and does not match the zone's project_id.
func testAccCheckTargetProjectNotAdmin(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("data source %s not found", n)
		}
		adminProj, ok := s.RootModule().Resources["openstack_dns_zone_v2.zone"]
		if !ok {
			return fmt.Errorf("DNS zone resource not found")
		}
		adminProjectID := adminProj.Primary.Attributes["project_id"]
		targetProjectID := ds.Primary.Attributes["target_project_id"]

		if targetProjectID == "" {
			return fmt.Errorf("target_project_id is empty in data source %s", n)
		}
		if targetProjectID == adminProjectID {
			return fmt.Errorf("target_project_id (%s) should not equal admin project id (%s)", targetProjectID, adminProjectID)
		}
		return nil
	}
}
