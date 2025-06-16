package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBlockStorageV3QosAssociation_importBasic(t *testing.T) {
	resourceName := "openstack_blockstorage_qos_association_v3.qos_association"

	qosName := "ACCPTTEST-" + acctest.RandString(5)

	vtName := "ACCPTTEST-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockstorageV3QosAssociationDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccBlockstorageV3QosAssociationBasic(qosName, vtName),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
