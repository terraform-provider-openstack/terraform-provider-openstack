package openstack

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/groups"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccIdentityV3Group_basic(t *testing.T) {
	var group groups.Group
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityV3GroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3Group_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3GroupExists("openstack_identity_group_v3.group_1", &group),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_group_v3.group_1", "name", &group.Name),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_group_v3.group_1", "description", &group.Description),
				),
			},
			{
				Config: testAccIdentityV3Group_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3GroupExists("openstack_identity_group_v3.group_1", &group),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_group_v3.group_1", "name", &group.Name),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_group_v3.group_1", "description", &group.Description),
				),
			},
		},
	})
}

func testAccCheckIdentityV3GroupDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	identityClient, err := config.IdentityV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_identity_group_v3" {
			continue
		}

		_, err := groups.Get(identityClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Group still exists")
		}
	}

	return nil
}

func testAccCheckIdentityV3GroupExists(n string, group *groups.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		identityClient, err := config.IdentityV3Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack identity client: %s", err)
		}

		found, err := groups.Get(identityClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Group not found")
		}

		*group = *found

		return nil
	}
}

const testAccIdentityV3Group_basic = `
resource "openstack_identity_group_v3" "group_1" {
	name = "group_1"
	description = "Terraform accept test 1"
}
`

const testAccIdentityV3Group_update = `
resource "openstack_identity_group_v3" "group_1" {
	name = "group_2"
	description = "Terraform accept test 2"
}
`
