package openstack

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccEndpointGroup_importBasic(t *testing.T) {
	resourceName := "openstack_vpnaas_endpoint_group_v2.group_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckVPN(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEndpointGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointGroupV2_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
