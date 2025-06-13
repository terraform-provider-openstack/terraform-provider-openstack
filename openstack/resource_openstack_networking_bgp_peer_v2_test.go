package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/bgp/peers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccNetworkingV2BGPPeer_basic(t *testing.T) {
	var bgpPeer peers.BGPPeer

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2BGPPeerDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2BGPPeerBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2BGPPeerExists(t.Context(), "openstack_networking_bgp_peer_v2.peer_1", &bgpPeer),
					resource.TestCheckResourceAttr(
						"openstack_networking_bgp_peer_v2.peer_1", "name", "peer_1"),
				),
			},
			{
				Config: testAccNetworkingV2BGPPeerUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_networking_bgp_peer_v2.peer_1", "name", ""),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2BGPPeerDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		networkingClient, err := config.NetworkingV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_networking_bgp_peer_v2" {
				continue
			}

			_, err := peers.Get(ctx, networkingClient, rs.Primary.ID).Extract()
			if err == nil {
				return errors.New("Network still exists")
			}
		}

		return nil
	}
}

func testAccCheckNetworkingV2BGPPeerExists(ctx context.Context, n string, peer *peers.BGPPeer) resource.TestCheckFunc {
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

		found, err := peers.Get(ctx, networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return errors.New("Network not found")
		}

		*peer = *found

		return nil
	}
}

const testAccNetworkingV2BGPPeerBasic = `
resource "openstack_networking_bgp_peer_v2" "peer_1" {
  name      = "peer_1"
  auth_type = "md5"
  password  = "123"
  remote_as = "1004"
  peer_ip  = "1.0.0.2"
}
`

const testAccNetworkingV2BGPPeerUpdate = `
resource "openstack_networking_bgp_peer_v2" "peer_1" {
  auth_type = "md5"
  password  = "456"
  remote_as = "1004"
  peer_ip  = "1.0.0.2"
}
`
