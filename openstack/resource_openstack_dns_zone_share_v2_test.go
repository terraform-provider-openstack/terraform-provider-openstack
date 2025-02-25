package openstack

import (
	"context"
	"fmt"
	"regexp"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/zones"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceDNSZoneShareV2_basic(t *testing.T) {
	var zone zones.Zone

	zoneName := fmt.Sprintf("ACPTTEST%s.com.", acctest.RandString(5))
	targetProjectName := fmt.Sprintf("ACPTTEST-Target-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckDNS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDNSZoneShareV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDNSZoneShareV2Config(zoneName, targetProjectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSZoneShareV2Exists("openstack_dns_zone_share_v2.share", &zone),
					resource.TestMatchResourceAttr("openstack_dns_zone_v2.zone", "name", regexp.MustCompile(`^ACPTTEST.+\.com\.$`)),
					resource.TestCheckResourceAttrPair("openstack_dns_zone_share_v2.share", "target_project_id", "openstack_identity_project_v3.target", "id"),
				),
			},
		},
	})
}

func testAccResourceDNSZoneShareV2Config(zoneName, targetProjectName string) string {
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
`, zoneName, targetProjectName)
}

func testAccCheckDNSZoneShareV2Exists(n string, zone *zones.Zone) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("DNS zone share resource not found: %s", n)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("DNS zone share resource has no ID set")
		}

		config := testAccProvider.Meta().(*Config)
		dnsClient, err := config.DNSV2Client(context.Background(), osRegionName)
		if err != nil {
			return fmt.Errorf("error creating DNS client: %s", err)
		}

		zoneID := rs.Primary.Attributes["zone_id"]
		_, shareID, err := parseDnsSharedZoneID(rs.Primary.ID)
		if err != nil {
			return err
		}
		ownerProjID := rs.Primary.Attributes["project_id"]

		shares, err := listZoneShares(context.Background(), dnsClient, zoneID, ownerProjID)
		if err != nil {
			return fmt.Errorf("error listing DNS zone shares: %s", err)
		}
		for _, s := range shares {
			if s.ID == shareID {
				return nil
			}
		}
		return fmt.Errorf("DNS zone share %s not found", shareID)
	}
}

func testAccCheckDNSZoneShareV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	dnsClient, err := config.DNSV2Client(context.Background(), osRegionName)
	if err != nil {
		return fmt.Errorf("error creating DNS client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		// Skip data sources and resources with no zone_id.
		if rs.Type == "data.openstack_dns_zone_share_v2" || rs.Primary.Attributes["zone_id"] == "" {
			continue
		}
		zoneID, shareID, err := parseDnsSharedZoneID(rs.Primary.ID)
		if err != nil {
			return err
		}
		// If zoneID is empty, skip the check.
		if zoneID == "" {
			continue
		}
		ownerProjID := rs.Primary.Attributes["project_id"]
		shares, err := listZoneShares(context.Background(), dnsClient, zoneID, ownerProjID)
		if err != nil {
			return fmt.Errorf("error listing DNS zone shares: %s", err)
		}
		for _, s := range shares {
			if s.ID == shareID {
				return fmt.Errorf("DNS zone share still exists: %s", shareID)
			}
		}
	}
	return nil
}
