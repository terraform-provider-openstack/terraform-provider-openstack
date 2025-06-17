package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBlockStorageV3VolumeTypeAccess_importBasic(t *testing.T) {
	resourceName := "openstack_blockstorage_volume_type_access_v3.volume_type_access"

	projectName := "ACCPTTEST-" + acctest.RandString(5)

	vtName := "ACCPTTEST-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockstorageV3VolumeTypeAccessDestroy(t.Context()),
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
