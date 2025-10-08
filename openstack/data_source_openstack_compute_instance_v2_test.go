package openstack

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
					testAccCheckComputeInstanceV2DataSourceID(
						"data.openstack_compute_instance_v2.source_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_instance_v2.source_1", "name", "instance_1"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_compute_instance_v2.source_1", "metadata", "openstack_compute_instance_v2.instance_1", "metadata"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_compute_instance_v2.source_1", "network.0.name"),
				),
			},
			{
				Config: testAccComputeV2InstanceDataSourceName(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeInstanceV2DataSourceID(
						"data.openstack_compute_instance_v2.source_2"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_instance_v2.source_2", "name", "instance_1"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_compute_instance_v2.source_2", "metadata", "openstack_compute_instance_v2.instance_1", "metadata"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_compute_instance_v2.source_2", "network.0.name"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_instance_v2.source_2", "tags.#", "2"),
					resource.TestCheckResourceAttr(
						"data.openstack_compute_instance_v2.source_2", "tags.0", "tag1"),
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
			return errors.New("Compute instance data source ID not set")
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
  tags = [
    "tag1",
	"tag2",
  ]
}
`, osNetworkID)
}

func testAccComputeV2InstanceDataSourceSource() string {
	return fmt.Sprintf(`
%s

data "openstack_compute_instance_v2" "source_1" {
  id = openstack_compute_instance_v2.instance_1.id
}
`, testAccComputeV2InstanceDataSourceBasic())
}

func testAccComputeV2InstanceDataSourceName() string {
	return fmt.Sprintf(`
%s

data "openstack_compute_instance_v2" "source_2" {
  name = "^instance.*$"

  tags_all = [
    "tag1",
	"tag2",
  ]

  tags_any = [
	"tag2",
  ]

  not_tags_all = [
	"tag3",
	"tag4",
  ]

  not_tags_any = [
	"tag5",
  ]

  depends_on = [
    openstack_compute_instance_v2.instance_1
  ]
}
`, testAccComputeV2InstanceDataSourceBasic())
}
