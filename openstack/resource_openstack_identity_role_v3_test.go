package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/roles"
)

func TestAccIdentityV3Role_basic(t *testing.T) {
	var role roles.Role
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIdentityV3RoleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3RoleBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3RoleExists("openstack_identity_role_v3.role_1", &role),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_role_v3.role_1", "name", &role.Name),
				),
			},
			{
				Config: testAccIdentityV3RoleUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3RoleExists("openstack_identity_role_v3.role_1", &role),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_role_v3.role_1", "name", &role.Name),
				),
			},
		},
	})
}

func testAccCheckIdentityV3RoleDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	identityClient, err := config.IdentityV3Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_identity_role_v3" {
			continue
		}

		_, err := roles.Get(identityClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Role still exists")
		}
	}

	return nil
}

func testAccCheckIdentityV3RoleExists(n string, role *roles.Role) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		identityClient, err := config.IdentityV3Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack identity client: %s", err)
		}

		found, err := roles.Get(identityClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Role not found")
		}

		*role = *found

		return nil
	}
}

const testAccIdentityV3RoleBasic = `
resource "openstack_identity_role_v3" "role_1" {
  name = "role_1"
}
`

const testAccIdentityV3RoleUpdate = `
resource "openstack_identity_role_v3" "role_1" {
  name = "role_2"
}
`
