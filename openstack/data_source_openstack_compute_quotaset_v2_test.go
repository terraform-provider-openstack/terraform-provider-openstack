package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccComputeV2QuotasetDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2QuotasetDataSourceBasic,
			},
			{
				Config: testAccComputeV2QuotasetDataSourceSource(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeQuotasetV2DataSourceID("data.openstack_compute_quotaset_v2.source"),
					resource.TestCheckResourceAttrSet("data.openstack_compute_quotaset_v2.source", "key_pairs"),
					resource.TestCheckResourceAttrSet("data.openstack_compute_quotaset_v2.source", "metadata_items"),
					resource.TestCheckResourceAttrSet("data.openstack_compute_quotaset_v2.source", "ram"),
					resource.TestCheckResourceAttrSet("data.openstack_compute_quotaset_v2.source", "cores"),
					resource.TestCheckResourceAttrSet("data.openstack_compute_quotaset_v2.source", "instances"),
					resource.TestCheckResourceAttrSet("data.openstack_compute_quotaset_v2.source", "server_groups"),
					resource.TestCheckResourceAttrSet("data.openstack_compute_quotaset_v2.source", "server_group_members"),
				),
			},
		},
	})
}

func testAccCheckComputeQuotasetV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find compute quotaset data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Compute quotaset data source ID not set")
		}

		return nil
	}
}

const testAccComputeV2QuotasetDataSourceBasic = `
resource "openstack_identity_project_v3" "project" {
  name = "test-quotaset-datasource"
}
`

func testAccComputeV2QuotasetDataSourceSource() string {
	return fmt.Sprintf(`
%s

data "openstack_compute_quotaset_v2" "source" {
  project_id = "${openstack_identity_project_v3.project.id}"
}
`, testAccComputeV2QuotasetDataSourceBasic)
}
