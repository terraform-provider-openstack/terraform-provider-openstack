package openstack

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccSFSV2Share_importBasic(t *testing.T) {
	resourceName := "openstack_sharedfilesystem_share_v2.share_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckSFS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSFSV2ShareDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSV2ShareConfig_basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
