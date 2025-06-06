package openstack

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccLBV2FlavorDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLBV2FlavorDataSourceBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2FlavorDataSourceID("data.openstack_lb_flavor_v2.flavor_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_flavor_v2.flavor_1", "name", "lb.acctest"),
				),
			},
		},
	})
}

func testAccCheckLBV2FlavorDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find flavor data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("Flavor data source ID not set")
		}

		return nil
	}
}

func testAccLBV2FlavorDataSourceBasic() string {
	return fmt.Sprintf(`
data "openstack_lb_flavor_v2" "flavor_1" {
  name = "%s"
}
`, osLbFlavorName)
}
