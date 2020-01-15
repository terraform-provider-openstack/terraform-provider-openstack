package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccKeyManagerSecretV1_importBasic(t *testing.T) {
	resourceName := "openstack_keymanager_secret_v1.secret_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckKeyManager(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecretV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerSecretV1_basic,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccKeyManagerSecretV1_importACLs(t *testing.T) {
	resourceName := "openstack_keymanager_secret_v1.secret_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckKeyManager(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecretV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerSecretV1_acls,
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"payload_content_encoding"},
			},
		},
	})
}
