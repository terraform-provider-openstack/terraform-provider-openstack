package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccResourceDNSZoneShareV2_importBasic(t *testing.T) {
	zoneName := fmt.Sprintf("ACPTTEST%s.com.", acctest.RandString(5))
	targetProjectName := "ACPTTEST-Target-" + acctest.RandString(5)

	zoneResourceName := "openstack_dns_zone_v2.zone"
	resourceName := "openstack_dns_zone_share_v2.share"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
			testAccPreCheckDNS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDNSZoneShareV2Destroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDNSZoneShareV2Config(zoneName, targetProjectName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAcckDNSZoneShareV2ImportID(zoneResourceName, resourceName),
			},
		},
	})
}

func TestAccResourceDNSZoneShareV2_importProjectID(t *testing.T) {
	zoneName := fmt.Sprintf("ACPTTEST%s.com.", acctest.RandString(5))
	targetProjectName := "ACPTTEST-Target-" + acctest.RandString(5)

	zoneResourceName := "openstack_dns_zone_v2.zone"
	resourceName := "openstack_dns_zone_share_v2.share"
	projectResourceName := "openstack_identity_project_v3.target"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
			testAccPreCheckDNS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDNSZoneShareV2Destroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDNSZoneShareV2Config(zoneName, targetProjectName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAcckDNSZoneShareProjectV2ImportID(zoneResourceName, resourceName, projectResourceName),
			},
		},
	})
}

func testAcckDNSZoneShareV2ImportID(zoneResource, shareResource string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		zone, ok := s.RootModule().Resources[zoneResource]
		if !ok {
			return "", fmt.Errorf("DNS zone not found: %s", zoneResource)
		}

		share, ok := s.RootModule().Resources[shareResource]
		if !ok {
			return "", fmt.Errorf("DNS zone share not found: %s", shareResource)
		}

		return fmt.Sprintf("%s/%s", zone.Primary.ID, share.Primary.ID), nil
	}
}

func testAcckDNSZoneShareProjectV2ImportID(zoneResource, shareResource, projectResource string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		zone, ok := s.RootModule().Resources[zoneResource]
		if !ok {
			return "", fmt.Errorf("DNS zone not found: %s", zoneResource)
		}

		share, ok := s.RootModule().Resources[shareResource]
		if !ok {
			return "", fmt.Errorf("DNS zone share not found: %s", shareResource)
		}

		project, ok := s.RootModule().Resources[projectResource]
		if !ok {
			return "", fmt.Errorf("Project not found: %s", projectResource)
		}

		return fmt.Sprintf("%s/%s/%s", zone.Primary.ID, share.Primary.ID, project.Primary.ID), nil
	}
}
