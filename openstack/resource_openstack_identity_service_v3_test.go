package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/services"
)

func TestAccIdentityV3Service_basic(t *testing.T) {
	var name, description string
	var service services.Service
	var serviceName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityV3ServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3Service_basic(serviceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ServiceExists("openstack_identity_service_v3.service_1", &service, &name, &description),
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
				Config: testAccIdentityV3Service_update(serviceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3ServiceExists("openstack_identity_service_v3.service_1", &service, &name, &description),
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

func testAccCheckIdentityV3ServiceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	identityClient, err := config.IdentityV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_identity_service_v3" {
			continue
		}

		_, err := services.Get(identityClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Service still exists")
		}
	}

	return nil
}

func testAccCheckIdentityV3ServiceExists(n string, service *services.Service, name *string, description *string) resource.TestCheckFunc {
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

		found, err := services.Get(identityClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Service not found")
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

func testAccIdentityV3Service_basic(serviceName string) string {
	return fmt.Sprintf(`
resource "openstack_identity_service_v3" "service_1" {
  name = "%s"
  type = "foo"
  description = "A service"
}
  `, serviceName)
}

func testAccIdentityV3Service_update(serviceName string) string {
	return fmt.Sprintf(`
resource "openstack_identity_service_v3" "service_1" {
  name = "%s"
  type = "bar"
  description = "A service"
  enabled = false
}
  `, serviceName)
}
