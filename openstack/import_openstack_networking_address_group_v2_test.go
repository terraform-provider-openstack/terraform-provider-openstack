package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNetworkingV2AddressGroup_importBasic(t *testing.T) {
	resourceName := "openstack_networking_address_group_v2.group_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2AddressGroupDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2AddressGroupBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// Some OpenStack deployments may not return a project_id.
				ImportStateVerifyIgnore: []string{"project_id"},
			},
		},
	})
}
