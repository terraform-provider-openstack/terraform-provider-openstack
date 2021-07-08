package openstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/db/v1/instances"
	"github.com/gophercloud/gophercloud/openstack/db/v1/users"
)

func TestAccDatabaseV1User_basic(t *testing.T) {
	var user users.User
	var instance instances.Instance

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckDatabase(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDatabaseV1UserDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseV1UserBasic(),
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
		DatabaseV1Client, err := config.DatabaseV1Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating cloud database client: %s", err)
		}

		pages, err := users.List(DatabaseV1Client, instance.ID).AllPages()
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

	DatabaseV1Client, err := config.DatabaseV1Client(osRegionName)
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

		pages, err := users.List(DatabaseV1Client, parts[0]).AllPages()
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

func testAccDatabaseV1UserBasic() string {
	return fmt.Sprintf(`
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
`, osDBDatastoreVersion, osDBDatastoreType, osNetworkID)
}
