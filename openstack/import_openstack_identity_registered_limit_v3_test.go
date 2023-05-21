package openstack

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIdentityV3RegisteredLimit_importBasic(t *testing.T) {
	_ = os.Setenv("OS_SYSTEM_SCOPE", "true")
	defer os.Unsetenv("OS_SYSTEM_SCOPE")
	resourceName := "openstack_identity_registered_limit_v3.limit_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIdentityV3RegisteredLimitDestroy,
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
