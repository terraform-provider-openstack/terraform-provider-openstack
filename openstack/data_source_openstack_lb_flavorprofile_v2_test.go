package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccLBV2FlavorProfileDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2FlavorProfileDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2FlavorProfileDataSourceID("data.openstack_lb_flavorprofile_v2.fp_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_flavorprofile_v2.fp_1", "name", "lb.acctest"),
				),
			},
		},
	})
}

func testAccCheckLBV2FlavorProfileDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find flavor data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Flavor data source ID not set")
		}

		return nil
	}
}

const testAccLBV2FlavorProfileDataSourceBasic = `
data "openstack_lb_flavorprofile_v2" "fp_1" {
  name = "lb.acctest"
}
`
