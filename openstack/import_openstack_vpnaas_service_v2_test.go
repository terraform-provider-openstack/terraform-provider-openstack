package openstack

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccServiceV2_importBasic(t *testing.T) {
	resourceName := "openstack_vpnaas_service_v2.service_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckVPN(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServiceV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceV2_basic,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
