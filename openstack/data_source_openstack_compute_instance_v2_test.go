package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccComputeV2InstanceDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceDataSourceBasic(),
			},
			{
				Config: testAccComputeV2InstanceDataSourceSource(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceV2DataSourceID("data.openstack_compute_instance_v2.source_1"),
					resource.TestCheckResourceAttr("data.openstack_compute_instance_v2.source_1", "name", "instance_1"),
					resource.TestCheckResourceAttrPair("data.openstack_compute_instance_v2.source_1", "metadata", "openstack_compute_instance_v2.instance_1", "metadata"),
					resource.TestCheckResourceAttrSet("data.openstack_compute_instance_v2.source_1", "network.0.name"),
				),
			},
		},
	})
}

func testAccCheckComputeInstanceV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find compute instance data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Compute instance data source ID not set")
		}

		return nil
	}
}

func testAccComputeV2InstanceDataSourceBasic() string {
	return fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  metadata = {
    foo = "bar"
  }
  network {
    uuid = "%s"
  }
}
`, osNetworkID)
}

func testAccComputeV2InstanceDataSourceSource() string {
	return fmt.Sprintf(`
%s

data "openstack_compute_instance_v2" "source_1" {
  id = "${openstack_compute_instance_v2.instance_1.id}"
}
`, testAccComputeV2InstanceDataSourceBasic())
}
