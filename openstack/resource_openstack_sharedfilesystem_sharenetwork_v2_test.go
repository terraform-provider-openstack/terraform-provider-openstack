package openstack

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/securityservices"
	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/sharenetworks"
)

func TestAccSFSV2Sharenetwork_basic(t *testing.T) {
	var sharenetwork1 sharenetworks.ShareNetwork
	var sharenetwork2 sharenetworks.ShareNetwork

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckSFS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSFSV2SharenetworkDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSFSV2SharenetworkConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSV2SharenetworkExists("openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", &sharenetwork1),
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
			resource.TestStep{
				Config: testAccSFSV2SharenetworkConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSV2SharenetworkExists("openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", &sharenetwork2),
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
					testAccCheckSFSV2SharenetworkNetDiffers(&sharenetwork1, &sharenetwork2),
				),
			},
		},
	})
}

func TestAccSFSV2Sharenetwork_secservice(t *testing.T) {
	var sharenetwork sharenetworks.ShareNetwork

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckSFS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSFSV2SharenetworkDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSFSV2SharenetworkConfig_secservice_1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSV2SharenetworkExists("openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", &sharenetwork),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "name", "test_sharenetwork_secure"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "description", "share the secure love"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "security_service_ids.#", "1"),
					testAccCheckSFSV2SharenetworkSecSvcExists("openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1"),
				),
			},
			resource.TestStep{
				Config: testAccSFSV2SharenetworkConfig_secservice_2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSV2SharenetworkExists("openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", &sharenetwork),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "name", "test_sharenetwork_secure"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "description", "share the secure love"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "security_service_ids.#", "2"),
					testAccCheckSFSV2SharenetworkSecSvcExists("openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1"),
				),
			},
			resource.TestStep{
				Config: testAccSFSV2SharenetworkConfig_secservice_3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSV2SharenetworkExists("openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", &sharenetwork),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "name", "test_sharenetwork_secure"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "description", "share the secure love"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "security_service_ids.#", "1"),
					testAccCheckSFSV2SharenetworkSecSvcExists("openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1"),
				),
			},
			resource.TestStep{
				Config: testAccSFSV2SharenetworkConfig_secservice_4,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSV2SharenetworkExists("openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", &sharenetwork),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "name", "test_sharenetwork"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "description", "share the love"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "security_service_ids.#", "0"),
					testAccCheckSFSV2SharenetworkSecSvcExists("openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1"),
				),
			},
		},
	})
}

func testAccCheckSFSV2SharenetworkDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	sfsClient, err := config.sharedfilesystemV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}
	mvSet, ouErr := setManilaMicroversion(sfsClient)
	if !mvSet && ouErr != nil {
		return ouErr
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

func testAccCheckSFSV2SharenetworkExists(n string, sharenetwork *sharenetworks.ShareNetwork) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		sfsClient, err := config.sharedfilesystemV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
		}
		mvSet, ouErr := setManilaMicroversion(sfsClient)
		if !mvSet && ouErr != nil {
			return ouErr
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

func testAccCheckSFSV2SharenetworkSecSvcExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		sfsClient, err := config.sharedfilesystemV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
		}
		mvSet, ouErr := setManilaMicroversion(sfsClient)
		if !mvSet && ouErr != nil {
			return ouErr
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

		api_security_service_ids := []string{}
		if len(securityServiceList) > 0 {
			api_security_service_ids = resourceSharedfilesystemSharenetworkSecurityServices2IDsV2(&securityServiceList)
		}

		tf_security_service_ids := []string{}
		for k, v := range rs.Primary.Attributes {
			if strings.HasPrefix(k, "security_service_ids.#") {
				continue
			}
			if strings.HasPrefix(k, "security_service_ids.") {
				tf_security_service_ids = append(tf_security_service_ids, v)
			}
		}

		sort.Strings(api_security_service_ids)
		sort.Strings(tf_security_service_ids)

		if !reflect.DeepEqual(api_security_service_ids, tf_security_service_ids) {
			return fmt.Errorf("API and Terraform security service IDs don't correspond: %#v != %#v", api_security_service_ids, tf_security_service_ids)
		}

		return nil
	}
}

func testAccCheckSFSV2SharenetworkNetDiffers(sharenetwork1, sharenetwork2 *sharenetworks.ShareNetwork) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if sharenetwork1.NeutronNetID != sharenetwork2.NeutronNetID && sharenetwork1.NeutronSubnetID != sharenetwork2.NeutronSubnetID {
			return nil
		}
		return fmt.Errorf("Underlying neutron network should differ")
	}
}

const testAccSFSV2SharenetworkConfig_basic = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
  lifecycle {
    ignore_changes = ["tags"]
  }
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_sharedfilesystem_sharenetwork_v2" "sharenetwork_1" {
  name                = "test_sharenetwork"
  description         = "share the love"
  neutron_net_id      = "${openstack_networking_network_v2.network_1.id}"
  neutron_subnet_id   = "${openstack_networking_subnet_v2.subnet_1.id}"
}
`

const testAccSFSV2SharenetworkConfig_update = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
  lifecycle {
    ignore_changes = ["tags"]
  }
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_networking_network_v2" "network_2" {
  name = "network_2"
  admin_state_up = "true"
  lifecycle {
    ignore_changes = ["tags"]
  }
}

resource "openstack_networking_subnet_v2" "subnet_2" {
  name = "subnet_2"
  cidr = "192.168.198.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_sharedfilesystem_sharenetwork_v2" "sharenetwork_1" {
  name                = "test_sharenetwork_new_net"
  description         = ""
  neutron_net_id      = "${openstack_networking_network_v2.network_2.id}"
  neutron_subnet_id   = "${openstack_networking_subnet_v2.subnet_2.id}"
}
`

const testAccSFSV2SharenetworkConfig_secservice_1 = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
  lifecycle {
    ignore_changes = ["tags"]
  }
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

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

resource "openstack_sharedfilesystem_sharenetwork_v2" "sharenetwork_1" {
  name                = "test_sharenetwork_secure"
  description         = "share the secure love"
  neutron_net_id      = "${openstack_networking_network_v2.network_1.id}"
  neutron_subnet_id   = "${openstack_networking_subnet_v2.subnet_1.id}"
  security_service_ids = [
    "${openstack_sharedfilesystem_securityservice_v2.securityservice_1.id}",
  ]
}
`

const testAccSFSV2SharenetworkConfig_secservice_2 = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
  lifecycle {
    ignore_changes = ["tags"]
  }
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

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
`

const testAccSFSV2SharenetworkConfig_secservice_3 = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
  lifecycle {
    ignore_changes = ["tags"]
  }
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

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

resource "openstack_sharedfilesystem_sharenetwork_v2" "sharenetwork_1" {
  name                = "test_sharenetwork_secure"
  description         = "share the secure love"
  neutron_net_id      = "${openstack_networking_network_v2.network_1.id}"
  neutron_subnet_id   = "${openstack_networking_subnet_v2.subnet_1.id}"
  security_service_ids = [
    "${openstack_sharedfilesystem_securityservice_v2.securityservice_2.id}",
  ]
}
`

const testAccSFSV2SharenetworkConfig_secservice_4 = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
  lifecycle {
    ignore_changes = ["tags"]
  }
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

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
  type        = "kerberos"
  server      = "192.168.199.11"
  dns_ip      = "192.168.199.11"
}

resource "openstack_sharedfilesystem_sharenetwork_v2" "sharenetwork_1" {
  name                = "test_sharenetwork"
  description         = "share the love"
  neutron_net_id      = "${openstack_networking_network_v2.network_1.id}"
  neutron_subnet_id   = "${openstack_networking_subnet_v2.subnet_1.id}"
}
`
