package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOpenStackIdentityV3EndpointDataSource_basic(t *testing.T) {
	serviceName := "keystone"
	serviceType := "identity"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackIdentityEndpointV3DataSourceBasic(serviceName, "public"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityEndpointV3DataSourceID("data.openstack_identity_endpoint_v3.endpoint_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_identity_endpoint_v3.endpoint_1", "service_name", serviceName),
					resource.TestCheckResourceAttr(
						"data.openstack_identity_endpoint_v3.endpoint_1", "service_type", serviceType),
				),
			},
		},
	})
}

func testAccCheckIdentityEndpointV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find endpoint data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Endpoint data source ID not set")
		}

		return nil
	}
}

func testAccOpenStackIdentityEndpointV3DataSourceBasic(name, iface string) string {
	return fmt.Sprintf(`
	data "openstack_identity_endpoint_v3" "endpoint_1" {
      service_name = "%s"
      interface = "%s"
	}
`, name, iface)
}
