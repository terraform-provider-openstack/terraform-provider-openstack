package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/gophercloud/gophercloud/openstack/db/v1/configurations"
)

func TestAccDatabaseV1Configuration_basic(t *testing.T) {
	var configuration configurations.Config

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckDatabase(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDatabaseV1ConfigurationBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseV1ConfigurationExists(
						"openstack_db_configuration_v1.basic", &configuration),
					resource.TestCheckResourceAttr(
						"openstack_db_configuration_v1.basic", "name", "basic"),
					resource.TestCheckResourceAttr(
						"openstack_db_configuration_v1.basic", "configuration.2.name", "max_connections"),
					resource.TestCheckResourceAttr(
						"openstack_db_configuration_v1.basic", "configuration.2.value", "200"),
				),
			},
		},
	})
}

func testAccCheckDatabaseV1ConfigurationExists(n string, configuration *configurations.Config) resource.TestCheckFunc {
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

		found, err := configurations.Get(databaseV1Client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Configuration not found")
		}

		*configuration = *found

		return nil
	}
}

var testAccDatabaseV1ConfigurationBasic = fmt.Sprintf(`
resource "openstack_db_configuration_v1" "basic" {
  name        = "basic"
  description = "test"

  datastore {
    version = "%s"
    type    = "%s"
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
`, OS_DB_DATASTORE_VERSION, OS_DB_DATASTORE_TYPE)
