package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/floatingips"
)

func TestAccComputeV2FloatingIPAssociate_basic(t *testing.T) {
	var instance servers.Server
	var fip floatingips.FloatingIP

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2FloatingIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2FloatingIPAssociateBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					testAccCheckNetworkingV2FloatingIPExists("openstack_networking_floatingip_v2.fip_1", &fip),
					testAccCheckComputeV2FloatingIPAssociateAssociated(&fip, &instance, 1),
				),
			},
			{
				Config: testAccComputeV2FloatingIPAssociateUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					testAccCheckNetworkingV2FloatingIPExists("openstack_networking_floatingip_v2.fip_1", &fip),
					testAccCheckComputeV2FloatingIPAssociateAssociated(&fip, &instance, 1),
				),
			},
		},
	})
}

func TestAccComputeV2FloatingIPAssociate_fixedIP(t *testing.T) {
	var instance servers.Server
	var fip floatingips.FloatingIP

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2FloatingIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2FloatingIPAssociateFixedIP(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					testAccCheckNetworkingV2FloatingIPExists("openstack_networking_floatingip_v2.fip_1", &fip),
					testAccCheckComputeV2FloatingIPAssociateAssociated(&fip, &instance, 1),
				),
			},
		},
	})
}

// Note: test disabled due to SDK V2 fails with the following error:
// Error: Invalid index
//
//   on terraform_plugin_test.tf line 17, in resource "openstack_compute_floatingip_associate_v2" "fip_1":
//   17:   fixed_ip = "${openstack_compute_instance_v2.instance_1.network.0.fixed_ip_v4}"
//     ├────────────────
//     │ openstack_compute_instance_v2.instance_1.network is empty list of object
//
// The given key does not identify an element in this collection value: the
// collection has no elements.
//     TestAccComputeV2FloatingIPAssociate_attachToFirstNetwork: resource_openstack_compute_floatingip_associate_v2_test.go:75: Step 1/1 error: Error running apply: exit status 1
//
//         Error: Invalid index
//
//           on terraform_plugin_test.tf line 17, in resource "openstack_compute_floatingip_associate_v2" "fip_1":
//           17:   fixed_ip = "${openstack_compute_instance_v2.instance_1.network.0.fixed_ip_v4}"
//             ├────────────────
//             │ openstack_compute_instance_v2.instance_1.network is empty list of object
//
//         The given key does not identify an element in this collection value: the
//         collection has no elements.
//func TestAccComputeV2FloatingIPAssociate_attachToFirstNetwork(t *testing.T) {
//	var instance servers.Server
//	var fip floatingips.FloatingIP
//
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() {
//			testAccPreCheck(t)
//			testAccPreCheckNonAdminOnly(t)
//		},
//		ProviderFactories: testAccProviders,
//		CheckDestroy:      testAccCheckComputeV2FloatingIPAssociateDestroy,
//		Steps: []resource.TestStep{
//			{
//				Config: testAccComputeV2FloatingIPAssociateAttachToFirstNetwork(),
//				Check: resource.ComposeTestCheckFunc(
//					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
//					testAccCheckNetworkingV2FloatingIPExists("openstack_networking_floatingip_v2.fip_1", &fip),
//					testAccCheckComputeV2FloatingIPAssociateAssociated(&fip, &instance, 1),
//				),
//			},
//		},
//	})
//}

func TestAccComputeV2FloatingIPAssociate_attachNew(t *testing.T) {
	var instance servers.Server
	var floatingIP1 floatingips.FloatingIP
	var floatingIP2 floatingips.FloatingIP

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2FloatingIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2FloatingIPAssociateAttachNew1(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					testAccCheckNetworkingV2FloatingIPExists("openstack_networking_floatingip_v2.fip_1", &floatingIP1),
					testAccCheckNetworkingV2FloatingIPExists("openstack_networking_floatingip_v2.fip_2", &floatingIP2),
					testAccCheckComputeV2FloatingIPAssociateAssociated(&floatingIP1, &instance, 1),
				),
			},
			{
				Config: testAccComputeV2FloatingIPAssociateAttachNew2(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					testAccCheckNetworkingV2FloatingIPExists("openstack_networking_floatingip_v2.fip_1", &floatingIP1),
					testAccCheckNetworkingV2FloatingIPExists("openstack_networking_floatingip_v2.fip_2", &floatingIP2),
					testAccCheckComputeV2FloatingIPAssociateAssociated(&floatingIP2, &instance, 1),
				),
			},
		},
	})
}

func TestAccComputeV2FloatingIPAssociate_waitUntilAssociated(t *testing.T) {
	var instance servers.Server
	var fip floatingips.FloatingIP

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckComputeV2FloatingIPAssociateDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccComputeV2FloatingIPAssociateWaitUntilAssociated(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckComputeV2InstanceExists("openstack_compute_instance_v2.instance_1", &instance),
					testAccCheckNetworkingV2FloatingIPExists("openstack_networking_floatingip_v2.fip_1", &fip),
					testAccCheckComputeV2FloatingIPAssociateAssociated(&fip, &instance, 1),
				),
			},
		},
	})
}

func testAccCheckComputeV2FloatingIPAssociateDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	computeClient, err := config.ComputeV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_compute_floatingip_associate_v2" {
			continue
		}

		floatingIP, instanceID, _, err := parseComputeFloatingIPAssociateID(rs.Primary.ID)
		if err != nil {
			return err
		}

		instance, err := servers.Get(computeClient, instanceID).Extract()
		if err != nil {
			// If the error is a 404, then the instance does not exist,
			// and therefore the floating IP cannot be associated to it.
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return nil
			}
			return err
		}

		// But if the instance still exists, then walk through its known addresses
		// and see if there's a floating IP.
		for _, networkAddresses := range instance.Addresses {
			for _, element := range networkAddresses.([]interface{}) {
				address := element.(map[string]interface{})
				if address["OS-EXT-IPS:type"] == "floating" {
					return fmt.Errorf("Floating IP %s is still attached to instance %s", floatingIP, instanceID)
				}
			}
		}
	}

	return nil
}

func testAccCheckComputeV2FloatingIPAssociateAssociated(
	fip *floatingips.FloatingIP, instance *servers.Server, n int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)
		computeClient, err := config.ComputeV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack compute client: %s", err)
		}

		newInstance, err := servers.Get(computeClient, instance.ID).Extract()
		if err != nil {
			return err
		}

		// Walk through the instance's addresses and find the match
		i := 0
		for _, networkAddresses := range newInstance.Addresses {
			i++
			if i != n {
				continue
			}
			for _, element := range networkAddresses.([]interface{}) {
				address := element.(map[string]interface{})
				if address["OS-EXT-IPS:type"] == "floating" && address["addr"] == fip.FloatingIP {
					return nil
				}
			}
		}
		return fmt.Errorf("Floating IP %s was not attached to instance %s", fip.FloatingIP, instance.ID)
	}
}

func testAccComputeV2FloatingIPAssociateBasic() string {
	return fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
}

resource "openstack_networking_floatingip_v2" "fip_1" {
}

resource "openstack_compute_floatingip_associate_v2" "fip_1" {
  floating_ip = "${openstack_networking_floatingip_v2.fip_1.address}"
  instance_id = "${openstack_compute_instance_v2.instance_1.id}"
}
`, osNetworkID)
}

func testAccComputeV2FloatingIPAssociateUpdate() string {
	return fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
}

resource "openstack_networking_floatingip_v2" "fip_1" {
  description = "test"
}

resource "openstack_compute_floatingip_associate_v2" "fip_1" {
  floating_ip = "${openstack_networking_floatingip_v2.fip_1.address}"
  instance_id = "${openstack_compute_instance_v2.instance_1.id}"
}
`, osNetworkID)
}

func testAccComputeV2FloatingIPAssociateFixedIP() string {
	return fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
}

resource "openstack_networking_floatingip_v2" "fip_1" {
}

resource "openstack_compute_floatingip_associate_v2" "fip_1" {
  floating_ip = "${openstack_networking_floatingip_v2.fip_1.address}"
  instance_id = "${openstack_compute_instance_v2.instance_1.id}"
  fixed_ip = "${openstack_compute_instance_v2.instance_1.access_ip_v4}"
}
`, osNetworkID)
}

//func testAccComputeV2FloatingIPAssociateAttachToFirstNetwork() string {
//	return fmt.Sprintf(`
//resource "openstack_compute_instance_v2" "instance_1" {
//  name = "instance_1"
//  security_groups = ["default"]
//
//  network {
//    uuid = "%s"
//  }
//}
//
//resource "openstack_networking_floatingip_v2" "fip_1" {
//}
//
//resource "openstack_compute_floatingip_associate_v2" "fip_1" {
//  floating_ip = "${openstack_networking_floatingip_v2.fip_1.address}"
//  instance_id = "${openstack_compute_instance_v2.instance_1.id}"
//  fixed_ip = "${openstack_compute_instance_v2.instance_1.network.0.fixed_ip_v4}"
//}
//`, osNetworkID)
//}

func testAccComputeV2FloatingIPAssociateAttachNew1() string {
	return fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
}

resource "openstack_networking_floatingip_v2" "fip_1" {
}

resource "openstack_networking_floatingip_v2" "fip_2" {
}

resource "openstack_compute_floatingip_associate_v2" "fip_1" {
  floating_ip = "${openstack_networking_floatingip_v2.fip_1.address}"
  instance_id = "${openstack_compute_instance_v2.instance_1.id}"
}
`, osNetworkID)
}

func testAccComputeV2FloatingIPAssociateAttachNew2() string {
	return fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
}

resource "openstack_networking_floatingip_v2" "fip_1" {
}

resource "openstack_networking_floatingip_v2" "fip_2" {
}

resource "openstack_compute_floatingip_associate_v2" "fip_1" {
  floating_ip = "${openstack_networking_floatingip_v2.fip_2.address}"
  instance_id = "${openstack_compute_instance_v2.instance_1.id}"
}
`, osNetworkID)
}

func testAccComputeV2FloatingIPAssociateWaitUntilAssociated() string {
	return fmt.Sprintf(`
resource "openstack_compute_instance_v2" "instance_1" {
  name = "instance_1"
  security_groups = ["default"]
  network {
    uuid = "%s"
  }
}

resource "openstack_networking_floatingip_v2" "fip_1" {
}

resource "openstack_compute_floatingip_associate_v2" "fip_1" {
  floating_ip = "${openstack_networking_floatingip_v2.fip_1.address}"
  instance_id = "${openstack_compute_instance_v2.instance_1.id}"

  wait_until_associated = true
}
`, osNetworkID)
}
