package openstack

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/addressgroups"
)

func TestAccNetworkingV2AddressGroup_basic(t *testing.T) {
	var addressGroup addressgroups.AddressGroup

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2AddressGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2AddressGroupBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2AddressGroupExists("openstack_networking_address_group_v2.group_1", &addressGroup),
					resource.TestCheckResourceAttrPtr("openstack_networking_address_group_v2.group_1", "id", &addressGroup.ID),
					resource.TestCheckResourceAttr("openstack_networking_address_group_v2.group_1", "name", "group_1"),
					resource.TestCheckResourceAttr("openstack_networking_address_group_v2.group_1", "description", "test"),
					resource.TestCheckResourceAttr("openstack_networking_address_group_v2.group_1", "addresses.#", "2"),
				),
			},
			{
				Config: testAccNetworkingV2AddressGroupUpdate1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr("openstack_networking_address_group_v2.group_1", "id", &addressGroup.ID),
					resource.TestCheckResourceAttr("openstack_networking_address_group_v2.group_1", "name", ""),
					resource.TestCheckResourceAttr("openstack_networking_address_group_v2.group_1", "description", ""),
					resource.TestCheckResourceAttr("openstack_networking_address_group_v2.group_1", "addresses.#", "1"),
				),
			},
			{
				Config: testAccNetworkingV2AddressGroupUpdate2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPtr("openstack_networking_address_group_v2.group_1", "id", &addressGroup.ID),
					resource.TestCheckResourceAttr("openstack_networking_address_group_v2.group_1", "name", "update"),
					resource.TestCheckResourceAttr("openstack_networking_address_group_v2.group_1", "description", "test update"),
					resource.TestCheckResourceAttr("openstack_networking_address_group_v2.group_1", "addresses.#", "2"),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2AddressGroupDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.NetworkingV2Client(context.TODO(), osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_networking_address_group_v2" {
			continue
		}

		_, err := addressgroups.Get(context.TODO(), networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Security address group still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2AddressGroupExists(n string, ag *addressgroups.AddressGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.NetworkingV2Client(context.TODO(), osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %s", err)
		}

		found, err := addressgroups.Get(context.TODO(), networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Security group not found")
		}

		*ag = *found

		return nil
	}
}

const testAccNetworkingV2AddressGroupBasic = `
resource "openstack_networking_address_group_v2" "group_1" {
  name        = "group_1"
  description = "test"
  addresses = [
    "192.168.0.1/32",
    "192.168.0.2/32",
  ]
}
`

const testAccNetworkingV2AddressGroupUpdate1 = `
resource "openstack_networking_address_group_v2" "group_1" {
  addresses = [
    "192.168.0.2/32",
  ]
}
`

const testAccNetworkingV2AddressGroupUpdate2 = `
resource "openstack_networking_address_group_v2" "group_1" {
  name        = "update"
  description = "test update"
  addresses = [
    "2001:db8::/32",
    "192.168.0.2/32",
  ]
}
`
