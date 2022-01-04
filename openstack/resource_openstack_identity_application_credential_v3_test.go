package openstack

import (
	"fmt"
	"reflect"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/applicationcredentials"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
)

func TestAccIdentityV3ApplicationCredential_basic(t *testing.T) {
	var applicationCredential applicationcredentials.ApplicationCredential

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIdentityV3ApplicationCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3ApplicationCredentialBasic,
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
				Config: testAccIdentityV3ApplicationCredentialCustomSecret,
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

func TestAccIdentityV3ApplicationCredential_access_rules(t *testing.T) {
	var ac1, ac2 applicationcredentials.ApplicationCredential

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIdentityV3ApplicationCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3ApplicationCredentialAccessRules,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ApplicationCredentialExists("openstack_identity_application_credential_v3.app_cred_1", &ac1),
					testAccCheckIdentityV3ApplicationCredentialExists("openstack_identity_application_credential_v3.app_cred_1", &ac2),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_application_credential_v3.app_cred_1", "name", &ac1.Name),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_application_credential_v3.app_cred_1", "description", &ac2.Description),
					resource.TestCheckResourceAttr(
						"openstack_identity_application_credential_v3.app_cred_1", "unrestricted", "false"),
					resource.TestCheckResourceAttrSet(
						"openstack_identity_application_credential_v3.app_cred_1", "secret"),
					resource.TestCheckResourceAttrSet(
						"openstack_identity_application_credential_v3.app_cred_1", "project_id"),
					resource.TestCheckResourceAttr(
						"openstack_identity_application_credential_v3.app_cred_1", "expires_at", "2219-02-13T12:12:12Z"),
					resource.TestCheckResourceAttr(
						"openstack_identity_application_credential_v3.app_cred_1", "access_rules.#", "3"),
					resource.TestCheckResourceAttr(
						"openstack_identity_application_credential_v3.app_cred_2", "access_rules.#", "3"),
					testAccCheckIdentityV3ApplicationCredentialAccessRulesEqual(&ac1, &ac2),
				),
			},
		},
	})
}

func testAccCheckIdentityV3ApplicationCredentialDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	identityClient, err := config.IdentityV3Client(osRegionName)
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
		identityClient, err := config.IdentityV3Client(osRegionName)
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

func testAccCheckIdentityV3ApplicationCredentialAccessRulesEqual(ac1, ac2 *applicationcredentials.ApplicationCredential) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if !reflect.DeepEqual(ac1.AccessRules, ac2.AccessRules) {
			return fmt.Errorf("AccessRules are not equal: %v != %v", ac1.AccessRules, ac2.AccessRules)
		}
		return nil
	}
}

const testAccIdentityV3ApplicationCredentialBasic = `
resource "openstack_identity_application_credential_v3" "app_cred_1" {
  name        = "monitoring"
  description = "read-only technical user"
  roles       = ["reader"]
  expires_at  = "2219-02-13T12:12:12Z"
}
`

const testAccIdentityV3ApplicationCredentialCustomSecret = `
resource "openstack_identity_application_credential_v3" "app_cred_1" {
  name         = "super-admin"
  description  = "wheel technical user"
  secret       = "foo"
  unrestricted = true
}
`

const testAccIdentityV3ApplicationCredentialAccessRules = `
resource "openstack_identity_application_credential_v3" "app_cred_1" {
  name        = "monitoring"
  roles       = ["reader"]
  expires_at  = "2219-02-13T12:12:12Z"

  access_rules {
    path    = "/v2.0/metrics"
    service = "monitoring"
    method  = "GET"
  }

  access_rules {
    path    = "/v2.0/metrics"
    service = "monitoring"
    method  = "POST"
  }

  access_rules {
    path    = "/v2.0/metrics"
    service = "monitoring"
    method  = "PUT"
  }
}

resource "openstack_identity_application_credential_v3" "app_cred_2" {
  depends_on  = [ openstack_identity_application_credential_v3.app_cred_1 ]
  name        = "monitoring2"
  roles       = ["reader"]
  expires_at  = "2219-02-13T12:12:12Z"

  dynamic "access_rules" {
    for_each = [for rule in openstack_identity_application_credential_v3.app_cred_1.access_rules : {
      path = rule.path
      service = rule.service
      method = rule.method
    }]

    content {
      path = "${access_rules.value.path}"
      service = "${access_rules.value.service}"
      method = "${access_rules.value.method}"
    }
  }
}
`
