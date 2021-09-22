package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/portsecurity"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/qos/policies"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
)

type testPortWithExtensions struct {
	ports.Port
	portsecurity.PortSecurityExt
	policies.QoSPolicyExt
}

func TestAccNetworkingV2Port_basic(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_noIP(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortNoIP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2PortCountFixedIPs(&port, 1),
				),
			},
			{
				Config: testAccNetworkingV2PortNoIPEmptyUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2PortCountFixedIPs(&port, 1),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_multipleNoIP(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortMultipleNoIP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2PortCountFixedIPs(&port, 3),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_allowedAddressPairs(t *testing.T) {
	var network networks.Network
	var subnet subnets.Subnet
	var vrrpPort1, vrrpPort2, instancePort ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortAllowedAddressPairs1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.vrrp_subnet", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.vrrp_network", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.vrrp_port_1", &vrrpPort1),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.vrrp_port_2", &vrrpPort2),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.instance_port", &instancePort),
					testAccCheckNetworkingV2PortCountAllowedAddressPairs(&instancePort, 2),
					resource.TestCheckResourceAttr("openstack_networking_port_v2.vrrp_port_1", "description", "test vrrp port"),
				),
			},
			{
				Config: testAccNetworkingV2PortAllowedAddressPairs2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.vrrp_subnet", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.vrrp_network", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.vrrp_port_1", &vrrpPort1),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.vrrp_port_2", &vrrpPort2),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.instance_port", &instancePort),
					testAccCheckNetworkingV2PortCountAllowedAddressPairs(&instancePort, 2),
					resource.TestCheckResourceAttr("openstack_networking_port_v2.vrrp_port_1", "description", ""),
				),
			},
			{
				Config: testAccNetworkingV2PortAllowedAddressPairs3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.vrrp_subnet", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.vrrp_network", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.vrrp_port_1", &vrrpPort1),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.vrrp_port_2", &vrrpPort2),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.instance_port", &instancePort),
					testAccCheckNetworkingV2PortCountAllowedAddressPairs(&instancePort, 2),
				),
			},
			{
				Config: testAccNetworkingV2PortAllowedAddressPairs4,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.vrrp_subnet", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.vrrp_network", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.vrrp_port_1", &vrrpPort1),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.vrrp_port_2", &vrrpPort2),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.instance_port", &instancePort),
					testAccCheckNetworkingV2PortCountAllowedAddressPairs(&instancePort, 1),
				),
			},
			{
				Config: testAccNetworkingV2PortAllowedAddressPairs5,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.vrrp_subnet", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.vrrp_network", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.vrrp_port_1", &vrrpPort1),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.vrrp_port_2", &vrrpPort2),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.instance_port", &instancePort),
					testAccCheckNetworkingV2PortCountAllowedAddressPairs(&instancePort, 0),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_allowedAddressPairsNoMAC(t *testing.T) {
	var network networks.Network
	var subnet subnets.Subnet
	var vrrpPort1, vrrpPort2, instancePort ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortAllowedAddressPairsNoMAC,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.vrrp_subnet", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.vrrp_network", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.vrrp_port_1", &vrrpPort1),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.vrrp_port_2", &vrrpPort2),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.instance_port", &instancePort),
					testAccCheckNetworkingV2PortCountAllowedAddressPairs(&instancePort, 2),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_multipleFixedIPs(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortMultipleFixedIPs,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2PortCountFixedIPs(&port, 3),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_timeout(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortTimeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_fixedIPs(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortFixedIPs,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "all_fixed_ips.0", "192.168.199.23"),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "all_fixed_ips.1", "192.168.199.24"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_updateSecurityGroups(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var secgroup1, secgroup2 groups.SecGroup
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortUpdateSecurityGroups1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_1", &secgroup1),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_2", &secgroup2),
					testAccCheckNetworkingV2PortCountSecurityGroups(&port, 1),
				),
			},
			{
				Config: testAccNetworkingV2PortUpdateSecurityGroups2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_1", &secgroup1),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_2", &secgroup2),
					testAccCheckNetworkingV2PortCountSecurityGroups(&port, 1),
				),
			},
			{
				Config: testAccNetworkingV2PortUpdateSecurityGroups3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_1", &secgroup1),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_2", &secgroup2),
					testAccCheckNetworkingV2PortCountSecurityGroups(&port, 2),
				),
			},
			{
				Config: testAccNetworkingV2PortUpdateSecurityGroups4,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_1", &secgroup1),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_2", &secgroup2),
					testAccCheckNetworkingV2PortCountSecurityGroups(&port, 1),
				),
			},
			{
				Config: testAccNetworkingV2PortUpdateSecurityGroups5,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_1", &secgroup1),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_2", &secgroup2),
					testAccCheckNetworkingV2PortCountSecurityGroups(&port, 0),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_noSecurityGroups(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var secgroup1, secgroup2 groups.SecGroup
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortNoSecurityGroups1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_1", &secgroup1),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_2", &secgroup2),
					testAccCheckNetworkingV2PortCountSecurityGroups(&port, 0),
				),
			},
			{
				Config: testAccNetworkingV2PortNoSecurityGroups2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_1", &secgroup1),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_2", &secgroup2),
					testAccCheckNetworkingV2PortCountSecurityGroups(&port, 1),
				),
			},
			{
				Config: testAccNetworkingV2PortNoSecurityGroups3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_1", &secgroup1),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_2", &secgroup2),
					testAccCheckNetworkingV2PortCountSecurityGroups(&port, 2),
				),
			},
			{
				Config: testAccNetworkingV2PortNoSecurityGroups4,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_1", &secgroup1),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_2", &secgroup2),
					testAccCheckNetworkingV2PortCountSecurityGroups(&port, 0),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_noFixedIP(t *testing.T) {
	var port ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortNoFixedIP1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "all_fixed_ips.#", "0"),
				),
			},
			{
				Config: testAccNetworkingV2PortNoFixedIP2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "all_fixed_ips.#", "1"),
				),
			},
			{
				Config: testAccNetworkingV2PortNoFixedIP1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "all_fixed_ips.#", "0"),
				),
			},
			{
				Config: testAccNetworkingV2PortNoFixedIP3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "all_fixed_ips.#", "2"),
				),
			},
			{
				Config: testAccNetworkingV2PortNoFixedIP1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "all_fixed_ips.#", "0"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_createExtraDHCPOpts(t *testing.T) {
	var network networks.Network
	var subnet subnets.Subnet
	var port ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortCreateExtraDhcpOpts,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "extra_dhcp_option.#", "2"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_updateExtraDHCPOpts(t *testing.T) {
	var network networks.Network
	var subnet subnets.Subnet
	var port ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
				),
			},
			{
				Config: testAccNetworkingV2PortUpdateExtraDhcpOpts1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "extra_dhcp_option.#", "1"),
				),
			},
			{
				Config: testAccNetworkingV2PortUpdateExtraDhcpOpts2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "extra_dhcp_option.#", "2"),
				),
			},
			{
				Config: testAccNetworkingV2PortUpdateExtraDhcpOpts3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "extra_dhcp_option.#", "2"),
				),
			},
			{
				Config: testAccNetworkingV2PortUpdateExtraDhcpOpts4,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "extra_dhcp_option.#", "2"),
				),
			},
			{
				Config: testAccNetworkingV2PortUpdateExtraDhcpOpts5,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "extra_dhcp_option.#", "2"),
				),
			},
			{
				Config: testAccNetworkingV2PortUpdateExtraDhcpOpts6,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckNoResourceAttr(
						"openstack_networking_port_v2.port_1", "extra_dhcp_option"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_adminStateUp_omit(t *testing.T) {
	var port ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortAdminStateUpOmit,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "admin_state_up", "true"),
					testAccCheckNetworkingV2PortAdminStateUp(&port, true),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_adminStateUp_true(t *testing.T) {
	var port ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortAdminStateUpTrue,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "admin_state_up", "true"),
					testAccCheckNetworkingV2PortAdminStateUp(&port, true),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_adminStateUp_false(t *testing.T) {
	var port ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortAdminStateUpFalse,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "admin_state_up", "false"),
					testAccCheckNetworkingV2PortAdminStateUp(&port, false),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_adminStateUp_update(t *testing.T) {
	var port ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortAdminStateUpOmit,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "admin_state_up", "true"),
					testAccCheckNetworkingV2PortAdminStateUp(&port, true),
				),
			},
			{
				Config: testAccNetworkingV2PortAdminStateUpFalse,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "admin_state_up", "false"),
					testAccCheckNetworkingV2PortAdminStateUp(&port, false),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_portSecurity_omit(t *testing.T) {
	var port testPortWithExtensions

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortSecurityOmit,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortWithExtensionsExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "port_security_enabled", "true"),
					testAccCheckNetworkingV2PortPortSecurityEnabled(&port, true),
				),
			},
			{
				Config: testAccNetworkingV2PortSecurityDisabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortWithExtensionsExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "port_security_enabled", "false"),
					testAccCheckNetworkingV2PortPortSecurityEnabled(&port, false),
				),
			},
			{
				Config: testAccNetworkingV2PortSecurityEnabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortWithExtensionsExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "port_security_enabled", "true"),
					testAccCheckNetworkingV2PortPortSecurityEnabled(&port, true),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_portSecurity_disabled(t *testing.T) {
	var port testPortWithExtensions

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortSecurityDisabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortWithExtensionsExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "port_security_enabled", "false"),
					testAccCheckNetworkingV2PortPortSecurityEnabled(&port, false),
				),
			},
			{
				Config: testAccNetworkingV2PortSecurityEnabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortWithExtensionsExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "port_security_enabled", "true"),
					testAccCheckNetworkingV2PortPortSecurityEnabled(&port, true),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_portSecurity_enabled(t *testing.T) {
	var port testPortWithExtensions

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortSecurityEnabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortWithExtensionsExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "port_security_enabled", "true"),
					testAccCheckNetworkingV2PortPortSecurityEnabled(&port, true),
				),
			},
			{
				Config: testAccNetworkingV2PortSecurityDisabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortWithExtensionsExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "port_security_enabled", "false"),
					testAccCheckNetworkingV2PortPortSecurityEnabled(&port, false),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_portBinding_create(t *testing.T) {
	var port ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortCreatePortBinding,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "extra_dhcp_option.#", "2"),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "binding.#", "1"),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "binding.0.vnic_type", "normal"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_portBinding_update(t *testing.T) {
	var port ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "admin_state_up", "true"),
					testAccCheckNetworkingV2PortAdminStateUp(&port, true),
				),
			},
			{
				Config: testAccNetworkingV2PortCreatePortBinding,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "extra_dhcp_option.#", "2"),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "binding.#", "1"),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "binding.0.vnic_type", "normal"),
				),
			},
			{
				Config: testAccNetworkingV2PortUpdatePortBinding1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "extra_dhcp_option.#", "1"),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "binding.#", "1"),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "binding.0.vnic_type", "normal"),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "binding.0.host_id", "localhost"),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "binding.0.profile", "{\"local_link_information\":[{\"port_id\":\"Ethernet3/4\",\"switch_id\":\"12:34:56:78:9A:BC\",\"switch_info\":\"info1\"},{\"port_id\":\"Ethernet3/4\",\"switch_id\":\"12:34:56:78:9A:BD\",\"switch_info\":\"info2\"}],\"vlan_type\":\"allowed\"}"),
				),
			},
			{
				Config: testAccNetworkingV2PortUpdatePortBinding2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "extra_dhcp_option.#", "0"),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "binding.#", "1"),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "binding.0.vnic_type", "baremetal"),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "binding.0.host_id", "localhost"),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "binding.0.profile", "{}"),
				),
			},
			{
				Config: testAccNetworkingV2PortUpdatePortBinding3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "extra_dhcp_option.#", "0"),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "binding.#", "1"),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "binding.0.vnic_type", "normal"),
				),
			},
			{
				Config: testAccNetworkingV2PortUpdatePortBinding4,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "extra_dhcp_option.#", "1"),
					// default computed values are in place
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "binding.#", "1"),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "binding.0.vnic_type", "normal"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_qos_policy_create(t *testing.T) {
	var (
		port      testPortWithExtensions
		qosPolicy policies.Policy
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortQosPolicy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortWithExtensionsExists(
						"openstack_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2QoSPolicyExists(
						"openstack_networking_qos_policy_v2.qos_policy_1", &qosPolicy),
					resource.TestCheckResourceAttrSet(
						"openstack_networking_port_v2.port_1", "qos_policy_id"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_qos_policy_update(t *testing.T) {
	var (
		port      testPortWithExtensions
		qosPolicy policies.Policy
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2PortBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortWithExtensionsExists("openstack_networking_port_v2.port_1", &port),
				),
			},
			{
				Config: testAccNetworkingV2PortQosPolicy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortWithExtensionsExists(
						"openstack_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2QoSPolicyExists(
						"openstack_networking_qos_policy_v2.qos_policy_1", &qosPolicy),
					resource.TestCheckResourceAttrSet(
						"openstack_networking_port_v2.port_1", "qos_policy_id"),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2PortDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.NetworkingV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_networking_port_v2" {
			continue
		}

		_, err := ports.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Port still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2PortExists(n string, port *ports.Port) resource.TestCheckFunc {
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

		found, err := ports.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Port not found")
		}

		*port = *found

		return nil
	}
}

func testAccCheckNetworkingV2PortWithExtensionsExists(
	n string, port *testPortWithExtensions) resource.TestCheckFunc {
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

		var p testPortWithExtensions
		err = ports.Get(networkingClient, rs.Primary.ID).ExtractInto(&p)
		if err != nil {
			return err
		}

		if p.ID != rs.Primary.ID {
			return fmt.Errorf("Port not found")
		}

		*port = p

		return nil
	}
}

func testAccCheckNetworkingV2PortCountFixedIPs(port *ports.Port, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(port.FixedIPs) != expected {
			return fmt.Errorf("Expected %d Fixed IPs, got %d", expected, len(port.FixedIPs))
		}

		return nil
	}
}

func testAccCheckNetworkingV2PortCountSecurityGroups(port *ports.Port, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(port.SecurityGroups) != expected {
			return fmt.Errorf("Expected %d Security Groups, got %d", expected, len(port.SecurityGroups))
		}

		return nil
	}
}

func testAccCheckNetworkingV2PortCountAllowedAddressPairs(
	port *ports.Port, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(port.AllowedAddressPairs) != expected {
			return fmt.Errorf("Expected %d Allowed Address Pairs, got %d", expected, len(port.AllowedAddressPairs))
		}

		return nil
	}
}

func testAccCheckNetworkingV2PortAdminStateUp(port *ports.Port, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if port.AdminStateUp != expected {
			return fmt.Errorf("Port has wrong admin_state_up. Expected %t, got %t", expected, port.AdminStateUp)
		}

		return nil
	}
}

func testAccCheckNetworkingV2PortPortSecurityEnabled(
	port *testPortWithExtensions, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if port.PortSecurityEnabled != expected {
			return fmt.Errorf("Port has wrong port_security_enabled. Expected %t, got %t", expected, port.PortSecurityEnabled)
		}

		return nil
	}
}

const testAccNetworkingV2PortBasic = `
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
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingV2PortNoIP = `
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
`

const testAccNetworkingV2PortNoIPEmptyUpdate = `
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
}
`

const testAccNetworkingV2PortMultipleNoIP = `
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
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
  }

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
  }

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
  }
}
`

const testAccNetworkingV2PortAllowedAddressPairs1 = `
resource "openstack_networking_network_v2" "vrrp_network" {
  name = "vrrp_network"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "vrrp_subnet" {
  name = "vrrp_subnet"
  cidr = "10.0.0.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.vrrp_network.id}"

  allocation_pools {
    start = "10.0.0.2"
    end = "10.0.0.200"
  }
}

resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group acceptance test"
}

resource "openstack_networking_router_v2" "vrrp_router" {
  name = "vrrp_router"
}

resource "openstack_networking_router_interface_v2" "vrrp_interface" {
  router_id = "${openstack_networking_router_v2.vrrp_router.id}"
  subnet_id = "${openstack_networking_subnet_v2.vrrp_subnet.id}"
}

resource "openstack_networking_port_v2" "vrrp_port_1" {
  name = "vrrp_port_1"
  description = "test vrrp port"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.vrrp_network.id}"

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.vrrp_subnet.id}"
    ip_address = "10.0.0.202"
  }
}

resource "openstack_networking_port_v2" "vrrp_port_2" {
  name = "vrrp_port_2"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.vrrp_network.id}"

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.vrrp_subnet.id}"
    ip_address = "10.0.0.201"
  }
}

resource "openstack_networking_port_v2" "instance_port" {
  name = "instance_port"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.vrrp_network.id}"

  allowed_address_pairs {
    ip_address = "${openstack_networking_port_v2.vrrp_port_1.fixed_ip.0.ip_address}"
    mac_address = "${openstack_networking_port_v2.vrrp_port_1.mac_address}"
  }

  allowed_address_pairs {
    ip_address = "${openstack_networking_port_v2.vrrp_port_2.fixed_ip.0.ip_address}"
    mac_address = "${openstack_networking_port_v2.vrrp_port_2.mac_address}"
  }
}
`

const testAccNetworkingV2PortAllowedAddressPairs2 = `
resource "openstack_networking_network_v2" "vrrp_network" {
  name = "vrrp_network"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "vrrp_subnet" {
  name = "vrrp_subnet"
  cidr = "10.0.0.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.vrrp_network.id}"

  allocation_pools {
    start = "10.0.0.2"
    end = "10.0.0.200"
  }
}

resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group acceptance test"
}

resource "openstack_networking_router_v2" "vrrp_router" {
  name = "vrrp_router"
}

resource "openstack_networking_router_interface_v2" "vrrp_interface" {
  router_id = "${openstack_networking_router_v2.vrrp_router.id}"
  subnet_id = "${openstack_networking_subnet_v2.vrrp_subnet.id}"
}

resource "openstack_networking_port_v2" "vrrp_port_1" {
  name = "vrrp_port_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.vrrp_network.id}"

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.vrrp_subnet.id}"
    ip_address = "10.0.0.202"
  }
}

resource "openstack_networking_port_v2" "vrrp_port_2" {
  name = "vrrp_port_2"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.vrrp_network.id}"

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.vrrp_subnet.id}"
    ip_address = "10.0.0.201"
  }
}

resource "openstack_networking_port_v2" "instance_port" {
  name = "instance_port"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.vrrp_network.id}"

  allowed_address_pairs {
    ip_address = "${openstack_networking_port_v2.vrrp_port_1.fixed_ip.0.ip_address}"
    mac_address = "${openstack_networking_port_v2.vrrp_port_1.mac_address}"
  }

  allowed_address_pairs {
    ip_address = "${openstack_networking_port_v2.vrrp_port_2.fixed_ip.0.ip_address}"
    mac_address = "${openstack_networking_port_v2.vrrp_port_2.mac_address}"
  }
}
`

const testAccNetworkingV2PortAllowedAddressPairs3 = `
resource "openstack_networking_network_v2" "vrrp_network" {
  name = "vrrp_network"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "vrrp_subnet" {
  name = "vrrp_subnet"
  cidr = "10.0.0.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.vrrp_network.id}"

  allocation_pools {
    start = "10.0.0.2"
    end = "10.0.0.200"
  }
}

resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group acceptance test"
}

resource "openstack_networking_router_v2" "vrrp_router" {
  name = "vrrp_router"
}

resource "openstack_networking_router_interface_v2" "vrrp_interface" {
  router_id = "${openstack_networking_router_v2.vrrp_router.id}"
  subnet_id = "${openstack_networking_subnet_v2.vrrp_subnet.id}"
}

resource "openstack_networking_port_v2" "vrrp_port_1" {
  name = "vrrp_port_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.vrrp_network.id}"

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.vrrp_subnet.id}"
    ip_address = "10.0.0.202"
  }
}

resource "openstack_networking_port_v2" "vrrp_port_2" {
  name = "vrrp_port_2"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.vrrp_network.id}"

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.vrrp_subnet.id}"
    ip_address = "10.0.0.201"
  }
}

resource "openstack_networking_port_v2" "instance_port" {
  name = "instance_port"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.vrrp_network.id}"
  security_group_ids = ["${openstack_networking_secgroup_v2.secgroup_1.id}"]

  allowed_address_pairs {
    ip_address = "${openstack_networking_port_v2.vrrp_port_1.fixed_ip.0.ip_address}"
    mac_address = "${openstack_networking_port_v2.vrrp_port_1.mac_address}"
  }

  allowed_address_pairs {
    ip_address = "${openstack_networking_port_v2.vrrp_port_2.fixed_ip.0.ip_address}"
    mac_address = "${openstack_networking_port_v2.vrrp_port_2.mac_address}"
  }
}
`

const testAccNetworkingV2PortAllowedAddressPairs4 = `
resource "openstack_networking_network_v2" "vrrp_network" {
  name = "vrrp_network"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "vrrp_subnet" {
  name = "vrrp_subnet"
  cidr = "10.0.0.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.vrrp_network.id}"

  allocation_pools {
    start = "10.0.0.2"
    end = "10.0.0.200"
  }
}

resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group acceptance test"
}

resource "openstack_networking_router_v2" "vrrp_router" {
  name = "vrrp_router"
}

resource "openstack_networking_router_interface_v2" "vrrp_interface" {
  router_id = "${openstack_networking_router_v2.vrrp_router.id}"
  subnet_id = "${openstack_networking_subnet_v2.vrrp_subnet.id}"
}

resource "openstack_networking_port_v2" "vrrp_port_1" {
  name = "vrrp_port_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.vrrp_network.id}"

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.vrrp_subnet.id}"
    ip_address = "10.0.0.202"
  }
}

resource "openstack_networking_port_v2" "vrrp_port_2" {
  name = "vrrp_port_2"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.vrrp_network.id}"

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.vrrp_subnet.id}"
    ip_address = "10.0.0.201"
  }
}

resource "openstack_networking_port_v2" "instance_port" {
  name = "instance_port"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.vrrp_network.id}"
  security_group_ids = ["${openstack_networking_secgroup_v2.secgroup_1.id}"]

  allowed_address_pairs {
    ip_address = "${openstack_networking_port_v2.vrrp_port_1.fixed_ip.0.ip_address}"
    mac_address = "${openstack_networking_port_v2.vrrp_port_1.mac_address}"
  }
}
`

const testAccNetworkingV2PortAllowedAddressPairs5 = `
resource "openstack_networking_network_v2" "vrrp_network" {
  name = "vrrp_network"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "vrrp_subnet" {
  name = "vrrp_subnet"
  cidr = "10.0.0.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.vrrp_network.id}"

  allocation_pools {
    start = "10.0.0.2"
    end = "10.0.0.200"
  }
}

resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group acceptance test"
}

resource "openstack_networking_router_v2" "vrrp_router" {
  name = "vrrp_router"
}

resource "openstack_networking_router_interface_v2" "vrrp_interface" {
  router_id = "${openstack_networking_router_v2.vrrp_router.id}"
  subnet_id = "${openstack_networking_subnet_v2.vrrp_subnet.id}"
}

resource "openstack_networking_port_v2" "vrrp_port_1" {
  name = "vrrp_port_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.vrrp_network.id}"

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.vrrp_subnet.id}"
    ip_address = "10.0.0.202"
  }
}

resource "openstack_networking_port_v2" "vrrp_port_2" {
  name = "vrrp_port_2"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.vrrp_network.id}"

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.vrrp_subnet.id}"
    ip_address = "10.0.0.201"
  }
}

resource "openstack_networking_port_v2" "instance_port" {
  name = "instance_port"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.vrrp_network.id}"
}
`

const testAccNetworkingV2PortMultipleFixedIPs = `
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
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.20"
  }

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.40"
  }
}
`

const testAccNetworkingV2PortTimeout = `
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
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`

const testAccNetworkingV2PortFixedIPs = `
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
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.24"
  }

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingV2PortUpdateSecurityGroups1 = `
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

resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group acceptance test"
}

resource "openstack_networking_secgroup_v2" "secgroup_2" {
  name = "secgroup_2"
  description = "terraform security group acceptance test"
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingV2PortUpdateSecurityGroups2 = `
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

resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group acceptance test"
}

resource "openstack_networking_secgroup_v2" "secgroup_2" {
  name = "secgroup_2"
  description = "terraform security group acceptance test"
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  security_group_ids = ["${openstack_networking_secgroup_v2.secgroup_1.id}"]

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingV2PortUpdateSecurityGroups3 = `
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

resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name = "security_group_1"
  description = "terraform security group acceptance test"
}

resource "openstack_networking_secgroup_v2" "secgroup_2" {
  name = "secgroup_2"
  description = "terraform security group acceptance test"
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  security_group_ids = [
    "${openstack_networking_secgroup_v2.secgroup_1.id}",
    "${openstack_networking_secgroup_v2.secgroup_2.id}"
  ]

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingV2PortUpdateSecurityGroups4 = `
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

resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name = "security_group"
  description = "terraform security group acceptance test"
}

resource "openstack_networking_secgroup_v2" "secgroup_2" {
  name = "secgroup_2"
  description = "terraform security group acceptance test"
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  security_group_ids = ["${openstack_networking_secgroup_v2.secgroup_2.id}"]

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingV2PortUpdateSecurityGroups5 = `
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

resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name = "security_group"
  description = "terraform security group acceptance test"
}

resource "openstack_networking_secgroup_v2" "secgroup_2" {
  name = "secgroup_2"
  description = "terraform security group acceptance test"
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  security_group_ids = []

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingV2PortNoSecurityGroups1 = `
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

resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group acceptance test"
}

resource "openstack_networking_secgroup_v2" "secgroup_2" {
  name = "secgroup_2"
  description = "terraform security group acceptance test"
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  no_security_groups = true

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingV2PortNoSecurityGroups2 = `
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

resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group acceptance test"
}

resource "openstack_networking_secgroup_v2" "secgroup_2" {
  name = "secgroup_2"
  description = "terraform security group acceptance test"
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  no_security_groups = false
  security_group_ids = ["${openstack_networking_secgroup_v2.secgroup_1.id}"]

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingV2PortNoSecurityGroups3 = `
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

resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group acceptance test"
}

resource "openstack_networking_secgroup_v2" "secgroup_2" {
  name = "secgroup_2"
  description = "terraform security group acceptance test"
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  no_security_groups = false
  security_group_ids = [
    "${openstack_networking_secgroup_v2.secgroup_1.id}",
    "${openstack_networking_secgroup_v2.secgroup_2.id}"
  ]

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingV2PortNoSecurityGroups4 = `
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

resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group acceptance test"
}

resource "openstack_networking_secgroup_v2" "secgroup_2" {
  name = "secgroup_2"
  description = "terraform security group acceptance test"
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  no_security_groups = true

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingV2PortAllowedAddressPairsNoMAC = `
resource "openstack_networking_network_v2" "vrrp_network" {
  name = "vrrp_network"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "vrrp_subnet" {
  name = "vrrp_subnet"
  cidr = "10.0.0.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.vrrp_network.id}"

  allocation_pools {
    start = "10.0.0.2"
    end = "10.0.0.200"
  }
}

resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group acceptance test"
}

resource "openstack_networking_router_v2" "vrrp_router" {
  name = "vrrp_router"
}

resource "openstack_networking_router_interface_v2" "vrrp_interface" {
  router_id = "${openstack_networking_router_v2.vrrp_router.id}"
  subnet_id = "${openstack_networking_subnet_v2.vrrp_subnet.id}"
}

resource "openstack_networking_port_v2" "vrrp_port_1" {
  name = "vrrp_port_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.vrrp_network.id}"

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.vrrp_subnet.id}"
    ip_address = "10.0.0.202"
  }
}

resource "openstack_networking_port_v2" "vrrp_port_2" {
  name = "vrrp_port_2"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.vrrp_network.id}"

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.vrrp_subnet.id}"
    ip_address = "10.0.0.201"
  }
}

resource "openstack_networking_port_v2" "instance_port" {
  name = "instance_port"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.vrrp_network.id}"

  allowed_address_pairs {
    ip_address = "${openstack_networking_port_v2.vrrp_port_1.fixed_ip.0.ip_address}"
  }

  allowed_address_pairs {
    ip_address = "${openstack_networking_port_v2.vrrp_port_2.fixed_ip.0.ip_address}"
  }
}
`

const testAccNetworkingV2PortNoFixedIP1 = `
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
  no_fixed_ip = true
}
`

const testAccNetworkingV2PortNoFixedIP2 = `
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
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingV2PortNoFixedIP3 = `
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
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.24"
  }
}
`

const testAccNetworkingV2PortCreateExtraDhcpOpts = `
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
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }

  extra_dhcp_option {
    name = "optionA"
    value = "valueA"
  }

  extra_dhcp_option {
    name = "optionB"
    value = "valueB"
  }
}
`

const testAccNetworkingV2PortUpdateExtraDhcpOpts1 = `
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
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }

  extra_dhcp_option {
    name = "optionC"
    value = "valueC"
  }
}
`

const testAccNetworkingV2PortUpdateExtraDhcpOpts2 = `
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
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }

  extra_dhcp_option {
    name = "optionC"
    value = "valueC"
  }

  extra_dhcp_option {
    name = "optionD"
    value = "valueD"
  }
}
`

const testAccNetworkingV2PortUpdateExtraDhcpOpts3 = `
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
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }

  extra_dhcp_option {
    name = "optionD"
    value = "valueD"
  }

  extra_dhcp_option {
    name = "optionE"
    value = "valueE"
  }
}
`

const testAccNetworkingV2PortUpdateExtraDhcpOpts4 = `
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
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }

  extra_dhcp_option {
    name = "optionD"
    value = "valueD"
  }

  extra_dhcp_option {
    name = "optionE"
    value = "valueEE"
  }
}
`

const testAccNetworkingV2PortUpdateExtraDhcpOpts5 = `
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
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }

  extra_dhcp_option {
    name = "optionD"
    value = "valueDD"
  }

  extra_dhcp_option {
    name = "optionE"
    value = "valueEE"
  }
}
`

const testAccNetworkingV2PortUpdateExtraDhcpOpts6 = `
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
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingV2PortAdminStateUpOmit = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  network_id = "${openstack_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingV2PortAdminStateUpTrue = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
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
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingV2PortAdminStateUpFalse = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "false"
  network_id = "${openstack_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingV2PortSecurityOmit = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  no_security_groups = true
  network_id = "${openstack_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingV2PortSecurityDisabled = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  no_security_groups = true
  port_security_enabled = false

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingV2PortSecurityEnabled = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  no_security_groups = true
  port_security_enabled = true

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }
}
`

const testAccNetworkingV2PortCreatePortBinding = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "false"
  network_id = "${openstack_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }

  binding {
    vnic_type = "normal"
  }

  extra_dhcp_option {
    name = "optionA"
    value = "valueA"
  }

  extra_dhcp_option {
    name = "optionB"
    value = "valueB"
  }
}
`

const testAccNetworkingV2PortUpdatePortBinding1 = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "false"
  network_id = "${openstack_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }

  binding {
    host_id = "localhost"
    profile = <<EOF
{
  "local_link_information": [
    {
      "switch_info": "info1",
      "port_id": "Ethernet3/4",
      "switch_id": "12:34:56:78:9A:BC"
    },
    {
      "switch_info": "info2",
      "port_id": "Ethernet3/4",
      "switch_id": "12:34:56:78:9A:BD"
    }
  ],
  "vlan_type": "allowed"
}
EOF
  }

  extra_dhcp_option {
    name = "optionA"
    value = "valueA"
  }
}
`

const testAccNetworkingV2PortUpdatePortBinding2 = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "false"
  network_id = "${openstack_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }

  binding {
    host_id = "localhost"
    vnic_type = "baremetal"
  }
}
`

const testAccNetworkingV2PortUpdatePortBinding3 = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "false"
  network_id = "${openstack_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }

  binding {
	vnic_type = "normal"
  }
}
`

const testAccNetworkingV2PortUpdatePortBinding4 = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "false"
  network_id = "${openstack_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }

  extra_dhcp_option {
    name = "optionA"
    value = "valueA"
  }
}
`

const testAccNetworkingV2PortQosPolicy = `
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

resource "openstack_networking_qos_policy_v2" "qos_policy_1" {
  name = "qos_policy_1"
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address = "192.168.199.23"
  }

  qos_policy_id  = "${openstack_networking_qos_policy_v2.qos_policy_1.id}"
}
`
