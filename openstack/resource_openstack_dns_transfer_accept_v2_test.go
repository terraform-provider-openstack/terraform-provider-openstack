package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/dns/v2/transfer/accept"
)

func TestAccDNSV2TransferAccept_basic(t *testing.T) {
	var transferAccept accept.TransferAccept
	var zoneName = fmt.Sprintf("ACPTTEST%s.com.", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckDNS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDNSV2TransferAcceptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2TransferAcceptBasic(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2TransferAcceptExists(
						"openstack_dns_transfer_accept_v2.accept_1", &transferAccept),
				),
			},
		},
	})
}

func TestAccDNSV2TransferAccept_ignoreStatusCheck(t *testing.T) {
	var transferAccept accept.TransferAccept
	var zoneName = fmt.Sprintf("ACPTTEST%s.com.", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckDNS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDNSV2TransferAcceptDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2TransferAcceptDisableCheck(zoneName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDNSV2TransferAcceptExists(
						"openstack_dns_transfer_accept_v2.accept_1", &transferAccept),
					resource.TestCheckResourceAttr(
						"openstack_dns_transfer_accept_v2.accept_1", "disable_status_check", "true"),
				),
			},
		},
	})
}

func testAccCheckDNSV2TransferAcceptExists(n string, transferAccept *accept.TransferAccept) resource.TestCheckFunc {
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

		found, err := accept.Get(dnsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Transfer accept not found")
		}

		*transferAccept = *found

		return nil
	}
}

func testAccCheckDNSV2TransferAcceptDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	dnsClient, err := config.DNSV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack DNS client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_dns_transfer_accept_v2" {
			continue
		}

		_, err := accept.Get(dnsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Transfer accept still exists")
		}
	}

	return nil
}

func testAccDNSV2TransferAcceptBasic(zoneName string) string {
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

		resource "openstack_dns_transfer_accept_v2" "accept_1" {
			zone_transfer_request_id = "${openstack_dns_transfer_request_v2.request_1.id}"
			key = "${openstack_dns_transfer_request_v2.request_1.key}"
        }
	`, zoneName)
}

func testAccDNSV2TransferAcceptDisableCheck(zoneName string) string {
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

		resource "openstack_dns_transfer_accept_v2" "accept_1" {
			zone_transfer_request_id = "${openstack_dns_transfer_request_v2.request_1.id}"
			key = "${openstack_dns_transfer_request_v2.request_1.key}"
			disable_status_check = true
        }
	`, zoneName)
}
