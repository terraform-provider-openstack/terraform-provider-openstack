package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccObjectStorageV1Account_importBasic(t *testing.T) {
	resourceName := "openstack_objectstorage_account_v1.account_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckSwift(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckObjectStorageV1AccountDestroy,
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
