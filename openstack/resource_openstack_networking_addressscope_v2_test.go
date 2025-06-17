package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/addressscopes"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
		CheckDestroy:      testAccCheckNetworkingV2AddressScopeDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2AddressScopeBasic(name),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2AddressScopeExists(t.Context(), "openstack_networking_addressscope_v2.addressscope_1", &addressScope),
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

func testAccCheckNetworkingV2AddressScopeExists(ctx context.Context, n string, addressScope *addressscopes.AddressScope) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		networkingClient, err := config.NetworkingV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %w", err)
		}

		found, err := addressscopes.Get(ctx, networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return errors.New("Address-scope not found")
		}

		*addressScope = *found

		return nil
	}
}

func testAccCheckNetworkingV2AddressScopeDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		networkingClient, err := config.NetworkingV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_networking_addressscope_v2" {
				continue
			}

			_, err := addressscopes.Get(ctx, networkingClient, rs.Primary.ID).Extract()
			if err == nil {
				return errors.New("Address-scope still exists")
			}
		}

		return nil
	}
}

func testAccNetworkingV2AddressScopeBasic(name string) string {
	return fmt.Sprintf(`
resource "openstack_networking_addressscope_v2" "addressscope_1" {
  name       = "%s"
  ip_version = 4
}
`, name)
}
