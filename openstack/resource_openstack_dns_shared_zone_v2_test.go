package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceDNSZoneShareV2_basic(t *testing.T) {
	zoneID := "dummy-zone-id"
	projectID := "dummy-project-id"
	sudoProjectID := "dummy-sudo-project-id"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: convertTestAccProviders(testAccProviders),
		CheckDestroy: func(state *terraform.State) error {
			// Implement a check to ensure the resource is destroyed
			return nil
		},
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDNSZoneShareV2Config(zoneID, projectID, sudoProjectID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openstack_dns_shared_zone_v2.test", "zone_id", zoneID),
					resource.TestCheckResourceAttr("openstack_dns_shared_zone_v2.test", "project_id", projectID),
					resource.TestCheckResourceAttr("openstack_dns_shared_zone_v2.test", "x_auth_sudo_project_id", sudoProjectID),
				),
			},
		},
	})
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

func testAccResourceDNSZoneShareV2Config(zoneID, projectID, sudoProjectID string) string {
	return fmt.Sprintf(`
resource "openstack_dns_shared_zone_v2" "test" {
  zone_id                 = "%s"
  project_id              = "%s"
  x_auth_sudo_project_id  = "%s"
}
`, zoneID, projectID, sudoProjectID)
}
