package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIdentityV3User_importBasic(t *testing.T) {
	resourceName := "openstack_identity_user_v3.user_1"

	userName := "ACCPTTEST-" + acctest.RandString(5)

	projectName := "ACCPTTEST-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIdentityV3UserDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3UserBasic(projectName, userName),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
				},
			},
		},
	})
}
