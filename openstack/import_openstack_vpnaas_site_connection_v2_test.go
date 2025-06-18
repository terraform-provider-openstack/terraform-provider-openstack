package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPNaaSV2SiteConnectionV2_importBasic(t *testing.T) {
	resourceName := "openstack_vpnaas_site_connection_v2.conn_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckVPN(t)
			t.Skip("Currently failing in GH-A")
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSiteConnectionV2Destroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccSiteConnectionV2Basic(),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
