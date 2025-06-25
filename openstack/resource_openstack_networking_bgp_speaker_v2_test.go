package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/bgp/speakers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccNetworkingV2BGPSpeaker_basic(t *testing.T) {
	var bgpSpeaker speakers.BGPSpeaker

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2BGPSpeakerDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2BGPSpeakerBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2BGPSpeakerExists(t.Context(), "openstack_networking_bgp_speaker_v2.speaker_1", &bgpSpeaker),
					resource.TestCheckResourceAttr(
						"openstack_networking_bgp_speaker_v2.speaker_1", "name", "speaker_1"),
					resource.TestCheckResourceAttr(
						"openstack_networking_bgp_speaker_v2.speaker_1", "local_as", "1004"),
					resource.TestCheckResourceAttr(
						"openstack_networking_bgp_speaker_v2.speaker_1", "peers.#", "1"),
					resource.TestCheckResourceAttr(
						"openstack_networking_bgp_speaker_v2.speaker_1", "networks.#", "1"),
					resource.TestCheckResourceAttr(
						"openstack_networking_bgp_speaker_v2.speaker_1", "ip_version", "4"),
					resource.TestCheckResourceAttr(
						"openstack_networking_bgp_speaker_v2.speaker_1", "advertise_floating_ip_host_routes", "true"),
					resource.TestCheckResourceAttr(
						"openstack_networking_bgp_speaker_v2.speaker_1", "advertise_tenant_networks", "true"),
				),
			},
			{
				Config: testAccNetworkingV2BGPSpeakerUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_networking_bgp_speaker_v2.speaker_1", "name", ""),
					resource.TestCheckResourceAttr(
						"openstack_networking_bgp_speaker_v2.speaker_1", "local_as", "1004"),
					resource.TestCheckResourceAttr(
						"openstack_networking_bgp_speaker_v2.speaker_1", "peers.#", "0"),
					resource.TestCheckResourceAttr(
						"openstack_networking_bgp_speaker_v2.speaker_1", "networks.#", "2"),
					resource.TestCheckResourceAttr(
						"openstack_networking_bgp_speaker_v2.speaker_1", "advertise_floating_ip_host_routes", "false"),
					resource.TestCheckResourceAttr(
						"openstack_networking_bgp_speaker_v2.speaker_1", "advertise_tenant_networks", "false"),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2BGPSpeakerDestroy(ctx context.Context) resource.TestCheckFunc {
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

			_, err := speakers.Get(ctx, networkingClient, rs.Primary.ID).Extract()
			if err == nil {
				return errors.New("Network still exists")
			}
		}

		return nil
	}
}

func testAccCheckNetworkingV2BGPSpeakerExists(ctx context.Context, n string, speaker *speakers.BGPSpeaker) resource.TestCheckFunc {
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

		found, err := speakers.Get(ctx, networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return errors.New("Network not found")
		}

		*speaker = *found

		return nil
	}
}

const testAccNetworkingV2BGPSpeakerBasic = `
resource "openstack_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = true
}

resource "openstack_networking_bgp_peer_v2" "peer_1" {
  name      = "test"
  auth_type = "md5"
  password  = "123"
  remote_as = "1001"
  peer_ip  = "1.0.0.2"
}

resource "openstack_networking_bgp_speaker_v2" "speaker_1" {
  name     = "speaker_1"
  local_as = "1004"

  peers = [
    openstack_networking_bgp_peer_v2.peer_1.id,
  ]

  networks = [
    openstack_networking_network_v2.network_1.id,
  ]
}
`

const testAccNetworkingV2BGPSpeakerUpdate = `
resource "openstack_networking_network_v2" "network_2" {
  name           = "network_1"
  admin_state_up = true
}

resource "openstack_networking_network_v2" "network_3" {
  name           = "network_1"
  admin_state_up = true
}

resource "openstack_networking_bgp_speaker_v2" "speaker_1" {
  advertise_floating_ip_host_routes = false
  advertise_tenant_networks         = false

  local_as = "1004"

  networks = [
    openstack_networking_network_v2.network_2.id,
    openstack_networking_network_v2.network_3.id,
  ]
}
`
