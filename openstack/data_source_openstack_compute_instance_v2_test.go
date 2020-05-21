package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccComputeV2InstanceDataSource(t *testing.T) {

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InstanceDataSource_basic,
			},
			{
				Config: testAccComputeV2InstanceDataSource_source,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceV2DataSourceId("data.openstack_compute_instance_v2.source_1"),
					resource.TestCheckResourceAttr("data.openstack_compute_instance_v2.source_1", "name", "instance_1"),
					resource.TestCheckResourceAttrPair("data.openstack_compute_instance_v2.source_1", "metadata", "openstack_compute_instance_v2.instance_1", "metadata"),
					resource.TestCheckResourceAttrSet("data.openstack_compute_instance_v2.source_1", "network.0.name"),
				),
			},
		},
	})
}

func testAccCheckComputeInstanceV2DataSourceId(n string) resource.TestCheckFunc {
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

var testAccComputeV2InstanceDataSource_basic = fmt.Sprintf(`
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
`, OS_NETWORK_ID)

var testAccComputeV2InstanceDataSource_source = fmt.Sprintf(`
%s

data "openstack_compute_instance_v2" "source_1" {
  id = "${openstack_compute_instance_v2.instance_1.id}"
}
`, testAccComputeV2InstanceDataSource_basic)
