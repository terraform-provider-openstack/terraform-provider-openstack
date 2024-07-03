package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/pools"
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
		CheckDestroy:      testAccCheckLBV2PoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccLbV2PoolConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2PoolExists("openstack_lb_pool_v2.pool_1", &pool),
					testAccCheckLBV2PoolHasTag("openstack_lb_pool_v2.pool_1", "foo"),
					testAccCheckLBV2PoolTagCount("openstack_lb_pool_v2.pool_1", 1),
				),
			},
			{
				Config: TestAccLbV2PoolConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openstack_lb_pool_v2.pool_1", "name", "pool_1_updated"),
					testAccCheckLBV2PoolHasTag("openstack_lb_pool_v2.pool_1", "bar"),
					testAccCheckLBV2PoolTagCount("openstack_lb_pool_v2.pool_1", 1),
				),
			},
		},
	})
}

func testAccCheckLBV2PoolHasTag(n, tag string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		lbClient, err := config.LoadBalancerV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack load balancing client: %s", err)
		}

		found, err := pools.Get(lbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Pool not found")
		}

		for _, v := range found.Tags {
			if tag == v {
				return nil
			}
		}

		return fmt.Errorf("Tag not found: %s", tag)
	}
}

func testAccCheckLBV2PoolTagCount(n string, expected int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		lbClient, err := config.LoadBalancerV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack load balancing client: %s", err)
		}

		found, err := pools.Get(lbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Pool not found")
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
		CheckDestroy:      testAccCheckLBV2PoolDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccLbV2PoolConfigOctaviaUDP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2PoolExists("openstack_lb_pool_v2.pool_1", &pool),
					resource.TestCheckResourceAttr("openstack_lb_pool_v2.pool_1", "protocol", "UDP"),
				),
			},
		},
	})
}

func testAccCheckLBV2PoolDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	lbClient, err := config.LoadBalancerV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack load balancing client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_lb_pool_v2" {
			continue
		}

		_, err := pools.Get(lbClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Pool still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV2PoolExists(n string, pool *pools.Pool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		lbClient, err := config.LoadBalancerV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack load balancing client: %s", err)
		}

		found, err := pools.Get(lbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Member not found")
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
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"

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
  loadbalancer_id = "${openstack_lb_loadbalancer_v2.loadbalancer_1.id}"
}

resource "openstack_lb_pool_v2" "pool_1" {
  name = "pool_1"
  protocol = "HTTP"
  lb_method = "ROUND_ROBIN"
  listener_id = "${openstack_lb_listener_v2.listener_1.id}"
  tags = ["foo"]

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
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"

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
  loadbalancer_id = "${openstack_lb_loadbalancer_v2.loadbalancer_1.id}"
}

resource "openstack_lb_pool_v2" "pool_1" {
  name = "pool_1_updated"
  protocol = "HTTP"
  lb_method = "LEAST_CONNECTIONS"
  admin_state_up = "true"
  listener_id = "${openstack_lb_listener_v2.listener_1.id}"
  tags = ["bar"]

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
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"

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
  loadbalancer_id = "${openstack_lb_loadbalancer_v2.loadbalancer_1.id}"

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`
