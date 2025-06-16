package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/roles"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccIdentityV3Role_basic(t *testing.T) {
	var role roles.Role

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIdentityV3RoleDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3RoleBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3RoleExists(t.Context(), "openstack_identity_role_v3.role_1", &role),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_role_v3.role_1", "name", &role.Name),
				),
			},
			{
				Config: testAccIdentityV3RoleUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3RoleExists(t.Context(), "openstack_identity_role_v3.role_1", &role),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_role_v3.role_1", "name", &role.Name),
				),
			},
		},
	})
}

func testAccCheckIdentityV3RoleDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		identityClient, err := config.IdentityV3Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack identity client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_identity_role_v3" {
				continue
			}

			_, err := roles.Get(ctx, identityClient, rs.Primary.ID).Extract()
			if err == nil {
				return errors.New("Role still exists")
			}
		}

		return nil
	}
}

func testAccCheckIdentityV3RoleExists(ctx context.Context, n string, role *roles.Role) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		identityClient, err := config.IdentityV3Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack identity client: %w", err)
		}

		found, err := roles.Get(ctx, identityClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return errors.New("Role not found")
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
