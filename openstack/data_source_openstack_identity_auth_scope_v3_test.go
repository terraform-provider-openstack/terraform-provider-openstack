package openstack

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOpenStackIdentityAuthScopeV3DataSource_basic(t *testing.T) {
	userName := os.Getenv("OS_USERNAME")
	projectName := os.Getenv("OS_PROJECT_NAME")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackIdentityAuthScopeV3DataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityAuthScopeV3DataSourceID("data.openstack_identity_auth_scope_v3.token"),
					resource.TestCheckResourceAttr(
						"data.openstack_identity_auth_scope_v3.token", "user_name", userName),
					resource.TestCheckResourceAttr(
						"data.openstack_identity_auth_scope_v3.token", "project_name", projectName),
					resource.TestCheckNoResourceAttr(
						"data.openstack_identity_auth_scope_v3.token", "token_id"),
				),
			},
		},
	})
}

func TestAccOpenStackIdentityAuthScopeV3DataSource_token_id(t *testing.T) {
	userName := os.Getenv("OS_USERNAME")
	projectName := os.Getenv("OS_PROJECT_NAME")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackIdentityAuthScopeV3DataSourceTokenID,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityAuthScopeV3DataSourceID("data.openstack_identity_auth_scope_v3.token"),
					resource.TestCheckResourceAttr(
						"data.openstack_identity_auth_scope_v3.token", "user_name", userName),
					resource.TestCheckResourceAttr(
						"data.openstack_identity_auth_scope_v3.token", "project_name", projectName),
					resource.TestCheckResourceAttrSet("data.openstack_identity_auth_scope_v3.token", "token_id"),
				),
			},
		},
	})
}

func testAccCheckIdentityAuthScopeV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find token data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Token data source ID not set")
		}

		return nil
	}
}

const testAccOpenStackIdentityAuthScopeV3DataSourceBasic = `
data "openstack_identity_auth_scope_v3" "token" {
  name = "my_token"
}
`

const testAccOpenStackIdentityAuthScopeV3DataSourceTokenID = `
data "openstack_identity_auth_scope_v3" "token" {
  name         = "my_token"
  set_token_id = true
}
`
