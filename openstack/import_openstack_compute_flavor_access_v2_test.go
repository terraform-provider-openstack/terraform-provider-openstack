package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccComputeV2FlavorAccess_importBasic(t *testing.T) {
	resourceName := "openstack_compute_flavor_access_v2.access_1"

	flavorName := fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))
	projectName := fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2FlavorAccessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2FlavorAccess_basic(flavorName, projectName),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
