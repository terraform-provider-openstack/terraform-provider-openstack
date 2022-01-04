package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/attachinterfaces"
)

func TestAccComputeV2InterfaceAttach_basic(t *testing.T) {
	var ai attachinterfaces.Interface

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2InterfaceAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InterfaceAttachBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InterfaceAttachExists("openstack_compute_interface_attach_v2.ai_1", &ai),
				),
			},
		},
	})
}

func TestAccComputeV2InterfaceAttach_IP(t *testing.T) {
	var ai attachinterfaces.Interface

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2InterfaceAttachDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2InterfaceAttachIP(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InterfaceAttachExists("openstack_compute_interface_attach_v2.ai_1", &ai),
					testAccCheckComputeV2InterfaceAttachIP(&ai, "192.168.1.100"),
				),
			},
		},
	})
}

func testAccCheckComputeV2InterfaceAttachDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	computeClient, err := config.ComputeV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_compute_interface_attach_v2" {
			continue
		}

		instanceID, portID, err := computeInterfaceAttachV2ParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = attachinterfaces.Get(computeClient, instanceID, portID).Extract()
		if err == nil {
			return fmt.Errorf("Volume attachment still exists")
		}
	}

	return nil
}

func testAccCheckComputeV2InterfaceAttachExists(n string, ai *attachinterfaces.Interface) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		computeClient, err := config.ComputeV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack compute client: %s", err)
		}

		instanceID, portID, err := computeInterfaceAttachV2ParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		found, err := attachinterfaces.Get(computeClient, instanceID, portID).Extract()
		if err != nil {
			return err
		}

		//if found.instanceID != instanceID || found.PortID != portID {
		if found.PortID != portID {
			return fmt.Errorf("InterfaceAttach not found")
		}

		*ai = *found

		return nil
	}
}

func testAccCheckComputeV2InterfaceAttachIP(
	ai *attachinterfaces.Interface, ip string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, i := range ai.FixedIPs {
			if i.IPAddress == ip {
				return nil
			}
		}
		return fmt.Errorf("Requested ip (%s) does not exist on port", ip)
	}
}

func testAccComputeV2InterfaceAttachBasic() string {
	return fmt.Sprintf(`
resource "openstack_networking_port_v2" "port_1" {
  name = "port_1"
  network_id = "%s"
  admin_state_up = "true"
}

resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
}

resource "openstack_compute_interface_attach_v2" "ai_1" {
  instance_id = "${openstack_compute_instance_v2.instance_1.id}"
  port_id = "${openstack_networking_port_v2.port_1.id}"
}
`, osNetworkID, osNetworkID)
}

func testAccComputeV2InterfaceAttachIP() string {
	return fmt.Sprintf(`
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  cidr = "192.168.1.0/24"
  ip_version = 4
  enable_dhcp = true
  no_gateway = true
}

resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
}

resource "openstack_compute_interface_attach_v2" "ai_1" {
  instance_id = "${openstack_compute_instance_v2.instance_1.id}"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  fixed_ip = "192.168.1.100"
}
`, osNetworkID)
}
