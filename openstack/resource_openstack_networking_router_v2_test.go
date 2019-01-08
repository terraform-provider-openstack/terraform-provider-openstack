package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/layer3/routers"
)

func TestAccNetworkingV2Router_basic(t *testing.T) {
	var router routers.Router

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2RouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Router_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2RouterExists("openstack_networking_router_v2.router_1", &router),
					resource.TestCheckResourceAttr(
						"openstack_networking_router_v2.router_1", "description", "router description"),
				),
			},
			{
				Config: testAccNetworkingV2Router_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_networking_router_v2.router_1", "name", "router_2"),
					resource.TestCheckResourceAttr(
						"openstack_networking_router_v2.router_1", "description", ""),
				),
			},
		},
	})
}

func TestAccNetworkingV2Router_updateExternalGateway(t *testing.T) {
	var router routers.Router

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2RouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Router_updateExternalGateway1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2RouterExists("openstack_networking_router_v2.router_1", &router),
				),
			},
			{
				Config: testAccNetworkingV2Router_updateExternalGateway2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_networking_router_v2.router_1", "external_network_id", OS_EXTGW_ID),
				),
			},
		},
	})
}

func TestAccNetworkingV2Router_vendor_opts(t *testing.T) {
	var router routers.Router

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2RouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Router_vendor_opts,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2RouterExists("openstack_networking_router_v2.router_1", &router),
					resource.TestCheckResourceAttr(
						"openstack_networking_router_v2.router_1", "external_gateway", OS_EXTGW_ID),
				),
			},
		},
	})
}

func TestAccNetworkingV2Router_vendor_opts_no_snat(t *testing.T) {
	var router routers.Router

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkingV2RouterDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2Router_vendor_opts_no_snat,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2RouterExists("openstack_networking_router_v2.router_1", &router),
					resource.TestCheckResourceAttr(
						"openstack_networking_router_v2.router_1", "external_gateway", OS_EXTGW_ID),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2RouterDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_networking_router_v2" {
			continue
		}

		_, err := routers.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Router still exists")
		}
	}

	return nil
}

func testAccCheckNetworkingV2RouterExists(n string, router *routers.Router) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %s", err)
		}

		found, err := routers.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Router not found")
		}

		*router = *found

		return nil
	}
}

const testAccNetworkingV2Router_basic = `
resource "openstack_networking_router_v2" "router_1" {
  name = "router_1"
  description = "router description"
  admin_state_up = "true"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`

const testAccNetworkingV2Router_update = `
resource "openstack_networking_router_v2" "router_1" {
  name = "router_2"
  admin_state_up = "true"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`

var testAccNetworkingV2Router_vendor_opts = fmt.Sprintf(`
resource "openstack_networking_router_v2" "router_1" {
	name = "router_1"
	admin_state_up = "true"
	external_network_id = "%s"
	vendor_options {
		set_router_gateway_after_create = true
	}
}
`, OS_EXTGW_ID)

var testAccNetworkingV2Router_vendor_opts_no_snat = fmt.Sprintf(`
resource "openstack_networking_router_v2" "router_1" {
        name = "router_1"
        admin_state_up = "true"
        distributed = "false"
        external_network_id = "%s"
        enable_snat = "false"
        vendor_options {
                set_router_gateway_after_create = true
        }
}
`, OS_EXTGW_ID)

const testAccNetworkingV2Router_updateExternalGateway1 = `
resource "openstack_networking_router_v2" "router_1" {
	name = "router"
	admin_state_up = "true"
}
`

var testAccNetworkingV2Router_updateExternalGateway2 = fmt.Sprintf(`
resource "openstack_networking_router_v2" "router_1" {
	name = "router"
	admin_state_up = "true"
	external_network_id = "%s"
}
`, OS_EXTGW_ID)
