package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/ec2credentials"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/tokens"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccIdentityV3Ec2Credential_basic(t *testing.T) {
	var Ec2Credential ec2credentials.Credential

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIdentityV3Ec2CredentialDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3Ec2CredentialBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3Ec2CredentialExists(t.Context(), "openstack_identity_ec2_credential_v3.ec2_cred_1", &Ec2Credential),
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

func testAccCheckIdentityV3Ec2CredentialDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		identityClient, err := config.IdentityV3Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack identity client: %w", err)
		}

		token := tokens.Get(ctx, identityClient, config.OsClient.TokenID)
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

			_, err := ec2credentials.Get(ctx, identityClient, user.ID, rs.Primary.ID).Extract()
			if err == nil {
				return errors.New("Ec2Credential still exists")
			}
		}

		return nil
	}
}

func testAccCheckIdentityV3Ec2CredentialExists(ctx context.Context, n string, ec2Credential *ec2credentials.Credential) resource.TestCheckFunc {
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

		token := tokens.Get(ctx, identityClient, config.OsClient.TokenID)
		if token.Err != nil {
			return token.Err
		}

		user, err := token.ExtractUser()
		if err != nil {
			return err
		}

		found, err := ec2credentials.Get(ctx, identityClient, user.ID, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.Access != rs.Primary.ID {
			return errors.New("Ec2Credential not found")
		}

		*ec2Credential = *found

		return nil
	}
}

const testAccIdentityV3Ec2CredentialBasic = `
resource "openstack_identity_ec2_credential_v3" "ec2_cred_1" {}
`
