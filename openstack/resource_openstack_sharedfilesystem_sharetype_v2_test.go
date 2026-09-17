package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/sharetypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccSFSV2ShareType_basic(t *testing.T) {
	var sharetype sharetypes.ShareType

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
			testAccPreCheckSFS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSFSV2ShareTypeDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccSFSV2ShareTypeConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSV2ShareTypeExists(t.Context(), "openstack_sharedfilesystem_sharetype_v2.sharetype_1", &sharetype),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharetype_v2.sharetype_1", "name", "sharetype_1"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharetype_v2.sharetype_1", "description", "created by terraform"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharetype_v2.sharetype_1", "is_public", "true"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharetype_v2.sharetype_1", "extra_specs.%", "1"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharetype_v2.sharetype_1", "extra_specs.driver_handles_share_servers", "false"),
				),
			},
			{
				Config: testAccSFSV2ShareTypeConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSV2ShareTypeExists(t.Context(), "openstack_sharedfilesystem_sharetype_v2.sharetype_1", &sharetype),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharetype_v2.sharetype_1", "name", "sharetype_1_renamed"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharetype_v2.sharetype_1", "description", "updated by terraform"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharetype_v2.sharetype_1", "is_public", "false"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharetype_v2.sharetype_1", "extra_specs.%", "3"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharetype_v2.sharetype_1", "extra_specs.driver_handles_share_servers", "true"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharetype_v2.sharetype_1", "extra_specs.snapshot_support", "true"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharetype_v2.sharetype_1", "extra_specs.my_custom_spec", "some_value"),
				),
			},
		},
	})
}

func TestAccSFSV2ShareType_EndpointCheck(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
			testAccPreCheckSFS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSFSV2ShareTypeDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				// register the sharev2 service and endpoint
				Config: testAccSFSV2ShareTypeConfigUpdateEndpointCheck,
			},
			{
				// test endpoint locator to pick up sharev2
				Config: testAccSFSV2ShareTypeConfigUpdateEndpointCheck + testAccSFSV2ShareTypeConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharetype_v2.sharetype_1", "name", "sharetype_1"),
				),
			},
		},
	})
}

func testAccCheckSFSV2ShareTypeDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		sfsClient, err := config.SharedfilesystemV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_sharedfilesystem_sharetype_v2" {
				continue
			}

			_, err := sharetypes.Get(ctx, sfsClient, rs.Primary.ID).Extract()
			if err == nil {
				return errors.New("Manila share type still exists")
			}
		}

		return nil
	}
}

func testAccCheckSFSV2ShareTypeExists(ctx context.Context, n string, sharetype *sharetypes.ShareType) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		sfsClient, err := config.SharedfilesystemV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %w", err)
		}

		found, err := sharetypes.Get(ctx, sfsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return errors.New("Share type not found")
		}

		*sharetype = *found

		return nil
	}
}

const testAccSFSV2ShareTypeConfigBasic = `
resource "openstack_sharedfilesystem_sharetype_v2" "sharetype_1" {
  name        = "sharetype_1"
  description = "created by terraform"
  is_public   = true

  extra_specs = {
    driver_handles_share_servers = "false"
  }
}
`

const testAccSFSV2ShareTypeConfigUpdate = `
resource "openstack_sharedfilesystem_sharetype_v2" "sharetype_1" {
  name        = "sharetype_1_renamed"
  description = "updated by terraform"
  is_public   = false

  extra_specs = {
    driver_handles_share_servers = "true"
    snapshot_support             = "true"
    my_custom_spec               = "some_value"
  }
}
`

const testAccSFSV2ShareTypeConfigUpdateEndpointCheck = `
resource "openstack_identity_service_v3" "service_1" {
  name = "manilav2"
  type = "sharev2"
}

resource "openstack_identity_endpoint_v3" "endpoint_1" {
  name            = "sharev2"
  service_id      = openstack_identity_service_v3.service_1.id
  endpoint_region = openstack_identity_service_v3.service_1.region
  url             = "http://my-endpoint"
}
`
