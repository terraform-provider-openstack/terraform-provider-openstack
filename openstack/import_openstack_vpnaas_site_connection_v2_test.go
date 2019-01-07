package openstack

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccSiteConnectionV2_importBasic(t *testing.T) {
	resourceName := "openstack_vpnaas_site_connection_v2.conn_1"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckVPN(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSiteConnectionV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteConnectionV2_basic,
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
