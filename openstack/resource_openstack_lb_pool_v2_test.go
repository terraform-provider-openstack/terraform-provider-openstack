package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/pools"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccLBV2Pool_basic(t *testing.T) {
	var pool pools.Pool

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV2PoolDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: TestAccLbV2PoolConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2PoolExists(t.Context(), "openstack_lb_pool_v2.pool_1", &pool),
					testAccCheckLBV2PoolHasTag(t.Context(), "openstack_lb_pool_v2.pool_1", "foo"),
					testAccCheckLBV2PoolTagCount(t.Context(), "openstack_lb_pool_v2.pool_1", 1),
					resource.TestCheckResourceAttr("openstack_lb_pool_v2.pool_1", "persistence.#", "1"),
					resource.TestCheckResourceAttr("openstack_lb_pool_v2.pool_1", "persistence.0.type", "APP_COOKIE"),
					resource.TestCheckResourceAttr("openstack_lb_pool_v2.pool_1", "persistence.0.cookie_name", "testCookie"),
					resource.TestCheckResourceAttr("openstack_lb_pool_v2.pool_1", "tls_enabled", "true"),
					resource.TestCheckResourceAttr("openstack_lb_pool_v2.pool_1", "tls_versions.#", "1"),
					resource.TestCheckResourceAttr("openstack_lb_pool_v2.pool_1", "alpn_protocols.#", "1"),
					resource.TestCheckResourceAttr("openstack_lb_pool_v2.pool_1", "tls_ciphers", "TLS13-CHACHA20-POLY1305-SHA256"),
				),
			},
			{
				Config: TestAccLbV2PoolConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openstack_lb_pool_v2.pool_1", "name", "pool_1_updated"),
					testAccCheckLBV2PoolHasTag(t.Context(), "openstack_lb_pool_v2.pool_1", "bar"),
					testAccCheckLBV2PoolTagCount(t.Context(), "openstack_lb_pool_v2.pool_1", 1),
					resource.TestCheckResourceAttr("openstack_lb_pool_v2.pool_1", "persistence.#", "0"),
					resource.TestCheckResourceAttr("openstack_lb_pool_v2.pool_1", "tls_enabled", "false"),
					// tls_versions reset to Octavia default value
					// resource.TestCheckResourceAttr("openstack_lb_pool_v2.pool_1", "tls_versions.#", "0"),
					resource.TestCheckResourceAttr("openstack_lb_pool_v2.pool_1", "alpn_protocols.#", "0"),
					// even though we unset this value, it should still be present in the state
					// because it is a computed field
					resource.TestCheckResourceAttr("openstack_lb_pool_v2.pool_1", "tls_ciphers", "TLS13-CHACHA20-POLY1305-SHA256"),
				),
			},
		},
	})
}

func testAccCheckLBV2PoolHasTag(ctx context.Context, n, tag string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		lbClient, err := config.LoadBalancerV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack load balancing client: %w", err)
		}

		found, err := pools.Get(ctx, lbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return errors.New("Pool not found")
		}

		for _, v := range found.Tags {
			if tag == v {
				return nil
			}
		}

		return fmt.Errorf("Tag not found: %s", tag)
	}
}

func testAccCheckLBV2PoolTagCount(ctx context.Context, n string, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		lbClient, err := config.LoadBalancerV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack load balancing client: %w", err)
		}

		found, err := pools.Get(ctx, lbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return errors.New("Pool not found")
		}

		if len(found.Tags) != expected {
			return fmt.Errorf("Expecting %d tags, found %d", expected, len(found.Tags))
		}

		return nil
	}
}

func TestAccLBV2Pool_octavia_udp(t *testing.T) {
	var pool pools.Pool

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV2PoolDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: TestAccLbV2PoolConfigOctaviaUDP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2PoolExists(t.Context(), "openstack_lb_pool_v2.pool_1", &pool),
					resource.TestCheckResourceAttr("openstack_lb_pool_v2.pool_1", "protocol", "UDP"),
				),
			},
		},
	})
}

func testAccCheckLBV2PoolDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		lbClient, err := config.LoadBalancerV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack load balancing client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_lb_pool_v2" {
				continue
			}

			_, err := pools.Get(ctx, lbClient, rs.Primary.ID).Extract()
			if err == nil {
				return fmt.Errorf("Pool still exists: %s", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccCheckLBV2PoolExists(ctx context.Context, n string, pool *pools.Pool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		lbClient, err := config.LoadBalancerV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack load balancing client: %w", err)
		}

		found, err := pools.Get(ctx, lbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return errors.New("Member not found")
		}

		*pool = *found

		return nil
	}
}

const TestAccLbV2PoolConfigBasic = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = openstack_networking_network_v2.network_1.id
}

resource "openstack_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = openstack_networking_subnet_v2.subnet_1.id

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}

resource "openstack_lb_listener_v2" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  loadbalancer_id = openstack_lb_loadbalancer_v2.loadbalancer_1.id
}

resource "openstack_lb_pool_v2" "pool_1" {
  name = "pool_1"
  protocol = "HTTP"
  lb_method = "ROUND_ROBIN"
  listener_id = openstack_lb_listener_v2.listener_1.id
  tags = ["foo"]

  tls_enabled    = true
  tls_versions   = ["TLSv1.2"]
  alpn_protocols = ["http/1.1"]
  tls_ciphers    = "TLS13-CHACHA20-POLY1305-SHA256"

  persistence {
    type        = "APP_COOKIE"
    cookie_name = "testCookie"
  }

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const TestAccLbV2PoolConfigUpdate = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = openstack_networking_network_v2.network_1.id
}

resource "openstack_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = openstack_networking_subnet_v2.subnet_1.id

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}

resource "openstack_lb_listener_v2" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  loadbalancer_id = openstack_lb_loadbalancer_v2.loadbalancer_1.id
}

resource "openstack_lb_pool_v2" "pool_1" {
  name = "pool_1_updated"
  protocol = "HTTP"
  lb_method = "LEAST_CONNECTIONS"
  admin_state_up = "true"
  listener_id = openstack_lb_listener_v2.listener_1.id
  tags = ["bar"]

  alpn_protocols = []

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const TestAccLbV2PoolConfigOctaviaUDP = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = openstack_networking_network_v2.network_1.id
}

resource "openstack_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = openstack_networking_subnet_v2.subnet_1.id

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}

resource "openstack_lb_pool_v2" "pool_1" {
  name = "pool_1"
  protocol = "UDP"
  lb_method = "ROUND_ROBIN"
  loadbalancer_id = openstack_lb_loadbalancer_v2.loadbalancer_1.id

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`
