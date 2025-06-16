package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDNSV2Zone_importBasic(t *testing.T) {
	zoneName := fmt.Sprintf("ACPTTEST%s.com.", acctest.RandString(5))

	resourceName := "openstack_dns_zone_v2.zone_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckDNS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDNSV2ZoneDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccDNSV2ZoneBasic(zoneName),
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
