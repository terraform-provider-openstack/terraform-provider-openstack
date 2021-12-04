package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccNetworkingV2QuotaDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2QuotaDataSourceBasic,
			},
			{
				Config: testAccNetworkingV2QuotaDataSourceSource(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingQuotaV2DataSourceID("data.openstack_networking_quota_v2.source"),
					resource.TestCheckResourceAttrSet("data.openstack_networking_quota_v2.source", "floatingip"),
					resource.TestCheckResourceAttrSet("data.openstack_networking_quota_v2.source", "network"),
					resource.TestCheckResourceAttrSet("data.openstack_networking_quota_v2.source", "port"),
					resource.TestCheckResourceAttrSet("data.openstack_networking_quota_v2.source", "rbac_policy"),
					resource.TestCheckResourceAttrSet("data.openstack_networking_quota_v2.source", "router"),
					resource.TestCheckResourceAttrSet("data.openstack_networking_quota_v2.source", "security_group"),
					resource.TestCheckResourceAttrSet("data.openstack_networking_quota_v2.source", "security_group_rule"),
					resource.TestCheckResourceAttrSet("data.openstack_networking_quota_v2.source", "subnet"),
					resource.TestCheckResourceAttrSet("data.openstack_networking_quota_v2.source", "subnetpool"),
				),
			},
		},
	})
}

func testAccCheckNetworkingQuotaV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find networking quota data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Networking quota data source ID not set")
		}

		return nil
	}
}

const testAccNetworkingV2QuotaDataSourceBasic = `
resource "openstack_identity_project_v3" "project" {
  name = "test-quota-datasource"
}
`

func testAccNetworkingV2QuotaDataSourceSource() string {
	return fmt.Sprintf(`
%s

data "openstack_networking_quota_v2" "source" {
  project_id = "${openstack_identity_project_v3.project.id}"
}
`, testAccNetworkingV2QuotaDataSourceBasic)
}
