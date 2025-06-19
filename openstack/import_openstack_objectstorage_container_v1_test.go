package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccObjectStorageV1Container_importBasic(t *testing.T) {
	resourceName := "openstack_objectstorage_container_v1.container_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckSwift(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckObjectStorageV1ContainerDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccObjectStorageV1ContainerComplete,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"force_destroy",
					"content_type",
					"metadata",
				},
			},
		},
	})
}
