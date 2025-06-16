package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBlockStorageQosV3_importBasic(t *testing.T) {
	resourceName := "openstack_blockstorage_qos_v3.qos"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockStorageQosV3Destroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageQosV3Basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
