package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/gophercloud/gophercloud/openstack/sharedfilesystems/v2/securityservices"
)

func TestAccSFSV2SecurityService_basic(t *testing.T) {
	var securityservice securityservices.SecurityService

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckSFS(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckSFSV2SecurityServiceDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSFSV2SecurityServiceConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSV2SecurityServiceExists("openstack_sharedfilesystem_securityservice_v2.securityservice_1", &securityservice),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_securityservice_v2.securityservice_1", "name", "security"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_securityservice_v2.securityservice_1", "description", "created by terraform"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_securityservice_v2.securityservice_1", "type", "active_directory"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_securityservice_v2.securityservice_1", "server", "192.168.199.10"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_securityservice_v2.securityservice_1", "dns_ip", "192.168.199.10"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_securityservice_v2.securityservice_1", "domain", "example.com"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_securityservice_v2.securityservice_1", "ou", "CN=Computers,DC=example,DC=com"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_securityservice_v2.securityservice_1", "user", "joinDomainUser"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_securityservice_v2.securityservice_1", "password", "s8cret"),
				),
			},
			{
				Config: testAccSFSV2SecurityServiceConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSV2SecurityServiceExists("openstack_sharedfilesystem_securityservice_v2.securityservice_1", &securityservice),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_securityservice_v2.securityservice_1", "name", "security_through_obscurity"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_securityservice_v2.securityservice_1", "description", ""),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_securityservice_v2.securityservice_1", "type", "kerberos"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_securityservice_v2.securityservice_1", "server", "192.168.199.11"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_securityservice_v2.securityservice_1", "dns_ip", "192.168.199.11"),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_securityservice_v2.securityservice_1", "domain", ""),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_securityservice_v2.securityservice_1", "ou", ""),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_securityservice_v2.securityservice_1", "user", ""),
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_securityservice_v2.securityservice_1", "password", ""),
				),
			},
		},
	})
}

func testAccCheckSFSV2SecurityServiceDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	sfsClient, err := config.sharedfilesystemV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_sharedfilesystem_securityservice_v2" {
			continue
		}

		_, err := securityservices.Get(sfsClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Manila securityservice still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckSFSV2SecurityServiceExists(n string, securityservice *securityservices.SecurityService) resource.TestCheckFunc {
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

		found, err := securityservices.Get(sfsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Member not found")
		}

		*securityservice = *found

		return nil
	}
}

const testAccSFSV2SecurityServiceConfig_basic = `
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
`

const testAccSFSV2SecurityServiceConfig_update = `
resource "openstack_sharedfilesystem_securityservice_v2" "securityservice_1" {
  name        = "security_through_obscurity"
  description = ""
  type        = "kerberos"
  server      = "192.168.199.11"
  dns_ip      = "192.168.199.11"
}
`
