package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccComputeQuotasetV2_importBasic(t *testing.T) {
	resourceName := "openstack_compute_quotaset_v2.quotaset_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityV3ProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeQuotasetV2_basic,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
