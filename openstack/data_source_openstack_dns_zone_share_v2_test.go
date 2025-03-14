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
	targetProjectName := fmt.Sprintf("ACPTTEST-Target-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckAdminOnly(t)
			testAccPreCheckDNS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDNSZoneShareV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDNSZoneShareV2Config(zoneName, targetProjectName),
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
					resource.TestCheckResourceAttrPair("data.openstack_dns_zone_share_v2.zone_project_id_and_target", "target_project_id", "openstack_identity_project_v3.target", "id"),
				),
			},
		},
	})
}

func testAccDataSourceDNSZoneShareV2Config(zoneName, targetProjectName string) string {
	return fmt.Sprintf(`
resource "openstack_dns_zone_v2" "zone" {
  name  = "%s"
  email = "admin@example.com"
  ttl   = 3000
  type  = "PRIMARY"
}

resource "openstack_identity_project_v3" "target" {
  name        = "%s"
  description = "The target project with which we share the zone"
}

resource "openstack_dns_zone_share_v2" "share" {
  zone_id           = openstack_dns_zone_v2.zone.id
  target_project_id = openstack_identity_project_v3.target.id
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
  target_project_id = openstack_identity_project_v3.target.id
  
  depends_on = [ openstack_dns_zone_share_v2.share ]
}
`, zoneName, targetProjectName)
}

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
