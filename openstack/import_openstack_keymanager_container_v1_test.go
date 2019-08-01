package openstack

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccKeyManagerContainerV1_importBasic(t *testing.T) {
	resourceName := "openstack_keymanager_container_v1.container_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckKeyManager(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckContainerV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKeyManagerContainerV1_basic,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
