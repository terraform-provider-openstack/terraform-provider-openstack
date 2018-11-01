package openstack

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccNetworkingV2Port_basic(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2Port_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_noip(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2Port_noip,
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
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2Port_multipleNoIP,
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
	var vrrp_port_1, vrrp_port_2, instance_port ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2Port_allowedAddressPairs_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.vrrp_subnet", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.vrrp_network", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.vrrp_port_1", &vrrp_port_1),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.vrrp_port_2", &vrrp_port_2),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.instance_port", &instance_port),
					testAccCheckNetworkingV2PortCountAllowedAddressPairs(&instance_port, 2),
				),
			},
			resource.TestStep{
				Config: testAccNetworkingV2Port_allowedAddressPairs_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.vrrp_subnet", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.vrrp_network", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.vrrp_port_1", &vrrp_port_1),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.vrrp_port_2", &vrrp_port_2),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.instance_port", &instance_port),
					testAccCheckNetworkingV2PortCountAllowedAddressPairs(&instance_port, 2),
				),
			},
			resource.TestStep{
				Config: testAccNetworkingV2Port_allowedAddressPairs_3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.vrrp_subnet", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.vrrp_network", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.vrrp_port_1", &vrrp_port_1),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.vrrp_port_2", &vrrp_port_2),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.instance_port", &instance_port),
					testAccCheckNetworkingV2PortCountAllowedAddressPairs(&instance_port, 2),
				),
			},
			resource.TestStep{
				Config: testAccNetworkingV2Port_allowedAddressPairs_4,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.vrrp_subnet", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.vrrp_network", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.vrrp_port_1", &vrrp_port_1),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.vrrp_port_2", &vrrp_port_2),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.instance_port", &instance_port),
					testAccCheckNetworkingV2PortCountAllowedAddressPairs(&instance_port, 1),
				),
			},
			resource.TestStep{
				Config: testAccNetworkingV2Port_allowedAddressPairs_5,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.vrrp_subnet", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.vrrp_network", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.vrrp_port_1", &vrrp_port_1),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.vrrp_port_2", &vrrp_port_2),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.instance_port", &instance_port),
					testAccCheckNetworkingV2PortCountAllowedAddressPairs(&instance_port, 0),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_allowedAddressPairsNoMAC(t *testing.T) {
	var network networks.Network
	var subnet subnets.Subnet
	var vrrp_port_1, vrrp_port_2, instance_port ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2Port_allowedAddressPairsNoMAC,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.vrrp_subnet", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.vrrp_network", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.vrrp_port_1", &vrrp_port_1),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.vrrp_port_2", &vrrp_port_2),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.instance_port", &instance_port),
					testAccCheckNetworkingV2PortCountAllowedAddressPairs(&instance_port, 2),
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
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2Port_multipleFixedIPs,
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
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2Port_timeout,
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
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2Port_fixedIPs,
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
	var secgroup_1, secgroup_2 groups.SecGroup
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2Port_updateSecurityGroups_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_1", &secgroup_1),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_2", &secgroup_2),
					testAccCheckNetworkingV2PortCountSecurityGroups(&port, 1),
				),
			},
			resource.TestStep{
				Config: testAccNetworkingV2Port_updateSecurityGroups_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_1", &secgroup_1),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_2", &secgroup_2),
					testAccCheckNetworkingV2PortCountSecurityGroups(&port, 1),
				),
			},
			resource.TestStep{
				Config: testAccNetworkingV2Port_updateSecurityGroups_3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_1", &secgroup_1),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_2", &secgroup_2),
					testAccCheckNetworkingV2PortCountSecurityGroups(&port, 2),
				),
			},
			resource.TestStep{
				Config: testAccNetworkingV2Port_updateSecurityGroups_4,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_1", &secgroup_1),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_2", &secgroup_2),
					testAccCheckNetworkingV2PortCountSecurityGroups(&port, 1),
				),
			},
			resource.TestStep{
				Config: testAccNetworkingV2Port_updateSecurityGroups_5,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_1", &secgroup_1),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_2", &secgroup_2),
					testAccCheckNetworkingV2PortCountSecurityGroups(&port, 0),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_noSecurityGroups(t *testing.T) {
	var network networks.Network
	var port ports.Port
	var secgroup_1, secgroup_2 groups.SecGroup
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2Port_noSecurityGroups_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_1", &secgroup_1),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_2", &secgroup_2),
					testAccCheckNetworkingV2PortCountSecurityGroups(&port, 0),
				),
			},
			resource.TestStep{
				Config: testAccNetworkingV2Port_noSecurityGroups_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_1", &secgroup_1),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_2", &secgroup_2),
					testAccCheckNetworkingV2PortCountSecurityGroups(&port, 1),
				),
			},
			resource.TestStep{
				Config: testAccNetworkingV2Port_noSecurityGroups_3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_1", &secgroup_1),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_2", &secgroup_2),
					testAccCheckNetworkingV2PortCountSecurityGroups(&port, 2),
				),
			},
			resource.TestStep{
				Config: testAccNetworkingV2Port_noSecurityGroups_4,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_1", &secgroup_1),
					testAccCheckNetworkingV2SecGroupExists(
						"openstack_networking_secgroup_v2.secgroup_2", &secgroup_2),
					testAccCheckNetworkingV2PortCountSecurityGroups(&port, 0),
				),
			},
		},
	})
}

func TestAccNetworkingV2Port_noFixedIP(t *testing.T) {
	var port ports.Port

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2Port_noFixedIP_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "all_fixed_ips.#", "0"),
				),
			},
			resource.TestStep{
				Config: testAccNetworkingV2Port_noFixedIP_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "all_fixed_ips.#", "1"),
				),
			},
			resource.TestStep{
				Config: testAccNetworkingV2Port_noFixedIP_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "all_fixed_ips.#", "0"),
				),
			},
			resource.TestStep{
				Config: testAccNetworkingV2Port_noFixedIP_3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "all_fixed_ips.#", "2"),
				),
			},
			resource.TestStep{
				Config: testAccNetworkingV2Port_noFixedIP_1,
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
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2Port_createExtraDHCPOpts,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "extra_dhcp_opts.0.opt_name", "optionA"),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "extra_dhcp_opts.0.opt_value", "valueA"),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "extra_dhcp_opts.1.opt_name", "optionB"),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "extra_dhcp_opts.1.opt_value", "valueB"),
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
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2PortDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkingV2Port_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
				),
			},
			resource.TestStep{
				Config: testAccNetworkingV2Port_updateExtraDHCPOpts_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "extra_dhcp_opts.0.opt_name", "optionC"),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "extra_dhcp_opts.0.opt_value", "valueC"),
				),
			},
			resource.TestStep{
				Config: testAccNetworkingV2Port_updateExtraDHCPOpts_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "extra_dhcp_opts.0.opt_name", "optionC"),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "extra_dhcp_opts.0.opt_value", "valueC"),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "extra_dhcp_opts.1.opt_name", "optionD"),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "extra_dhcp_opts.1.opt_value", "valueD"),
				),
			},
			resource.TestStep{
				Config: testAccNetworkingV2Port_updateExtraDHCPOpts_3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "extra_dhcp_opts.0.opt_name", "optionD"),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "extra_dhcp_opts.0.opt_value", "valueD"),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "extra_dhcp_opts.1.opt_name", "optionE"),
					resource.TestCheckResourceAttr(
						"openstack_networking_port_v2.port_1", "extra_dhcp_opts.1.opt_value", "valueE"),
				),
			},
			resource.TestStep{
				Config: testAccNetworkingV2Port_updateExtraDHCPOpts_4,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					resource.TestCheckNoResourceAttr(
						"openstack_networking_port_v2.port_1", "extra_dhcp_opts.0.opt_name"),
					resource.TestCheckNoResourceAttr(
						"openstack_networking_port_v2.port_1", "extra_dhcp_opts.0.opt_value"),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2PortDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
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
		networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
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

const testAccNetworkingV2Port_basic = `
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

const testAccNetworkingV2Port_noip = `
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
}
`

const testAccNetworkingV2Port_multipleNoIP = `
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

const testAccNetworkingV2Port_allowedAddressPairs_1 = `
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

const testAccNetworkingV2Port_allowedAddressPairs_2 = `
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

const testAccNetworkingV2Port_allowedAddressPairs_3 = `
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

const testAccNetworkingV2Port_allowedAddressPairs_4 = `
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

const testAccNetworkingV2Port_allowedAddressPairs_5 = `
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

const testAccNetworkingV2Port_multipleFixedIPs = `
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

const testAccNetworkingV2Port_timeout = `
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

const testAccNetworkingV2Port_fixedIPs = `
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

const testAccNetworkingV2Port_updateSecurityGroups_1 = `
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

const testAccNetworkingV2Port_updateSecurityGroups_2 = `
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

const testAccNetworkingV2Port_updateSecurityGroups_3 = `
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

const testAccNetworkingV2Port_updateSecurityGroups_4 = `
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

const testAccNetworkingV2Port_updateSecurityGroups_5 = `
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

const testAccNetworkingV2Port_noSecurityGroups_1 = `
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

const testAccNetworkingV2Port_noSecurityGroups_2 = `
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

const testAccNetworkingV2Port_noSecurityGroups_3 = `
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

const testAccNetworkingV2Port_noSecurityGroups_4 = `
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

const testAccNetworkingV2Port_allowedAddressPairsNoMAC = `
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

const testAccNetworkingV2Port_noFixedIP_1 = `
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

const testAccNetworkingV2Port_noFixedIP_2 = `
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

const testAccNetworkingV2Port_noFixedIP_3 = `
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

const testAccNetworkingV2Port_createExtraDHCPOpts = `
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
  extra_dhcp_opts {
    opt_name = "optionA"
    opt_value = "valueA"
  }
  extra_dhcp_opts {
    opt_name = "optionB"
    opt_value = "valueB"
  }
}
`

const testAccNetworkingV2Port_updateExtraDHCPOpts_1 = `
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
  extra_dhcp_opts {
    opt_name = "optionC"
    opt_value = "valueC"
  }
}
`

const testAccNetworkingV2Port_updateExtraDHCPOpts_2 = `
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
  extra_dhcp_opts {
    opt_name = "optionC"
    opt_value = "valueC"
  }
  extra_dhcp_opts {
    opt_name = "optionD"
    opt_value = "valueD"
  }
}
`

const testAccNetworkingV2Port_updateExtraDHCPOpts_3 = `
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
  extra_dhcp_opts {
    opt_name = "optionD"
    opt_value = "valueD"
  }
  extra_dhcp_opts {
    opt_name = "optionE"
    opt_value = "valueE"
  }
}
`

const testAccNetworkingV2Port_updateExtraDHCPOpts_4 = `
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
