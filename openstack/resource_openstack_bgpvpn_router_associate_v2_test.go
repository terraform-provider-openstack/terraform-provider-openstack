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

func TestAccBGPVPNRouterAssociateV2_basic(t *testing.T) {
	var ra bgpvpns.RouterAssociation

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckBGPVPN(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBGPVPNRouterAssociateV2Destroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccBGPVPNRouterAssociateV2Config,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBGPVPNRouterAssociateV2Exists(t.Context(),
						"openstack_bgpvpn_router_associate_v2.association_1", &ra),
					resource.TestCheckResourceAttrPtr("openstack_bgpvpn_router_associate_v2.association_1", "router_id", &ra.RouterID),
					resource.TestCheckResourceAttr("openstack_bgpvpn_router_associate_v2.association_1", "advertise_extra_routes", "true"),
				),
			},
			{
				Config: testAccBGPVPNRouterAssociateV2ConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBGPVPNRouterAssociateV2Exists(t.Context(),
						"openstack_bgpvpn_router_associate_v2.association_1", &ra),
					resource.TestCheckResourceAttrPtr("openstack_bgpvpn_router_associate_v2.association_1", "router_id", &ra.RouterID),
					resource.TestCheckResourceAttr("openstack_bgpvpn_router_associate_v2.association_1", "advertise_extra_routes", "false"),
				),
			},
		},
	})
}

func TestAccBGPVPNRouterAssociateV2_no_routes_advertise(t *testing.T) {
	var ra bgpvpns.RouterAssociation

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckBGPVPN(t)
			t.Skip(`bug in OpenStack which ignores {"advertise_extra_routes": false} on POST request, while handles on PUT request.`)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckBGPVPNRouterAssociateV2Destroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccBGPVPNRouterAssociateV2ConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBGPVPNRouterAssociateV2Exists(t.Context(),
						"openstack_bgpvpn_router_associate_v2.association_1", &ra),
					resource.TestCheckResourceAttrPtr("openstack_bgpvpn_router_associate_v2.association_1", "router_id", &ra.RouterID),
					resource.TestCheckResourceAttr("openstack_bgpvpn_router_associate_v2.association_1", "advertise_extra_routes", "false"),
				),
			},
		},
	})
}

func testAccCheckBGPVPNRouterAssociateV2Destroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		networkingClient, err := config.NetworkingV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_bgpvpn_router_associate_v2" {
				continue
			}

			bgpvpnID, id, err := parsePairedIDs(rs.Primary.ID, "openstack_bgpvpn_router_associate_v2")
			if err != nil {
				return err
			}

			_, err = bgpvpns.GetRouterAssociation(ctx, networkingClient, bgpvpnID, id).Extract()
			if err == nil {
				return fmt.Errorf("BGP VPN router association (%s) still exists", id)
			}

			if !gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return err
			}
		}

		return nil
	}
}

func testAccCheckBGPVPNRouterAssociateV2Exists(ctx context.Context, n string, ra *bgpvpns.RouterAssociation) resource.TestCheckFunc {
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

		bgpvpnID, id, err := parsePairedIDs(rs.Primary.ID, "openstack_bgpvpn_router_associate_v2")
		if err != nil {
			return err
		}

		found, err := bgpvpns.GetRouterAssociation(ctx, networkingClient, bgpvpnID, id).Extract()
		if err != nil {
			return err
		}

		if found.ID != id {
			return errors.New("BGP VPN router association not found")
		}

		*ra = *found

		return nil
	}
}

const testAccBGPVPNRouterAssociateV2Config = `
resource "openstack_bgpvpn_v2" "bgpvpn_1" {
  name = "bgpvpn_1"
}

resource "openstack_networking_router_v2" "router_1" {
  name = "router_1"
}

resource "openstack_bgpvpn_router_associate_v2" "association_1" {
  bgpvpn_id = openstack_bgpvpn_v2.bgpvpn_1.id
  router_id = openstack_networking_router_v2.router_1.id
}
`

const testAccBGPVPNRouterAssociateV2ConfigUpdate = `
resource "openstack_bgpvpn_v2" "bgpvpn_1" {
  name = "bgpvpn_1"
}

resource "openstack_networking_router_v2" "router_1" {
  name = "router_1"
}

resource "openstack_bgpvpn_router_associate_v2" "association_1" {
  bgpvpn_id = openstack_bgpvpn_v2.bgpvpn_1.id
  router_id = openstack_networking_router_v2.router_1.id
  advertise_extra_routes = false
}
`
