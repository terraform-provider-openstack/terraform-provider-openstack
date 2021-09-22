package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccIdentityV3ApplicationCredential_importBasic(t *testing.T) {
	resourceName := "openstack_identity_application_credential_v3.app_cred_1"

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
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIdentityV3ApplicationCredentialDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3ApplicationCredentialCustomSecret,
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
