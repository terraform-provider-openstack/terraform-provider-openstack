package openstack

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccLoadBalancerV2FlavorDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLoadBalancerV2FlavorDataSourceBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLoadBalancerV2FlavorDataSourceID("data.openstack_loadbalancer_flavor_v2.flavor_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_loadbalancer_flavor_v2.flavor_1", "name", "lb.acctest"),
				),
			},
		},
	})
}

func testAccCheckLoadBalancerV2FlavorDataSourceID(n string) resource.TestCheckFunc {
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

func testAccLoadBalancerV2FlavorDataSourceBasic() string {
	return fmt.Sprintf(`
data "openstack_loadbalancer_flavor_v2" "flavor_1" {
  name = "%s"
}
`, osLbFlavorName)
}
