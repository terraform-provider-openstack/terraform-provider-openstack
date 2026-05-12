package openstack

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/bgp/peers"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOpenStackNetworkingBGPSpeakerV2DataSource_basic(t *testing.T) {
	var peer peers.BGPPeer

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingBGPSpeakerV2DataSourceBGPSpeakerBasic(),
			},
			{
				Config: testAccOpenStackNetworkingBGPSpeakerV2DataSourceBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.openstack_networking_bgp_speaker_v2.speaker_acc"),
					testAccCheckNetworkingV2BGPPeerExists(t.Context(), "openstack_networking_bgp_peer_v2.peer_acc", &peer),
					resource.TestCheckResourceAttr("data.openstack_networking_bgp_speaker_v2.speaker_acc", "name", "speaker_acc"),
					resource.TestCheckResourceAttr("data.openstack_networking_bgp_speaker_v2.speaker_acc", "ip_version", "4"),
					resource.TestCheckResourceAttr("data.openstack_networking_bgp_speaker_v2.speaker_acc", "local_as", "65001"),
					resource.TestCheckResourceAttr("data.openstack_networking_bgp_speaker_v2.speaker_acc", "advertise_floating_ip_host_routes", "true"),
					resource.TestCheckResourceAttr("data.openstack_networking_bgp_speaker_v2.speaker_acc", "advertise_tenant_networks", "true"),
					resource.TestCheckResourceAttr("data.openstack_networking_bgp_speaker_v2.speaker_acc", "networks.0", osExtGwID),
					resource.TestCheckResourceAttrPtr("data.openstack_networking_bgp_speaker_v2.speaker_acc", "peers.0", &peer.ID),
				),
			},
		},
	})
}

func TestAccOpenStackNetworkingBGPSpeakerV2DataSource_speakerID(t *testing.T) {
	var peer peers.BGPPeer

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingBGPSpeakerV2DataSourceBGPSpeakerID(),
			},
			{
				Config: testAccOpenStackNetworkingBGPSpeakerV2DataSourceID(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.openstack_networking_bgp_speaker_v2.speaker_acc"),
					testAccCheckNetworkingV2BGPPeerExists(t.Context(), "openstack_networking_bgp_peer_v2.peer_acc", &peer),
					resource.TestCheckResourceAttr("data.openstack_networking_bgp_speaker_v2.speaker_acc", "name", "speaker_acc2"),
					resource.TestCheckResourceAttr("data.openstack_networking_bgp_speaker_v2.speaker_acc", "ip_version", "4"),
					resource.TestCheckResourceAttr("data.openstack_networking_bgp_speaker_v2.speaker_acc", "local_as", "65002"),
					resource.TestCheckResourceAttr("data.openstack_networking_bgp_speaker_v2.speaker_acc", "advertise_floating_ip_host_routes", "true"),
					resource.TestCheckResourceAttr("data.openstack_networking_bgp_speaker_v2.speaker_acc", "advertise_tenant_networks", "true"),
					resource.TestCheckResourceAttr("data.openstack_networking_bgp_speaker_v2.speaker_acc", "networks.0", osExtGwID),
					resource.TestCheckResourceAttrPtr("data.openstack_networking_bgp_speaker_v2.speaker_acc", "peers.0", &peer.ID),
				),
			},
		},
	})
}

func TestAccOpenStackNetworkingBGPSpeakerV2DataSource_name(t *testing.T) {
	var peer peers.BGPPeer

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingBGPSpeakerV2DataSourceBGPSpeakerName(),
			},
			{
				Config: testAccOpenStackNetworkingBGPSpeakerV2DataSourceName(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.openstack_networking_bgp_speaker_v2.speaker_acc"),
					testAccCheckNetworkingV2BGPPeerExists(t.Context(), "openstack_networking_bgp_peer_v2.peer_acc", &peer),
					resource.TestCheckResourceAttr("data.openstack_networking_bgp_speaker_v2.speaker_acc", "name", "speaker_acc2"),
					resource.TestCheckResourceAttr("data.openstack_networking_bgp_speaker_v2.speaker_acc", "ip_version", "4"),
					resource.TestCheckResourceAttr("data.openstack_networking_bgp_speaker_v2.speaker_acc", "local_as", "65002"),
					resource.TestCheckResourceAttr("data.openstack_networking_bgp_speaker_v2.speaker_acc", "advertise_floating_ip_host_routes", "true"),
					resource.TestCheckResourceAttr("data.openstack_networking_bgp_speaker_v2.speaker_acc", "advertise_tenant_networks", "true"),
					resource.TestCheckResourceAttr("data.openstack_networking_bgp_speaker_v2.speaker_acc", "networks.0", osExtGwID),
					resource.TestCheckResourceAttrPtr("data.openstack_networking_bgp_speaker_v2.speaker_acc", "peers.0", &peer.ID),
				),
			},
		},
	})
}

const testAccOpenStackNetworkingBGPSpeakerV2DataSourcePeer = `
resource "openstack_networking_bgp_peer_v2" "peer_acc" {
  name = "peer_acc"
  peer_ip = "127.0.0.1"
  remote_as = 65001
  auth_type = "md5"
  password = "secret"
}
`

func testAccOpenStackNetworkingBGPSpeakerV2DataSourceBGPSpeakerBasic() string {
	return fmt.Sprintf(`
%s

resource "openstack_networking_bgp_speaker_v2" "speaker_acc" {
  name       = "speaker_acc"
  ip_version = 4
  local_as   = 65001
  advertise_floating_ip_host_routes = true
  advertise_tenant_networks = true

  networks = [
    "%s",
  ]

  peers = [
    openstack_networking_bgp_peer_v2.peer_acc.id,
  ]
}`, testAccOpenStackNetworkingBGPSpeakerV2DataSourcePeer, osExtGwID)
}

func testAccOpenStackNetworkingBGPSpeakerV2DataSourceBasic() string {
	return fmt.Sprintf(`
%s

data "openstack_networking_bgp_speaker_v2" "speaker_acc" {
}
`, testAccOpenStackNetworkingBGPSpeakerV2DataSourceBGPSpeakerBasic())
}

func testAccOpenStackNetworkingBGPSpeakerV2DataSourceBGPSpeakerID() string {
	return fmt.Sprintf(`
%s

resource "openstack_networking_bgp_speaker_v2" "speaker_acc1" {
  name       = "speaker_acc1"
  local_as   = 65001
}

resource "openstack_networking_bgp_speaker_v2" "speaker_acc2" {
  name       = "speaker_acc2"
  ip_version = 4
  local_as   = 65002
  advertise_floating_ip_host_routes = true
  advertise_tenant_networks = true

  networks = [
    "%s",
  ]

  peers = [
    openstack_networking_bgp_peer_v2.peer_acc.id,
  ]
}`, testAccOpenStackNetworkingBGPSpeakerV2DataSourcePeer, osExtGwID)
}

func testAccOpenStackNetworkingBGPSpeakerV2DataSourceID() string {
	return fmt.Sprintf(`
%s

data "openstack_networking_bgp_speaker_v2" "speaker_acc" {
  speaker_id = resource.openstack_networking_bgp_speaker_v2.speaker_acc2.id
}
`, testAccOpenStackNetworkingBGPSpeakerV2DataSourceBGPSpeakerName())
}

func testAccOpenStackNetworkingBGPSpeakerV2DataSourceBGPSpeakerName() string {
	return fmt.Sprintf(`
%s

resource "openstack_networking_bgp_speaker_v2" "speaker_acc1" {
  name       = "speaker_acc1"
  local_as   = 65001
}

resource "openstack_networking_bgp_speaker_v2" "speaker_acc2" {
  name       = "speaker_acc2"
  ip_version = 4
  local_as   = 65002
  advertise_floating_ip_host_routes = true
  advertise_tenant_networks = true

  networks = [
    "%s",
  ]

  peers = [
    openstack_networking_bgp_peer_v2.peer_acc.id,
  ]
}`, testAccOpenStackNetworkingBGPSpeakerV2DataSourcePeer, osExtGwID)
}

func testAccOpenStackNetworkingBGPSpeakerV2DataSourceName() string {
	return fmt.Sprintf(`
%s

data "openstack_networking_bgp_speaker_v2" "speaker_acc" {
  name = "speaker_acc2"
}
`, testAccOpenStackNetworkingBGPSpeakerV2DataSourceBGPSpeakerName())
}
