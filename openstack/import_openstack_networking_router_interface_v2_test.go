package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccNetworkingV2RouterInterface_importBasic_port(t *testing.T) {
	resourceName := "openstack_networking_router_interface_v2.int_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2RouterInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterInterfaceBasicPort,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetworkingV2RouterInterface_importBasic_subnet(t *testing.T) {
	resourceName := "openstack_networking_router_interface_v2.int_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2RouterInterfaceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterInterfaceBasicSubnet,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
