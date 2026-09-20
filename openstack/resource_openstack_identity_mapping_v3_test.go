package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/federation"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccIdentityV3Mapping_basic(t *testing.T) {
	var mapping federation.Mapping

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIdentityV3MappingDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3MappingBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3MappingExists(t.Context(), "openstack_identity_mapping_v3.mapping_1", &mapping),
					resource.TestCheckResourceAttr("openstack_identity_mapping_v3.mapping_1", "mapping_id", "terraform-test"),
					resource.TestCheckResourceAttr("openstack_identity_mapping_v3.mapping_1", "rules", `[{"local":[{"user":{"name":"{0}"}}],"remote":[{"type":"UserName"}]}]`),
					resource.TestCheckResourceAttrPair(
						"data.openstack_identity_mapping_v3.mapping_1", "rules",
						"openstack_identity_mapping_v3.mapping_1", "rules",
					),
				),
			},
			{
				Config: testAccIdentityV3MappingUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3MappingExists(t.Context(), "openstack_identity_mapping_v3.mapping_1", &mapping),
					resource.TestCheckResourceAttr("openstack_identity_mapping_v3.mapping_1", "mapping_id", "terraform-test"),
					resource.TestCheckResourceAttr("openstack_identity_mapping_v3.mapping_1", "rules", `[{"local":[{"user":{"name":"{0}"}}],"remote":[{"type":"UserName"},{"type":"orgPersonType","any_one_of":["Employee"]}]}]`),
				),
			},
		},
	})
}

func testAccCheckIdentityV3MappingDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		identityClient, err := config.IdentityV3Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack identity client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_identity_mapping_v3" {
				continue
			}

			_, err := federation.GetMapping(ctx, identityClient, rs.Primary.ID).Extract()
			if err == nil {
				return errors.New("Identity mapping still exists")
			}
		}

		return nil
	}
}

func testAccCheckIdentityV3MappingExists(ctx context.Context, n string, mapping *federation.Mapping) resource.TestCheckFunc {
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

		found, err := federation.GetMapping(ctx, identityClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return errors.New("Identity mapping not found")
		}

		*mapping = *found

		return nil
	}
}

const testAccIdentityV3MappingBasic = `
resource "openstack_identity_mapping_v3" "mapping_1" {
  mapping_id = "terraform-test"
  rules = jsonencode([
    {
      local = [
        {
          user = {
            name = "{0}"
          }
        }
      ]
      remote = [
        {
          type = "UserName"
        }
      ]
    }
  ])
}

data "openstack_identity_mapping_v3" "mapping_1" {
  mapping_id = openstack_identity_mapping_v3.mapping_1.id
}
`

const testAccIdentityV3MappingUpdate = `
resource "openstack_identity_mapping_v3" "mapping_1" {
  mapping_id = "terraform-test"
  rules = jsonencode([
    {
      local = [
        {
          user = {
            name = "{0}"
          }
        }
      ]
      remote = [
        {
          type = "UserName"
        },
        {
          type       = "orgPersonType"
          any_one_of = ["Employee"]
        }
      ]
    }
  ])
}
`
