package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccBlockStorageV3QosAssociation_importBasic(t *testing.T) {
	resourceName := "openstack_blockstorage_qos_association_v3.qos_association"

	var qosName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))
	var vtName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBlockstorageV3QosAssociationDestroy,
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
