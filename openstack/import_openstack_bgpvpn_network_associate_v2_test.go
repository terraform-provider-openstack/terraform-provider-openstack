package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBGPVPNNetworkAssociateV2_Import(t *testing.T) {
	resourceName := "openstack_bgpvpn_network_associate_v2.association_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
			testAccPreCheckBGPVPN(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBGPVPNNetworkAssociateV2Destroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccBGPVPNNetworkAssociateV2Config,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
