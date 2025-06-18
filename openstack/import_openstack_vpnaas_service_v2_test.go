package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPNaaSV2ServiceV2_importBasic(t *testing.T) {
	resourceName := "openstack_vpnaas_service_v2.service_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckVPN(t)
			t.Skip("Currently failing in GH-A")
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckServiceV2Destroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceV2Basic(),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
