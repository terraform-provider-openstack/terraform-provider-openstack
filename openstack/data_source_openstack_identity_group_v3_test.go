package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOpenStackIdentityV3GroupDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackIdentityV3GroupDataSource_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3DataSourceID("data.openstack_identity_group_v3.group_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_identity_group_v3.group_1", "name", "admins"),
				),
			},
		},
	})
}

func testAccCheckIdentityV3GroupDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find group data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Group data source ID not set")
		}

		return nil
	}
}

const testAccOpenStackIdentityV3GroupDataSource_basic = `
data "openstack_identity_group_v3" "group_1" {
    name = "admins"
}
`
