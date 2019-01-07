package openstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/db/v1/instances"
	"github.com/gophercloud/gophercloud/openstack/db/v1/users"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDatabaseV1User_basic(t *testing.T) {
	var user users.User
	var instance instances.Instance

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckDatabase(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDatabaseV1UserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseV1UserBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseV1InstanceExists(
						"openstack_db_instance_v1.basic", &instance),
					testAccCheckDatabaseV1UserExists(
						"openstack_db_user_v1.basic", &instance, &user),
					resource.TestCheckResourceAttrPtr(
						"openstack_db_user_v1.basic", "name", &user.Name),
				),
			},
		},
	})
}

func testAccCheckDatabaseV1UserExists(n string, instance *instances.Instance, user *users.User) resource.TestCheckFunc {

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
			return fmt.Errorf("Malformed user name: %s", rs.Primary.ID)
		}

		config := testAccProvider.Meta().(*Config)
		databaseV1Client, err := config.databaseV1Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating cloud database client: %s", err)
		}

		pages, err := users.List(databaseV1Client, instance.ID).AllPages()
		if err != nil {
			return fmt.Errorf("Unable to retrieve users: %s", err)
		}

		allUsers, err := users.ExtractUsers(pages)
		if err != nil {
			return fmt.Errorf("Unable to extract users: %s", err)
		}

		for _, u := range allUsers {
			if u.Name == parts[1] {
				*user = u
				return nil
			}
		}

		return fmt.Errorf("User %s does not exist", n)
	}
}

func testAccCheckDatabaseV1UserDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	databaseV1Client, err := config.databaseV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating cloud database client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_db_user_v1" {
			continue
		}

		parts := strings.SplitN(rs.Primary.ID, "/", 2)
		if len(parts) != 2 {
			return fmt.Errorf("Malformed username: %s", rs.Primary.ID)
		}

		pages, err := users.List(databaseV1Client, parts[0]).AllPages()
		if err != nil {
			return nil
		}

		allUsers, err := users.ExtractUsers(pages)
		if err != nil {
			return fmt.Errorf("Unable to extract users: %s", err)
		}

		var exists bool
		for _, v := range allUsers {
			if v.Name == parts[1] {
				exists = true
			}
		}

		if exists {
			return fmt.Errorf("User still exists")
		}
	}

	return nil
}

var testAccDatabaseV1UserBasic = fmt.Sprintf(`
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

resource "openstack_db_user_v1" "basic" {
  name        = "basic"
  instance_id = "${openstack_db_instance_v1.basic.id}"
  password    = "password"
  databases   = ["testdb"]
}
`, OS_DB_DATASTORE_VERSION, OS_DB_DATASTORE_TYPE, OS_NETWORK_ID)
