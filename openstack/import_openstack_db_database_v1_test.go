package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatabaseV1Database_importBasic(t *testing.T) {
	resourceName := "openstack_db_database_v1.basic"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckDatabase(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDatabaseV1DatabaseDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseV1DatabaseBasic(),
			},

			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"region",
				},
			},
		},
	})
}
