package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccSFSV2ShareNetworkDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckSFS(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSV2ShareNetworkDataSourceBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "id",
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "id"),
					resource.TestCheckResourceAttr(
						"data.openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "security_service_ids.#", "2"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "neutron_net_id",
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "neutron_net_id"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "ip_version",
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1", "ip_version"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_2", "id",
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_2", "id"),
					resource.TestCheckResourceAttr(
						"data.openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_2", "security_service_ids.#", "1"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_2", "neutron_net_id",
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_2", "neutron_net_id"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_2", "ip_version",
						"openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_2", "ip_version"),
				),
			},
		},
	})
}

const testAccSFSV2ShareNetworkDataSourceBasic = `
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

resource "openstack_sharedfilesystem_sharenetwork_v2" "sharenetwork_2" {
  name                = "test_sharenetwork_secure"
  description         = "share the less secure love"
  neutron_net_id      = "${openstack_networking_network_v2.network_1.id}"
  neutron_subnet_id   = "${openstack_networking_subnet_v2.subnet_1.id}"
  security_service_ids = [
    "${openstack_sharedfilesystem_securityservice_v2.securityservice_1.id}",
  ]
}

data "openstack_sharedfilesystem_sharenetwork_v2" "sharenetwork_1" {
  name                = "${openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1.name}"
  security_service_id = "${openstack_sharedfilesystem_securityservice_v2.securityservice_2.id}"
  ip_version          = "${openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_1.ip_version}"
}

data "openstack_sharedfilesystem_sharenetwork_v2" "sharenetwork_2" {
  name                = "test_sharenetwork_secure"
  description         = "${openstack_sharedfilesystem_sharenetwork_v2.sharenetwork_2.description}"
  security_service_id = "${openstack_sharedfilesystem_securityservice_v2.securityservice_1.id}"
}
`
