package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/gophercloud/gophercloud/openstack/db/v1/instances"
)

func TestAccDatabaseV1Instance_basic(t *testing.T) {
	var instance instances.Instance

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckDatabase(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDatabaseV1InstanceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseV1InstanceExists(
						"openstack_db_instance_v1.basic", &instance),
					resource.TestCheckResourceAttr(
						"openstack_db_instance_v1.basic", "name", "basic"),
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
		databaseV1Client, err := config.databaseV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack compute client: %s", err)
		}

		found, err := instances.Get(databaseV1Client, rs.Primary.ID).Extract()
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

var testAccDatabaseV1InstanceBasic = fmt.Sprintf(`
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
`, OS_DB_DATASTORE_VERSION, OS_DB_DATASTORE_TYPE, OS_NETWORK_ID)
