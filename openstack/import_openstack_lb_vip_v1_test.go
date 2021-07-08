package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccLBV1VIP_importBasic(t *testing.T) {
	resourceName := "openstack_lb_vip_v1.vip_1"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckDeprecated(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV1VIPDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbV1VIPBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
