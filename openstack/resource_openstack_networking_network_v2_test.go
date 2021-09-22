package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/secgroups"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/portsecurity"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/qos/policies"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/subnets"
)

type testNetworkWithExtensions struct {
	networks.Network
	portsecurity.PortSecurityExt
	policies.QoSPolicyExt
}

func TestAccNetworkingV2Network_basic(t *testing.T) {
	var network networks.Network

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2NetworkBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					resource.TestCheckResourceAttr(
						"openstack_networking_network_v2.network_1", "name", "network_1"),
					resource.TestCheckResourceAttr(
						"openstack_networking_network_v2.network_1", "description", "my network description"),
				),
			},
			{
				Config: testAccNetworkingV2NetworkUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_networking_network_v2.network_1", "name", "network_2"),
					resource.TestCheckResourceAttr(
						"openstack_networking_network_v2.network_1", "description", ""),
				),
			},
		},
	})
}

func TestAccNetworkingV2Network_netstack(t *testing.T) {
	var network networks.Network
	var subnet subnets.Subnet
	var router routers.Router

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2NetworkNetstack,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckNetworkingV2RouterExists("openstack_networking_router_v2.router_1", &router),
					testAccCheckNetworkingV2RouterInterfaceExists(
						"openstack_networking_router_interface_v2.ri_1"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Network_fullstack(t *testing.T) {
	var instance servers.Server
	var network networks.Network
	var port ports.Port
	var secgroup secgroups.SecurityGroup
	var subnet subnets.Subnet

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2NetworkFullstack,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2SubnetExists("openstack_networking_subnet_v2.subnet_1", &subnet),
					testAccCheckComputeV2SecGroupExists("openstack_compute_secgroup_v2.secgroup_1", &secgroup),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.port_1", &port),
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
				),
			},
		},
	})
}

func TestAccNetworkingV2Network_timeout(t *testing.T) {
	var network networks.Network

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2NetworkTimeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
				),
			},
		},
	})
}

func TestAccNetworkingV2Network_multipleSegmentMappings(t *testing.T) {
	var network networks.Network

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2NetworkMultipleSegmentMappings,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
				),
			},
		},
	})
}

func TestAccNetworkingV2Network_externalCreate(t *testing.T) {
	var network networks.Network

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2NetworkExternal,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					resource.TestCheckResourceAttr(
						"openstack_networking_network_v2.network_1", "external", "true"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Network_externalUpdate(t *testing.T) {
	var network networks.Network

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2NetworkBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
				),
			},
			{
				Config: testAccNetworkingV2NetworkExternal,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					resource.TestCheckResourceAttr(
						"openstack_networking_network_v2.network_1", "external", "true"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Network_transparent_vlan_Create(t *testing.T) {
	var network networks.Network

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckTransparentVLAN(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2NetworkTransparentVlan,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					resource.TestCheckResourceAttr(
						"openstack_networking_network_v2.network_1", "transparent_vlan", "true"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Network_adminStateUp_omit(t *testing.T) {
	var network networks.Network

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2NetworkAdminStateUpOmit,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					resource.TestCheckResourceAttr(
						"openstack_networking_network_v2.network_1", "admin_state_up", "true"),
					testAccCheckNetworkingV2NetworkAdminStateUp(&network, true),
				),
			},
		},
	})
}

func TestAccNetworkingV2Network_adminStateUp_true(t *testing.T) {
	var network networks.Network

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2NetworkAdminStateUpTrue,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					resource.TestCheckResourceAttr(
						"openstack_networking_network_v2.network_1", "admin_state_up", "true"),
					testAccCheckNetworkingV2NetworkAdminStateUp(&network, true),
				),
			},
		},
	})
}

func TestAccNetworkingV2Network_adminStateUp_false(t *testing.T) {
	var network networks.Network

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2NetworkAdminStateUpFalse,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					resource.TestCheckResourceAttr(
						"openstack_networking_network_v2.network_1", "admin_state_up", "false"),
					testAccCheckNetworkingV2NetworkAdminStateUp(&network, false),
				),
			},
		},
	})
}

func TestAccNetworkingV2Network_adminStateUp_update(t *testing.T) {
	var network networks.Network

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2NetworkAdminStateUpOmit,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					resource.TestCheckResourceAttr(
						"openstack_networking_network_v2.network_1", "admin_state_up", "true"),
					testAccCheckNetworkingV2NetworkAdminStateUp(&network, true),
				),
			},
			{
				Config: testAccNetworkingV2NetworkAdminStateUpFalse,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists("openstack_networking_network_v2.network_1", &network),
					resource.TestCheckResourceAttr(
						"openstack_networking_network_v2.network_1", "admin_state_up", "false"),
					testAccCheckNetworkingV2NetworkAdminStateUp(&network, false),
				),
			},
		},
	})
}

func TestAccNetworkingV2Network_portSecurity_omit(t *testing.T) {
	var network testNetworkWithExtensions

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2NetworkAdminStateUpOmit,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkWithExtensionsExists(
						"openstack_networking_network_v2.network_1", &network),
					resource.TestCheckResourceAttr(
						"openstack_networking_network_v2.network_1", "port_security_enabled", "true"),
					testAccCheckNetworkingV2NetworkPortSecurityEnabled(&network, true),
				),
			},
			{
				Config: testAccNetworkingV2NetworkPortSecurityDisabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkWithExtensionsExists(
						"openstack_networking_network_v2.network_1", &network),
					resource.TestCheckResourceAttr(
						"openstack_networking_network_v2.network_1", "port_security_enabled", "false"),
					testAccCheckNetworkingV2NetworkPortSecurityEnabled(&network, false),
				),
			},
			{
				Config: testAccNetworkingV2NetworkPortSecurityEnabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkWithExtensionsExists(
						"openstack_networking_network_v2.network_1", &network),
					resource.TestCheckResourceAttr(
						"openstack_networking_network_v2.network_1", "port_security_enabled", "true"),
					testAccCheckNetworkingV2NetworkPortSecurityEnabled(&network, true),
				),
			},
		},
	})
}

func TestAccNetworkingV2Network_portSecurity_disabled(t *testing.T) {
	var network testNetworkWithExtensions

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2NetworkPortSecurityDisabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkWithExtensionsExists(
						"openstack_networking_network_v2.network_1", &network),
					resource.TestCheckResourceAttr(
						"openstack_networking_network_v2.network_1", "port_security_enabled", "false"),
					testAccCheckNetworkingV2NetworkPortSecurityEnabled(&network, false),
				),
			},
			{
				Config: testAccNetworkingV2NetworkPortSecurityEnabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkWithExtensionsExists(
						"openstack_networking_network_v2.network_1", &network),
					resource.TestCheckResourceAttr(
						"openstack_networking_network_v2.network_1", "port_security_enabled", "true"),
					testAccCheckNetworkingV2NetworkPortSecurityEnabled(&network, true),
				),
			},
		},
	})
}

func TestAccNetworkingV2Network_portSecurity_enabled(t *testing.T) {
	var network testNetworkWithExtensions

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2NetworkPortSecurityEnabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkWithExtensionsExists(
						"openstack_networking_network_v2.network_1", &network),
					resource.TestCheckResourceAttr(
						"openstack_networking_network_v2.network_1", "port_security_enabled", "true"),
					testAccCheckNetworkingV2NetworkPortSecurityEnabled(&network, true),
				),
			},
			{
				Config: testAccNetworkingV2NetworkPortSecurityDisabled,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkWithExtensionsExists(
						"openstack_networking_network_v2.network_1", &network),
					resource.TestCheckResourceAttr(
						"openstack_networking_network_v2.network_1", "port_security_enabled", "false"),
					testAccCheckNetworkingV2NetworkPortSecurityEnabled(&network, false),
				),
			},
		},
	})
}

func TestAccNetworkingV2Network_qos_policy_create(t *testing.T) {
	var (
		network   testNetworkWithExtensions
		qosPolicy policies.Policy
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2NetworkQosPolicy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkWithExtensionsExists(
						"openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2QoSPolicyExists(
						"openstack_networking_qos_policy_v2.qos_policy_1", &qosPolicy),
					resource.TestCheckResourceAttrSet(
						"openstack_networking_network_v2.network_1", "qos_policy_id"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Network_qos_policy_update(t *testing.T) {
	var (
		network   testNetworkWithExtensions
		qosPolicy policies.Policy
	)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2NetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2NetworkBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkWithExtensionsExists(
						"openstack_networking_network_v2.network_1", &network),
					resource.TestCheckResourceAttr(
						"openstack_networking_network_v2.network_1", "name", "network_1"),
					resource.TestCheckResourceAttr(
						"openstack_networking_network_v2.network_1", "description", "my network description"),
				),
			},
			{
				Config: testAccNetworkingV2NetworkQosPolicy,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkWithExtensionsExists(
						"openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingV2QoSPolicyExists(
						"openstack_networking_qos_policy_v2.qos_policy_1", &qosPolicy),
					resource.TestCheckResourceAttrSet(
						"openstack_networking_network_v2.network_1", "qos_policy_id"),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2NetworkDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.NetworkingV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_networking_network_v2" {
			continue
		}

		_, err := networks.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Network still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2NetworkExists(n string, network *networks.Network) resource.TestCheckFunc {
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

		found, err := networks.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Network not found")
		}

		*network = *found

		return nil
	}
}

func testAccCheckNetworkingV2NetworkWithExtensionsExists(n string, network *testNetworkWithExtensions) resource.TestCheckFunc {
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

		var n testNetworkWithExtensions
		err = networks.Get(networkingClient, rs.Primary.ID).ExtractInto(&n)
		if err != nil {
			return err
		}

		if n.ID != rs.Primary.ID {
			return fmt.Errorf("Network not found")
		}

		*network = n

		return nil
	}
}

func testAccCheckNetworkingV2NetworkAdminStateUp(network *networks.Network, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if network.AdminStateUp != expected {
			return fmt.Errorf("Network has wrong admin_state_up. Expected %t, got %t", expected, network.AdminStateUp)
		}

		return nil
	}
}

func testAccCheckNetworkingV2NetworkPortSecurityEnabled(network *testNetworkWithExtensions, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if network.PortSecurityEnabled != expected {
			return fmt.Errorf("Network has wrong port_security_enabled. Expected %t, got %t", expected, network.PortSecurityEnabled)
		}

		return nil
	}
}

const testAccNetworkingV2NetworkBasic = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  description = "my network description"
  admin_state_up = "true"
}
`

const testAccNetworkingV2NetworkUpdate = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_2"
  admin_state_up = "true"
}
`

const testAccNetworkingV2NetworkNetstack = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.10.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_router_v2" "router_1" {
  name = "router_1"
}

resource "openstack_networking_router_interface_v2" "ri_1" {
  router_id = "${openstack_networking_router_v2.router_1.id}"
  subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
}
`

const testAccNetworkingV2NetworkFullstack = `
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

resource "openstack_compute_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "a security group"
  rule {
    from_port = 22
    to_port = 22
    ip_protocol = "tcp"
    cidr = "0.0.0.0/0"
  }
}

resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  admin_state_up = "true"
  security_group_ids = ["${openstack_compute_secgroup_v2.secgroup_1.id}"]
  network_id = "${openstack_networking_network_v2.network_1.id}"

  fixed_ip {
    subnet_id =  "${openstack_networking_subnet_v2.subnet_1.id}"
    ip_address =  "192.168.199.23"
  }
}

resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["${openstack_compute_secgroup_v2.secgroup_1.name}"]

  network {
    port = "${openstack_networking_port_v2.port_1.id}"
  }
}
`

const testAccNetworkingV2NetworkTimeout = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`

const testAccNetworkingV2NetworkMultipleSegmentMappings = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"

  segments {
    segmentation_id = 2
    network_type = "vxlan"
  }
}
`

const testAccNetworkingV2NetworkExternal = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
	admin_state_up = "true"
	external = "true"
}
`

const testAccNetworkingV2NetworkTransparentVlan = `
resource "openstack_networking_network_v2" "network_1" {
  name             = "network_1"
	admin_state_up   = "true"
	transparent_vlan = "true"
}
`

const testAccNetworkingV2NetworkAdminStateUpOmit = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
}
`

const testAccNetworkingV2NetworkAdminStateUpTrue = `
resource "openstack_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}
`

const testAccNetworkingV2NetworkAdminStateUpFalse = `
resource "openstack_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "false"
}
`

const testAccNetworkingV2NetworkPortSecurityDisabled = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  port_security_enabled = "false"
}
`

const testAccNetworkingV2NetworkPortSecurityEnabled = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  port_security_enabled = "true"
}
`

const testAccNetworkingV2NetworkQosPolicy = `
resource "openstack_networking_qos_policy_v2" "qos_policy_1" {
  name = "qos_policy_1"
}

resource "openstack_networking_network_v2" "network_1" {
  name           = "network_1"
  description    = "my network description"
  admin_state_up = "true"
  qos_policy_id  = "${openstack_networking_qos_policy_v2.qos_policy_1.id}"
}
`
