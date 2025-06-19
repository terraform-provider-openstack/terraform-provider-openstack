package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccObjectStorageV1Account_importBasic(t *testing.T) {
	resourceName := "openstack_objectstorage_account_v1.account_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckSwift(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckObjectStorageV1AccountDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccObjectStorageV1AccountBasic,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"headers",
				},
			},
		},
	})
}
