package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/services"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccIdentityV3Service_basic(t *testing.T) {
	var name, description string

	var service services.Service

	serviceName := "ACCPTTEST-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIdentityV3ServiceDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3ServiceBasic(serviceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ServiceExists(t.Context(), "openstack_identity_service_v3.service_1", &service, &name, &description),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_service_v3.service_1", "name", &name),
					resource.TestCheckResourceAttr(
						"openstack_identity_service_v3.service_1", "type", "foo"),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_service_v3.service_1", "description", &description),
					resource.TestCheckResourceAttr(
						"openstack_identity_service_v3.service_1", "enabled", "true"),
				),
			},
			{
				Config: testAccIdentityV3ServiceUpdate(serviceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ServiceExists(t.Context(), "openstack_identity_service_v3.service_1", &service, &name, &description),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_service_v3.service_1", "name", &name),
					resource.TestCheckResourceAttr(
						"openstack_identity_service_v3.service_1", "type", "bar"),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_service_v3.service_1", "description", &description),
					resource.TestCheckResourceAttr(
						"openstack_identity_service_v3.service_1", "enabled", "false"),
				),
			},
		},
	})
}

func testAccCheckIdentityV3ServiceDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		identityClient, err := config.IdentityV3Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack identity client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_identity_service_v3" {
				continue
			}

			_, err := services.Get(ctx, identityClient, rs.Primary.ID).Extract()
			if err == nil {
				return errors.New("Service still exists")
			}
		}

		return nil
	}
}

func testAccCheckIdentityV3ServiceExists(ctx context.Context, n string, service *services.Service, name *string, description *string) resource.TestCheckFunc {
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

		found, err := services.Get(ctx, identityClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return errors.New("Service not found")
		}

		if v, ok := found.Extra["name"]; ok {
			*name = v.(string)
		}

		if v, ok := found.Extra["description"]; ok {
			*description = v.(string)
		}

		*service = *found

		return nil
	}
}

func testAccIdentityV3ServiceBasic(serviceName string) string {
	return fmt.Sprintf(`
resource "openstack_identity_service_v3" "service_1" {
  name = "%s"
  type = "foo"
  description = "A service"
}
  `, serviceName)
}

func testAccIdentityV3ServiceUpdate(serviceName string) string {
	return fmt.Sprintf(`
resource "openstack_identity_service_v3" "service_1" {
  name = "%s"
  type = "bar"
  description = "A service"
  enabled = false
}
  `, serviceName)
}
