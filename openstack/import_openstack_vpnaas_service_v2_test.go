package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccServiceV2_importBasic(t *testing.T) {
	resourceName := "openstack_vpnaas_service_v2.service_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckVPN(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServiceV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccServiceV2Basic,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
