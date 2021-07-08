package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBlockStorageV2Quotaset_importBasic(t *testing.T) {
	resourceName := "openstack_blockstorage_quotaset_v2.quotaset_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
			testAccPreCheckBlockStorageV2(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIdentityV3ProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageQuotasetV2Basic,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
