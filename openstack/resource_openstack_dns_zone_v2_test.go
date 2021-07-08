package openstack

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/dns/v2/zones"
)

func TestAccDNSV2Zone_basic(t *testing.T) {
	var zone zones.Zone
	var zoneName = fmt.Sprintf("ACPTTEST%s.com.", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckDNS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDNSV2ZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2ZoneBasic(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2ZoneExists("openstack_dns_zone_v2.zone_1", &zone),
					resource.TestCheckResourceAttr(
						"openstack_dns_zone_v2.zone_1", "description", "a zone"),
				),
			},
			{
				Config: testAccDNSV2ZoneUpdate(zoneName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openstack_dns_zone_v2.zone_1", "name", zoneName),
					resource.TestCheckResourceAttr("openstack_dns_zone_v2.zone_1", "email", "email2@example.com"),
					resource.TestCheckResourceAttr("openstack_dns_zone_v2.zone_1", "ttl", "6000"),
					resource.TestCheckResourceAttr("openstack_dns_zone_v2.zone_1", "type", "PRIMARY"),
					resource.TestCheckResourceAttr(
						"openstack_dns_zone_v2.zone_1", "description", "an updated zone"),
				),
			},
		},
	})
}

func TestAccDNSV2Zone_ignoreStatusCheck(t *testing.T) {
	var zone zones.Zone
	var zoneName = fmt.Sprintf("ACPTTEST%s.com.", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckDNS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDNSV2ZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2ZoneDisableCheck(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2ZoneExists("openstack_dns_zone_v2.zone_1", &zone),
					resource.TestCheckResourceAttr(
						"openstack_dns_zone_v2.zone_1", "disable_status_check", "true"),
				),
			},
			{
				Config: testAccDNSV2ZoneBasic(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2ZoneExists("openstack_dns_zone_v2.zone_1", &zone),
					resource.TestCheckResourceAttr(
						"openstack_dns_zone_v2.zone_1", "disable_status_check", "false"),
				),
			},
		},
	})
}

func TestAccDNSV2Zone_readTTL(t *testing.T) {
	var zone zones.Zone
	var zoneName = fmt.Sprintf("ACPTTEST%s.com.", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckDNS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDNSV2ZoneDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2ZoneReadTTL(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2ZoneExists("openstack_dns_zone_v2.zone_1", &zone),
					resource.TestCheckResourceAttr("openstack_dns_zone_v2.zone_1", "type", "PRIMARY"),
					resource.TestMatchResourceAttr(
						"openstack_dns_zone_v2.zone_1", "ttl", regexp.MustCompile("^[0-9]+$")),
				),
			},
		},
	})
}

func testAccCheckDNSV2ZoneDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	dnsClient, err := config.DNSV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack DNS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_dns_zone_v2" {
			continue
		}

		_, err := zones.Get(dnsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Zone still exists")
		}
	}

	return nil
}

func testAccCheckDNSV2ZoneExists(n string, zone *zones.Zone) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		dnsClient, err := config.DNSV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack DNS client: %s", err)
		}

		found, err := zones.Get(dnsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Zone not found")
		}

		*zone = *found

		return nil
	}
}

func testAccDNSV2ZoneBasic(zoneName string) string {
	return fmt.Sprintf(`
		resource "openstack_dns_zone_v2" "zone_1" {
			name = "%s"
			email = "email1@example.com"
			description = "a zone"
			ttl = 3000
			type = "PRIMARY"
		}
	`, zoneName)
}

func testAccDNSV2ZoneUpdate(zoneName string) string {
	return fmt.Sprintf(`
		resource "openstack_dns_zone_v2" "zone_1" {
			name = "%s"
			email = "email2@example.com"
			description = "an updated zone"
			ttl = 6000
			type = "PRIMARY"
		}
	`, zoneName)
}

func testAccDNSV2ZoneReadTTL(zoneName string) string {
	return fmt.Sprintf(`
		resource "openstack_dns_zone_v2" "zone_1" {
			name = "%s"
			email = "email1@example.com"
		}
	`, zoneName)
}

func testAccDNSV2ZoneDisableCheck(zoneName string) string {
	return fmt.Sprintf(`
		resource "openstack_dns_zone_v2" "zone_1" {
			name = "%s"
			email = "email1@example.com"
			description = "a zone"
			ttl = 3000
			type = "PRIMARY"
			disable_status_check = true
		}
	`, zoneName)
}
