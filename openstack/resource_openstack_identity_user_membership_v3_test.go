package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/groups"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/users"
)

func TestAccIdentityV3UserMembership_basic(t *testing.T) {
	var group groups.Group
	var groupName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))

	var user users.User
	var userName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIdentityV3UserMembershipDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3UserMembershipBasic(groupName, userName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3UserExists("openstack_identity_user_v3.user_1", &user),
					testAccCheckIdentityV3GroupExists("openstack_identity_group_v3.group_1", &group),
					testAccCheckIdentityV3UserMembershipExists("openstack_identity_user_membership_v3.user_membership_1"),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_user_membership_v3.user_membership_1", "user_id", &user.ID),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_user_membership_v3.user_membership_1", "group_id", &group.ID),
				),
			},
		},
	})
}

func testAccCheckIdentityV3UserMembershipDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	identityClient, err := config.IdentityV3Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_identity_user_membership_v3" {
			continue
		}

		uid, gid, err := parseUserMembershipID(rs.Primary.ID)
		if err != nil {
			return err
		}

		um, err := users.IsMemberOfGroup(identityClient, gid, uid).Extract()
		if err == nil && um {
			return fmt.Errorf("User membership still exists")
		}
	}

	return nil
}

func testAccCheckIdentityV3UserMembershipExists(n string) resource.TestCheckFunc {
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

		uid, gid, err := parseUserMembershipID(rs.Primary.ID)
		if err != nil {
			return err
		}

		um, err := users.IsMemberOfGroup(identityClient, gid, uid).Extract()
		if err != nil || !um {
			return fmt.Errorf("User membership not found")
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
	user_id = "${openstack_identity_user_v3.user_1.id}"
	group_id = "${openstack_identity_group_v3.group_1.id}"
	}
    `, groupName, userName)
}
