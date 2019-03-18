package openstack

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccSecretV1_importBasic(t *testing.T) {
	resourceName := "openstack_keymanager_secret_v1.secret_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckKeyManager(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSecretV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeymanagerSecretV1_basic,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"payload",
					"payload_content_type",
				},
			},
		},
	})
}
