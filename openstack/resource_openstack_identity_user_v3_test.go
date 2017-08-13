package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/users"
)

func TestAccIdentityV3User_basic(t *testing.T) {
	var project projects.Project
	var projectName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))
	var user users.User
	var userName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityV3UserDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccIdentityV3User_basic(projectName, userName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3UserExists("openstack_identity_user_v3.user_1", &user),
					testAccCheckIdentityV3ProjectExists("openstack_identity_project_v3.project_1", &project),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "name", userName),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "description", "A user"),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "enabled", "true"),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "ignore_change_password_upon_first_use", "true"),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "multi_factor_auth_enabled", "true"),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "multi_factor_auth_rule.#", "2"),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "multi_factor_auth_rule.0.rule.0", "password"),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "multi_factor_auth_rule.0.rule.1", "totp"),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "multi_factor_auth_rule.1.rule.0", "password"),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "multi_factor_auth_rule.1.rule.1", "custom-auth-method"),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "extra.email", "jdoe@example.com"),
				),
			},
			resource.TestStep{
				Config: testAccIdentityV3User_update(projectName, userName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3UserExists("openstack_identity_user_v3.user_1", &user),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "name", userName),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "description", "Some user"),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "enabled", "false"),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "ignore_change_password_upon_first_use", "false"),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "multi_factor_auth_rule.#", "1"),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "multi_factor_auth_rule.0.rule.0", "password"),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "multi_factor_auth_rule.0.rule.1", "totp"),
					resource.TestCheckResourceAttr(
						"openstack_identity_user_v3.user_1", "extra.email", "jdoe@foobar.com"),
				),
			},
		},
	})
}

func testAccCheckIdentityV3UserDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	identityClient, err := config.identityV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_identity_user_v3" {
			continue
		}

		_, err := users.Get(identityClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("User still exists")
		}
	}

	return nil
}

func testAccCheckIdentityV3UserExists(n string, user *users.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		identityClient, err := config.identityV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack identity client: %s", err)
		}

		found, err := users.Get(identityClient, rs.Primary.ID).Extract()
		if err == nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("User not found")
		}

		*user = *found

		return nil
	}
}

func testAccIdentityV3User_basic(projectName, userName string) string {
	return fmt.Sprintf(`
    resource "openstack_identity_project_v3" "project_1" {
      name = "%s"
    }

    resource "openstack_identity_user_v3" "user_1" {
      default_project_id = "${openstack_identity_project_v3.project_1.id}"
      name = "%s"
      description = "A user"
      password = "password123"
      ignore_change_password_upon_first_use = true
      multi_factor_auth_enabled = true

      multi_factor_auth_rule {
        rule = ["password", "totp"]
      }

      multi_factor_auth_rule {
        rule = ["password", "custom-auth-method"]
      }

      extra {
        email = "jdoe@example.com"
      }
    }
  `, projectName, userName)
}

func testAccIdentityV3User_update(projectName, userName string) string {
	return fmt.Sprintf(`
    resource "openstack_identity_project_v3" "project_1" {
      name = "%s"
    }

    resource "openstack_identity_user_v3" "user_1" {
      default_project_id = "${openstack_identity_project_v3.project_1.id}"
      name = "%s"
      description = "Some user"
      enabled = false
      password = "password123"
      ignore_change_password_upon_first_use = false
      multi_factor_auth_enabled = true

      multi_factor_auth_rule {
        rule = ["password", "totp"]
      }

      extra {
        email = "jdoe@foobar.com"
      }
    }
  `, projectName, userName)
}
