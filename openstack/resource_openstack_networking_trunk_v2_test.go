package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/trunks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
)

func TestAccNetworkingV2Trunk_nosubports(t *testing.T) {
	var port_1 ports.Port
	var trunk_1 trunks.Trunk

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2TrunkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Trunk_noSubports,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.parent_port_1", &port_1),
					testAccCheckNetworkingV2TrunkExists("openstack_networking_trunk_v2.trunk_1", []string{}, &trunk_1),
					resource.TestCheckResourceAttr(
						"openstack_networking_trunk_v2.trunk_1", "name", "trunk_1"),
					resource.TestCheckResourceAttr(
						"openstack_networking_trunk_v2.trunk_1", "description", "trunk_1 description"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Trunk_subports(t *testing.T) {
	var parent_port_1, subport_1, subport_2 ports.Port
	var trunk_1 trunks.Trunk

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2TrunkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Trunk_subports,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.parent_port_1", &parent_port_1),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.subport_1", &subport_1),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.subport_2", &subport_2),
					testAccCheckNetworkingV2TrunkExists("openstack_networking_trunk_v2.trunk_1", []string{"openstack_networking_port_v2.subport_1", "openstack_networking_port_v2.subport_2"}, &trunk_1, &subport_1, &subport_2),
				),
			},
		},
	})
}

func TestAccNetworkingV2Trunk_tags(t *testing.T) {
	var parent_port_1 ports.Port
	var trunk_1 trunks.Trunk

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2TrunkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Trunk_tags_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.parent_port_1", &parent_port_1),
					testAccCheckNetworkingV2TrunkExists("openstack_networking_trunk_v2.trunk_1", []string{}, &trunk_1),
					testAccCheckNetworkingV2Tags("openstack_networking_trunk_v2.trunk_1", []string{"a", "b", "c"}),
				),
			},
			{
				Config: testAccNetworkingV2Trunk_tags_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.parent_port_1", &parent_port_1),
					testAccCheckNetworkingV2TrunkExists("openstack_networking_trunk_v2.trunk_1", []string{}, &trunk_1),
					testAccCheckNetworkingV2Tags("openstack_networking_trunk_v2.trunk_1", []string{"c", "d", "e"}),
				),
			},
		},
	})
}

func TestAccNetworkingV2Trunk_trunkUpdateSubports(t *testing.T) {
	var parent_port_1, subport_1, subport_2, subport_3, subport_4 ports.Port
	var trunk_1 trunks.Trunk

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2TrunkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Trunk_updateSubports_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.parent_port_1", &parent_port_1),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.subport_1", &subport_1),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.subport_2", &subport_2),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.subport_3", &subport_3),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.subport_4", &subport_4),
					testAccCheckNetworkingV2TrunkExists("openstack_networking_trunk_v2.trunk_1", []string{"openstack_networking_port_v2.subport_1", "openstack_networking_port_v2.subport_2"}, &trunk_1, &subport_1, &subport_2),
					resource.TestCheckResourceAttr(
						"openstack_networking_trunk_v2.trunk_1", "description", "trunk_1 description"),
				),
			},
			{
				Config: testAccNetworkingV2Trunk_updateSubports_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.parent_port_1", &parent_port_1),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.subport_1", &subport_1),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.subport_2", &subport_2),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.subport_3", &subport_3),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.subport_4", &subport_4),
					testAccCheckNetworkingV2TrunkExists("openstack_networking_trunk_v2.trunk_1", []string{"openstack_networking_port_v2.subport_1", "openstack_networking_port_v2.subport_3", "openstack_networking_port_v2.subport_4"}, &trunk_1, &subport_1, &subport_3, &subport_4),
					resource.TestCheckResourceAttr(
						"openstack_networking_trunk_v2.trunk_1", "description", ""),
				),
			},
			{
				Config: testAccNetworkingV2Trunk_updateSubports_3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.parent_port_1", &parent_port_1),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.subport_1", &subport_1),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.subport_2", &subport_2),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.subport_3", &subport_3),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.subport_4", &subport_4),
					testAccCheckNetworkingV2TrunkExists("openstack_networking_trunk_v2.trunk_1", []string{"openstack_networking_port_v2.subport_1", "openstack_networking_port_v2.subport_3", "openstack_networking_port_v2.subport_4"}, &trunk_1, &subport_1, &subport_3, &subport_4),
					resource.TestCheckResourceAttr(
						"openstack_networking_trunk_v2.trunk_1", "description", ""),
				),
			},
			{
				Config: testAccNetworkingV2Trunk_updateSubports_4,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.parent_port_1", &parent_port_1),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.subport_1", &subport_1),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.subport_2", &subport_2),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.subport_3", &subport_3),
					testAccCheckNetworkingV2PortExists("openstack_networking_port_v2.subport_4", &subport_4),
					testAccCheckNetworkingV2TrunkExists("openstack_networking_trunk_v2.trunk_1", []string{}, &trunk_1),
					resource.TestCheckResourceAttr(
						"openstack_networking_trunk_v2.trunk_1", "description", "trunk_1 updated description"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Trunk_computeInstance(t *testing.T) {
	var instance_1 servers.Server
	var parent_port_1, subport_1 ports.Port
	var trunk_1 trunks.Trunk

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2InstanceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Trunk_computeInstance,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance_1),
					testAccCheckNetworkingV2PortExists(
						"openstack_networking_port_v2.parent_port_1", &parent_port_1),
					testAccCheckNetworkingV2PortExists(
						"openstack_networking_port_v2.subport_1", &subport_1),
					testAccCheckNetworkingV2TrunkExists("openstack_networking_trunk_v2.trunk_1", []string{"openstack_networking_port_v2.subport_1"}, &trunk_1, &subport_1),
					resource.TestCheckResourceAttrPtr(
						"openstack_compute_instance_v2.instance_1", "network.0.port", &trunk_1.PortID),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2TrunkDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	client, err := config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_networking_trunk_v2" {
			continue
		}

		_, err := trunks.Get(client, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Trunk still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2TrunkExists(n string, subportResourceNames []string, trunk *trunks.Trunk, subports ...*ports.Port) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Trunk not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Trunk ID is not set")
		}

		var subportResources map[string]bool
		if len(subports) > 0 {
			if len(subportResourceNames) != len(subports) {
				return fmt.Errorf("Amount of subport resource names and subports do not match")
			}

			subportResources = make(map[string]bool)
			for i, subport := range subports {
				if subportResource, ok := s.RootModule().Resources[subportResourceNames[i]]; ok {
					subportResources[subportResource.Primary.ID] = true
				} else {
					return fmt.Errorf("Subport not found: %s", subport.ID)
				}
			}
		}

		config := testAccProvider.Meta().(*Config)
		client, err := config.networkingV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %s", err)
		}

		found, err := trunks.Get(client, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if len(found.Subports) != len(subports) {
			return fmt.Errorf("The amount of retrieved trunk subports and trunk subports to check does not match")
		}

		if len(subports) > 0 {
			for _, subport := range found.Subports {
				if _, ok := subportResources[subport.PortID]; !ok {
					return fmt.Errorf("Trunk Subport not found: %s", subport.PortID)
				}
			}
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Trunk not found")
		}

		if found.Name != rs.Primary.Attributes["name"] {
			return fmt.Errorf("Trunk name does not match")
		}

		*trunk = *found

		return nil
	}
}

const testAccNetworkingV2Trunk_noSubports = `
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

resource "openstack_networking_port_v2" "parent_port_1" {
  name = "parent_port_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_trunk_v2" "trunk_1" {
  name = "trunk_1"
  description = "trunk_1 description"
  port_id = "${openstack_networking_port_v2.parent_port_1.id}"
  admin_state_up = "true"
}
`

const testAccNetworkingV2Trunk_subports = `
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

resource "openstack_networking_port_v2" "parent_port_1" {
  name = "parent_port_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "subport_1" {
  name = "subport_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "subport_2" {
  name = "subport_2"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_trunk_v2" "trunk_1" {
  name = "trunk_1"
  description = "trunk_1 description"
  port_id = "${openstack_networking_port_v2.parent_port_1.id}"
  admin_state_up = "true"

  sub_port {
	  port_id = "${openstack_networking_port_v2.subport_1.id}"
	  segmentation_id = 1
	  segmentation_type = "vlan"
  }

  sub_port {
	  port_id = "${openstack_networking_port_v2.subport_2.id}"
	  segmentation_id = 2
	  segmentation_type = "vlan"
  }
}
`

const testAccNetworkingV2Trunk_updateSubports_1 = `
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

resource "openstack_networking_port_v2" "parent_port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "subport_1" {
  name = "subport_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "subport_2" {
  name = "subport_2"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "subport_3" {
  name = "subport_3"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "subport_4" {
  name = "subport_4"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_trunk_v2" "trunk_1" {
  name = "trunk_1"
  description = "trunk_1 description"
  admin_state_up = "true"
  port_id = "${openstack_networking_port_v2.parent_port_1.id}"

  sub_port {
	  port_id = "${openstack_networking_port_v2.subport_1.id}"
	  segmentation_id = 1
	  segmentation_type = "vlan"
  }

  sub_port {
	  port_id = "${openstack_networking_port_v2.subport_2.id}"
	  segmentation_id = 2
	  segmentation_type = "vlan"
  }
}
`

const testAccNetworkingV2Trunk_updateSubports_2 = `
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

resource "openstack_networking_port_v2" "parent_port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "subport_1" {
  name = "subport_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "subport_2" {
  name = "subport_2"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "subport_3" {
  name = "subport_3"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "subport_4" {
  name = "subport_4"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_trunk_v2" "trunk_1" {
  name = "update_trunk_1"
  admin_state_up = "true"
  port_id = "${openstack_networking_port_v2.parent_port_1.id}"

  sub_port {
	  port_id = "${openstack_networking_port_v2.subport_1.id}"
	  segmentation_id = 1
	  segmentation_type = "vlan"
  }

  sub_port {
	  port_id = "${openstack_networking_port_v2.subport_3.id}"
	  segmentation_id = 3
	  segmentation_type = "vlan"
  }

  sub_port {
	  port_id = "${openstack_networking_port_v2.subport_4.id}"
	  segmentation_id = 4
	  segmentation_type = "vlan"
  }
}
`

const testAccNetworkingV2Trunk_updateSubports_3 = `
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

resource "openstack_networking_port_v2" "parent_port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "subport_1" {
  name = "subport_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "subport_2" {
  name = "subport_2"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "subport_3" {
  name = "subport_3"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "subport_4" {
  name = "subport_4"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_trunk_v2" "trunk_1" {
  name = "trunk_1"
  description = ""
  admin_state_up = "true"
  port_id = "${openstack_networking_port_v2.parent_port_1.id}"

  sub_port {
	  port_id = "${openstack_networking_port_v2.subport_1.id}"
	  segmentation_id = 1
	  segmentation_type = "vlan"
  }

  sub_port {
	  port_id = "${openstack_networking_port_v2.subport_3.id}"
	  segmentation_id = 3
	  segmentation_type = "vlan"
  }

  sub_port {
	  port_id = "${openstack_networking_port_v2.subport_4.id}"
	  segmentation_id = 4
	  segmentation_type = "vlan"
  }
}
`

const testAccNetworkingV2Trunk_updateSubports_4 = `
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

resource "openstack_networking_port_v2" "parent_port_1" {
  name = "port_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "subport_1" {
  name = "subport_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "subport_2" {
  name = "subport_2"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "subport_3" {
  name = "subport_3"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_port_v2" "subport_4" {
  name = "subport_4"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_trunk_v2" "trunk_1" {
  name = "trunk_1"
  description = "trunk_1 updated description"
  port_id = "${openstack_networking_port_v2.parent_port_1.id}"
  admin_state_up = "true"
}
`

const testAccNetworkingV2Trunk_computeInstance = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  cidr = "192.168.1.0/24"
  ip_version = 4
  enable_dhcp = true
  no_gateway = true
}

resource "openstack_networking_port_v2" "parent_port_1" {
  depends_on = [
    "openstack_networking_subnet_v2.subnet_1",
  ]

  name = "parent_port_1"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  admin_state_up = "true"
}

resource "openstack_networking_port_v2" "subport_1" {
  depends_on = [
    "openstack_networking_subnet_v2.subnet_1",
  ]

  name = "subport_1"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  admin_state_up = "true"
}

resource "openstack_networking_trunk_v2" "trunk_1" {
  name = "trunk_1"
  admin_state_up = "true"
  port_id = "${openstack_networking_port_v2.parent_port_1.id}"

  sub_port {
	  port_id = "${openstack_networking_port_v2.subport_1.id}"
	  segmentation_id = 1
	  segmentation_type = "vlan"
  }
}

resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]

  network {
    port = "${openstack_networking_trunk_v2.trunk_1.port_id}"
  }
}
`

const testAccNetworkingV2Trunk_tags_1 = `
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

resource "openstack_networking_port_v2" "parent_port_1" {
  name = "parent_port_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_trunk_v2" "trunk_1" {
  name = "trunk_1"
  port_id = "${openstack_networking_port_v2.parent_port_1.id}"
  admin_state_up = "true"

  tags = ["a", "b", "c"]
}
`

const testAccNetworkingV2Trunk_tags_2 = `
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

resource "openstack_networking_port_v2" "parent_port_1" {
  name = "parent_port_1"
  admin_state_up = "true"
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_trunk_v2" "trunk_1" {
  name = "trunk_1"
  port_id = "${openstack_networking_port_v2.parent_port_1.id}"
  admin_state_up = "true"

  tags = ["c", "d", "e"]
}
`
