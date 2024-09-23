package openstack

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/bgpvpns"
)

func TestAccBGPVPNNetworkAssociateV2_basic(t *testing.T) {
	var na bgpvpns.NetworkAssociation
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckBGPVPN(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBGPVPNNetworkAssociateV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBGPVPNNetworkAssociateV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBGPVPNNetworkAssociateV2Exists(
						"openstack_bgpvpn_network_associate_v2.association_1", &na),
					resource.TestCheckResourceAttrPtr("openstack_bgpvpn_network_associate_v2.association_1", "network_id", &na.NetworkID),
				),
			},
		},
	})
}

func testAccCheckBGPVPNNetworkAssociateV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.NetworkingV2Client(context.TODO(), osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_bgpvpn_network_associate_v2" {
			continue
		}

		bgpvpnID, id, err := parsePairedIDs(rs.Primary.ID, "openstack_bgpvpn_network_associate_v2")
		if err != nil {
			return err
		}

		_, err = bgpvpns.GetNetworkAssociation(context.TODO(), networkingClient, bgpvpnID, id).Extract()
		if err == nil {
			return fmt.Errorf("BGP VPN network association (%s) still exists", id)
		}
		if !gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
			return err
		}
	}
	return nil
}

func testAccCheckBGPVPNNetworkAssociateV2Exists(n string, na *bgpvpns.NetworkAssociation) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.NetworkingV2Client(context.TODO(), osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %s", err)
		}

		bgpvpnID, id, err := parsePairedIDs(rs.Primary.ID, "openstack_bgpvpn_network_associate_v2")
		if err != nil {
			return err
		}

		found, err := bgpvpns.GetNetworkAssociation(context.TODO(), networkingClient, bgpvpnID, id).Extract()
		if err != nil {
			return err
		}

		if found.ID != id {
			return fmt.Errorf("BGP VPN network association not found")
		}

		*na = *found

		return nil
	}
}

const testAccBGPVPNNetworkAssociateV2Config = `
resource "openstack_bgpvpn_v2" "bgpvpn_1" {
  name = "bgpvpn_1"
}

resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
}

resource "openstack_bgpvpn_network_associate_v2" "association_1" {
  bgpvpn_id = openstack_bgpvpn_v2.bgpvpn_1.id
  network_id = openstack_networking_network_v2.network_1.id
}
`
