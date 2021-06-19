package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIdentityV3Ec2Credential_importBasic(t *testing.T) {
	resourceName := "openstack_identity_ec2_credential_v3.ec2_cred_1"

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
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret", "project_id"},
			},
		},
	})
}
