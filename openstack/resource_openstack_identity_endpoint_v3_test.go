package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/endpoints"
	"github.com/gophercloud/gophercloud/v2/pagination"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccIdentityV3Endpoint_basic(t *testing.T) {
	var endpoint endpoints.Endpoint

	endpointName := "ACCPTTEST-" + acctest.RandString(5)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckIdentityV3EndpointDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccIdentityV3EndpointBasic(endpointName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3EndpointExists(t.Context(), "openstack_identity_endpoint_v3.endpoint_1", &endpoint),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_endpoint_v3.endpoint_1", "name", &endpoint.Name),
					resource.TestCheckResourceAttrPair(
						"openstack_identity_service_v3.service_1", "id",
						"openstack_identity_endpoint_v3.endpoint_1", "service_id"),
					resource.TestCheckResourceAttrPair(
						"openstack_identity_service_v3.service_1", "region",
						"openstack_identity_endpoint_v3.endpoint_1", "endpoint_region"),
					resource.TestCheckResourceAttr(
						"openstack_identity_endpoint_v3.endpoint_1", "url", "http://myservice.local/v1.0/%(tenant_id)s"),
				),
			},
			{
				Config: testAccIdentityV3EndpointUpdate(endpointName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3EndpointExists(t.Context(), "openstack_identity_endpoint_v3.endpoint_1", &endpoint),
					resource.TestCheckResourceAttrPtr(
						"openstack_identity_endpoint_v3.endpoint_1", "name", &endpoint.Name),
					resource.TestCheckResourceAttrPair(
						"openstack_identity_service_v3.service_1", "id",
						"openstack_identity_endpoint_v3.endpoint_1", "service_id"),
					resource.TestCheckResourceAttr(
						"openstack_identity_endpoint_v3.endpoint_1", "endpoint_region", "interstate76"),
					resource.TestCheckResourceAttr(
						"openstack_identity_endpoint_v3.endpoint_1", "url", "http://my-new-service/v1.0/%(tenant_id)s"),
				),
			},
		},
	})
}

func testAccCheckIdentityV3EndpointDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		identityClient, err := config.IdentityV3Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack identity client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_identity_endpoint_v3" {
				continue
			}

			var endpoint endpoints.Endpoint

			err = endpoints.List(identityClient, nil).EachPage(ctx, func(_ context.Context, page pagination.Page) (bool, error) {
				endpointList, err := endpoints.ExtractEndpoints(page)
				if err != nil {
					return false, err
				}

				for _, v := range endpointList {
					if v.ID == rs.Primary.ID {
						endpoint = v

						break
					}
				}

				return true, nil
			})
			if err != nil {
				return fmt.Errorf("Error retrieving OpenStack identity endpoints: %w", err)
			}

			if endpoint != (endpoints.Endpoint{}) {
				return errors.New("Endpoint still exists")
			}
		}

		return nil
	}
}

func testAccCheckIdentityV3EndpointExists(ctx context.Context, n string, endpoint *endpoints.Endpoint) resource.TestCheckFunc {
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

		var found *endpoints.Endpoint

		err = endpoints.List(identityClient, nil).EachPage(ctx, func(_ context.Context, page pagination.Page) (bool, error) {
			endpointList, err := endpoints.ExtractEndpoints(page)
			if err != nil {
				return false, err
			}

			for _, ep := range endpointList {
				e := ep
				if e.ID == rs.Primary.ID {
					found = &e

					break
				}
			}

			return true, nil
		})

		if err != nil || *found == (endpoints.Endpoint{}) {
			return errors.New("Endpoint not found")
		}

		if found.ID != rs.Primary.ID {
			return errors.New("Endpoint not found")
		}

		*endpoint = *found

		return nil
	}
}

func testAccIdentityV3EndpointBasic(endpointName string) string {
	return fmt.Sprintf(`
resource "openstack_identity_service_v3" "service_1" {
  name = "foo"
  type = "bar"
}

resource "openstack_identity_endpoint_v3" "endpoint_1" {
  name = "%s"
  service_id = openstack_identity_service_v3.service_1.id
  endpoint_region = openstack_identity_service_v3.service_1.region
  url = "http://myservice.local/v1.0/%%(tenant_id)s"
}
  `, endpointName)
}

func testAccIdentityV3EndpointUpdate(endpointName string) string {
	return fmt.Sprintf(`
resource "openstack_identity_service_v3" "service_1" {
  name = "baz"
  type = "qux"
}

resource "openstack_identity_endpoint_v3" "endpoint_1" {
  name = "%s"
  service_id = openstack_identity_service_v3.service_1.id
  endpoint_region = "interstate76"
  url = "http://my-new-service/v1.0/%%(tenant_id)s"
}
  `, endpointName)
}
