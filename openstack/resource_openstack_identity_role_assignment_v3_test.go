package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/roles"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/users"
	"github.com/gophercloud/gophercloud/v2/pagination"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccIdentityV3RoleAssignment_basic(t *testing.T) {
	var role roles.Role

	var user users.User

	var project projects.Project

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIdentityV3RoleAssignmentDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3RoleAssignmentBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3RoleAssignmentExists(t.Context(), "openstack_identity_role_assignment_v3.role_assignment_1", &role, &user, &project),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_role_assignment_v3.role_assignment_1", "project_id", &project.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_role_assignment_v3.role_assignment_1", "user_id", &user.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_role_assignment_v3.role_assignment_1", "role_id", &role.ID),
				),
			},
		},
	})
}

func testAccCheckIdentityV3RoleAssignmentDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		identityClient, err := config.IdentityV3Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack identity client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_identity_role_assignment_v3" {
				continue
			}

			_, err := roles.Get(ctx, identityClient, rs.Primary.ID).Extract()
			if err == nil {
				return errors.New("Role assignment still exists")
			}
		}

		return nil
	}
}

func testAccCheckIdentityV3RoleAssignmentExists(ctx context.Context, n string, role *roles.Role, user *users.User, project *projects.Project) resource.TestCheckFunc {
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
			return fmt.Errorf("Error determining openstack_identity_role_assignment_v3 ID: %w", err)
		}

		opts := roles.ListAssignmentsOpts{
			GroupID:        groupID,
			ScopeDomainID:  domainID,
			ScopeProjectID: projectID,
			UserID:         userID,
		}

		pager := roles.ListAssignments(identityClient, opts)

		var assignment roles.RoleAssignment

		err = pager.EachPage(ctx, func(_ context.Context, page pagination.Page) (bool, error) {
			assignmentList, err := roles.ExtractRoleAssignments(page)
			if err != nil {
				return false, err
			}

			for _, a := range assignmentList {
				if a.Role.ID == roleID {
					assignment = a

					return false, nil
				}
			}

			return true, nil
		})
		if err != nil {
			return err
		}

		p, err := projects.Get(ctx, identityClient, assignment.Scope.Project.ID).Extract()
		if err != nil {
			return errors.New("Project not found")
		}

		*project = *p

		u, err := users.Get(ctx, identityClient, assignment.User.ID).Extract()
		if err != nil {
			return errors.New("User not found")
		}

		*user = *u

		r, err := roles.Get(ctx, identityClient, assignment.Role.ID).Extract()
		if err != nil {
			return errors.New("Role not found")
		}

		*role = *r

		return nil
	}
}

const testAccIdentityV3RoleAssignmentBasic = `
resource "openstack_identity_project_v3" "project_1" {
  name = "project_1"
}

resource "openstack_identity_user_v3" "user_1" {
  name = "user_1"
  default_project_id = openstack_identity_project_v3.project_1.id
}

resource "openstack_identity_role_v3" "role_1" {
  name = "role_1"
}

resource "openstack_identity_role_assignment_v3" "role_assignment_1" {
  user_id = openstack_identity_user_v3.user_1.id
  project_id = openstack_identity_project_v3.project_1.id
  role_id = openstack_identity_role_v3.role_1.id
}
`
