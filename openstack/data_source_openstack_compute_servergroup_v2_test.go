package openstack

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccComputeServerGroupV2DataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeServerGroupV2DataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeServerGroupV2DataSourceID("data.openstack_compute_servergroup_v2.server_group_1"),
					resource.TestCheckResourceAttr("data.openstack_compute_servergroup_v2.server_group_1", "name", "my-servergroup"),
					resource.TestCheckResourceAttr("data.openstack_compute_servergroup_v2.server_group_1", "policy", "anti-affinity"),
					resource.TestCheckResourceAttr("data.openstack_compute_servergroup_v2.server_group_1", "rules.max_server_per_host", strconv.Itoa(3)),
				),
			},
		},
	})
}

func testAccCheckComputeServerGroupV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("ServerGroup data source ID not set")
		}

		return nil
	}
}

const testAccComputeServerGroupV2DataSourceBasic = `
resource "openstack_compute_servergroup_v2" "server_group_1" {
  name     = "my-servergroup"
  policies = ["anti-affinity"]
  rules {
    max_server_per_host = 3
  }
}

data "openstack_compute_servergroup_v2" "server_group_1" {
  name = "${openstack_compute_servergroup_v2.server_group_1.name}"
}
`
