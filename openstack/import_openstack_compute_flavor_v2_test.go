package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccComputeV2Flavor_importBasic(t *testing.T) {
	resourceName := "openstack_compute_flavor_v2.flavor_1"

	flavorName := acctest.RandomWithPrefix("tf-acc-flavor")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2FlavorDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2FlavorBasic(flavorName),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
