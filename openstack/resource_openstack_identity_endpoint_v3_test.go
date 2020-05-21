package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/endpoints"
	"github.com/gophercloud/gophercloud/pagination"
)

func TestAccIdentityV3Endpoint_basic(t *testing.T) {
	var endpoint endpoints.Endpoint
	var endpointName = fmt.Sprintf("ACCPTTEST-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIdentityV3EndpointDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3Endpoint_basic(endpointName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3EndpointExists("openstack_identity_endpoint_v3.endpoint_1", &endpoint),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_endpoint_v3.endpoint_1", "name", &endpoint.Name),
					resource.TestCheckResourceAttrPair(
						"openstack_identity_service_v3.service_1", "id",
						"openstack_identity_endpoint_v3.endpoint_1", "service_id"),
					resource.TestCheckResourceAttrPair(
						"openstack_identity_service_v3.service_1", "region",
						"openstack_identity_endpoint_v3.endpoint_1", "endpoint_region"),
					resource.TestCheckResourceAttr(
						"openstack_identity_endpoint_v3.endpoint_1", "url", "http://myservice.local"),
				),
			},
			{
				Config: testAccIdentityV3Endpoint_update(endpointName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3EndpointExists("openstack_identity_endpoint_v3.endpoint_1", &endpoint),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_endpoint_v3.endpoint_1", "name", &endpoint.Name),
					resource.TestCheckResourceAttrPair(
						"openstack_identity_service_v3.service_1", "id",
						"openstack_identity_endpoint_v3.endpoint_1", "service_id"),
					resource.TestCheckResourceAttr(
						"openstack_identity_endpoint_v3.endpoint_1", "endpoint_region", "interstate76"),
					resource.TestCheckResourceAttr(
						"openstack_identity_endpoint_v3.endpoint_1", "url", "http://my-new-service.local"),
				),
			},
		},
	})
}

func testAccCheckIdentityV3EndpointDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	identityClient, err := config.IdentityV3Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack identity client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_identity_endpoint_v3" {
			continue
		}

		var endpoint endpoints.Endpoint
		endpoints.List(identityClient, nil).EachPage(func(page pagination.Page) (bool, error) {
			if endpointList, err := endpoints.ExtractEndpoints(page); err != nil {
				return false, err
			} else {
				for _, v := range endpointList {
					if v.ID == rs.Primary.ID {
						endpoint = v
						break
					}
				}
			}
			return true, nil
		})

		if endpoint != (endpoints.Endpoint{}) {
			return fmt.Errorf("Endpoint still exists")
		}
	}

	return nil
}

func testAccCheckIdentityV3EndpointExists(n string, endpoint *endpoints.Endpoint) resource.TestCheckFunc {
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

		var found *endpoints.Endpoint
		err = endpoints.List(identityClient, nil).EachPage(func(page pagination.Page) (bool, error) {
			if endpointList, err := endpoints.ExtractEndpoints(page); err != nil {
				return false, err
			} else {
				for _, v := range endpointList {
					if v.ID == rs.Primary.ID {
						found = &v
						break
					}
				}
			}
			return true, nil
		})

		if err != nil || *found == (endpoints.Endpoint{}) {
			return fmt.Errorf("Endpoint not found")
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Endpoint not found")
		}

		*endpoint = *found

		return nil
	}
}

func testAccIdentityV3Endpoint_basic(endpointName string) string {
	return fmt.Sprintf(`
resource "openstack_identity_service_v3" "service_1" {
  name = "foo"
  type = "bar"
}

resource "openstack_identity_endpoint_v3" "endpoint_1" {
  name = "%s"
  service_id = "${openstack_identity_service_v3.service_1.id}"
  endpoint_region = "${openstack_identity_service_v3.service_1.region}"
  url = "http://myservice.local"
}
  `, endpointName)
}

func testAccIdentityV3Endpoint_update(endpointName string) string {
	return fmt.Sprintf(`
resource "openstack_identity_service_v3" "service_1" {
  name = "baz"
  type = "qux"
}

resource "openstack_identity_endpoint_v3" "endpoint_1" {
  name = "%s"
  service_id = "${openstack_identity_service_v3.service_1.id}"
  endpoint_region = "interstate76"
  url = "http://my-new-service.local"
}
  `, endpointName)
}
