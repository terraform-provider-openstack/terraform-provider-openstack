package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccIdentityV3Ec2Credential_importBasic(t *testing.T) {
	resourceName := "openstack_identity_ec2_credential_v3.ec2_cred_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityV3Ec2CredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3Ec2Credential_basic,
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
