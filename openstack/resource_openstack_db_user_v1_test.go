package openstack

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/instances"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/users"
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

		_, userName, err := parsePairedIDs(rs.Primary.ID, "openstack_db_user_v1")
		if err != nil {
			return err
		}

		config := testAccProvider.Meta().(*Config)
		DatabaseV1Client, err := config.DatabaseV1Client(context.TODO(), osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating cloud database client: %s", err)
		}

		pages, err := users.List(DatabaseV1Client, instance.ID).AllPages(context.TODO())
		if err != nil {
			return fmt.Errorf("Unable to retrieve users: %s", err)
		}

		allUsers, err := users.ExtractUsers(pages)
		if err != nil {
			return fmt.Errorf("Unable to extract users: %s", err)
		}

		for _, u := range allUsers {
			if u.Name == userName {
				*user = u
				return nil
			}
		}

		return fmt.Errorf("User %s does not exist", n)
	}
}

func testAccCheckDatabaseV1UserDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)

	DatabaseV1Client, err := config.DatabaseV1Client(context.TODO(), osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating cloud database client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_db_user_v1" {
			continue
		}

		id, userName, err := parsePairedIDs(rs.Primary.ID, "openstack_db_user_v1")
		if err != nil {
			return err
		}

		pages, err := users.List(DatabaseV1Client, id).AllPages(context.TODO())
		if err != nil {
			return nil
		}

		allUsers, err := users.ExtractUsers(pages)
		if err != nil {
			return fmt.Errorf("Unable to extract users: %s", err)
		}

		var exists bool
		for _, v := range allUsers {
			if v.Name == userName {
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
