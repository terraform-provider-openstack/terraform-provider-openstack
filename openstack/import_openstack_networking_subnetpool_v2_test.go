package openstack

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccNetworkingV2SubnetPoolImportBasic(t *testing.T) {
	resourceName := "openstack_networking_subnetpool_v2.subnetpool_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2SubnetPoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SubnetPoolBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
