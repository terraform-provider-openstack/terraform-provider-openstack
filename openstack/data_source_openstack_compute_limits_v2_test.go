package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccComputeV2LimitsDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2LimitsDataSourceBasic,
			},
			{
				Config: testAccComputeV2LimitsDataSourceSource(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeLimitsV2DataSourceID("data.openstack_compute_limits_v2.source"),
					resource.TestCheckResourceAttrSet("data.openstack_compute_limits_v2.source", "max_total_cores"),
					resource.TestCheckResourceAttrSet("data.openstack_compute_limits_v2.source", "max_image_meta"),
					resource.TestCheckResourceAttrSet("data.openstack_compute_limits_v2.source", "max_server_meta"),
					resource.TestCheckResourceAttrSet("data.openstack_compute_limits_v2.source", "max_personality"),
					resource.TestCheckResourceAttrSet("data.openstack_compute_limits_v2.source", "max_personality_size"),
					resource.TestCheckResourceAttrSet("data.openstack_compute_limits_v2.source", "max_total_keypairs"),
					resource.TestCheckResourceAttrSet("data.openstack_compute_limits_v2.source", "max_security_groups"),
					resource.TestCheckResourceAttrSet("data.openstack_compute_limits_v2.source", "max_security_group_rules"),
					resource.TestCheckResourceAttrSet("data.openstack_compute_limits_v2.source", "max_server_groups"),
					resource.TestCheckResourceAttrSet("data.openstack_compute_limits_v2.source", "max_server_group_members"),
					resource.TestCheckResourceAttrSet("data.openstack_compute_limits_v2.source", "max_total_floating_ips"),
					resource.TestCheckResourceAttrSet("data.openstack_compute_limits_v2.source", "max_total_instances"),
					resource.TestCheckResourceAttrSet("data.openstack_compute_limits_v2.source", "max_total_ram_size"),
					resource.TestCheckResourceAttrSet("data.openstack_compute_limits_v2.source", "total_cores_used"),
					resource.TestCheckResourceAttrSet("data.openstack_compute_limits_v2.source", "total_instances_used"),
					resource.TestCheckResourceAttrSet("data.openstack_compute_limits_v2.source", "total_floating_ips_used"),
					resource.TestCheckResourceAttrSet("data.openstack_compute_limits_v2.source", "total_ram_used"),
					resource.TestCheckResourceAttrSet("data.openstack_compute_limits_v2.source", "total_security_groups_used"),
					resource.TestCheckResourceAttrSet("data.openstack_compute_limits_v2.source", "total_server_groups_used"),
				),
			},
		},
	})
}

func testAccCheckComputeLimitsV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find compute limits data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Compute limits data source ID not set")
		}

		return nil
	}
}

const testAccComputeV2LimitsDataSourceBasic = `
resource "openstack_identity_project_v3" "project" {
  name = "test-limits-datasource"
}
`

func testAccComputeV2LimitsDataSourceSource() string {
	return fmt.Sprintf(`
%s

data "openstack_compute_limits_v2" "source" {
  project_id = "${openstack_identity_project_v3.project.id}"
}
`, testAccComputeV2LimitsDataSourceBasic)
}
