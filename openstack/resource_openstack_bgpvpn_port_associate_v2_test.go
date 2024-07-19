package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/bgpvpns"
)

func TestAccBGPVPNPortAssociateV2_basic(t *testing.T) {
	var pa bgpvpns.PortAssociation
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckBGPVPN(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBGPVPNPortAssociateV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBGPVPNPortAssociateV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBGPVPNPortAssociateV2Exists(
						"openstack_bgpvpn_port_associate_v2.association_1", &pa),
					resource.TestCheckResourceAttrPtr("openstack_bgpvpn_port_associate_v2.association_1", "port_id", &pa.PortID),
					resource.TestCheckResourceAttr("openstack_bgpvpn_port_associate_v2.association_1", "advertise_fixed_ips", "true"),
				),
			},
			{
				Config: testAccBGPVPNPortAssociateV2ConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBGPVPNPortAssociateV2Exists(
						"openstack_bgpvpn_port_associate_v2.association_1", &pa),
					resource.TestCheckResourceAttrPtr("openstack_bgpvpn_port_associate_v2.association_1", "port_id", &pa.PortID),
					resource.TestCheckResourceAttr("openstack_bgpvpn_port_associate_v2.association_1", "advertise_fixed_ips", "false"),
				),
			},
		},
	})
}

func TestAccBGPVPNPortAssociateV2_no_fixed_ips_advertise(t *testing.T) {
	var pa bgpvpns.PortAssociation
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckBGPVPN(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBGPVPNPortAssociateV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBGPVPNPortAssociateV2ConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBGPVPNPortAssociateV2Exists(
						"openstack_bgpvpn_port_associate_v2.association_1", &pa),
					resource.TestCheckResourceAttrPtr("openstack_bgpvpn_port_associate_v2.association_1", "port_id", &pa.PortID),
					resource.TestCheckResourceAttr("openstack_bgpvpn_port_associate_v2.association_1", "advertise_fixed_ips", "false"),
				),
			},
		},
	})
}

func testAccCheckBGPVPNPortAssociateV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.NetworkingV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_bgpvpn_port_associate_v2" {
			continue
		}

		bgpvpnID, id, err := parsePairedIDs(rs.Primary.ID, "openstack_bgpvpn_port_associate_v2")
		if err != nil {
			return err
		}

		_, err = bgpvpns.GetPortAssociation(networkingClient, bgpvpnID, id).Extract()
		if err == nil {
			return fmt.Errorf("BGP VPN port association (%s) still exists", id)
		}
		if _, ok := err.(gophercloud.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckBGPVPNPortAssociateV2Exists(n string, pa *bgpvpns.PortAssociation) resource.TestCheckFunc {
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

		bgpvpnID, id, err := parsePairedIDs(rs.Primary.ID, "openstack_bgpvpn_port_associate_v2")
		if err != nil {
			return err
		}

		found, err := bgpvpns.GetPortAssociation(networkingClient, bgpvpnID, id).Extract()
		if err != nil {
			return err
		}

		if found.ID != id {
			return fmt.Errorf("BGP VPN port association not found")
		}

		*pa = *found

		return nil
	}
}

const testAccBGPVPNPortAssociateV2Config = `
resource "openstack_bgpvpn_v2" "bgpvpn_1" {
  name = "bgpvpn_1"
}

resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
  }
}

resource "openstack_bgpvpn_port_associate_v2" "association_1" {
  bgpvpn_id = openstack_bgpvpn_v2.bgpvpn_1.id
  port_id = openstack_networking_port_v2.port_1.id
}
`

const testAccBGPVPNPortAssociateV2ConfigUpdate = `
resource "openstack_bgpvpn_v2" "bgpvpn_1" {
  name = "bgpvpn_1"
}

resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
  }
}

resource "openstack_bgpvpn_port_associate_v2" "association_1" {
  bgpvpn_id = openstack_bgpvpn_v2.bgpvpn_1.id
  port_id = openstack_networking_port_v2.port_1.id
  advertise_fixed_ips = false
}
`
