package openstack

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func testAccHypervisorDataSource() string {
	return fmt.Sprintf(`
data "openstack_compute_hypervisor_v2" "host01" {
  hostname = "%s"
}
    `, osHypervisorEnvironment)
}

func TestAccComputeHypervisorV2DataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccHypervisorDataSource(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeHypervisorV2DataSourceID("data.openstack_compute_hypervisor_v2.host01"),
					resource.TestCheckResourceAttrSet("data.openstack_compute_hypervisor_v2.host01", "hostname"),
				),
			},
		},
	})
}

func testAccCheckComputeHypervisorV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("Data source ID not set")
		}

		return nil
	}
}
