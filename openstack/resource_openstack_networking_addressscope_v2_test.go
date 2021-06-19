package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/addressscopes"
)

func TestAccNetworkingV2AddressScope_basic(t *testing.T) {
	var addressScope addressscopes.AddressScope

	name := acctest.RandomWithPrefix("tf-acc-addrscope")
	newName := acctest.RandomWithPrefix("tf-acc-addrscope")

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2AddressScopeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2AddressScopeBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2AddressScopeExists("openstack_networking_addressscope_v2.addressscope_1", &addressScope),
					resource.TestCheckResourceAttr("openstack_networking_addressscope_v2.addressscope_1", "name", name),
					resource.TestCheckResourceAttr("openstack_networking_addressscope_v2.addressscope_1", "ip_version", "4"),
					resource.TestCheckResourceAttr("openstack_networking_addressscope_v2.addressscope_1", "shared", "false"),
				),
			},
			{
				Config: testAccNetworkingV2AddressScopeBasic(newName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openstack_networking_addressscope_v2.addressscope_1", "name", newName),
					resource.TestCheckResourceAttr("openstack_networking_addressscope_v2.addressscope_1", "ip_version", "4"),
					resource.TestCheckResourceAttr("openstack_networking_addressscope_v2.addressscope_1", "shared", "false"),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2AddressScopeExists(n string, addressScope *addressscopes.AddressScope) resource.TestCheckFunc {
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

		found, err := addressscopes.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Address-scope not found")
		}

		*addressScope = *found

		return nil
	}
}

func testAccCheckNetworkingV2AddressScopeDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.NetworkingV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_networking_addressscope_v2" {
			continue
		}

		_, err := addressscopes.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Address-scope still exists")
		}
	}

	return nil
}

func testAccNetworkingV2AddressScopeBasic(name string) string {
	return fmt.Sprintf(`
resource "openstack_networking_addressscope_v2" "addressscope_1" {
  name       = "%s"
  ip_version = 4
}
`, name)
}
