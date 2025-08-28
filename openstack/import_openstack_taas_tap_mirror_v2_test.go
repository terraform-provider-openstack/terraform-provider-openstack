package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccTaasTapMirrorV2_importBasic(t *testing.T) {
	resourceName := "openstack_taas_tap_mirror_v2.tap_mirror_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckTaas(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckTapMirrorV2Destroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccTapMirrorV2Basic,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
