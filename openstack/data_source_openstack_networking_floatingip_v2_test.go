package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOpenStackNetworkingFloatingIPV2DataSource_address(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingFloatingIPV2DataSource_fip,
			},
			{
				Config: testAccOpenStackNetworkingFloatingIPV2DataSource_address,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingFloatingIPV2DataSourceID("data.openstack_networking_floatingip_v2.fip"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_networking_floatingip_v2.fip", "address"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_networking_floatingip_v2.fip", "pool"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_networking_floatingip_v2.fip", "status"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_networking_floatingip_v2.fip", "description"),
				),
			},
		},
	})
}

func testAccCheckNetworkingFloatingIPV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find floating IP data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Floating IP data source ID not set")
		}

		return nil
	}
}

var testAccOpenStackNetworkingFloatingIPV2DataSource_fip = fmt.Sprintf(`
resource "openstack_networking_floatingip_v2" "fip" {
  pool = "%s"
  description = "test fip"
}
`, OS_POOL_NAME)

var testAccOpenStackNetworkingFloatingIPV2DataSource_address = fmt.Sprintf(`
%s

data "openstack_networking_floatingip_v2" "fip" {
   address = "${openstack_networking_floatingip_v2.fip.address}"
   description = "test fip"
}
`, testAccOpenStackNetworkingFloatingIPV2DataSource_fip)
