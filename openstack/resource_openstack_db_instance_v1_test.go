package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/db/v1/configurations"
	"github.com/gophercloud/gophercloud/openstack/db/v1/instances"
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
		CheckDestroy:      testAccCheckDatabaseV1InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseV1InstanceBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseV1InstanceExists(
						"openstack_db_instance_v1.basic", &instance),
					resource.TestCheckResourceAttrPtr(
						"openstack_db_instance_v1.basic", "name", &instance.Name),
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
					testAccCheckDatabaseV1ConfigurationExists(
						"openstack_db_configuration_v1.basic", &configuration),
				),
			},
		},
	})
}

func testAccCheckDatabaseV1InstanceExists(n string, instance *instances.Instance) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		DatabaseV1Client, err := config.DatabaseV1Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack compute client: %s", err)
		}

		found, err := instances.Get(DatabaseV1Client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Instance not found")
		}

		*instance = *found

		return nil
	}
}

func testAccCheckDatabaseV1InstanceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	DatabaseV1Client, err := config.DatabaseV1Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_db_instance_v1" {
			continue
		}

		_, err := instances.Get(DatabaseV1Client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Instance still exists")
		}
	}

	return nil
}

func testAccDatabaseV1InstanceBasic() string {
	return fmt.Sprintf(`
resource "openstack_db_instance_v1" "basic" {
  name             = "basic"
  configuration_id = "${openstack_db_configuration_v1.basic.id}"

  datastore {
    version = "%[1]s"
    type    = "%[2]s"
  }

  network {
    uuid = "%[3]s"
  }

  size = 10

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
`, osDBDatastoreVersion, osDBDatastoreType, osNetworkID)
}
