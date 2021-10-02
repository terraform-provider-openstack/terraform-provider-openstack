package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/dns/v2/transfer/request"
)

func TestAccDNSV2TransferRequest_basic(t *testing.T) {
	var transferRequest request.TransferRequest
	var zoneName = fmt.Sprintf("ACPTTEST%s.com.", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckDNS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDNSV2TransferRequestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2TransferRequestBasic(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2TransferRequestExists(
						"openstack_dns_transfer_request_v2.request_1", &transferRequest),
					resource.TestCheckResourceAttr(
						"openstack_dns_transfer_request_v2.request_1", "description", "a transfer request"),
				),
			},
			{
				Config: testAccDNSV2TransferRequestUpdate(zoneName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_dns_transfer_request_v2.request_1", "description", "an updated transfer request"),
				),
			},
		},
	})
}

func TestAccDNSV2TransferRequest_ignoreStatusCheck(t *testing.T) {
	var transferRequest request.TransferRequest
	var zoneName = fmt.Sprintf("ACPTTEST%s.com.", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckDNS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDNSV2TransferRequestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2TransferRequestDisableCheck(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2TransferRequestExists("openstack_dns_transfer_request_v2.request_1", &transferRequest),
					resource.TestCheckResourceAttr(
						"openstack_dns_transfer_request_v2.request_1", "disable_status_check", "true"),
				),
			},
			{
				Config: testAccDNSV2TransferRequestBasic(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2TransferRequestExists("openstack_dns_transfer_request_v2.request_1", &transferRequest),
					resource.TestCheckResourceAttr(
						"openstack_dns_transfer_request_v2.request_1", "disable_status_check", "false"),
				),
			},
		},
	})
}

func testAccCheckDNSV2TransferRequestDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	dnsClient, err := config.DNSV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack DNS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_dns_transfer_request_v2" {
			continue
		}

		_, err := request.Get(dnsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Transfer request still exists")
		}
	}

	return nil
}

func testAccCheckDNSV2TransferRequestExists(n string, transferRequest *request.TransferRequest) resource.TestCheckFunc {
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

		found, err := request.Get(dnsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Transfer request not found")
		}

		*transferRequest = *found

		return nil
	}
}

func testAccDNSV2TransferRequestBasic(zoneName string) string {
	return fmt.Sprintf(`
		resource "openstack_dns_zone_v2" "zone_1" {
			name = "%s"
			email = "email1@example.com"
			description = "a zone"
			ttl = 3000
			type = "PRIMARY"
		}

		resource "openstack_dns_transfer_request_v2" "request_1" {
			zone_id = "${openstack_dns_zone_v2.zone_1.id}"
			target_project_id = "${openstack_dns_zone_v2.zone_1.project_id}"
			description = "a transfer request"
        }
	`, zoneName)
}

func testAccDNSV2TransferRequestUpdate(zoneName string) string {
	return fmt.Sprintf(`
		resource "openstack_dns_zone_v2" "zone_1" {
			name = "%s"
			email = "email1@example.com"
			description = "a zone"
			ttl = 3000
			type = "PRIMARY"
		}

		resource "openstack_dns_transfer_request_v2" "request_1" {
			zone_id = "${openstack_dns_zone_v2.zone_1.id}"
			target_project_id = "${openstack_dns_zone_v2.zone_1.project_id}"
			description = "an updated transfer request"
        }
	`, zoneName)
}

func testAccDNSV2TransferRequestDisableCheck(zoneName string) string {
	return fmt.Sprintf(`
		resource "openstack_dns_zone_v2" "zone_1" {
			name = "%s"
			email = "email1@example.com"
			description = "a zone"
			ttl = 3000
			type = "PRIMARY"
			disable_status_check = true
		}

		resource "openstack_dns_transfer_request_v2" "request_1" {
			zone_id = "${openstack_dns_zone_v2.zone_1.id}"
			target_project_id = "${openstack_dns_zone_v2.zone_1.project_id}"
			description = "a transfer request"
			disable_status_check = true
        }
	`, zoneName)
}
