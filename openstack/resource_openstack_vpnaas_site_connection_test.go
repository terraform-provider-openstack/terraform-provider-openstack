package openstack

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/siteconnections"
)

func TestAccSiteConnectionV2_basic(t *testing.T) {
	var conn siteconnections.Connection
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckVPN(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSiteConnectionV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSiteConnectionV2Basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSiteConnectionV2Exists(
						"openstack_vpnaas_site_connection_v2.conn_1", &conn),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_site_connection_v2.conn_1", "ikepolicy_id", &conn.IKEPolicyID),
					resource.TestCheckResourceAttr("openstack_vpnaas_site_connection_v2.conn_1", "admin_state_up", strconv.FormatBool(conn.AdminStateUp)),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_site_connection_v2.conn_1", "psk", &conn.PSK),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_site_connection_v2.conn_1", "ipsecpolicy_id", &conn.IPSecPolicyID),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_site_connection_v2.conn_1", "vpnservice_id", &conn.VPNServiceID),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_site_connection_v2.conn_1", "local_ep_group_id", &conn.LocalEPGroupID),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_site_connection_v2.conn_1", "local_id", &conn.LocalID),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_site_connection_v2.conn_1", "peer_ep_group_id", &conn.PeerEPGroupID),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_site_connection_v2.conn_1", "name", &conn.Name),
				),
			},
		},
	})
}

func testAccCheckSiteConnectionV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.NetworkingV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_vpnaas_site_connection" {
			continue
		}
		_, err = siteconnections.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Site connection (%s) still exists", rs.Primary.ID)
		}
		if _, ok := err.(gophercloud.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckSiteConnectionV2Exists(n string, conn *siteconnections.Connection) resource.TestCheckFunc {
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

		var found *siteconnections.Connection

		found, err = siteconnections.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*conn = *found

		return nil
	}
}

func testAccSiteConnectionV2Basic() string {
	return fmt.Sprintf(`
	resource "openstack_networking_network_v2" "network_1" {
		name           = "tf_test_network"
  		admin_state_up = "true"
	}

	resource "openstack_networking_subnet_v2" "subnet_1" {
  		network_id = "${openstack_networking_network_v2.network_1.id}"
  		cidr       = "192.168.199.0/24"
  		ip_version = 4
	}

	resource "openstack_networking_router_v2" "router_1" {
  		name             = "my_router"
  		external_network_id = "%s"
	}

	resource "openstack_networking_router_interface_v2" "router_interface_1" {
  		router_id = "${openstack_networking_router_v2.router_1.id}"
  		subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
	}

	resource "openstack_vpnaas_service_v2" "service_1" {
		router_id = "${openstack_networking_router_v2.router_1.id}",
		admin_state_up = "false"
	}

	resource "openstack_vpnaas_ipsec_policy_v2" "policy_1" {
	}

	resource "openstack_vpnaas_ike_policy_v2" "policy_2" {
	}

	resource "openstack_vpnaas_endpoint_group_v2" "group_1" {
		type = "cidr"
		endpoints = ["10.0.0.24/24", "10.0.0.25/24"]
	}
	resource "openstack_vpnaas_endpoint_group_v2" "group_2" {
		type = "subnet"
		endpoints = [ "${openstack_networking_subnet_v2.subnet_1.id}" ]
	}

	resource "openstack_vpnaas_site_connection_v2" "conn_1" {
		name = "connection_1"
		ikepolicy_id = "${openstack_vpnaas_ike_policy_v2.policy_2.id}"
		ipsecpolicy_id = "${openstack_vpnaas_ipsec_policy_v2.policy_1.id}"
		vpnservice_id = "${openstack_vpnaas_service_v2.service_1.id}"
		psk = "secret"
		peer_address = "192.168.10.1"
		peer_id = "192.168.10.1"
		local_ep_group_id = "${openstack_vpnaas_endpoint_group_v2.group_2.id}"
		peer_ep_group_id = "${openstack_vpnaas_endpoint_group_v2.group_1.id}"
		depends_on = ["openstack_networking_router_interface_v2.router_interface_1"]
	}
	`, osExtGwID)
}
