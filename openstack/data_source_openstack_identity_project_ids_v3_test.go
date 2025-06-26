package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOpenStackIdentityProjectIDsV3DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackIdentityProjectIDsV3DataSourceEmpty(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.openstack_identity_project_ids_v3.projects_empty", "ids.#", "0"),
				),
			},
			{
				Config: testAccOpenStackIdentityProjectIDsV3DataSourceName(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.openstack_identity_project_ids_v3.projects_by_name", "ids.#", "1"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_identity_project_ids_v3.projects_by_name", "ids.0",
						"openstack_identity_project_v3.project_1", "id"),
				),
			},
			{
				Config: testAccOpenStackIdentityProjectIDsV3DataSourceRegex(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.openstack_identity_project_ids_v3.projects_by_name_regex", "ids.#", "1"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_identity_project_ids_v3.projects_by_name_regex", "ids.0",
						"openstack_identity_project_v3.project_2", "id"),
				),
			},
			{
				Config: testAccOpenStackIdentityProjectIDsV3DataSourceTags(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.openstack_identity_project_ids_v3.projects_by_tag", "ids.#", "2"),
				),
			},
		},
	})
}

const testAccOpenStackIdentityProjectIDsV3DataSourceProjects = `
	resource "openstack_identity_project_v3" "project_1" {
	  name = "project_1"
	  tags = ["dev", "product1"]
	}

	resource "openstack_identity_project_v3" "project_2" {
	  name = "project_2"
	  tags = ["prod", "product1"]
	}
`

func testAccOpenStackIdentityProjectIDsV3DataSourceEmpty() string {
	return fmt.Sprintf(`
%s

data "openstack_identity_project_ids_v3" "projects_empty" {
    name = "non-existed-project"
}
`, testAccOpenStackIdentityProjectIDsV3DataSourceProjects)
}

func testAccOpenStackIdentityProjectIDsV3DataSourceName() string {
	return fmt.Sprintf(`
%s

data "openstack_identity_project_ids_v3" "projects_by_name" {
	name = openstack_identity_project_v3.project_1.name
}
`, testAccOpenStackIdentityProjectIDsV3DataSourceProjects)
}

func testAccOpenStackIdentityProjectIDsV3DataSourceRegex() string {
	return fmt.Sprintf(`
%s

data "openstack_identity_project_ids_v3" "projects_by_name_regex" {
	name_regex = "^.+_2$"
}
`, testAccOpenStackIdentityProjectIDsV3DataSourceProjects)
}

func testAccOpenStackIdentityProjectIDsV3DataSourceTags() string {
	return fmt.Sprintf(`
%s

data "openstack_identity_project_ids_v3" "projects_by_tag" {
	tags = ["product1"]
}
`, testAccOpenStackIdentityProjectIDsV3DataSourceProjects)
}
