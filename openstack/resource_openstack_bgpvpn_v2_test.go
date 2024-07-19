package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/bgpvpns"
)

func TestAccBGPVPNV2_basic(t *testing.T) {
	var bgpvpn bgpvpns.BGPVPN
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckBGPVPN(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBGPVPNV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBGPVPNV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBGPVPNV2Exists(
						"openstack_bgpvpn_v2.bgpvpn_1", &bgpvpn),
					resource.TestCheckResourceAttrPtr("openstack_bgpvpn_v2.bgpvpn_1", "name", &bgpvpn.Name),
					resource.TestCheckResourceAttr("openstack_bgpvpn_v2.bgpvpn_1", "type", "l3"),
				),
			},
		},
	})
}

func testAccCheckBGPVPNV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.NetworkingV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_bgpvpn_v2" {
			continue
		}
		_, err = bgpvpns.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("BGPVPN (%s) still exists", rs.Primary.ID)
		}
		if _, ok := err.(gophercloud.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckBGPVPNV2Exists(n string, bgpvpn *bgpvpns.BGPVPN) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.NetworkingV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %s", err)
		}

		var found *bgpvpns.BGPVPN
		found, err = bgpvpns.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("BGP VPN not found")
		}

		*bgpvpn = *found

		return nil
	}
}

const testAccBGPVPNV2Config = `
resource "openstack_bgpvpn_v2" "bgpvpn_1" {
  name = "bgpvpn_1"
}
`
