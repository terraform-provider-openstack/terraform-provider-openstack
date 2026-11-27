package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccIdentityV3Mapping_importBasic(t *testing.T) {
	resourceName := "openstack_identity_mapping_v3.mapping_1"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIdentityV3MappingDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3MappingBasic,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
