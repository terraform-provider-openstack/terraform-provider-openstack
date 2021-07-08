package openstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/db/v1/databases"
	"github.com/gophercloud/gophercloud/openstack/db/v1/instances"
)

func TestAccDatabaseV1Database_basic(t *testing.T) {
	var db databases.Database
	var instance instances.Instance

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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseV1InstanceExists(
						"openstack_db_instance_v1.basic", &instance),
					testAccCheckDatabaseV1DatabaseExists(
						"openstack_db_database_v1.basic", &instance, &db),
					resource.TestCheckResourceAttrPtr(
						"openstack_db_database_v1.basic", "name", &db.Name),
				),
			},
		},
	})
}

func testAccCheckDatabaseV1DatabaseExists(n string, instance *instances.Instance, db *databases.Database) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		parts := strings.SplitN(rs.Primary.ID, "/", 2)
		if len(parts) != 2 {
			return fmt.Errorf("Malformed database name: %s", rs.Primary.ID)
		}

		config := testAccProvider.Meta().(*Config)
		DatabaseV1Client, err := config.DatabaseV1Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack compute client: %s", err)
		}

		pages, err := databases.List(DatabaseV1Client, instance.ID).AllPages()
		if err != nil {
			return fmt.Errorf("Unable to retrieve databases: %s", err)
		}

		allDatabases, err := databases.ExtractDBs(pages)
		if err != nil {
			return fmt.Errorf("Unable to extract databases: %s", err)
		}

		for _, v := range allDatabases {
			if v.Name == parts[1] {
				*db = v
				return nil
			}
		}

		return fmt.Errorf("database %s does not exist", n)
	}
}

func testAccCheckDatabaseV1DatabaseDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	DatabaseV1Client, err := config.DatabaseV1Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_db_database_v1" {
			continue
		}

		parts := strings.SplitN(rs.Primary.ID, "/", 2)
		if len(parts) != 2 {
			return fmt.Errorf("Malformed database name: %s", rs.Primary.ID)
		}

		pages, err := databases.List(DatabaseV1Client, parts[0]).AllPages()
		if err != nil {
			return nil
		}

		allDatabases, err := databases.ExtractDBs(pages)
		if err != nil {
			return fmt.Errorf("Unable to extract databases: %s", err)
		}

		var exists bool
		for _, v := range allDatabases {
			if v.Name == parts[1] {
				exists = true
			}
		}

		if exists {
			return fmt.Errorf("database still exists")
		}
	}

	return nil
}

func testAccDatabaseV1DatabaseBasic() string {
	return fmt.Sprintf(`
resource "openstack_db_instance_v1" "basic" {
  name = "basic"
  size = 10

  datastore {
    version = "%s"
    type    = "%s"
  }

  network {
    uuid = "%s"
  }
}

resource "openstack_db_database_v1" "basic" {
  name        = "basic"
  instance_id = "${openstack_db_instance_v1.basic.id}"
}
`, osDBDatastoreVersion, osDBDatastoreType, osNetworkID)
}
