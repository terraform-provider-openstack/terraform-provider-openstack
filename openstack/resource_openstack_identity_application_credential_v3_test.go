package openstack

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/applicationcredentials"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
)

func TestAccIdentityV3ApplicationCredential_basic(t *testing.T) {
	var applicationCredential applicationcredentials.ApplicationCredential

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityV3ApplicationCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3ApplicationCredential_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ApplicationCredentialExists("openstack_identity_application_credential_v3.app_cred_1", &applicationCredential),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_application_credential_v3.app_cred_1", "name", &applicationCredential.Name),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_application_credential_v3.app_cred_1", "description", &applicationCredential.Description),
					resource.TestCheckResourceAttr(
						"openstack_identity_application_credential_v3.app_cred_1", "unrestricted", "false"),
					resource.TestCheckResourceAttrSet(
						"openstack_identity_application_credential_v3.app_cred_1", "secret"),
					resource.TestCheckResourceAttrSet(
						"openstack_identity_application_credential_v3.app_cred_1", "project_id"),
					resource.TestCheckResourceAttr(
						"openstack_identity_application_credential_v3.app_cred_1", "expires_at", "2219-02-13T12:12:12Z"),
					resource.TestCheckResourceAttr(
						"openstack_identity_application_credential_v3.app_cred_1", "roles.#", "1"),
					testAccCheckIdentityV3ApplicationCredentialRoleNameExists("reader", &applicationCredential),
				),
			},
			{
				Config: testAccIdentityV3ApplicationCredential_custom_secret,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ApplicationCredentialExists("openstack_identity_application_credential_v3.app_cred_1", &applicationCredential),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_application_credential_v3.app_cred_1", "name", &applicationCredential.Name),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_application_credential_v3.app_cred_1", "description", &applicationCredential.Description),
					resource.TestCheckResourceAttr(
						"openstack_identity_application_credential_v3.app_cred_1", "unrestricted", "true"),
					resource.TestCheckResourceAttr(
						"openstack_identity_application_credential_v3.app_cred_1", "secret", "foo"),
					resource.TestCheckResourceAttrSet(
						"openstack_identity_application_credential_v3.app_cred_1", "project_id"),
					resource.TestCheckResourceAttr(
						"openstack_identity_application_credential_v3.app_cred_1", "expires_at", ""),
					resource.TestMatchResourceAttr(
						"openstack_identity_application_credential_v3.app_cred_1", "roles.#", regexp.MustCompile("^[2-9]\\d*")),
				),
			},
		},
	})
}

func testAccCheckIdentityV3ApplicationCredentialDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	identityClient, err := config.identityV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	token := tokens.Get(identityClient, config.OsClient.TokenID)
	if token.Err != nil {
		return token.Err
	}

	user, err := token.ExtractUser()
	if err != nil {
		return err
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_identity_application_credential_v3" {
			continue
		}

		_, err := applicationcredentials.Get(identityClient, user.ID, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("ApplicationCredential still exists")
		}
	}

	return nil
}

func testAccCheckIdentityV3ApplicationCredentialExists(n string, applicationCredential *applicationcredentials.ApplicationCredential) resource.TestCheckFunc {
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

		token := tokens.Get(identityClient, config.OsClient.TokenID)
		if token.Err != nil {
			return token.Err
		}

		user, err := token.ExtractUser()
		if err != nil {
			return err
		}

		found, err := applicationcredentials.Get(identityClient, user.ID, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("ApplicationCredential not found")
		}

		*applicationCredential = *found

		return nil
	}
}

func testAccCheckIdentityV3ApplicationCredentialRoleNameExists(role string, applicationCredential *applicationcredentials.ApplicationCredential) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		roles := flattenIdentityApplicationCredentialRolesV3(applicationCredential.Roles)
		exists := strSliceContains(roles, role)
		if exists {
			return nil
		}
		return fmt.Errorf("The %s role was not found in %+q", role, roles)
	}
}

const testAccIdentityV3ApplicationCredential_basic = `
resource "openstack_identity_application_credential_v3" "app_cred_1" {
  name        = "monitoring"
  description = "read-only technical user"
  roles       = ["reader"]
  expires_at  = "2219-02-13T12:12:12Z"
}
`

const testAccIdentityV3ApplicationCredential_custom_secret = `
resource "openstack_identity_application_credential_v3" "app_cred_1" {
  name         = "super-admin"
  description  = "wheel technical user"
  secret       = "foo"
  unrestricted = true
}
`
