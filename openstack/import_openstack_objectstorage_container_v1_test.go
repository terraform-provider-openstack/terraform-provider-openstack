package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccObjectStorageV1Container_importBasic(t *testing.T) {
	resourceName := "openstack_objectstorage_container_v1.container_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckSwift(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckObjectStorageV1ContainerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccObjectStorageV1Container_complete,
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
