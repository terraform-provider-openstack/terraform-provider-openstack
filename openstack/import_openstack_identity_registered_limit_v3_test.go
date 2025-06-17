package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIdentityV3RegisteredLimit_importBasic(t *testing.T) {
	resourceName := "openstack_identity_registered_limit_v3.limit_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIdentityV3RegisteredLimitDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3RegisteredLimitBasic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
