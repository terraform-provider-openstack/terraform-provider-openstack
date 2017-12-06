package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/gophercloud/gophercloud/openstack/db/v1/databases"
)

func TestAccDatabaseV1Database_basic(t *testing.T) {
	var db databases.Database

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckDatabase(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDatabaseV1DatabaseBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseV1DatabaseExists(
						"openstack_db_database_v1.basic", &db),
					resource.TestCheckResourceAttr(
						"openstack_db_database_v1.basic", "name", "basic"),
				),
			},
		},
	})
}

func testAccCheckDatabaseV1DatabaseExists(n string, db *databases.Database) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		databaseV1Client, err := config.databaseV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack compute client: %s", err)
		}

		dbInstance := s.RootModule().Resources["openstack_db_instance_v1.basic"]

		pages, err := databases.List(databaseV1Client, dbInstance.Primary.ID).AllPages()
		if err != nil {
			return fmt.Errorf("Unable to retrieve databases, pages: %s", err)
		}

		allDatabases, err := databases.ExtractDBs(pages)
		if err != nil {
			return fmt.Errorf("Unable to retrieve databases, extract: %s", err)
		}

		for _, v := range allDatabases {
			if v.Name == db.Name {
				return nil
			}
		}
		return nil
	}
}

var testAccDatabaseV1DatabaseBasic = fmt.Sprintf(`
resource "openstack_db_instance_v1" "basic" {
  name = "basic"
  datastore {
    version = "%s"
    type    = "%s"
  }

  network {
    uuid = "%s"
  }
  size = 10

}

resource "openstack_db_database_v1" "basic" {
  name     = "basic"
  instance = "${openstack_db_instance_v1.basic.id}"
}
`, OS_DB_DATASTORE_VERSION, OS_DB_DATASTORE_TYPE, OS_NETWORK_ID)
