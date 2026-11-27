package openstack

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/configurations"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/instances"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
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

func TestAccDatabaseV1Instance_resize(t *testing.T) {
	var (
		instance   instances.Instance
		instanceID string
	)

	instanceName := acctest.RandomWithPrefix("tf-acc-db-resize")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckDatabase(t)

			if osFlavorID == "" || osFlavorIDResize == "" {
				t.Skip("OS_FLAVOR_ID and OS_FLAVOR_ID_RESIZE must be set to resize a database instance")
			}
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDatabaseV1InstanceDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseV1InstanceResize(instanceName, osFlavorID, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseV1InstanceExists(t.Context(),
						"openstack_db_instance_v1.resize", &instance),
					testAccCheckDatabaseV1InstanceSameID(&instance, &instanceID),
					testAccCheckDatabaseV1InstanceFlavorAndSize(&instance, osFlavorID, 1),
					resource.TestCheckResourceAttr(
						"openstack_db_instance_v1.resize", "flavor_id", osFlavorID),
					resource.TestCheckResourceAttr(
						"openstack_db_instance_v1.resize", "size", "1"),
				),
			},
			{
				// Flavor only.
				Config: testAccDatabaseV1InstanceResize(instanceName, osFlavorIDResize, 1),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("openstack_db_instance_v1.resize", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseV1InstanceExists(t.Context(),
						"openstack_db_instance_v1.resize", &instance),
					testAccCheckDatabaseV1InstanceSameID(&instance, &instanceID),
					testAccCheckDatabaseV1InstanceFlavorAndSize(&instance, osFlavorIDResize, 1),
					resource.TestCheckResourceAttr(
						"openstack_db_instance_v1.resize", "flavor_id", osFlavorIDResize),
				),
			},
			{
				// Volume only.
				Config: testAccDatabaseV1InstanceResize(instanceName, osFlavorIDResize, 2),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("openstack_db_instance_v1.resize", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseV1InstanceExists(t.Context(),
						"openstack_db_instance_v1.resize", &instance),
					testAccCheckDatabaseV1InstanceSameID(&instance, &instanceID),
					testAccCheckDatabaseV1InstanceFlavorAndSize(&instance, osFlavorIDResize, 2),
					resource.TestCheckResourceAttr(
						"openstack_db_instance_v1.resize", "size", "2"),
				),
			},
			{
				// Flavor and volume together.
				Config: testAccDatabaseV1InstanceResize(instanceName, osFlavorID, 3),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("openstack_db_instance_v1.resize", plancheck.ResourceActionUpdate),
					},
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseV1InstanceExists(t.Context(),
						"openstack_db_instance_v1.resize", &instance),
					testAccCheckDatabaseV1InstanceSameID(&instance, &instanceID),
					testAccCheckDatabaseV1InstanceFlavorAndSize(&instance, osFlavorID, 3),
					resource.TestCheckResourceAttr(
						"openstack_db_instance_v1.resize", "flavor_id", osFlavorID),
					resource.TestCheckResourceAttr(
						"openstack_db_instance_v1.resize", "size", "3"),
				),
			},
			{
				// A smaller volume is refused when planning.
				Config:      testAccDatabaseV1InstanceResize(instanceName, osFlavorID, 2),
				ExpectError: regexp.MustCompile(`cannot be decreased from 3 to 2 GB`),
			},
		},
	})
}

// testAccCheckDatabaseV1InstanceSameID records the ID of the instance on the
// first call and fails if a later call finds another ID, i.e. if the
// instance was replaced.
func testAccCheckDatabaseV1InstanceSameID(instance *instances.Instance, instanceID *string) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if *instanceID == "" {
			*instanceID = instance.ID

			return nil
		}

		if instance.ID != *instanceID {
			return fmt.Errorf("Instance was replaced: ID is %s, expected %s", instance.ID, *instanceID)
		}

		return nil
	}
}

func testAccCheckDatabaseV1InstanceFlavorAndSize(instance *instances.Instance, flavorID string, size int) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		if instance.Flavor.ID != flavorID {
			return fmt.Errorf("Instance has flavor %s, expected %s", instance.Flavor.ID, flavorID)
		}

		if instance.Volume.Size != size {
			return fmt.Errorf("Instance has a volume of %d GB, expected %d GB", instance.Volume.Size, size)
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

  size = 10
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
`, osDBDatastoreVersion, osDBDatastoreType, osNetworkID)
}

func testAccDatabaseV1InstanceResize(name, flavorID string, size int) string {
	return fmt.Sprintf(`
resource "openstack_db_instance_v1" "resize" {
  name      = "%[1]s"
  flavor_id = "%[2]s"
  size      = %[3]d

  # size requires volume_type: empty means the default volume type
  volume_type = ""

  datastore {
    version = "%[4]s"
    type    = "%[5]s"
  }

  network {
    uuid = "%[6]s"
  }

  database {
    name = "testdb"
  }

  user {
    name      = "testuser"
    password  = "testpassword"
    databases = ["testdb"]
    host      = "%%"
  }
}
`, name, flavorID, size, osDBDatastoreVersion, osDBDatastoreType, osNetworkID)
}
