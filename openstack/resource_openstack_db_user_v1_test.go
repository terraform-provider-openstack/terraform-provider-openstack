package openstack

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/db/v1/users"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDatabaseV1User_basic(t *testing.T) {
	var user users.User

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheckDatabase(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccDatabaseV1UserBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatabaseV1UserExists(
						"openstack_db_user_v1.basic", &user),
					resource.TestCheckResourceAttr(
						"openstack_db_user_v1.basic", "name", "basic"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckDatabaseV1UserExists(n string, user *users.User) resource.TestCheckFunc {
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

		dbInstance := s.RootModule().Resources["openstack_db_instance_v1.basic"]

		pages, err := users.List(databaseV1Client, dbInstance.Primary.ID).AllPages()
		if err != nil {
			return fmt.Errorf("Unable to retrieve users, pages: %s", err)
		}
		allUsers, err := users.ExtractUsers(pages)
		if err != nil {
			return fmt.Errorf("Unable to retrieve users, extract: %s", err)
		}

		for _, v := range allUsers {
			if v.Name == user.Name {
				return nil
			}
		}

		return nil
	}
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
  name      = "basic"
  instance  = "${openstack_db_instance_v1.basic.id}"
  password  = "password"
  databases = ["testdb"]
}
`, OS_DB_DATASTORE_VERSION, OS_DB_DATASTORE_TYPE, OS_NETWORK_ID)
