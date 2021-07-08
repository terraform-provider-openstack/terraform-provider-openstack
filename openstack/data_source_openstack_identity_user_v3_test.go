package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOpenStackIdentityV3UserDataSource_basic(t *testing.T) {
	userName := fmt.Sprintf("tf_test_%s", acctest.RandString(5))
	userPassword := acctest.RandString(20)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackIdentityUserV3DataSourceUser(userName, userPassword),
			},
			{
				Config: testAccOpenStackIdentityUserV3DataSourceBasic(userName, userPassword),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityUserV3DataSourceID("data.openstack_identity_user_v3.user_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_identity_user_v3.user_1", "name", userName),
					resource.TestCheckResourceAttr(
						"data.openstack_identity_user_v3.user_1", "domain_id", "default"),
					resource.TestCheckResourceAttr(
						"data.openstack_identity_user_v3.user_1", "enabled", "true"),
					testAccCheckIdentityUserV3DataSourceDefaultProjectID(
						"data.openstack_identity_user_v3.user_1", "openstack_identity_project_v3.project_1"),
				),
			},
		},
	})
}

func testAccCheckIdentityUserV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find user data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("User data source ID not set")
		}

		return nil
	}
}

func testAccCheckIdentityUserV3DataSourceDefaultProjectID(n1, n2 string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ds1, ok := s.RootModule().Resources[n1]
		if !ok {
			return fmt.Errorf("Can't find user data source: %s", n1)
		}

		if ds1.Primary.ID == "" {
			return fmt.Errorf("User data source ID not set")
		}

		rs2, ok := s.RootModule().Resources[n2]
		if !ok {
			return fmt.Errorf("Can't find project resource: %s", n2)
		}

		if rs2.Primary.ID == "" {
			return fmt.Errorf("Project resource ID not set")
		}

		if rs2.Primary.ID != ds1.Primary.Attributes["default_project_id"] {
			return fmt.Errorf("Project id and user default_project_id don't match")
		}

		return nil
	}
}

func testAccOpenStackIdentityUserV3DataSourceUser(name, password string) string {
	return fmt.Sprintf(`
	%s

	resource "openstack_identity_user_v3" "user_1" {
	  name = "%s"
	  password = "%s"
	  default_project_id = "${openstack_identity_project_v3.project_1.id}"
	}
`, testAccOpenStackIdentityProjectV3DataSourceProject(fmt.Sprintf("%s_project", name), acctest.RandString(20), "tag1", "tag2"), name, password)
}

func testAccOpenStackIdentityUserV3DataSourceBasic(name, password string) string {
	return fmt.Sprintf(`
	%s

	data "openstack_identity_user_v3" "user_1" {
      name = "${openstack_identity_user_v3.user_1.name}"
	}
`, testAccOpenStackIdentityUserV3DataSourceUser(name, password))
}
