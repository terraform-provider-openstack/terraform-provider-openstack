package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
)

func TestAccIdentityV3Project_basic(t *testing.T) {
	var project projects.Project
	var projectName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIdentityV3ProjectDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3ProjectBasic(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectExists("openstack_identity_project_v3.project_1", &project),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_project_v3.project_1", "name", &project.Name),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_project_v3.project_1", "description", &project.Description),
					resource.TestCheckResourceAttr(
						"openstack_identity_project_v3.project_1", "domain_id", "default"),
					resource.TestCheckResourceAttr(
						"openstack_identity_project_v3.project_1", "enabled", "true"),
					resource.TestCheckResourceAttr(
						"openstack_identity_project_v3.project_1", "is_domain", "false"),
				),
			},
			{
				Config: testAccIdentityV3ProjectUpdate(projectName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ProjectExists("openstack_identity_project_v3.project_1", &project),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_project_v3.project_1", "name", &project.Name),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_project_v3.project_1", "description", &project.Description),
					resource.TestCheckResourceAttr(
						"openstack_identity_project_v3.project_1", "domain_id", "default"),
					resource.TestCheckResourceAttr(
						"openstack_identity_project_v3.project_1", "enabled", "false"),
					resource.TestCheckResourceAttr(
						"openstack_identity_project_v3.project_1", "is_domain", "false"),
					testAccCheckIdentityV3ProjectHasTag("openstack_identity_project_v3.project_1", "tag1"),
					testAccCheckIdentityV3ProjectHasTag("openstack_identity_project_v3.project_1", "tag2"),
					testAccCheckIdentityV3ProjectTagCount("openstack_identity_project_v3.project_1", 2),
				),
			},
		},
	})
}

func testAccCheckIdentityV3ProjectDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	identityClient, err := config.IdentityV3Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_identity_project_v3" {
			continue
		}

		_, err := projects.Get(identityClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Project still exists")
		}
	}

	return nil
}

func testAccCheckIdentityV3ProjectExists(n string, project *projects.Project) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		identityClient, err := config.IdentityV3Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack identity client: %s", err)
		}

		found, err := projects.Get(identityClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Project not found")
		}

		*project = *found

		return nil
	}
}

func testAccCheckIdentityV3ProjectHasTag(n, tag string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		identityClient, err := config.IdentityV3Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack identity client: %s", err)
		}

		found, err := projects.Get(identityClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Project not found")
		}

		for _, v := range found.Tags {
			if tag == v {
				return nil
			}
		}

		return fmt.Errorf("Tag not found: %s", tag)
	}
}

func testAccCheckIdentityV3ProjectTagCount(n string, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		identityClient, err := config.IdentityV3Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack identity client: %s", err)
		}

		found, err := projects.Get(identityClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Project not found")
		}

		if len(found.Tags) != expected {
			return fmt.Errorf("Expecting %d tags, found %d", expected, len(found.Tags))
		}

		return nil
	}
}

func testAccIdentityV3ProjectBasic(projectName string) string {
	return fmt.Sprintf(`
    resource "openstack_identity_project_v3" "project_1" {
      name = "%s"
      description = "A project"
    }
  `, projectName)
}

func testAccIdentityV3ProjectUpdate(projectName string) string {
	return fmt.Sprintf(`
    resource "openstack_identity_project_v3" "project_1" {
      name = "%s"
      description = "Some project"
	  enabled = false
	  tags = ["tag1","tag2"]
    }
  `, projectName)
}
