package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOpenStackNetworkingFloatingIPV2DataSource_address(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingFloatingIPV2DataSourceFloatingIP(),
			},
			{
				Config: testAccOpenStackNetworkingFloatingIPV2DataSourceAddress(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingFloatingIPV2DataSourceID("data.openstack_networking_floatingip_v2.fip_1"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_networking_floatingip_v2.fip_1", "address"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_networking_floatingip_v2.fip_1", "pool"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_networking_floatingip_v2.fip_1", "status"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_networking_floatingip_v2.fip_1", "description"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_floatingip_v2.fip_1", "tags.#", "1"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_floatingip_v2.fip_1", "all_tags.#", "2"),
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

func testAccOpenStackNetworkingFloatingIPV2DataSourceFloatingIP() string {
	return fmt.Sprintf(`
resource "openstack_networking_floatingip_v2" "fip_1" {
  pool = "%s"
  description = "test fip"
  tags = [
    "foo",
    "bar",
  ]
}
`, osPoolName)
}

func testAccOpenStackNetworkingFloatingIPV2DataSourceAddress() string {
	return fmt.Sprintf(`
%s

data "openstack_networking_floatingip_v2" "fip_1" {
  address = "${openstack_networking_floatingip_v2.fip_1.address}"
  description = "test fip"
  tags = [
    "foo",
  ]
}
`, testAccOpenStackNetworkingFloatingIPV2DataSourceFloatingIP())
}
