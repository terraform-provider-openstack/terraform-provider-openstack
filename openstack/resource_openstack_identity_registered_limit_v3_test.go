package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/registeredlimits"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/services"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccIdentityV3RegisteredLimit_basic(t *testing.T) {
	var service services.Service

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIdentityV3RegisteredLimitDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3RegisteredLimitBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3RegisteredLimitExists(t.Context(), "openstack_identity_registered_limit_v3.limit_1", &service),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_registered_limit_v3.limit_1", "service_id", &service.ID),
					resource.TestCheckResourceAttr(
						"openstack_identity_registered_limit_v3.limit_1", "resource_name", "instances"),
					resource.TestCheckResourceAttr(
						"openstack_identity_registered_limit_v3.limit_1", "default_limit", "10"),
					resource.TestCheckResourceAttr(
						"openstack_identity_registered_limit_v3.limit_1", "description", "foo"),
				),
			},
			{
				Config: testAccIdentityV3RegisteredLimitUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3RegisteredLimitExists(t.Context(), "openstack_identity_registered_limit_v3.limit_1", &service),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_registered_limit_v3.limit_1", "service_id", &service.ID),
					resource.TestCheckResourceAttr(
						"openstack_identity_registered_limit_v3.limit_1", "resource_name", "instances"),
					resource.TestCheckResourceAttr(
						"openstack_identity_registered_limit_v3.limit_1", "default_limit", "100"),
					resource.TestCheckResourceAttr(
						"openstack_identity_registered_limit_v3.limit_1", "description", "bar"),
				),
			},
		},
	})
}

func testAccCheckIdentityV3RegisteredLimitDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		identityClient, err := config.IdentityV3Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack identity client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_identity_registered_limit_v3" {
				continue
			}

			_, err := registeredlimits.Get(ctx, identityClient, rs.Primary.ID).Extract()
			if err == nil {
				return errors.New("Registered limit still exists")
			}
		}

		return nil
	}
}

func testAccCheckIdentityV3RegisteredLimitExists(ctx context.Context, n string, service *services.Service) resource.TestCheckFunc {
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

		found, err := registeredlimits.Get(ctx, identityClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return errors.New("Registered limit not found")
		}

		svc, err := services.Get(ctx, identityClient, found.ServiceID).Extract()
		if err != nil {
			return fmt.Errorf("Error retrieving OpenStack service %s: %w", found.ServiceID, err)
		}

		*service = *svc

		return nil
	}
}

const testAccIdentityV3RegisteredLimitBasic = `
data "openstack_identity_service_v3" "nova" {
  name = "nova"
}

resource "openstack_identity_registered_limit_v3" "limit_1" {
  service_id = data.openstack_identity_service_v3.nova.id
  resource_name = "instances"
  default_limit = 10
  description = "foo"
}
`

const testAccIdentityV3RegisteredLimitUpdate = `
data "openstack_identity_service_v3" "nova" {
  name = "nova"
}

resource "openstack_identity_registered_limit_v3" "limit_1" {
  service_id = data.openstack_identity_service_v3.nova.id
  resource_name = "instances"
  default_limit = 100
  description = "bar"
}
`
