package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOpenStackIdentityV3EndpointDataSource_basic(t *testing.T) {
	endpointName := "identity"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackIdentityEndpointV3DataSource_basic(endpointName, "public"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityEndpointV3DataSourceID("data.openstack_identity_endpoint_v3.endpoint_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_identity_endpoint_v3.endpoint_1", "service_name", endpointName),
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

func testAccOpenStackIdentityEndpointV3DataSource_basic(name, iface string) string {
	return fmt.Sprintf(`
	data "openstack_identity_endpoint_v3" "endpoint_1" {
      service_name = "%s"
      interface = "%s"
	}
`, name, iface)
}
