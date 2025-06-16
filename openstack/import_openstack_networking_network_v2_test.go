package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccNetworkingV2Network_importBasic(t *testing.T) {
	resourceName := "openstack_networking_network_v2.network_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2NetworkDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2NetworkBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccNetworkingV2Network_importSegments(t *testing.T) {
	resourceName := "openstack_networking_network_v2.network_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2NetworkDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2NetworkMultipleSegmentMappings,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
