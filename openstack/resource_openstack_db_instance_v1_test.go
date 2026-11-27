package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/configurations"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/instances"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDatabaseV1Instance_basic(t *testing.T) {
	var instance instances.Instance

	var configuration configurations.Config

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckDatabase(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDatabaseV1InstanceDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseV1InstanceBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseV1InstanceExists(t.Context(),
						"openstack_db_instance_v1.basic", &instance),
					resource.TestCheckResourceAttrPtr(
						"openstack_db_instance_v1.basic", "name", &instance.Name),
					resource.TestCheckResourceAttr(
						"openstack_db_instance_v1.basic", "volume_type", "lvmdriver-1"),
					resource.TestCheckResourceAttr(
						"openstack_db_instance_v1.basic", "user.0.name", "testuser"),
					resource.TestCheckResourceAttr(
						"openstack_db_instance_v1.basic", "user.0.password", "testpassword"),
					resource.TestCheckResourceAttr(
						"openstack_db_instance_v1.basic", "database.0.name", "testdb1"),
					resource.TestCheckResourceAttr(
						"openstack_db_instance_v1.basic", "database.0.charset", "utf8"),
					resource.TestCheckResourceAttr(
						"openstack_db_instance_v1.basic", "database.0.collate", "utf8_general_ci"),
					resource.TestCheckResourceAttr(
						"openstack_db_instance_v1.basic", "database.1.name", "testdb2"),
					resource.TestCheckResourceAttr(
						"openstack_db_instance_v1.basic", "database.1.charset", "utf8"),
					resource.TestCheckResourceAttr(
						"openstack_db_instance_v1.basic", "database.1.collate", "utf8_general_ci"),
					resource.TestCheckResourceAttrSet(
						"openstack_db_instance_v1.basic", "configuration_id"),
					testAccCheckDatabaseV1ConfigurationExists(t.Context(),
						"openstack_db_configuration_v1.basic", &configuration),
				),
			},
		},
	})
}

func TestAccDatabaseV1Instance_resizeVolume(t *testing.T) {
	var (
		instance   instances.Instance
		instanceID string
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckDatabase(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDatabaseV1InstanceDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseV1InstanceBasicWithSize(10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseV1InstanceExists(t.Context(),
						"openstack_db_instance_v1.basic", &instance),
					testAccCheckDatabaseV1InstanceID("openstack_db_instance_v1.basic", &instanceID),
					resource.TestCheckResourceAttr("openstack_db_instance_v1.basic", "size", "10"),
				),
			},
			{
				Config: testAccDatabaseV1InstanceBasicWithSize(11),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseV1InstanceExists(t.Context(),
						"openstack_db_instance_v1.basic", &instance),
					testAccCheckDatabaseV1InstanceID("openstack_db_instance_v1.basic", &instanceID),
					resource.TestCheckResourceAttr("openstack_db_instance_v1.basic", "size", "11"),
				),
			},
		},
	})
}

func testAccCheckDatabaseV1InstanceID(name string, expectedID *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resourceState, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		if *expectedID == "" {
			*expectedID = resourceState.Primary.ID

			return nil
		}

		if resourceState.Primary.ID != *expectedID {
			return fmt.Errorf("database instance was recreated: got ID %s, want %s", resourceState.Primary.ID, *expectedID)
		}

		return nil
	}
}

func testAccCheckDatabaseV1InstanceExists(ctx context.Context, n string, instance *instances.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		databaseV1Client, err := config.DatabaseV1Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack compute client: %w", err)
		}

		found, err := instances.Get(ctx, databaseV1Client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return errors.New("Instance not found")
		}

		*instance = *found

		return nil
	}
}

func testAccCheckDatabaseV1InstanceDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		databaseV1Client, err := config.DatabaseV1Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack compute client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_db_instance_v1" {
				continue
			}

			_, err := instances.Get(ctx, databaseV1Client, rs.Primary.ID).Extract()
			if err == nil {
				return errors.New("Instance still exists")
			}
		}

		return nil
	}
}

func testAccDatabaseV1InstanceBasic() string {
	return testAccDatabaseV1InstanceBasicWithSize(10)
}

func testAccDatabaseV1InstanceBasicWithSize(size int) string {
	return fmt.Sprintf(`
resource "openstack_db_instance_v1" "basic" {
  name             = "basic"
  configuration_id = openstack_db_configuration_v1.basic.id

  datastore {
    version = "%[1]s"
    type    = "%[2]s"
  }

  network {
    uuid = "%[3]s"
  }

  size = %[4]d
  volume_type = "lvmdriver-1"

  database {
    name    = "testdb1"
    charset = "utf8"
    collate = "utf8_general_ci"
  }

  database {
    name    = "testdb2"
    charset = "utf8"
    collate = "utf8_general_ci"
  }

  user {
    name      = "testuser"
    password  = "testpassword"
    databases = ["testdb1"]
    host      = "%%"
  }

}

resource "openstack_db_configuration_v1" "basic" {
  name        = "basic"
  description = "test"

  datastore {
    version = "%[1]s"
    type    = "%[2]s"
  }

  configuration {
    name  = "collation_server"
    value = "latin1_swedish_ci"
  }

  configuration {
    name  = "collation_database"
    value = "latin1_swedish_ci"
  }

  configuration {
    name  = "max_connections"
    value = 200
  }
}
`, osDBDatastoreVersion, osDBDatastoreType, osNetworkID, size)
}
