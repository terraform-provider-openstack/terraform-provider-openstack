package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccOpenStackIdentityV3DomainDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackIdentityV3DomainDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIdentityV3DomainDataSourceID("data.openstack_identity_domain_v3.domain_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_identity_domain_v3.domain_1", "name", "Default"),
				),
			},
		},
	})
}

func testAccCheckIdentityV3DomainDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find domain data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Domain data source ID not set")
		}

		return nil
	}
}

const testAccOpenStackIdentityV3DomainDataSourceBasic = `
data "openstack_identity_domain_v3" "domain_1" {
    name = "Default"
}
`
