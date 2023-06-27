package openstack

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/groups"
)

func TestAccFWGroupV2_basic(t *testing.T) {
	var policyID *string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckFW(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckFWGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWGroupV2Basic1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWGroupV2("openstack_fw_group_v2.group_1", "", "", policyID),
				),
			},
			{
				Config: testAccFWGroupV2Basic2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWGroupV2(
						"openstack_fw_group_v2.group_1", "group_1", "terraform acceptance test", policyID),
				),
			},
			{
				Config: testAccFWGroupV2Basic3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWGroupV2(
						"openstack_fw_group_v2.group_1", "new_name_group_1", "new description terraform acceptance test", policyID),
				),
			},
			{
				Config: testAccFWGroupV2Basic1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWGroupV2(
						"openstack_fw_group_v2.group_1", "", "", policyID),
				),
			},
		},
	})
}

func TestAccFWGroupV2_shared(t *testing.T) {
	var policyID *string

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
			testAccPreCheckFW(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckFWGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWGroupV2Shared,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWGroupV2("openstack_fw_group_v2.group_1", "", "", policyID),
					resource.TestCheckResourceAttr(
						"openstack_fw_group_v2.group_1", "shared", "true"),
				),
			},
		},
	})
}

func TestAccFWGroupV2_port(t *testing.T) {
	var group groups.Group

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckFW(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckFWGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWGroupV2Port,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWGroupV2Exists("openstack_fw_group_v2.group_1", &group),
					testAccCheckFWGroupPortCount(&group, 1),
				),
			},
		},
	})
}

func TestAccFWGroupV2_no_port(t *testing.T) {
	var group groups.Group

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckFW(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckFWGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWGroupV2NoPort,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWGroupV2Exists("openstack_fw_group_v2.group_1", &group),
					resource.TestCheckResourceAttr("openstack_fw_group_v2.group_1", "description", "firewall group port test"),
					testAccCheckFWGroupPortCount(&group, 0),
				),
			},
		},
	})
}

func TestAccFWGroupV2_port_update(t *testing.T) {
	var group groups.Group

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckFW(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckFWGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWGroupV2Port,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWGroupV2Exists("openstack_fw_group_v2.group_1", &group),
					testAccCheckFWGroupPortCount(&group, 1),
				),
			},
			{
				Config: testAccFWGroupV2PortAdd,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWGroupV2Exists("openstack_fw_group_v2.group_1", &group),
					testAccCheckFWGroupPortCount(&group, 2),
				),
			},
		},
	})
}

func TestAccFWGroupV2_port_remove(t *testing.T) {
	var group groups.Group

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckFW(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckFWGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccFWGroupV2Port,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWGroupV2Exists("openstack_fw_group_v2.group_1", &group),
					testAccCheckFWGroupPortCount(&group, 1),
				),
			},
			{
				Config: testAccFWGroupV2PortRemove,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckFWGroupV2Exists("openstack_fw_group_v2.group_1", &group),
					testAccCheckFWGroupPortCount(&group, 0),
				),
			},
		},
	})
}

func testAccCheckFWGroupV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.NetworkingV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_fw_group_v2" {
			continue
		}

		_, err = groups.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Firewall group (%s) still exists", rs.Primary.ID)
		}
		if _, ok := err.(gophercloud.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckFWGroupV2Exists(n string, group *groups.Group) resource.TestCheckFunc {
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
			return fmt.Errorf("Exists) Error creating OpenStack networking client: %s", err)
		}

		found, err := groups.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Firewall group not found. Expected %s, got %s", rs.Primary.ID, found.ID)
		}

		*group = *found

		return nil
	}
}

func testAccCheckFWGroupPortCount(group *groups.Group, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(group.Ports) != expected {
			return fmt.Errorf("Expected %d Ports, got %d", expected, len(group.Ports))
		}

		return nil
	}
}

func testAccCheckFWGroupV2(n, expectedName, expectedDescription string, IngressFirewallPolicyID *string) resource.TestCheckFunc {
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
			return fmt.Errorf("Exists) Error creating OpenStack networking client: %s", err)
		}

		var found *groups.Group
		for i := 0; i < 5; i++ {
			// Firewall group creation is asynchronous. Retry some times
			// if we get a 404 error. Fail on any other error.
			found, err = groups.Get(networkingClient, rs.Primary.ID).Extract()
			if err != nil {
				if _, ok := err.(gophercloud.ErrDefault404); ok {
					time.Sleep(time.Second)
					continue
				}
				return err
			}
			break
		}

		switch {
		case found.Name != expectedName:
			err = fmt.Errorf("Expected Name to be <%s> but found <%s>", expectedName, found.Name)
		case found.Description != expectedDescription:
			err = fmt.Errorf("Expected Description to be <%s> but found <%s>",
				expectedDescription, found.Description)
		case found.IngressFirewallPolicyID == "":
			err = fmt.Errorf("Policy should not be empty")
		}

		if err != nil {
			return err
		}

		IngressFirewallPolicyID = &found.IngressFirewallPolicyID

		return nil
	}
}

const testAccFWGroupV2Basic1 = `
resource "openstack_fw_policy_v2" "ingress_firewall_policy_1" {
  name = "ingress_firewall_policy_1"
}

resource "openstack_fw_policy_v2" "egress_firewall_policy_1" {
  name = "egress_firewall_policy_1"
}

resource "openstack_fw_group_v2" "group_1" {
  ingress_firewall_policy_id = "${openstack_fw_policy_v2.ingress_firewall_policy_1.id}"
  egress_firewall_policy_id  = "${openstack_fw_policy_v2.egress_firewall_policy_1.id}"
}
`

const testAccFWGroupV2Shared = `
resource "openstack_fw_policy_v2" "ingress_firewall_policy_1" {
  name = "ingress_firewall_policy_1"
}

resource "openstack_fw_policy_v2" "egress_firewall_policy_1" {
  name = "egress_firewall_policy_1"
}

resource "openstack_fw_group_v2" "group_1" {
  ingress_firewall_policy_id = "${openstack_fw_policy_v2.ingress_firewall_policy_1.id}"
  egress_firewall_policy_id  = "${openstack_fw_policy_v2.egress_firewall_policy_1.id}"
  shared                     = true
}
`

const testAccFWGroupV2Basic2 = `
resource "openstack_fw_policy_v2" "ingress_firewall_policy_1" {
  name = "ingress_firewall_policy_1"
}

resource "openstack_fw_policy_v2" "egress_firewall_policy_1" {
  name = "egress_firewall_policy_1"
}

resource "openstack_fw_group_v2" "group_1" {
  name                       = "group_1"
  description                = "terraform acceptance test"
  ingress_firewall_policy_id = "${openstack_fw_policy_v2.ingress_firewall_policy_1.id}"
  egress_firewall_policy_id  = "${openstack_fw_policy_v2.egress_firewall_policy_1.id}"
  admin_state_up             = true
}
`

const testAccFWGroupV2Basic3 = `
resource "openstack_fw_policy_v2" "ingress_firewall_policy_1" {
  name = "ingress_firewall_policy_1"
}

resource "openstack_fw_policy_v2" "egress_firewall_policy_1" {
  name = "egress_firewall_policy_1"
}

resource "openstack_fw_group_v2" "group_1" {
  name                       = "new_name_group_1"
  description                = "new description terraform acceptance test"
  ingress_firewall_policy_id = "${openstack_fw_policy_v2.ingress_firewall_policy_1.id}"
  egress_firewall_policy_id  = "${openstack_fw_policy_v2.egress_firewall_policy_1.id}"
  admin_state_up             = true
}
`

const testAccFWGroupV2Port = `
resource "openstack_networking_router_v2" "router_1" {
  name           = "router_1"
  admin_state_up = true
}

resource "openstack_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  network_id = "${openstack_networking_network_v2.network_1.id}"
  cidr       = "10.20.30.0/24"
  ip_version = 4
}

resource "openstack_networking_router_interface_v2" "router_interface_1" {
  router_id = "${openstack_networking_router_v2.router_1.id}"
  subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
}

resource "openstack_fw_policy_v2" "ingress_firewall_policy_1" {
  name = "ingress_firewall_policy_1"
}

resource "openstack_fw_group_v2" "group_1" {
  name        = "group_1"
  description = "firewall group port test"
  ports       = [
    "${openstack_networking_router_interface_v2.router_interface_1.port_id}",
  ]
}
`

const testAccFWGroupV2PortAdd = `
resource "openstack_networking_router_v2" "router_1" {
  name           = "router_1"
  admin_state_up = "true"
}

resource "openstack_networking_network_v2" "network_1" {
  name           = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  network_id = "${openstack_networking_network_v2.network_1.id}"
  cidr       = "10.20.30.0/24"
  ip_version = 4
}

resource "openstack_networking_router_interface_v2" "router_interface_1" {
  router_id = "${openstack_networking_router_v2.router_1.id}"
  subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
}

resource "openstack_networking_router_v2" "router_2" {
  name           = "router_2"
  admin_state_up = "true"
}

resource "openstack_networking_network_v2" "network_2" {
  name           = "network_2"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_2" {
  network_id = "${openstack_networking_network_v2.network_2.id}"
  cidr       = "20.30.40.0/24"
  ip_version = 4
}

resource "openstack_networking_router_interface_v2" "router_interface_2" {
  router_id = "${openstack_networking_router_v2.router_2.id}"
  subnet_id = "${openstack_networking_subnet_v2.subnet_2.id}"
}

resource "openstack_fw_policy_v2" "ingress_firewall_policy_1" {
  name = "ingress_firewall_policy_1"
}

resource "openstack_fw_group_v2" "group_1" {
  name                       = "group_1"
  description                = "firewall group port test"
  ingress_firewall_policy_id = "${openstack_fw_policy_v2.ingress_firewall_policy_1.id}"
  ports                      = [
    "${openstack_networking_router_interface_v2.router_interface_1.port_id}",
    "${openstack_networking_router_interface_v2.router_interface_2.port_id}",
  ]
}
`

const testAccFWGroupV2PortRemove = `
resource "openstack_fw_policy_v2" "ingress_firewall_policy_1" {
  name = "ingress_firewall_policy_1"
}

resource "openstack_fw_group_v2" "group_1" {
  name                       = "group_1"
  description                = "firewall group port test"
  ingress_firewall_policy_id = "${openstack_fw_policy_v2.ingress_firewall_policy_1.id}"
  ports                      = []
}
`

const testAccFWGroupV2NoPort = `
resource "openstack_fw_policy_v2" "ingress_firewall_policy_1" {
  name = "ingress_firewall_policy_1"
}

resource "openstack_fw_policy_v2" "egress_firewall_policy_1" {
  name = "egress_firewall_policy_1"
}

resource "openstack_fw_group_v2" "group_1" {
  name                       = "group_1"
  description                = "firewall group port test"
  ingress_firewall_policy_id = "${openstack_fw_policy_v2.ingress_firewall_policy_1.id}"
  egress_firewall_policy_id  = "${openstack_fw_policy_v2.egress_firewall_policy_1.id}"
}
`
