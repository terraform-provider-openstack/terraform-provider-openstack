package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBGPVPNV2_Import(t *testing.T) {
	resourceName := "openstack_bgpvpn_v2.bgpvpn_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
			testAccPreCheckBGPVPN(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBGPVPNV2Destroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccBGPVPNV2Config,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "name", "bgpvpn_1"),
					resource.TestCheckResourceAttr(resourceName, "type", "l3"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
