package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOpenStackNetworkingRouterV2DataSource_name(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingRouterV2DataSourceRouter,
			},
			{
				Config: testAccOpenStackNetworkingRouterV2DataSourceName(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingRouterV2DataSourceID("data.openstack_networking_router_v2.router"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_networking_router_v2.router", "name"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_networking_router_v2.router", "description"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_networking_router_v2.router", "admin_state_up"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_networking_router_v2.router", "status"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_router_v2.router", "tags.#", "1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_router_v2.router", "all_tags.#", "2"),
				),
			},
		},
	})
}

func testAccCheckNetworkingRouterV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find router data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Router data source ID not set")
		}

		return nil
	}
}

const testAccOpenStackNetworkingRouterV2DataSourceRouter = `
resource "openstack_networking_router_v2" "router" {
  name           = "router_tf"
  description    = "description"
  admin_state_up = "true"
  tags = [
    "foo",
    "bar",
  ]
}
`

func testAccOpenStackNetworkingRouterV2DataSourceName() string {
	return fmt.Sprintf(`
%s

data "openstack_networking_router_v2" "router" {
  name           = "${openstack_networking_router_v2.router.name}"
  description    = "description"
  admin_state_up = "true"
  tags = [
    "foo",
  ]
}
`, testAccOpenStackNetworkingRouterV2DataSourceRouter)
}
