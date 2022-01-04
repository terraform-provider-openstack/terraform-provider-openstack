package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/listeners"
)

func TestAccLBV2Listener_basic(t *testing.T) {
	var listener listeners.Listener

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV2ListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbV2ListenerConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2ListenerExists("openstack_lb_listener_v2.listener_1", &listener),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "connection_limit", "-1"),
				),
			},
			{
				Config: testAccLbV2ListenerConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "name", "listener_1_updated"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "connection_limit", "100"),
				),
			},
		},
	})
}

func TestAccLBV2Listener_octavia(t *testing.T) {
	var listener listeners.Listener

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
			testAccPreCheckUseOctavia(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV2ListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbV2ListenerConfigOctavia,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2ListenerExists("openstack_lb_listener_v2.listener_1", &listener),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "connection_limit", "5"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "timeout_client_data", "1000"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "timeout_member_connect", "2000"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "timeout_member_data", "3000"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "timeout_tcp_inspect", "4000"),
				),
			},
			{
				Config: testAccLbV2ListenerConfigOctaviaUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "name", "listener_1_updated"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "connection_limit", "100"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "timeout_client_data", "4000"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "timeout_member_connect", "3000"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "timeout_member_data", "2000"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "timeout_tcp_inspect", "1000"),
				),
			},
		},
	})
}

func TestAccLBV2Listener_octavia_udp(t *testing.T) {
	var listener listeners.Listener

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
			testAccPreCheckUseOctavia(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV2ListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbV2ListenerConfigOctaviaUDP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2ListenerExists("openstack_lb_listener_v2.listener_1", &listener),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "protocol", "UDP"),
				),
			},
		},
	})
}

func TestAccLBV2ListenerConfig_octavia_insert_headers(t *testing.T) {
	var listener listeners.Listener

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
			testAccPreCheckUseOctavia(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV2ListenerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbV2ListenerConfigOctaviaInsertHeaders1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2ListenerExists("openstack_lb_listener_v2.listener_1", &listener),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "insert_headers.X-Forwarded-For", "true"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "insert_headers.X-Forwarded-Port", "false"),
				),
			},
			{
				Config: testAccLbV2ListenerConfigOctaviaInsertHeaders2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2ListenerExists("openstack_lb_listener_v2.listener_1", &listener),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "insert_headers.X-Forwarded-For", "false"),
					resource.TestCheckResourceAttr(
						"openstack_lb_listener_v2.listener_1", "insert_headers.X-Forwarded-Port", "true"),
				),
			},
			{
				Config: testAccLbV2ListenerConfigOctavia,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2ListenerExists("openstack_lb_listener_v2.listener_1", &listener),
					resource.TestCheckNoResourceAttr(
						"openstack_lb_listener_v2.listener_1", "insert_headers.X-Forwarded-For"),
					resource.TestCheckNoResourceAttr(
						"openstack_lb_listener_v2.listener_1", "insert_headers.X-Forwarded-Port"),
				),
			},
		},
	})
}

func testAccCheckLBV2ListenerDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	lbClient, err := chooseLBV2AccTestClient(config, osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack load balancing client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_lb_listener_v2" {
			continue
		}

		_, err := listeners.Get(lbClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Listener still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV2ListenerExists(n string, listener *listeners.Listener) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		lbClient, err := chooseLBV2AccTestClient(config, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack load balancing client: %s", err)
		}

		found, err := listeners.Get(lbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Member not found")
		}

		*listener = *found

		return nil
	}
}

const testAccLbV2ListenerConfigBasic = `
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
  name            = "pool_1"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = "${openstack_lb_loadbalancer_v2.loadbalancer_1.id}"
}

resource "openstack_lb_listener_v2" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  default_pool_id = "${openstack_lb_pool_v2.pool_1.id}"
  loadbalancer_id = "${openstack_lb_loadbalancer_v2.loadbalancer_1.id}"

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const testAccLbV2ListenerConfigUpdate = `
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
  name = "listener_1_updated"
  protocol = "HTTP"
  protocol_port = 8080
  connection_limit = 100
  admin_state_up = "true"
  loadbalancer_id = "${openstack_lb_loadbalancer_v2.loadbalancer_1.id}"

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const testAccLbV2ListenerConfigOctavia = `
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
  name            = "pool_1"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = "${openstack_lb_loadbalancer_v2.loadbalancer_1.id}"
}

resource "openstack_lb_listener_v2" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  connection_limit = 5
  timeout_client_data = 1000
  timeout_member_connect = 2000
  timeout_member_data = 3000
  timeout_tcp_inspect = 4000
  default_pool_id = "${openstack_lb_pool_v2.pool_1.id}"
  loadbalancer_id = "${openstack_lb_loadbalancer_v2.loadbalancer_1.id}"

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const testAccLbV2ListenerConfigOctaviaUpdate = `
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
  name            = "pool_1"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = "${openstack_lb_loadbalancer_v2.loadbalancer_1.id}"
}

resource "openstack_lb_listener_v2" "listener_1" {
  name = "listener_1_updated"
  protocol = "HTTP"
  protocol_port = 8080
  connection_limit = 100
  timeout_client_data = 4000
  timeout_member_connect = 3000
  timeout_member_data = 2000
  timeout_tcp_inspect = 1000
  admin_state_up = "true"
  default_pool_id = "${openstack_lb_pool_v2.pool_1.id}"
  loadbalancer_id = "${openstack_lb_loadbalancer_v2.loadbalancer_1.id}"

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const testAccLbV2ListenerConfigOctaviaUDP = `
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
  name            = "pool_1"
  protocol        = "UDP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = "${openstack_lb_loadbalancer_v2.loadbalancer_1.id}"
}

resource "openstack_lb_listener_v2" "listener_1" {
  name = "listener_1"
  protocol = "UDP"
  protocol_port = 53
  default_pool_id = "${openstack_lb_pool_v2.pool_1.id}"
  loadbalancer_id = "${openstack_lb_loadbalancer_v2.loadbalancer_1.id}"

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const testAccLbV2ListenerConfigOctaviaInsertHeaders1 = `
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
  name            = "pool_1"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = "${openstack_lb_loadbalancer_v2.loadbalancer_1.id}"
}

resource "openstack_lb_listener_v2" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  default_pool_id = "${openstack_lb_pool_v2.pool_1.id}"
  loadbalancer_id = "${openstack_lb_loadbalancer_v2.loadbalancer_1.id}"

  insert_headers = {
    X-Forwarded-For = "true"
    X-Forwarded-Port = "false"
  }

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const testAccLbV2ListenerConfigOctaviaInsertHeaders2 = `
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
  name            = "pool_1"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = "${openstack_lb_loadbalancer_v2.loadbalancer_1.id}"
}

resource "openstack_lb_listener_v2" "listener_1" {
  name = "listener_1"
  protocol = "HTTP"
  protocol_port = 8080
  default_pool_id = "${openstack_lb_pool_v2.pool_1.id}"
  loadbalancer_id = "${openstack_lb_loadbalancer_v2.loadbalancer_1.id}"

  insert_headers = {
    X-Forwarded-For = "false"
    X-Forwarded-Port = "true"
  }

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`
