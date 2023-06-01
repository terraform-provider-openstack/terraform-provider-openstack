package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/osinherit"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/roles"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/users"
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
		CheckDestroy:      testAccCheckIdentityV3InheritRoleAssignmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3InheritRoleAssignmentBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3InheritRoleAssignmentExists("openstack_identity_inherit_role_assignment_v3.role_assignment_1", &role, &user),
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

func testAccCheckIdentityV3InheritRoleAssignmentDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	identityClient, err := config.IdentityV3Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_identity_inherit_role_assignment_v3" {
			continue
		}

		domainID, projectID, groupID, userID, roleID, err := identityRoleAssignmentV3ParseID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error determining openstack_identity_inherit_role_assignment_v3 ID: %s", err)
		}

		var opts = osinherit.ValidateOpts{
			GroupID:   groupID,
			DomainID:  domainID,
			ProjectID: projectID,
			UserID:    userID,
		}

		err = osinherit.Validate(identityClient, roleID, opts).ExtractErr()
		if err == nil {
			return fmt.Errorf("Inherit Role assignment still exists")
		}
	}

	return nil
}

func testAccCheckIdentityV3InheritRoleAssignmentExists(n string, role *roles.Role, user *users.User) resource.TestCheckFunc {
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

		domainID, projectID, groupID, userID, roleID, err := identityRoleAssignmentV3ParseID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error determining openstack_identity_inherit_role_assignment_v3 ID: %s", err)
		}

		var opts = osinherit.ValidateOpts{
			GroupID:   groupID,
			DomainID:  domainID,
			ProjectID: projectID,
			UserID:    userID,
		}

		err = osinherit.Validate(identityClient, roleID, opts).ExtractErr()
		if err != nil {
			return err
		}

		u, err := users.Get(identityClient, userID).Extract()
		if err != nil {
			return fmt.Errorf("User not found")
		}
		*user = *u
		r, err := roles.Get(identityClient, roleID).Extract()
		if err != nil {
			return fmt.Errorf("Role not found")
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
  user_id = "${openstack_identity_user_v3.user_1.id}"
  domain_id = "default"
  role_id = "${openstack_identity_role_v3.role_1.id}"
}
`
