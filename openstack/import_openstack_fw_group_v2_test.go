package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccFWGroupV2_importBasic(t *testing.T) {
	resourceName := "openstack_fw_group_v2.group_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckFW(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckFWGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWGroupV2Basic1,
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
