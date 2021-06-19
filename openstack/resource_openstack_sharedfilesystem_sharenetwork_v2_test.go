package openstack

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/securityservices"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/sharenetworks"
)

func TestAccSFSV2ShareNetwork_basic(t *testing.T) {
	var sharenetwork1 sharenetworks.ShareNetwork
	var sharenetwork2 sharenetworks.ShareNetwork

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckSFS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSFSV2ShareNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSV2ShareNetworkConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSV2ShareNetworkExists("openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", &sharenetwork1),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "name", "test_sharenetwork"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "description", "share the love"),
					resource.TestMatchResourceAttr(
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "neutron_net_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
					resource.TestMatchResourceAttr(
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "neutron_subnet_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
				),
			},
			{
				Config: testAccSFSV2ShareNetworkConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSV2ShareNetworkExists("openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", &sharenetwork2),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "name", "test_sharenetwork_new_net"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "description", ""),
					resource.TestMatchResourceAttr(
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "neutron_net_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
					resource.TestMatchResourceAttr(
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "neutron_subnet_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
					testAccCheckSFSV2ShareNetworkNetDiffers(&sharenetwork1, &sharenetwork2),
				),
			},
		},
	})
}

func TestAccSFSV2ShareNetwork_secservice(t *testing.T) {
	var sharenetwork sharenetworks.ShareNetwork

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckSFS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSFSV2ShareNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSV2ShareNetworkConfigSecService1(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSV2ShareNetworkExists("openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", &sharenetwork),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "name", "test_sharenetwork_secure"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "description", "share the secure love"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "security_service_ids.#", "1"),
					testAccCheckSFSV2ShareNetworkSecSvcExists("openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1"),
				),
			},
			{
				Config: testAccSFSV2ShareNetworkConfigSecService2(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSV2ShareNetworkExists("openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", &sharenetwork),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "name", "test_sharenetwork_secure"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "description", "share the secure love"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "security_service_ids.#", "2"),
					testAccCheckSFSV2ShareNetworkSecSvcExists("openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1"),
				),
			},
			{
				Config: testAccSFSV2ShareNetworkConfigSecService3(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSV2ShareNetworkExists("openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", &sharenetwork),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "name", "test_sharenetwork_secure"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "description", "share the secure love"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "security_service_ids.#", "1"),
					testAccCheckSFSV2ShareNetworkSecSvcExists("openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1"),
				),
			},
			{
				Config: testAccSFSV2ShareNetworkConfigSecService4(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSV2ShareNetworkExists("openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", &sharenetwork),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "name", "test_sharenetwork"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "description", "share the love"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "security_service_ids.#", "0"),
					testAccCheckSFSV2ShareNetworkSecSvcExists("openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1"),
				),
			},
		},
	})
}

func testAccCheckSFSV2ShareNetworkDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	sfsClient, err := config.SharedfilesystemV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_sharedfilesystem_securityservice_v2" {
			continue
		}

		_, err := sharenetworks.Get(sfsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Manila sharenetwork still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckSFSV2ShareNetworkExists(n string, sharenetwork *sharenetworks.ShareNetwork) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		sfsClient, err := config.SharedfilesystemV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
		}

		found, err := sharenetworks.Get(sfsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Member not found")
		}

		*sharenetwork = *found

		return nil
	}
}

func testAccCheckSFSV2ShareNetworkSecSvcExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		sfsClient, err := config.SharedfilesystemV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
		}

		securityServiceListOpts := securityservices.ListOpts{ShareNetworkID: rs.Primary.ID}
		securityServicePages, err := securityservices.List(sfsClient, securityServiceListOpts).AllPages()
		if err != nil {
			return err
		}
		securityServiceList, err := securityservices.ExtractSecurityServices(securityServicePages)
		if err != nil {
			return err
		}

		apiSecurityServiceIDs := resourceSharedFilesystemShareNetworkV2SecSvcToArray(&securityServiceList)

		var tfSecurityServiceIDs []string
		for k, v := range rs.Primary.Attributes {
			if strings.HasPrefix(k, "security_service_ids.#") {
				continue
			}
			if strings.HasPrefix(k, "security_service_ids.") {
				tfSecurityServiceIDs = append(tfSecurityServiceIDs, v)
			}
		}

		sort.Strings(apiSecurityServiceIDs)
		sort.Strings(tfSecurityServiceIDs)

		if !reflect.DeepEqual(apiSecurityServiceIDs, tfSecurityServiceIDs) {
			return fmt.Errorf("API and Terraform security service IDs don't correspond: %#v != %#v", apiSecurityServiceIDs, tfSecurityServiceIDs)
		}

		return nil
	}
}

func testAccCheckSFSV2ShareNetworkNetDiffers(sharenetwork1, sharenetwork2 *sharenetworks.ShareNetwork) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if sharenetwork1.NeutronNetID != sharenetwork2.NeutronNetID && sharenetwork1.NeutronSubnetID != sharenetwork2.NeutronSubnetID {
			return nil
		}
		return fmt.Errorf("Underlying neutron network should differ")
	}
}

const testAccSFSV2ShareNetworkConfig = `
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
`

func testAccSFSV2ShareNetworkConfigBasic() string {
	return fmt.Sprintf(`
%s

resource "openstack_sharedfilesystem_sharenetwork_v2" "sharenetwork_1" {
  name                = "test_sharenetwork"
  description         = "share the love"
  neutron_net_id      = "${openstack_networking_network_v2.network_1.id}"
  neutron_subnet_id   = "${openstack_networking_subnet_v2.subnet_1.id}"
}
`, testAccSFSV2ShareNetworkConfig)
}

func testAccSFSV2ShareNetworkConfigUpdate() string {
	return fmt.Sprintf(`
%s

resource "openstack_networking_network_v2" "network_2" {
  name = "network_2"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_2" {
  name = "subnet_2"
  cidr = "192.168.198.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_2.id}"
}

resource "openstack_sharedfilesystem_sharenetwork_v2" "sharenetwork_1" {
  name                = "test_sharenetwork_new_net"
  description         = ""
  neutron_net_id      = "${openstack_networking_network_v2.network_2.id}"
  neutron_subnet_id   = "${openstack_networking_subnet_v2.subnet_2.id}"
}
`, testAccSFSV2ShareNetworkConfig)
}

const testAccSFSV2ShareNetworkConfigSecService = `
resource "openstack_sharedfilesystem_securityservice_v2" "securityservice_1" {
  name        = "security"
  description = "created by terraform"
  type        = "active_directory"
  server      = "192.168.199.10"
  dns_ip      = "192.168.199.10"
  domain      = "example.com"
  ou          = "CN=Computers,DC=example,DC=com"
  user        = "joinDomainUser"
  password    = "s8cret"
}

resource "openstack_sharedfilesystem_securityservice_v2" "securityservice_2" {
  name        = "security_through_obscurity"
  description = ""
  type        = "kerberos"
  server      = "192.168.199.11"
  dns_ip      = "192.168.199.11"
}
`

func testAccSFSV2ShareNetworkConfigSecService1() string {
	return fmt.Sprintf(`
%s

%s

resource "openstack_sharedfilesystem_sharenetwork_v2" "sharenetwork_1" {
  name                = "test_sharenetwork_secure"
  description         = "share the secure love"
  neutron_net_id      = "${openstack_networking_network_v2.network_1.id}"
  neutron_subnet_id   = "${openstack_networking_subnet_v2.subnet_1.id}"
  security_service_ids = [
    "${openstack_sharedfilesystem_securityservice_v2.securityservice_1.id}",
  ]
}
`, testAccSFSV2ShareNetworkConfig, testAccSFSV2ShareNetworkConfigSecService)
}

func testAccSFSV2ShareNetworkConfigSecService2() string {
	return fmt.Sprintf(`
%s

%s

resource "openstack_sharedfilesystem_sharenetwork_v2" "sharenetwork_1" {
  name                = "test_sharenetwork_secure"
  description         = "share the secure love"
  neutron_net_id      = "${openstack_networking_network_v2.network_1.id}"
  neutron_subnet_id   = "${openstack_networking_subnet_v2.subnet_1.id}"
  security_service_ids = [
    "${openstack_sharedfilesystem_securityservice_v2.securityservice_1.id}",
    "${openstack_sharedfilesystem_securityservice_v2.securityservice_2.id}",
  ]
}
`, testAccSFSV2ShareNetworkConfig, testAccSFSV2ShareNetworkConfigSecService)
}

func testAccSFSV2ShareNetworkConfigSecService3() string {
	return fmt.Sprintf(`
%s

%s

resource "openstack_sharedfilesystem_sharenetwork_v2" "sharenetwork_1" {
  name                = "test_sharenetwork_secure"
  description         = "share the secure love"
  neutron_net_id      = "${openstack_networking_network_v2.network_1.id}"
  neutron_subnet_id   = "${openstack_networking_subnet_v2.subnet_1.id}"
  security_service_ids = [
    "${openstack_sharedfilesystem_securityservice_v2.securityservice_2.id}",
  ]
}
`, testAccSFSV2ShareNetworkConfig, testAccSFSV2ShareNetworkConfigSecService)
}

func testAccSFSV2ShareNetworkConfigSecService4() string {
	return fmt.Sprintf(`
%s

%s

resource "openstack_sharedfilesystem_sharenetwork_v2" "sharenetwork_1" {
  name                = "test_sharenetwork"
  description         = "share the love"
  neutron_net_id      = "${openstack_networking_network_v2.network_1.id}"
  neutron_subnet_id   = "${openstack_networking_subnet_v2.subnet_1.id}"
}
`, testAccSFSV2ShareNetworkConfig, testAccSFSV2ShareNetworkConfigSecService)
}
