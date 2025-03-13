package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccKeyManagerOrderV1_importBasic(t *testing.T) {
	resourceName := "openstack_keymanager_order_v1.test-acc-basic"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
			testAccPreCheckKeyManager(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckOrderV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerOrderV1Symmetric,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
