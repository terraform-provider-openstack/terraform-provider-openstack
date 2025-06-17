package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccComputeV2FlavorAccess_importBasic(t *testing.T) {
	resourceName := "openstack_compute_flavor_access_v2.access_1"

	flavorName := "ACCPTTEST-" + acctest.RandString(5)
	projectName := "ACCPTTEST-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2FlavorAccessDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2FlavorAccessBasic(flavorName, projectName),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
