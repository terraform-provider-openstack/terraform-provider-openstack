package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
		CheckDestroy:      testAccCheckDatabaseV1DatabaseDestroy(t.Context()),
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
