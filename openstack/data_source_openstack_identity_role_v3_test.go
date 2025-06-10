package openstack

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOpenStackIdentityV3RoleDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackIdentityV3RoleDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3RoleDataSourceID("data.openstack_identity_role_v3.role_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_identity_role_v3.role_1", "name", "admin"),
				),
			},
		},
	})
}

func testAccCheckIdentityV3RoleDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find role data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("Role data source ID not set")
		}

		return nil
	}
}

const testAccOpenStackIdentityV3RoleDataSourceBasic = `
data "openstack_identity_role_v3" "role_1" {
    name = "admin"
}
`
