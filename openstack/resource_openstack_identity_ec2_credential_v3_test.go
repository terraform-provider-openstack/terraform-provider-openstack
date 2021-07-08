package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/extensions/ec2credentials"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
)

func TestAccIdentityV3Ec2Credential_basic(t *testing.T) {
	var Ec2Credential ec2credentials.Credential

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIdentityV3Ec2CredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3Ec2CredentialBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3Ec2CredentialExists("openstack_identity_ec2_credential_v3.ec2_cred_1", &Ec2Credential),
					resource.TestCheckResourceAttrSet(
						"openstack_identity_ec2_credential_v3.ec2_cred_1", "secret"),
					resource.TestCheckResourceAttrSet(
						"openstack_identity_ec2_credential_v3.ec2_cred_1", "project_id"),
					resource.TestCheckResourceAttrSet(
						"openstack_identity_ec2_credential_v3.ec2_cred_1", "access"),
					resource.TestCheckResourceAttrSet(
						"openstack_identity_ec2_credential_v3.ec2_cred_1", "user_id"),
				),
			},
		},
	})
}

func testAccCheckIdentityV3Ec2CredentialDestroy(s *terraform.State) error {
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
		if rs.Type != "openstack_identity_ec2_credential_v3" {
			continue
		}

		_, err := ec2credentials.Get(identityClient, user.ID, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Ec2Credential still exists")
		}
	}

	return nil
}

func testAccCheckIdentityV3Ec2CredentialExists(n string, Ec2Credential *ec2credentials.Credential) resource.TestCheckFunc {
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

		found, err := ec2credentials.Get(identityClient, user.ID, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Access != rs.Primary.ID {
			return fmt.Errorf("Ec2Credential not found")
		}

		*Ec2Credential = *found

		return nil
	}
}

const testAccIdentityV3Ec2CredentialBasic = `
resource "openstack_identity_ec2_credential_v3" "ec2_cred_1" {}
`
