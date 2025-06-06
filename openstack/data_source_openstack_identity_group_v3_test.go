package openstack

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOpenStackIdentityV3GroupDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackIdentityV3GroupDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3GroupDataSourceID("data.openstack_identity_group_v3.group_1"),
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
			return errors.New("Group data source ID not set")
		}

		return nil
	}
}

const testAccOpenStackIdentityV3GroupDataSourceBasic = `
data "openstack_identity_group_v3" "group_1" {
    name = "admins"
}
`
