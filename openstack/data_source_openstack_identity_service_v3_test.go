package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccOpenStackIdentityV3ServiceDataSource_basic(t *testing.T) {
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
				Config: testAccOpenStackIdentityServiceV3DataSourceBasic(serviceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityServiceV3DataSourceID("data.openstack_identity_service_v3.service_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_identity_service_v3.service_1", "name", serviceName),
					resource.TestCheckResourceAttr(
						"data.openstack_identity_service_v3.service_1", "type", serviceType),
				),
			},
		},
	})
}

func testAccCheckIdentityServiceV3DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find service data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Service data source ID not set")
		}

		return nil
	}
}

func testAccOpenStackIdentityServiceV3DataSourceBasic(name string) string {
	return fmt.Sprintf(`
data "openstack_identity_service_v3" "service_1" {
  name = "%s"
}
`, name)
}
