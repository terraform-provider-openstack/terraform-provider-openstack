package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/attachinterfaces"
)

var FixedIP string = "10.119.100.100"

func TestAccComputeV2AttachInterface_basic(t *testing.T) {
	var ai attachinterfaces.Interface

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2AttachInterfaceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeV2AttachInterface_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2AttachInterfaceExists("openstack_compute_attach_interface_v2.ai_1", &ai),
				),
			},
		},
	})
}

func TestAccComputeV2AttachInterface_IP(t *testing.T) {
	var ai attachinterfaces.Interface

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2AttachInterfaceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeV2AttachInterface_IP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2AttachInterfaceExists("openstack_compute_attach_interface_v2.ai_1", &ai),
					testAccCheckComputeV2AttachInterfaceIP(&ai, FixedIP),
				),
			},
		},
	})
}

func TestAccComputeV2AttachInterface_timeout(t *testing.T) {
	var ai attachinterfaces.Interface

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckComputeV2AttachInterfaceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeV2AttachInterface_timeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2AttachInterfaceExists("openstack_compute_attach_interface_v2.ai_1", &ai),
				),
			},
		},
	})
}

func testAccCheckComputeV2AttachInterfaceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	computeClient, err := config.computeV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_compute_attach_interface_v2" {
			continue
		}

		instanceId, portId, err := parseComputeAttachInterfaceId(rs.Primary.ID)
		if err != nil {
			return err
		}

		_, err = attachinterfaces.Get(computeClient, instanceId, portId).Extract()
		if err == nil {
			return fmt.Errorf("Volume attachment still exists")
		}
	}

	return nil
}

func testAccCheckComputeV2AttachInterfaceExists(n string, ai *attachinterfaces.Interface) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		computeClient, err := config.computeV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack compute client: %s", err)
		}

		instanceId, portId, err := parseComputeAttachInterfaceId(rs.Primary.ID)
		if err != nil {
			return err
		}

		found, err := attachinterfaces.Get(computeClient, instanceId, portId).Extract()
		if err != nil {
			return err
		}

		//if found.instanceID != instanceID || found.PortID != portId {
		if found.PortID != portId {
			return fmt.Errorf("AttachInterface not found")
		}

		*ai = *found

		return nil
	}
}

func testAccCheckComputeV2AttachInterfaceIP(
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

var testAccComputeV2AttachInterface_basic = fmt.Sprintf(`
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

resource "openstack_compute_attach_interface_v2" "ai_1" {
  instance_id = "${openstack_compute_instance_v2.instance_1.id}"
  port_id = "${openstack_networking_port_v2.port_1.id}"
}
`, OS_NETWORK_ID, OS_NETWORK_ID)

var testAccComputeV2AttachInterface_IP = fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
}

resource "openstack_compute_attach_interface_v2" "ai_1" {
  instance_id = "${openstack_compute_instance_v2.instance_1.id}"
  network_id = "%s"
  fixed_ip = "%s"
}
`, OS_NETWORK_ID, OS_NETWORK_ID, FixedIP)

var testAccComputeV2AttachInterface_timeout = fmt.Sprintf(`
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

resource "openstack_compute_attach_interface_v2" "ai_1" {
  instance_id = "${openstack_compute_instance_v2.instance_1.id}"
  port_id = "${openstack_networking_port_v2.port_1.id}"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`, OS_NETWORK_ID, OS_NETWORK_ID)
