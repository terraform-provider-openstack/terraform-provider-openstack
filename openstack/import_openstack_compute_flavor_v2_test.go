package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccComputeV2Flavor_importBasic(t *testing.T) {
	resourceName := "openstack_compute_flavor_v2.flavor_1"
	var flavorName = acctest.RandomWithPrefix("tf-acc-flavor")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2FlavorDestroy,
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
