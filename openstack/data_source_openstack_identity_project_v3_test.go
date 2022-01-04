package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOpenStackIdentityV3ProjectDataSource_basic(t *testing.T) {
	projectName := fmt.Sprintf("tf_test_%s", acctest.RandString(5))
	projectDescription := acctest.RandString(20)
	projectTag1 := acctest.RandString(20)
	projectTag2 := acctest.RandString(20)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackIdentityProjectV3DataSourceProject(projectName, projectDescription, projectTag1, projectTag2),
			},
			{
				Config: testAccOpenStackIdentityProjectV3DataSourceBasic(projectName, projectDescription, projectTag1, projectTag2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectDataSourceID("data.openstack_identity_project_v3.project_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_identity_project_v3.project_1", "name", projectName),
					resource.TestCheckResourceAttr(
						"openstack_identity_project_v3.project_1", "description", projectDescription),
					resource.TestCheckResourceAttr(
						"openstack_identity_project_v3.project_1", "enabled", "true"),
					resource.TestCheckResourceAttr(
						"openstack_identity_project_v3.project_1", "is_domain", "false"),
					resource.TestCheckResourceAttr(
						"openstack_identity_project_v3.project_1", "tags.#", "2"),
				),
			},
		},
	})
}

func testAccCheckIdentityV3ProjectDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find project data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Project data source ID not set")
		}

		return nil
	}
}

func testAccOpenStackIdentityProjectV3DataSourceProject(name, description, tag1, tag2 string) string {
	return fmt.Sprintf(`
	resource "openstack_identity_project_v3" "project_1" {
	  name = "%s"
	  description = "%s"
	  tags = ["%s", "%s"]
	}
`, name, description, tag1, tag2)
}

func testAccOpenStackIdentityProjectV3DataSourceBasic(name, description, tag1, tag2 string) string {
	return fmt.Sprintf(`
	%s

	data "openstack_identity_project_v3" "project_1" {
      name = "${openstack_identity_project_v3.project_1.name}"
	}
`, testAccOpenStackIdentityProjectV3DataSourceProject(name, description, tag1, tag2))
}
