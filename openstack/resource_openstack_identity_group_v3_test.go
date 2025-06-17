package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/groups"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccIdentityV3Group_basic(t *testing.T) {
	var group groups.Group

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIdentityV3GroupDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3GroupBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3GroupExists(t.Context(), "openstack_identity_group_v3.group_1", &group),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_group_v3.group_1", "name", &group.Name),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_group_v3.group_1", "description", &group.Description),
				),
			},
			{
				Config: testAccIdentityV3GroupUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3GroupExists(t.Context(), "openstack_identity_group_v3.group_1", &group),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_group_v3.group_1", "name", &group.Name),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_group_v3.group_1", "description", &group.Description),
				),
			},
		},
	})
}

func testAccCheckIdentityV3GroupDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		identityClient, err := config.IdentityV3Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack identity client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_identity_group_v3" {
				continue
			}

			_, err := groups.Get(ctx, identityClient, rs.Primary.ID).Extract()
			if err == nil {
				return errors.New("Group still exists")
			}
		}

		return nil
	}
}

func testAccCheckIdentityV3GroupExists(ctx context.Context, n string, group *groups.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		identityClient, err := config.IdentityV3Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack identity client: %w", err)
		}

		found, err := groups.Get(ctx, identityClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return errors.New("Group not found")
		}

		*group = *found

		return nil
	}
}

const testAccIdentityV3GroupBasic = `
resource "openstack_identity_group_v3" "group_1" {
	name = "group_1"
	description = "Terraform accept test 1"
}
`

const testAccIdentityV3GroupUpdate = `
resource "openstack_identity_group_v3" "group_1" {
	name = "group_2"
	description = "Terraform accept test 2"
}
`
