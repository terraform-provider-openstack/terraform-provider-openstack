package openstack

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/dns/v2/recordsets"
)

func randomZoneName() string {
	return fmt.Sprintf("ACPTTEST-zone-%s.com.", acctest.RandString(5))
}

func TestAccDNSV2RecordSet_basic(t *testing.T) {
	var recordset recordsets.RecordSet
	zoneName := randomZoneName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckDNS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDNSV2RecordSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2RecordSetBasic(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2RecordSetExists("openstack_dns_recordset_v2.recordset_1", &recordset),
					resource.TestCheckResourceAttr(
						"openstack_dns_recordset_v2.recordset_1", "description", "a record set"),
					resource.TestCheckResourceAttr(
						"openstack_dns_recordset_v2.recordset_1", "records.0", "10.1.0.0"),
				),
			},
			{
				Config: testAccDNSV2RecordSetUpdate(zoneName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openstack_dns_recordset_v2.recordset_1", "name", zoneName),
					resource.TestCheckResourceAttr("openstack_dns_recordset_v2.recordset_1", "ttl", "6000"),
					resource.TestCheckResourceAttr("openstack_dns_recordset_v2.recordset_1", "type", "A"),
					resource.TestCheckResourceAttr(
						"openstack_dns_recordset_v2.recordset_1", "description", "an updated record set"),
					resource.TestCheckResourceAttr(
						"openstack_dns_recordset_v2.recordset_1", "records.0", "10.1.0.1"),
				),
			},
		},
	})
}

func TestAccDNSV2RecordSet_ipv6(t *testing.T) {
	var recordset recordsets.RecordSet
	zoneName := randomZoneName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckDNS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDNSV2RecordSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2RecordSetIPv6(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2RecordSetExists("openstack_dns_recordset_v2.recordset_1", &recordset),
					resource.TestCheckResourceAttr(
						"openstack_dns_recordset_v2.recordset_1", "description", "a record set"),
					resource.TestCheckResourceAttr(
						"openstack_dns_recordset_v2.recordset_1", "records.0", "fd2b:db7f:6ae:dd8d::1"),
					resource.TestCheckResourceAttr(
						"openstack_dns_recordset_v2.recordset_1", "records.1", "fd2b:db7f:6ae:dd8d::2"),
				),
			},
		},
	})
}

func TestAccDNSV2RecordSet_readTTL(t *testing.T) {
	var recordset recordsets.RecordSet
	zoneName := randomZoneName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckDNS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDNSV2RecordSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2RecordSetReadTTL(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2RecordSetExists("openstack_dns_recordset_v2.recordset_1", &recordset),
					resource.TestMatchResourceAttr(
						"openstack_dns_recordset_v2.recordset_1", "ttl", regexp.MustCompile("^[0-9]+$")),
				),
			},
		},
	})
}

func TestAccDNSV2RecordSet_ensureSameTTL(t *testing.T) {
	var recordset recordsets.RecordSet
	zoneName := randomZoneName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckDNS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDNSV2RecordSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2RecordSetEnsureSameTTL1(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2RecordSetExists("openstack_dns_recordset_v2.recordset_1", &recordset),
					resource.TestCheckResourceAttr(
						"openstack_dns_recordset_v2.recordset_1", "records.0", "10.1.0.1"),
					resource.TestCheckResourceAttr(
						"openstack_dns_recordset_v2.recordset_1", "ttl", "3000"),
				),
			},
			{
				Config: testAccDNSV2RecordSetEnsureSameTTL2(zoneName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_dns_recordset_v2.recordset_1", "records.0", "10.1.0.2"),
					resource.TestCheckResourceAttr(
						"openstack_dns_recordset_v2.recordset_1", "ttl", "3000"),
				),
			},
		},
	})
}

func TestAccDNSV2RecordSet_sudoProjectID(t *testing.T) {
	var recordset recordsets.RecordSet
	zoneName := randomZoneName()

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckDNS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDNSV2RecordSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2RecordSetSudoProjectID(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2RecordSetExists("openstack_dns_recordset_v2.recordset_1", &recordset),
					resource.TestCheckResourceAttr(
						"openstack_dns_recordset_v2.recordset_1", "records.0", "10.1.0.1"),
				),
			},
		},
	})
}

func testAccCheckDNSV2RecordSetDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	dnsClient, err := config.DNSV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack DNS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_dns_recordset_v2" {
			continue
		}

		zoneID, recordsetID, err := dnsRecordSetV2ParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = recordsets.Get(dnsClient, zoneID, recordsetID).Extract()
		if err == nil {
			return fmt.Errorf("Record set still exists")
		}
	}

	return nil
}

func testAccCheckDNSV2RecordSetExists(n string, recordset *recordsets.RecordSet) resource.TestCheckFunc {
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

		zoneID, recordsetID, err := dnsRecordSetV2ParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		found, err := recordsets.Get(dnsClient, zoneID, recordsetID).Extract()
		if err != nil {
			return err
		}

		if found.ID != recordsetID {
			return fmt.Errorf("Record set not found")
		}

		*recordset = *found

		return nil
	}
}

func TestAccDNSV2RecordSet_ignoreStatusCheck(t *testing.T) {
	var recordset recordsets.RecordSet
	zoneName := randomZoneName()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckDNS(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDNSV2RecordSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2RecordSetDisableCheck(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2RecordSetExists("openstack_dns_recordset_v2.recordset_1", &recordset),
					resource.TestCheckResourceAttr(
						"openstack_dns_recordset_v2.recordset_1", "disable_status_check", "true"),
				),
			},
			{
				Config: testAccDNSV2RecordSetBasic(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2RecordSetExists("openstack_dns_recordset_v2.recordset_1", &recordset),
					resource.TestCheckResourceAttr(
						"openstack_dns_recordset_v2.recordset_1", "disable_status_check", "false"),
				),
			},
		},
	})
}

func testAccDNSV2RecordSetBasic(zoneName string) string {
	return fmt.Sprintf(`
		resource "openstack_dns_zone_v2" "zone_1" {
			name = "%s"
			email = "email2@example.com"
			description = "a zone"
			ttl = 6000
			type = "PRIMARY"
		}

		resource "openstack_dns_recordset_v2" "recordset_1" {
			zone_id = "${openstack_dns_zone_v2.zone_1.id}"
			name = "%s"
			type = "A"
			description = "a record set"
			ttl = 3000
			records = ["10.1.0.0"]
		}
	`, zoneName, zoneName)
}

func testAccDNSV2RecordSetUpdate(zoneName string) string {
	return fmt.Sprintf(`
		resource "openstack_dns_zone_v2" "zone_1" {
			name = "%s"
			email = "email2@example.com"
			description = "an updated zone"
			ttl = 6000
			type = "PRIMARY"
		}

		resource "openstack_dns_recordset_v2" "recordset_1" {
			zone_id = "${openstack_dns_zone_v2.zone_1.id}"
			name = "%s"
			type = "A"
			description = "an updated record set"
			ttl = 6000
			records = ["10.1.0.1"]
		}
	`, zoneName, zoneName)
}

func testAccDNSV2RecordSetReadTTL(zoneName string) string {
	return fmt.Sprintf(`
		resource "openstack_dns_zone_v2" "zone_1" {
			name = "%s"
			email = "email2@example.com"
			description = "an updated zone"
			ttl = 6000
			type = "PRIMARY"
		}

		resource "openstack_dns_recordset_v2" "recordset_1" {
			zone_id = "${openstack_dns_zone_v2.zone_1.id}"
			name = "%s"
			type = "A"
			records = ["10.1.0.2"]
		}
	`, zoneName, zoneName)
}

func testAccDNSV2RecordSetIPv6(zoneName string) string {
	return fmt.Sprintf(`
		resource "openstack_dns_zone_v2" "zone_1" {
			name = "%s"
			email = "email2@example.com"
			description = "a zone"
			ttl = 6000
			type = "PRIMARY"
		}

		resource "openstack_dns_recordset_v2" "recordset_1" {
			zone_id = "${openstack_dns_zone_v2.zone_1.id}"
			name = "%s"
			type = "AAAA"
			description = "a record set"
			ttl = 3000
			records = [
				"[fd2b:db7f:6ae:dd8d::1]",
				"fd2b:db7f:6ae:dd8d::2"
			]
		}
	`, zoneName, zoneName)
}

func testAccDNSV2RecordSetEnsureSameTTL1(zoneName string) string {
	return fmt.Sprintf(`
		resource "openstack_dns_zone_v2" "zone_1" {
			name = "%s"
			email = "email2@example.com"
			ttl = 6000
			type = "PRIMARY"
		}

		resource "openstack_dns_recordset_v2" "recordset_1" {
			zone_id = "${openstack_dns_zone_v2.zone_1.id}"
			name = "%s"
			type = "A"
			ttl = 3000
			records = ["10.1.0.1"]
		}
	`, zoneName, zoneName)
}

func testAccDNSV2RecordSetEnsureSameTTL2(zoneName string) string {
	return fmt.Sprintf(`
		resource "openstack_dns_zone_v2" "zone_1" {
			name = "%s"
			email = "email2@example.com"
			ttl = 6000
			type = "PRIMARY"
		}

		resource "openstack_dns_recordset_v2" "recordset_1" {
			zone_id = "${openstack_dns_zone_v2.zone_1.id}"
			name = "%s"
			type = "A"
			ttl = 3000
			records = ["10.1.0.2"]
		}
	`, zoneName, zoneName)
}

func testAccDNSV2RecordSetDisableCheck(zoneName string) string {
	return fmt.Sprintf(`
		resource "openstack_dns_zone_v2" "zone_1" {
			name = "%s"
			email = "email2@example.com"
			description = "a zone"
			ttl = 6000
			type = "PRIMARY"
		}

		resource "openstack_dns_recordset_v2" "recordset_1" {
			zone_id = "${openstack_dns_zone_v2.zone_1.id}"
			name = "%s"
			type = "A"
			description = "a record set"
			ttl = 3000
			records = ["10.1.0.0"]
			disable_status_check = true
		}
	`, zoneName, zoneName)
}

func testAccDNSV2RecordSetSudoProjectID(zoneName string) string {
	return fmt.Sprintf(`
		resource "openstack_dns_zone_v2" "zone_1" {
			name = "%s"
			email = "email2@example.com"
			ttl = 6000
			type = "PRIMARY"
		}

		resource "openstack_dns_recordset_v2" "recordset_1" {
			zone_id = "${openstack_dns_zone_v2.zone_1.id}"
			name = "%s"
			type = "A"
			ttl = 3000
			records = ["10.1.0.1"]
            project_id = "${openstack_dns_zone_v2.zone_1.project_id}"
		}
	`, zoneName, zoneName)
}
