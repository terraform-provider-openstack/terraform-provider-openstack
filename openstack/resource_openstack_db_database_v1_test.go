package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/databases"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/instances"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
		CheckDestroy:      testAccCheckDatabaseV1DatabaseDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseV1DatabaseBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseV1InstanceExists(t.Context(),
						"openstack_db_instance_v1.basic", &instance),
					testAccCheckDatabaseV1DatabaseExists(t.Context(),
						"openstack_db_database_v1.basic", &instance, &db),
					resource.TestCheckResourceAttrPtr(
						"openstack_db_database_v1.basic", "name", &db.Name),
				),
			},
		},
	})
}

func testAccCheckDatabaseV1DatabaseExists(ctx context.Context, n string, instance *instances.Instance, db *databases.Database) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No ID is set")
		}

		_, userName, err := parsePairedIDs(rs.Primary.ID, "openstack_db_database_v1")
		if err != nil {
			return err
		}

		config := testAccProvider.Meta().(*Config)

		databaseV1Client, err := config.DatabaseV1Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack compute client: %w", err)
		}

		pages, err := databases.List(databaseV1Client, instance.ID).AllPages(ctx)
		if err != nil {
			return fmt.Errorf("Unable to retrieve databases: %w", err)
		}

		allDatabases, err := databases.ExtractDBs(pages)
		if err != nil {
			return fmt.Errorf("Unable to extract databases: %w", err)
		}

		for _, v := range allDatabases {
			if v.Name == userName {
				*db = v

				return nil
			}
		}

		return fmt.Errorf("database %s does not exist", n)
	}
}

func testAccCheckDatabaseV1DatabaseDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		databaseV1Client, err := config.DatabaseV1Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack compute client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_db_database_v1" {
				continue
			}

			id, userName, err := parsePairedIDs(rs.Primary.ID, "openstack_db_database_v1")
			if err != nil {
				return err
			}

			pages, err := databases.List(databaseV1Client, id).AllPages(ctx)
			if err != nil {
				return nil
			}

			allDatabases, err := databases.ExtractDBs(pages)
			if err != nil {
				return fmt.Errorf("Unable to extract databases: %w", err)
			}

			var exists bool

			for _, v := range allDatabases {
				if v.Name == userName {
					exists = true
				}
			}

			if exists {
				return errors.New("database still exists")
			}
		}

		return nil
	}
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
  instance_id = openstack_db_instance_v1.basic.id
}
`, osDBDatastoreVersion, osDBDatastoreType, osNetworkID)
}
