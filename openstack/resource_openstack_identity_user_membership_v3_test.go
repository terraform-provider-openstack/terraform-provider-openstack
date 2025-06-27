package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/groups"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/users"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccIdentityV3UserMembership_basic(t *testing.T) {
	var group groups.Group

	groupName := "ACCPTTEST-" + acctest.RandString(5)

	var user users.User

	userName := "ACCPTTEST-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIdentityV3UserMembershipDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3UserMembershipBasic(groupName, userName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3UserExists(t.Context(), "openstack_identity_user_v3.user_1", &user),
					testAccCheckIdentityV3GroupExists(t.Context(), "openstack_identity_group_v3.group_1", &group),
					testAccCheckIdentityV3UserMembershipExists(t.Context(), "openstack_identity_user_membership_v3.user_membership_1"),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_user_membership_v3.user_membership_1", "user_id", &user.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_user_membership_v3.user_membership_1", "group_id", &group.ID),
				),
			},
		},
	})
}

func testAccCheckIdentityV3UserMembershipDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		identityClient, err := config.IdentityV3Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack identity client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_identity_user_membership_v3" {
				continue
			}

			uid, gid, err := parsePairedIDs(rs.Primary.ID, "openstack_identity_user_membership_v3")
			if err != nil {
				return err
			}

			um, err := users.IsMemberOfGroup(ctx, identityClient, gid, uid).Extract()
			if err == nil && um {
				return errors.New("User membership still exists")
			}
		}

		return nil
	}
}

func testAccCheckIdentityV3UserMembershipExists(ctx context.Context, n string) resource.TestCheckFunc {
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

		uid, gid, err := parsePairedIDs(rs.Primary.ID, "openstack_identity_user_membership_v3")
		if err != nil {
			return err
		}

		um, err := users.IsMemberOfGroup(ctx, identityClient, gid, uid).Extract()
		if err != nil || !um {
			return errors.New("User membership not found")
		}

		return nil
	}
}

func testAccIdentityV3UserMembershipBasic(groupName, userName string) string {
	return fmt.Sprintf(`
	resource "openstack_identity_group_v3" "group_1" {
	name = "%s"
	}

	resource "openstack_identity_user_v3" "user_1" {
	name = "%s"
	}

	resource "openstack_identity_user_membership_v3" "user_membership_1" {
	user_id = openstack_identity_user_v3.user_1.id
	group_id = openstack_identity_group_v3.group_1.id
	}
    `, groupName, userName)
}
