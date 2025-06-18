package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccVPNaaSV2EndpointGroup_importBasic(t *testing.T) {
	resourceName := "openstack_vpnaas_endpoint_group_v2.group_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckVPN(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckEndpointGroupV2Destroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointGroupV2Basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
