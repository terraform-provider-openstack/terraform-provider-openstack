package openstack

import (
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccNetworkingV2PortExtraDHCPOptions_importBasic(t *testing.T) {
	resourceName := "openstack_networking_port_extradhcp_options_v2.opts_1"
	networkName := acctest.RandomWithPrefix("tf-acc-network")
	subnetName := acctest.RandomWithPrefix("tf-acc-subnet")
	portName := acctest.RandomWithPrefix("tf-acc-port")
	optsName := acctest.RandomWithPrefix("tf-acc-extradhcp-opts")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortExtraDHCPOptionsDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2PortExtraDHCPOptions_basic(networkName, subnetName, portName, optsName),
			},
			resource.TestStep{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"name",
				},
			},
		},
	})
}
