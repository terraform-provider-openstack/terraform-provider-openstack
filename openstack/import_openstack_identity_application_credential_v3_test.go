package openstack

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccIdentityV3ApplicationCredential_importBasic(t *testing.T) {
	resourceName := "openstack_identity_application_credential_v3.app_cred_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityV3ApplicationCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3ApplicationCredential_basic,
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret"},
			},
		},
	})
}

func TestAccIdentityV3ApplicationCredential_importCustomSecret(t *testing.T) {
	resourceName := "openstack_identity_application_credential_v3.app_cred_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityV3ApplicationCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3ApplicationCredential_custom_secret,
			},

			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"secret"},
			},
		},
	})
}
