package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDNSV2TransferRequest_importBasic(t *testing.T) {
	zoneName := randomZoneName()
	resourceName := "openstack_dns_transfer_request_v2.request_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckDNS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDNSV2TransferRequestDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2TransferRequestBasic(zoneName),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"disable_status_check",
				},
			},
		},
	})
}
