package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccEndpointGroup_importBasic(t *testing.T) {
	resourceName := "openstack_vpnaas_endpoint_group_v2.group_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckVPN(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckEndpointGroupV2Destroy,
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
