package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/securityservices"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccSFSV2SecurityService_basic(t *testing.T) {
	var securityservice securityservices.SecurityService

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckSFS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSFSV2SecurityServiceDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccSFSV2SecurityServiceConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSV2SecurityServiceExists(t.Context(), "openstack_sharedfilesystem_securityservice_v2.securityservice_1", &securityservice),
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
				Config: testAccSFSV2SecurityServiceConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSFSV2SecurityServiceExists(t.Context(), "openstack_sharedfilesystem_securityservice_v2.securityservice_1", &securityservice),
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

func TestAccSFSV2SecurityService_EndpointCheck(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
			testAccPreCheckSFS(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckSFSV2SecurityServiceDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				// register the volumev2 service and endpoint
				Config: testAccSFSV2SecurityServiceConfigUpdateEndpointCheck,
			},
			{
				// test endpoint locator to pick up sharev2
				Config: testAccSFSV2SecurityServiceConfigUpdateEndpointCheck + testAccSFSV2SecurityServiceConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_sharedfilesystem_securityservice_v2.securityservice_1", "name", "security"),
				),
			},
		},
	})
}

func testAccCheckSFSV2SecurityServiceDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		sfsClient, err := config.SharedfilesystemV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_sharedfilesystem_securityservice_v2" {
				continue
			}

			_, err := securityservices.Get(ctx, sfsClient, rs.Primary.ID).Extract()
			if err == nil {
				return fmt.Errorf("Manila securityservice still exists: %s", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccCheckSFSV2SecurityServiceExists(ctx context.Context, n string, securityservice *securityservices.SecurityService) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		sfsClient, err := config.SharedfilesystemV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack sharedfilesystem client: %w", err)
		}

		found, err := securityservices.Get(ctx, sfsClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return errors.New("Member not found")
		}

		*securityservice = *found

		return nil
	}
}

const testAccSFSV2SecurityServiceConfigBasic = `
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

const testAccSFSV2SecurityServiceConfigUpdate = `
resource "openstack_sharedfilesystem_securityservice_v2" "securityservice_1" {
  name        = "security_through_obscurity"
  description = ""
  type        = "kerberos"
  server      = "192.168.199.11"
  dns_ip      = "192.168.199.11"
}
`

const testAccSFSV2SecurityServiceConfigUpdateEndpointCheck = `
resource "openstack_identity_service_v3" "service_1" {
  name = "manilav2"
  type = "sharev2"
}

resource "openstack_identity_endpoint_v3" "endpoint_1" {
  name            = "sharev2"
  service_id      = openstack_identity_service_v3.service_1.id
  endpoint_region = openstack_identity_service_v3.service_1.region
  url             = "http://my-endpoint"
}
`
