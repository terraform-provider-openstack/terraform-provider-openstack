package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/zones"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceDNSZoneShareV2_basic(t *testing.T) {
	zoneName := fmt.Sprintf("ACPTTEST%s.com.", acctest.RandString(5))
	targetProjectName := "ACPTTEST-Target-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckAdminOnly(t)
			testAccPreCheckDNS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDNSZoneShareV2Destroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDNSZoneShareV2Config(zoneName, targetProjectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSZoneShareV2Exists(t.Context(), "openstack_dns_zone_share_v2.share"),
					resource.TestCheckResourceAttr("openstack_dns_zone_v2.zone", "name", zoneName),
					resource.TestCheckResourceAttrPair("openstack_dns_zone_share_v2.share", "target_project_id", "openstack_identity_project_v3.target", "id"),
				),
			},
		},
	})
}

func testAccCheckDNSZoneShareV2Exists(ctx context.Context, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("DNS zone share resource not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("DNS zone share resource has no ID set")
		}

		config := testAccProvider.Meta().(*Config)

		dnsClient, err := config.DNSV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating DNS client: %w", err)
		}

		zoneID := rs.Primary.Attributes["zone_id"]
		shareID := rs.Primary.ID
		dnsClient.MoreHeaders = map[string]string{
			headerAuthSudoTenantID: rs.Primary.Attributes["project_id"],
		}

		_, err = zones.GetShare(ctx, dnsClient, zoneID, shareID).Extract()
		if err != nil {
			return fmt.Errorf("Error getting DNS zone share: %w", err)
		}

		return nil
	}
}

func testAccCheckDNSZoneShareV2Destroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		dnsClient, err := config.DNSV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating DNS client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			// Skip data sources and resources with no zone_id.
			if rs.Type == "data.openstack_dns_zone_share_v2" || rs.Primary.Attributes["zone_id"] == "" {
				continue
			}

			zoneID := rs.Primary.Attributes["zone_id"]
			shareID := rs.Primary.ID
			dnsClient.MoreHeaders = map[string]string{
				headerAuthSudoTenantID: rs.Primary.Attributes["project_id"],
			}

			_, err = zones.GetShare(ctx, dnsClient, zoneID, shareID).Extract()
			if err == nil {
				return fmt.Errorf("DNS zone share still exists: %s", shareID)
			}

			return nil
		}

		return nil
	}
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
