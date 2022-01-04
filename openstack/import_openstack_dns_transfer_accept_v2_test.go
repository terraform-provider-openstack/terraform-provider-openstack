package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDNSV2TransferAccept_importBasic(t *testing.T) {
	zoneName := randomZoneName()
	resourceName := "openstack_dns_transfer_request_v2.request_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckDNS(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2TransferAcceptBasic(zoneName),
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
