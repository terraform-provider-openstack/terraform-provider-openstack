package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIdentityV3Project_importBasic(t *testing.T) {
	resourceName := "openstack_identity_project_v3.project_1"

	projectName := "ACCPTTEST-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIdentityV3ProjectDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3ProjectBasic(projectName),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
