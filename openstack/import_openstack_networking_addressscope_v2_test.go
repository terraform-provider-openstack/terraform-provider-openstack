package openstack

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccNetworkingV2AddressScopeImport_basic(t *testing.T) {
	resourceName := "openstack_networking_addressscope_v2.addressscope_1"
	name := acctest.RandomWithPrefix("tf-acc-addrscope")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2AddressScopeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2AddressScopeBasic(name),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
