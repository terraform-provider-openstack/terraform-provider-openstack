package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBlockStorageV3VolumeTypeAccess_importBasic(t *testing.T) {
	resourceName := "openstack_blockstorage_volume_type_access_v3.volume_type_access"

	var projectName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))
	var vtName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockstorageV3VolumeTypeAccessDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockstorageV3VolumeTypeAccessBasic(projectName, vtName),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
