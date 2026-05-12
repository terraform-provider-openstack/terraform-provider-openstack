package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccOpenStackNetworkingBGPPeerV2DataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingBGPPeerV2DataSourceBGPPeerBasic(),
			},
			{
				Config: testAccOpenStackNetworkingBGPPeerV2DataSourceBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.openstack_networking_bgp_peer_v2.peer_acc"),
					resource.TestCheckResourceAttr("data.openstack_networking_bgp_peer_v2.peer_acc", "name", "peer_acc1"),
					resource.TestCheckResourceAttr("data.openstack_networking_bgp_peer_v2.peer_acc", "peer_ip", "127.0.0.1"),
					resource.TestCheckResourceAttr("data.openstack_networking_bgp_peer_v2.peer_acc", "remote_as", "65001"),
					resource.TestCheckResourceAttr("data.openstack_networking_bgp_peer_v2.peer_acc", "auth_type", "md5"),
				),
			},
		},
	})
}

func TestAccOpenStackNetworkingBGPPeerV2DataSource_peerID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingBGPPeerV2DataSourceBGPPeerID(),
			},
			{
				Config: testAccOpenStackNetworkingBGPPeerV2DataSourceID(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.openstack_networking_bgp_peer_v2.peer_acc"),
					resource.TestCheckResourceAttr("data.openstack_networking_bgp_peer_v2.peer_acc", "name", "peer_acc2"),
					resource.TestCheckResourceAttr("data.openstack_networking_bgp_peer_v2.peer_acc", "peer_ip", "127.0.0.2"),
					resource.TestCheckResourceAttr("data.openstack_networking_bgp_peer_v2.peer_acc", "remote_as", "65002"),
					resource.TestCheckResourceAttr("data.openstack_networking_bgp_peer_v2.peer_acc", "auth_type", "md5"),
				),
			},
		},
	})
}

func TestAccOpenStackNetworkingBGPPeerV2DataSource_name(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingBGPPeerV2DataSourceBGPPeerName(),
			},
			{
				Config: testAccOpenStackNetworkingBGPPeerV2DataSourceName(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingNetworkV2DataSourceID("data.openstack_networking_bgp_peer_v2.peer_acc"),
					resource.TestCheckResourceAttr("data.openstack_networking_bgp_peer_v2.peer_acc", "name", "peer_acc2"),
					resource.TestCheckResourceAttr("data.openstack_networking_bgp_peer_v2.peer_acc", "peer_ip", "127.0.0.2"),
					resource.TestCheckResourceAttr("data.openstack_networking_bgp_peer_v2.peer_acc", "remote_as", "65002"),
					resource.TestCheckResourceAttr("data.openstack_networking_bgp_peer_v2.peer_acc", "auth_type", "md5"),
				),
			},
		},
	})
}

const testAccOpenStackNetworkingBGPPeerV2DataSourceBGPPeerPeer1 = `
resource "openstack_networking_bgp_peer_v2" "peer_acc1" {
  name = "peer_acc1"
  peer_ip = "127.0.0.1"
  remote_as = 65001
  auth_type = "md5"
  password = "secret"
}`

const testAccOpenStackNetworkingBGPPeerV2DataSourceBGPPeerPeer2 = `
resource "openstack_networking_bgp_peer_v2" "peer_acc2" {
  name = "peer_acc2"
  peer_ip = "127.0.0.2"
  remote_as = 65002
  auth_type = "md5"
  password = "secret"
}`

func testAccOpenStackNetworkingBGPPeerV2DataSourceBGPPeerBasic() string {
	return fmt.Sprintf(`
%s`, testAccOpenStackNetworkingBGPPeerV2DataSourceBGPPeerPeer1)
}

func testAccOpenStackNetworkingBGPPeerV2DataSourceBasic() string {
	return fmt.Sprintf(`
%s

data "openstack_networking_bgp_peer_v2" "peer_acc" {
}
`, testAccOpenStackNetworkingBGPPeerV2DataSourceBGPPeerPeer1)
}

func testAccOpenStackNetworkingBGPPeerV2DataSourceBGPPeerID() string {
	return fmt.Sprintf(`
%s

%s`, testAccOpenStackNetworkingBGPPeerV2DataSourceBGPPeerPeer1,
		testAccOpenStackNetworkingBGPPeerV2DataSourceBGPPeerPeer2)
}

func testAccOpenStackNetworkingBGPPeerV2DataSourceID() string {
	return fmt.Sprintf(`
%s

%s

data "openstack_networking_bgp_peer_v2" "peer_acc" {
  peer_id = resource.openstack_networking_bgp_peer_v2.peer_acc2.id
}
`, testAccOpenStackNetworkingBGPPeerV2DataSourceBGPPeerPeer1,
		testAccOpenStackNetworkingBGPPeerV2DataSourceBGPPeerPeer2)
}

func testAccOpenStackNetworkingBGPPeerV2DataSourceBGPPeerName() string {
	return fmt.Sprintf(`
%s

%s`, testAccOpenStackNetworkingBGPPeerV2DataSourceBGPPeerPeer1,
		testAccOpenStackNetworkingBGPPeerV2DataSourceBGPPeerPeer2)
}

func testAccOpenStackNetworkingBGPPeerV2DataSourceName() string {
	return fmt.Sprintf(`
%s

%s

data "openstack_networking_bgp_peer_v2" "peer_acc" {
  name = "peer_acc2"
}
`, testAccOpenStackNetworkingBGPPeerV2DataSourceBGPPeerPeer1,
		testAccOpenStackNetworkingBGPPeerV2DataSourceBGPPeerPeer2)
}
