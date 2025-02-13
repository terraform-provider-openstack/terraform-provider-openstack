package openstack

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceDNSZoneShareV2_basic(t *testing.T) {
	zoneID := "dummy-zone-id"
	targetProjectID := "dummy-target-project-id"
	projectID := "dummy-project-id"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: convertTestAccProviders(testAccProviders),
		CheckDestroy: func(state *terraform.State) error {
			return testAccCheckDNSZoneShareV2Destroy(state, zoneID, targetProjectID, projectID)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDNSZoneShareV2Config(zoneID, targetProjectID, projectID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openstack_dns_shared_zone_v2.test", "zone_id", zoneID),
					resource.TestCheckResourceAttr("openstack_dns_shared_zone_v2.test", "target_project_id", targetProjectID),
					resource.TestCheckResourceAttr("openstack_dns_shared_zone_v2.test", "project_id", projectID),
				),
			},
		},
	})
}

func testAccCheckDNSZoneShareV2Destroy(state *terraform.State, zoneID, targetProjectID, projectID string) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.DNSV2Client(context.Background(), "")
	if err != nil {
		return fmt.Errorf("error creating OpenStack DNS client: %s", err)
	}

	shares, err := listZoneShares(client, zoneID, projectID)
	if err != nil {
		return fmt.Errorf("error listing shares for DNS zone %s: %s", zoneID, err)
	}

	for _, s := range shares {
		if s.TargetProjectID == targetProjectID {
			return fmt.Errorf("DNS zone share still exists for zone %s and target project %s", zoneID, targetProjectID)
		}
	}

	return nil
}

func convertTestAccProviders(providers map[string]func() (*schema.Provider, error)) map[string]*schema.Provider {
	converted := make(map[string]*schema.Provider)
	for key, providerFunc := range providers {
		if provider, err := providerFunc(); err == nil {
			converted[key] = provider
		}
	}
	return converted
}

func testAccResourceDNSZoneShareV2Config(zoneID, targetProjectID, projectID string) string {
	return fmt.Sprintf(`
resource "openstack_dns_shared_zone_v2" "test" {
  zone_id            = "%s"
  target_project_id  = "%s"
  project_id         = "%s"
}
`, zoneID, targetProjectID, projectID)
}
