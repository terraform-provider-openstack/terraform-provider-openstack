package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/layer3/routers"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/qos/policies"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccNetworkingV2Router_basic(t *testing.T) {
	var router routers.Router

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2RouterDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2RouterExists(t.Context(), "openstack_networking_router_v2.router_1", &router),
					resource.TestCheckResourceAttr(
						"openstack_networking_router_v2.router_1", "description", "router description"),
				),
			},
			{
				Config: testAccNetworkingV2RouterUpdate,
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
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2RouterDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterUpdateExternalGateway1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2RouterExists(t.Context(), "openstack_networking_router_v2.router_1", &router),
				),
			},
			{
				Config: testAccNetworkingV2RouterUpdateExternalGateway2(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_networking_router_v2.router_1", "external_network_id", osExtGwID),
				),
			},
		},
	})
}

func TestAccNetworkingV2Router_vendor_opts(t *testing.T) {
	var router routers.Router

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2RouterDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterVendorOpts(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2RouterExists(t.Context(), "openstack_networking_router_v2.router_1", &router),
					resource.TestCheckResourceAttr(
						"openstack_networking_router_v2.router_1", "external_network_id", osExtGwID),
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
			// (rule:create_router and rule:create_router:distributed) is disallowed by policy
			testAccPreCheckAdminOnly(t)
			t.Skip("Currently failing in GH-A: Cannot enable DVR + OVN on devstack")
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2RouterDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterVendorOptsNoSnat(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2RouterExists(t.Context(), "openstack_networking_router_v2.router_1", &router),
					resource.TestCheckResourceAttr(
						"openstack_networking_router_v2.router_1", "external_network_id", osExtGwID),
				),
			},
		},
	})
}

func TestAccNetworkingV2Router_extFixedIPs(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2RouterDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterExtFixedIPs(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_networking_router_v2.router_2", "name", "router_2"),
					resource.TestCheckResourceAttr(
						"openstack_networking_router_v2.router_2", "external_fixed_ip.#", "2"),
					resource.TestCheckResourceAttr(
						"openstack_networking_router_v2.router_2", "enable_snat", "true"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Router_extSubnetIDs(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2RouterDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterExtSubnetIDs(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_networking_router_v2.router_2", "name", "router_2"),
					resource.TestCheckResourceAttr(
						"openstack_networking_router_v2.router_2", "external_fixed_ip.#", "1"),
					resource.TestCheckResourceAttr(
						"openstack_networking_router_v2.router_2", "enable_snat", "true"),
				),
			},
		},
	})
}

func TestAccNetworkingV2Router_extQoSPolicy(t *testing.T) {
	var policy policies.Policy

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2RouterDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterExtQoSPolicy(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2QoSPolicyExists(t.Context(),
						"openstack_networking_qos_policy_v2.qos_policy_1", &policy),
					resource.TestCheckResourceAttr(
						"openstack_networking_router_v2.router_1", "name", "router_1"),
					resource.TestCheckResourceAttr(
						"openstack_networking_router_v2.router_1", "external_fixed_ip.#", "2"),
					resource.TestCheckResourceAttr(
						"openstack_networking_router_v2.router_1", "enable_snat", "true"),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_router_v2.router_1", "external_qos_policy_id", &policy.ID),
				),
			},
		},
	})
}

func TestAccNetworkingV2Router_extIPAddress(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2RouterDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2RouterExtIPAddress(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_networking_router_v2.router_1", "name", "router_1"),
					resource.TestCheckResourceAttr(
						"openstack_networking_router_v2.router_1", "external_fixed_ip.#", "1"),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2RouterDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		networkingClient, err := config.NetworkingV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_networking_router_v2" {
				continue
			}

			_, err := routers.Get(ctx, networkingClient, rs.Primary.ID).Extract()
			if err == nil {
				return errors.New("Router still exists")
			}
		}

		return nil
	}
}

func testAccCheckNetworkingV2RouterExists(ctx context.Context, n string, router *routers.Router) resource.TestCheckFunc {
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

		found, err := routers.Get(ctx, networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return errors.New("Router not found")
		}

		*router = *found

		return nil
	}
}

const testAccNetworkingV2RouterBasic = `
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

const testAccNetworkingV2RouterUpdate = `
resource "openstack_networking_router_v2" "router_1" {
  name = "router_2"
  admin_state_up = "true"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`

func testAccNetworkingV2RouterVendorOpts() string {
	return fmt.Sprintf(`
resource "openstack_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
  external_network_id = "%s"
  vendor_options {
    set_router_gateway_after_create = true
  }
}
`, osExtGwID)
}

func testAccNetworkingV2RouterVendorOptsNoSnat() string {
	return fmt.Sprintf(`
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
`, osExtGwID)
}

const testAccNetworkingV2RouterUpdateExternalGateway1 = `
resource "openstack_networking_router_v2" "router_1" {
  name = "router"
  admin_state_up = "true"
}
`

func testAccNetworkingV2RouterUpdateExternalGateway2() string {
	return fmt.Sprintf(`
resource "openstack_networking_router_v2" "router_1" {
  name = "router"
  admin_state_up = "true"
  external_network_id = "%s"
}
`, osExtGwID)
}

func testAccNetworkingV2RouterExtFixedIPs() string {
	return fmt.Sprintf(`
resource "openstack_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
  external_network_id = "%s"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}

resource "openstack_networking_router_v2" "router_2" {
  name = "router_2"
  admin_state_up = "true"
  external_network_id = "%s"

  external_fixed_ip {
    subnet_id = openstack_networking_router_v2.router_1.external_fixed_ip.0.subnet_id
  }

  external_fixed_ip {
    subnet_id = openstack_networking_router_v2.router_1.external_fixed_ip.0.subnet_id
  }

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`, osExtGwID, osExtGwID)
}

func testAccNetworkingV2RouterExtSubnetIDs() string {
	return fmt.Sprintf(`
resource "openstack_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
  external_network_id = "%s"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}

resource "openstack_networking_router_v2" "router_2" {
  name = "router_2"
  admin_state_up = "true"
  external_network_id = "%s"

  external_subnet_ids = [
    "%s", # wrong UUID
    openstack_networking_router_v2.router_1.external_fixed_ip.0.subnet_id,
    "%s", # wrong UUID again
  ]

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`, osExtGwID, osExtGwID, osExtGwID, osExtGwID)
}

func testAccNetworkingV2RouterExtQoSPolicy() string {
	return fmt.Sprintf(`
resource "openstack_networking_qos_policy_v2" "qos_policy_1" {
  name = "qos_policy_1"
}

resource "openstack_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
  external_network_id = "%s"
  external_qos_policy_id = openstack_networking_qos_policy_v2.qos_policy_1.id

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`, osExtGwID)
}

func testAccNetworkingV2RouterExtIPAddress() string {
	return fmt.Sprintf(`
data "openstack_networking_subnet_v2" "subnet_1" {
  name = "public-subnet"
  network_id = "%s"
}

resource "openstack_networking_router_v2" "router_1" {
  name = "router_1"
  admin_state_up = "true"
  external_network_id = "%s"

  external_fixed_ip {
    ip_address = cidrhost(format("%%s/24", data.openstack_networking_subnet_v2.subnet_1.allocation_pools.0.end),-100)
  }

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`, osExtGwID, osExtGwID)
}
