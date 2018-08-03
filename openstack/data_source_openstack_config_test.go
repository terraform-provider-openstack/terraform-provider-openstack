package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccOpenStackOpenStackConfigDataSource_basic(t *testing.T) {
	authURL := "http://10.0.0.3:5000"

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccOpenStackOpenStackConfigDataSource_basic(authURL),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckOpenStackConfigDataSourceID("data.openstack_config.config_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_config.config_1", "auth_url", authURL),
					resource.TestCheckResourceAttr(
						"data.openstack_config.config_1", "swauth", "false"),
				),
			},
		},
	})
}

func testAccCheckOpenStackConfigDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find config data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Endpoint data source ID not set")
		}

		return nil
	}
}

func testAccOpenStackOpenStackConfigDataSource_basic(authURL string) string {
	return fmt.Sprintf(`
	data "openstack_config" "config_1" {
      auth_url = "%s"
	}
`, authURL)
}
