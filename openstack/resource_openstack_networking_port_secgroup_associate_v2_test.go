package openstack

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/networks"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
)

func TestAccNetworkingV2PortSecGroupAssociate_update(t *testing.T) {
	var port ports.Port

	if os.Getenv("TF_ACC") != "" {
		hiddenPort, err := testAccCheckNetworkingV2PortSecGroupCreatePort(t, "hidden_port", true)
		if err != nil {
			t.Fatal(err)
		}
		defer testAccCheckNetworkingV2PortSecGroupDeletePort(t, hiddenPort) //nolint:errcheck
	}

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			// enforce = false
			{ // step 0
				Config: testAccNetworkingV2PortSecGroupAssociateManifestUpdate0(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortSecGroupAssociateExists("openstack_networking_port_secgroup_associate_v2.port_1", &port),
					testAccCheckNetworkingV2PortSecGroupAssociateCountSecurityGroups(&port, 3),
				),
			},
			{ // step 1
				Config: testAccNetworkingV2PortSecGroupAssociateManifestUpdate1(), // unset user defined security groups only
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortSecGroupAssociateExists("data.openstack_networking_port_v2.hidden_port_1", &port),
					testAccCheckNetworkingV2PortSecGroupAssociateCountSecurityGroups(&port, 2),
				),
			},
			// enforce = true
			{ // step 2
				Config: testAccNetworkingV2PortSecGroupAssociateManifestUpdate2(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortSecGroupAssociateExists("openstack_networking_port_secgroup_associate_v2.port_1", &port),
					testAccCheckNetworkingV2PortSecGroupAssociateCountSecurityGroups(&port, 1),
				),
			},
			{ // step 3
				Config: testAccNetworkingV2PortSecGroupAssociateManifestUpdate3(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortSecGroupAssociateExists("openstack_networking_port_secgroup_associate_v2.port_1", &port),
					testAccCheckNetworkingV2PortSecGroupAssociateCountSecurityGroups(&port, 3),
				),
			},
			{ // step 4
				Config: testAccNetworkingV2PortSecGroupAssociateManifestUpdate4(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortSecGroupAssociateExists("openstack_networking_port_secgroup_associate_v2.port_1", &port),
					testAccCheckNetworkingV2PortSecGroupAssociateCountSecurityGroups(&port, 1),
				),
			},
			{ // step 5
				Config: testAccNetworkingV2PortSecGroupAssociateManifestUpdate5(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortSecGroupAssociateExists("openstack_networking_port_secgroup_associate_v2.port_1", &port),
					testAccCheckNetworkingV2PortSecGroupAssociateCountSecurityGroups(&port, 0),
				),
			},
			{ // step 6
				Config: testAccNetworkingV2PortSecGroupAssociateManifestUpdate6(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortSecGroupAssociateExists("openstack_networking_port_secgroup_associate_v2.port_1", &port),
					testAccCheckNetworkingV2PortSecGroupAssociateCountSecurityGroups(&port, 1),
				),
			},
			// enforce = false
			{ // step 7
				Config: testAccNetworkingV2PortSecGroupAssociateManifestUpdate7(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortSecGroupAssociateExists("openstack_networking_port_secgroup_associate_v2.port_1", &port),
					testAccCheckNetworkingV2PortSecGroupAssociateCountSecurityGroups(&port, 1),
				),
			},
			{ // step 8
				Config: testAccNetworkingV2PortSecGroupAssociateManifestUpdate8(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortSecGroupAssociateExists("openstack_networking_port_secgroup_associate_v2.port_1", &port),
					testAccCheckNetworkingV2PortSecGroupAssociateCountSecurityGroups(&port, 3),
				),
			},
			{ // step 9
				Config: testAccNetworkingV2PortSecGroupAssociateManifestUpdate9(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortSecGroupAssociateExists("openstack_networking_port_secgroup_associate_v2.port_1", &port),
					testAccCheckNetworkingV2PortSecGroupAssociateCountSecurityGroups(&port, 1),
				),
			},
			{ // step 10
				Config: testAccNetworkingV2PortSecGroupAssociateManifestUpdate10(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortSecGroupAssociateExists("openstack_networking_port_secgroup_associate_v2.port_1", &port),
					testAccCheckNetworkingV2PortSecGroupAssociateCountSecurityGroups(&port, 0),
				),
			},
			{ // step 11
				Config: testAccNetworkingV2PortSecGroupAssociateManifestUpdate11(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortSecGroupAssociateExists("openstack_networking_port_secgroup_associate_v2.port_1", &port),
					testAccCheckNetworkingV2PortSecGroupAssociateCountSecurityGroups(&port, 1),
				),
			},
			{ // step 12
				Config: testAccNetworkingV2PortSecGroupAssociateManifestUpdate12(), // cleanup all the ports
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2PortSecGroupAssociateExists("data.openstack_networking_port_v2.hidden_port_1", &port),
					testAccCheckNetworkingV2PortSecGroupAssociateCountSecurityGroups(&port, 0),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2PortSecGroupCreatePort(t *testing.T, portName string, defaultSecGroups bool) (*ports.Port, error) {
	config, err := testAccAuthFromEnv()
	if err != nil {
		return nil, err
	}

	client, err := config.NetworkingV2Client(osRegionName)
	if err != nil {
		return nil, err
	}

	createNetOpts := networks.CreateOpts{
		Name:         "test",
		AdminStateUp: gophercloud.Enabled,
	}

	network, err := networks.Create(client, createNetOpts).Extract()
	if err != nil {
		return nil, err
	}

	t.Logf("Network %s created", network.ID)

	var securityGroups []string
	if defaultSecGroups {
		// create default security groups
		createSecGroupOpts := groups.CreateOpts{
			Name: "default_1",
		}

		secGroup1, err := groups.Create(client, createSecGroupOpts).Extract()
		if err != nil {
			return nil, err
		}

		t.Logf("Default security group 1 %s created", secGroup1.ID)

		createSecGroupOpts.Name = "default_2"

		secGroup2, err := groups.Create(client, createSecGroupOpts).Extract()
		if err != nil {
			return nil, err
		}

		t.Logf("Default security group 2 %s created", secGroup2.ID)

		// reversed order, just in case
		securityGroups = append(securityGroups, secGroup2.ID)
		securityGroups = append(securityGroups, secGroup1.ID)
	}

	// create port with default security groups assigned
	createOpts := ports.CreateOpts{
		NetworkID:      network.ID,
		Name:           portName,
		SecurityGroups: &securityGroups,
		AdminStateUp:   gophercloud.Enabled,
	}

	port, err := ports.Create(client, createOpts).Extract()
	if err != nil {
		nErr := networks.Delete(client, network.ID).ExtractErr()
		if nErr != nil {
			return nil, fmt.Errorf("Unable to create port (%s) and delete network (%s: %s)", err, network.ID, nErr)
		}
		return nil, err
	}

	t.Logf("Port %s created", port.ID)

	return port, nil
}

func testAccCheckNetworkingV2PortSecGroupDeletePort(t *testing.T, port *ports.Port) error {
	config, err := testAccAuthFromEnv()
	if err != nil {
		return err
	}

	client, err := config.NetworkingV2Client(osRegionName)
	if err != nil {
		return err
	}

	err = ports.Delete(client, port.ID).ExtractErr()
	if err != nil {
		return err
	}

	t.Logf("Port %s deleted", port.ID)

	// delete default security groups
	for _, secGroupID := range port.SecurityGroups {
		err = groups.Delete(client, secGroupID).ExtractErr()
		if err != nil {
			return err
		}
		t.Logf("Default security group %s deleted", secGroupID)
	}

	err = networks.Delete(client, port.NetworkID).ExtractErr()
	if err != nil {
		return err
	}

	t.Logf("Network %s deleted", port.NetworkID)

	return nil
}

func testAccCheckNetworkingV2PortSecGroupAssociateExists(n string, port *ports.Port) resource.TestCheckFunc {
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

func testAccCheckNetworkingV2PortSecGroupAssociateCountSecurityGroups(port *ports.Port, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if len(port.SecurityGroups) != expected {
			return fmt.Errorf("Expected %d Security Groups, got %d", expected, len(port.SecurityGroups))
		}

		return nil
	}
}

const testAccNetworkingV2PortSecGroupAssociate = `
resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "terraform security group acceptance test"
}

resource "openstack_networking_secgroup_v2" "secgroup_2" {
  name = "secgroup_2"
  description = "terraform security group acceptance test"
}

data "openstack_networking_secgroup_v2" "default_1" {
  name = "default_1"
}

data "openstack_networking_secgroup_v2" "default_2" {
  name = "default_2"
}

data "openstack_networking_port_v2" "hidden_port_1" {
  name = "hidden_port"
}
`

func testAccNetworkingV2PortSecGroupAssociateManifestUpdate0() string {
	return fmt.Sprintf(`
%s

resource "openstack_networking_port_secgroup_associate_v2" "port_1" {
  port_id = "${data.openstack_networking_port_v2.hidden_port_1.id}"
  enforce = "false"
  security_group_ids = [
    "${openstack_networking_secgroup_v2.secgroup_1.id}",
  ]
}
`, testAccNetworkingV2PortSecGroupAssociate)
}

func testAccNetworkingV2PortSecGroupAssociateManifestUpdate1() string {
	return fmt.Sprintf(`
%s
`, testAccNetworkingV2PortSecGroupAssociate)
}

func testAccNetworkingV2PortSecGroupAssociateManifestUpdate2() string {
	return fmt.Sprintf(`
%s

resource "openstack_networking_port_secgroup_associate_v2" "port_1" {
  port_id = "${data.openstack_networking_port_v2.hidden_port_1.id}"
  enforce = "true"
  security_group_ids = [
    "${openstack_networking_secgroup_v2.secgroup_1.id}",
  ]
}
`, testAccNetworkingV2PortSecGroupAssociate)
}

func testAccNetworkingV2PortSecGroupAssociateManifestUpdate3() string {
	return fmt.Sprintf(`
%s

resource "openstack_networking_port_secgroup_associate_v2" "port_1" {
  port_id = "${data.openstack_networking_port_v2.hidden_port_1.id}"
  enforce = "true"
  security_group_ids = [
    "${openstack_networking_secgroup_v2.secgroup_1.id}",
    "${openstack_networking_secgroup_v2.secgroup_2.id}",
    "${data.openstack_networking_secgroup_v2.default_2.id}",
  ]
}
`, testAccNetworkingV2PortSecGroupAssociate)
}

func testAccNetworkingV2PortSecGroupAssociateManifestUpdate4() string {
	return fmt.Sprintf(`
%s

resource "openstack_networking_port_secgroup_associate_v2" "port_1" {
  port_id = "${data.openstack_networking_port_v2.hidden_port_1.id}"
  enforce = "true"
  security_group_ids = [
    "${openstack_networking_secgroup_v2.secgroup_2.id}",
  ]
}
`, testAccNetworkingV2PortSecGroupAssociate)
}

func testAccNetworkingV2PortSecGroupAssociateManifestUpdate5() string {
	return fmt.Sprintf(`
%s

data "openstack_networking_port_v2" "port_1" {
  port_id = "${openstack_networking_port_secgroup_associate_v2.port_1.id}"
}

resource "openstack_networking_port_secgroup_associate_v2" "port_1" {
  port_id = "${data.openstack_networking_port_v2.hidden_port_1.id}"
  enforce = "true"
  security_group_ids = []
}
`, testAccNetworkingV2PortSecGroupAssociate)
}

func testAccNetworkingV2PortSecGroupAssociateManifestUpdate6() string {
	return fmt.Sprintf(`
%s

resource "openstack_networking_port_secgroup_associate_v2" "port_1" {
  port_id = "${data.openstack_networking_port_v2.hidden_port_1.id}"
  enforce = "true"
  security_group_ids = [
    "${data.openstack_networking_secgroup_v2.default_2.id}",
  ]
}
`, testAccNetworkingV2PortSecGroupAssociate)
}

func testAccNetworkingV2PortSecGroupAssociateManifestUpdate7() string {
	return fmt.Sprintf(`
%s

resource "openstack_networking_port_secgroup_associate_v2" "port_1" {
  port_id = "${data.openstack_networking_port_v2.hidden_port_1.id}"
  enforce = "false"
  security_group_ids = [
    "${openstack_networking_secgroup_v2.secgroup_1.id}",
  ]
}
`, testAccNetworkingV2PortSecGroupAssociate)
}

func testAccNetworkingV2PortSecGroupAssociateManifestUpdate8() string {
	return fmt.Sprintf(`
%s

resource "openstack_networking_port_secgroup_associate_v2" "port_1" {
  port_id = "${data.openstack_networking_port_v2.hidden_port_1.id}"
  enforce = "false"
  security_group_ids = [
    "${openstack_networking_secgroup_v2.secgroup_1.id}",
    "${openstack_networking_secgroup_v2.secgroup_2.id}",
    "${data.openstack_networking_secgroup_v2.default_2.id}",
  ]
}
`, testAccNetworkingV2PortSecGroupAssociate)
}

func testAccNetworkingV2PortSecGroupAssociateManifestUpdate9() string {
	return fmt.Sprintf(`
%s

resource "openstack_networking_port_secgroup_associate_v2" "port_1" {
  port_id = "${data.openstack_networking_port_v2.hidden_port_1.id}"
  enforce = "false"
  security_group_ids = [
    "${openstack_networking_secgroup_v2.secgroup_2.id}",
  ]
}
`, testAccNetworkingV2PortSecGroupAssociate)
}

func testAccNetworkingV2PortSecGroupAssociateManifestUpdate10() string {
	return fmt.Sprintf(`
%s

resource "openstack_networking_port_secgroup_associate_v2" "port_1" {
  port_id = "${data.openstack_networking_port_v2.hidden_port_1.id}"
  enforce = "false"
  security_group_ids = []
}
`, testAccNetworkingV2PortSecGroupAssociate)
}

func testAccNetworkingV2PortSecGroupAssociateManifestUpdate11() string {
	return fmt.Sprintf(`
%s

resource "openstack_networking_port_secgroup_associate_v2" "port_1" {
  port_id = "${data.openstack_networking_port_v2.hidden_port_1.id}"
  enforce = "false"
  security_group_ids = [
    "${data.openstack_networking_secgroup_v2.default_2.id}",
  ]
}
`, testAccNetworkingV2PortSecGroupAssociate)
}

func testAccNetworkingV2PortSecGroupAssociateManifestUpdate12() string {
	return fmt.Sprintf(`
%s
`, testAccNetworkingV2PortSecGroupAssociate)
}
