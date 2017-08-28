package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccComputeV2Flavor_importBasic(t *testing.T) {
	resourceName := "openstack_compute_flavor_v2.flavor_1"
	var projectName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2FlavorDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeV2Flavor_basic(projectName),
			},

			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
