package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNetworkingV2FloatingIPAssociate_importBasic(t *testing.T) {
	resourceName := "openstack_networking_floatingip_associate_v2.fip_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2FloatingIPAssociateDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2FloatingIPAssociateBasic(),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
