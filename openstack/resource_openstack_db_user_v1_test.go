package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/instances"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/users"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
		CheckDestroy:      testAccCheckDatabaseV1UserDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseV1UserBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseV1InstanceExists(t.Context(),
						"openstack_db_instance_v1.basic", &instance),
					testAccCheckDatabaseV1UserExists(t.Context(),
						"openstack_db_user_v1.basic", &instance, &user),
					resource.TestCheckResourceAttrPtr(
						"openstack_db_user_v1.basic", "name", &user.Name),
				),
			},
		},
	})
}

func testAccCheckDatabaseV1UserExists(ctx context.Context, n string, instance *instances.Instance, user *users.User) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No ID is set")
		}

		_, userName, err := parsePairedIDs(rs.Primary.ID, "openstack_db_user_v1")
		if err != nil {
			return err
		}

		config := testAccProvider.Meta().(*Config)

		databaseV1Client, err := config.DatabaseV1Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating cloud database client: %w", err)
		}

		pages, err := users.List(databaseV1Client, instance.ID).AllPages(ctx)
		if err != nil {
			return fmt.Errorf("Unable to retrieve users: %w", err)
		}

		allUsers, err := users.ExtractUsers(pages)
		if err != nil {
			return fmt.Errorf("Unable to extract users: %w", err)
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

func testAccCheckDatabaseV1UserDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		databaseV1Client, err := config.DatabaseV1Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating cloud database client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_db_user_v1" {
				continue
			}

			id, userName, err := parsePairedIDs(rs.Primary.ID, "openstack_db_user_v1")
			if err != nil {
				return err
			}

			pages, err := users.List(databaseV1Client, id).AllPages(ctx)
			if err != nil {
				return nil
			}

			allUsers, err := users.ExtractUsers(pages)
			if err != nil {
				return fmt.Errorf("Unable to extract users: %w", err)
			}

			var exists bool

			for _, v := range allUsers {
				if v.Name == userName {
					exists = true
				}
			}

			if exists {
				return errors.New("User still exists")
			}
		}

		return nil
	}
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
  instance_id = openstack_db_instance_v1.basic.id
  password    = "password"
  databases   = ["testdb"]
}
`, osDBDatastoreVersion, osDBDatastoreType, osNetworkID)
}
