package openstack

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/bgpvpns"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccBGPVPNPortAssociateV2_basic(t *testing.T) {
	var pa bgpvpns.PortAssociation

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckBGPVPN(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBGPVPNPortAssociateV2Destroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccBGPVPNPortAssociateV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBGPVPNPortAssociateV2Exists(t.Context(),
						"openstack_bgpvpn_port_associate_v2.association_1", &pa),
					resource.TestCheckResourceAttrPtr("openstack_bgpvpn_port_associate_v2.association_1", "port_id", &pa.PortID),
					resource.TestCheckResourceAttr("openstack_bgpvpn_port_associate_v2.association_1", "advertise_fixed_ips", "true"),
				),
			},
			{
				Config: testAccBGPVPNPortAssociateV2ConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBGPVPNPortAssociateV2Exists(t.Context(),
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
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckBGPVPN(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBGPVPNPortAssociateV2Destroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccBGPVPNPortAssociateV2ConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBGPVPNPortAssociateV2Exists(t.Context(),
						"openstack_bgpvpn_port_associate_v2.association_1", &pa),
					resource.TestCheckResourceAttrPtr("openstack_bgpvpn_port_associate_v2.association_1", "port_id", &pa.PortID),
					resource.TestCheckResourceAttr("openstack_bgpvpn_port_associate_v2.association_1", "advertise_fixed_ips", "false"),
				),
			},
		},
	})
}

func testAccCheckBGPVPNPortAssociateV2Destroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		networkingClient, err := config.NetworkingV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_bgpvpn_port_associate_v2" {
				continue
			}

			bgpvpnID, id, err := parsePairedIDs(rs.Primary.ID, "openstack_bgpvpn_port_associate_v2")
			if err != nil {
				return err
			}

			_, err = bgpvpns.GetPortAssociation(ctx, networkingClient, bgpvpnID, id).Extract()
			if err == nil {
				return fmt.Errorf("BGP VPN port association (%s) still exists", id)
			}

			if !gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return err
			}
		}

		return nil
	}
}

func testAccCheckBGPVPNPortAssociateV2Exists(ctx context.Context, n string, pa *bgpvpns.PortAssociation) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		networkingClient, err := config.NetworkingV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %w", err)
		}

		bgpvpnID, id, err := parsePairedIDs(rs.Primary.ID, "openstack_bgpvpn_port_associate_v2")
		if err != nil {
			return err
		}

		found, err := bgpvpns.GetPortAssociation(ctx, networkingClient, bgpvpnID, id).Extract()
		if err != nil {
			return err
		}

		if found.ID != id {
			return errors.New("BGP VPN port association not found")
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
  network_id = openstack_networking_network_v2.network_1.id
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = openstack_networking_network_v2.network_1.id

  fixed_ip {
    subnet_id = openstack_networking_subnet_v2.subnet_1.id
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
  network_id = openstack_networking_network_v2.network_1.id
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = openstack_networking_network_v2.network_1.id

  fixed_ip {
    subnet_id = openstack_networking_subnet_v2.subnet_1.id
  }
}

resource "openstack_bgpvpn_port_associate_v2" "association_1" {
  bgpvpn_id = openstack_bgpvpn_v2.bgpvpn_1.id
  port_id = openstack_networking_port_v2.port_1.id
  advertise_fixed_ips = false
}
`
