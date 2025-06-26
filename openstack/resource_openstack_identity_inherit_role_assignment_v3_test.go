package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/osinherit"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/roles"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/users"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccIdentityV3InheritRoleAssignment_basic(t *testing.T) {
	var role roles.Role

	var user users.User

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIdentityV3InheritRoleAssignmentDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3InheritRoleAssignmentBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3InheritRoleAssignmentExists(t.Context(), "openstack_identity_inherit_role_assignment_v3.role_assignment_1", &role, &user),
					resource.TestCheckResourceAttr(
						"openstack_identity_inherit_role_assignment_v3.role_assignment_1", "domain_id", "default"),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_inherit_role_assignment_v3.role_assignment_1", "user_id", &user.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_inherit_role_assignment_v3.role_assignment_1", "role_id", &role.ID),
				),
			},
		},
	})
}

func testAccCheckIdentityV3InheritRoleAssignmentDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		identityClient, err := config.IdentityV3Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack identity client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_identity_inherit_role_assignment_v3" {
				continue
			}

			domainID, projectID, groupID, userID, roleID, err := identityRoleAssignmentV3ParseID(rs.Primary.ID)
			if err != nil {
				return fmt.Errorf("Error determining openstack_identity_inherit_role_assignment_v3 ID: %w", err)
			}

			opts := osinherit.ValidateOpts{
				GroupID:   groupID,
				DomainID:  domainID,
				ProjectID: projectID,
				UserID:    userID,
			}

			err = osinherit.Validate(ctx, identityClient, roleID, opts).ExtractErr()
			if err == nil {
				return errors.New("Inherit Role assignment still exists")
			}
		}

		return nil
	}
}

func testAccCheckIdentityV3InheritRoleAssignmentExists(ctx context.Context, n string, role *roles.Role, user *users.User) resource.TestCheckFunc {
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

		domainID, projectID, groupID, userID, roleID, err := identityRoleAssignmentV3ParseID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error determining openstack_identity_inherit_role_assignment_v3 ID: %w", err)
		}

		opts := osinherit.ValidateOpts{
			GroupID:   groupID,
			DomainID:  domainID,
			ProjectID: projectID,
			UserID:    userID,
		}

		err = osinherit.Validate(ctx, identityClient, roleID, opts).ExtractErr()
		if err != nil {
			return err
		}

		u, err := users.Get(ctx, identityClient, userID).Extract()
		if err != nil {
			return errors.New("User not found")
		}

		*user = *u

		r, err := roles.Get(ctx, identityClient, roleID).Extract()
		if err != nil {
			return errors.New("Role not found")
		}

		*role = *r

		return nil
	}
}

const testAccIdentityV3InheritRoleAssignmentBasic = `
resource "openstack_identity_user_v3" "user_1" {
  name = "user_1"
  domain_id = "default"
}

resource "openstack_identity_role_v3" "role_1" {
  name = "role_1"
  domain_id = "default"
}

resource "openstack_identity_inherit_role_assignment_v3" "role_assignment_1" {
  user_id = openstack_identity_user_v3.user_1.id
  domain_id = "default"
  role_id = openstack_identity_role_v3.role_1.id
}
`
