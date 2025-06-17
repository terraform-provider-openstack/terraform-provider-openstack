package openstack

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/zones"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceDNSZoneShareV2_basic(t *testing.T) {
	zoneName := fmt.Sprintf("ACPTTEST%s.com.", acctest.RandString(5))
	targetProjectName := "ACPTTEST-Target-" + acctest.RandString(5)

	var zone zones.Zone

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckAdminOnly(t)
			testAccPreCheckDNS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDNSZoneShareV2Destroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDNSZoneShareV2Config(zoneName, targetProjectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2ZoneExists(t.Context(), "openstack_dns_zone_v2.zone", &zone),
					resource.TestCheckResourceAttrPtr("data.openstack_dns_zone_share_v2.zone_only", "zone_id", &zone.ID),

					resource.TestCheckResourceAttrPtr("data.openstack_dns_zone_share_v2.zone_and_project_id", "zone_id", &zone.ID),
					resource.TestCheckResourceAttrPair("data.openstack_dns_zone_share_v2.zone_and_project_id", "project_id", "openstack_dns_zone_v2.zone", "project_id"),

					resource.TestCheckResourceAttrPtr("data.openstack_dns_zone_share_v2.zone_project_id_and_target", "zone_id", &zone.ID),
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
  share_id = openstack_dns_zone_share_v2.share.id
}

data "openstack_dns_zone_share_v2" "zone_and_project_id" {
  zone_id    = openstack_dns_zone_v2.zone.id
  project_id = openstack_dns_zone_v2.zone.project_id

  share_id = openstack_dns_zone_share_v2.share.id
}

data "openstack_dns_zone_share_v2" "zone_project_id_and_target" {
  zone_id           = openstack_dns_zone_v2.zone.id
  project_id        = openstack_dns_zone_v2.zone.project_id
  target_project_id = openstack_identity_project_v3.target.id

  share_id = openstack_dns_zone_share_v2.share.id
}
`, zoneName, targetProjectName)
}
