package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/monitors"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccLBV2Monitor_basic(t *testing.T) {
	var monitor monitors.Monitor

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV2MonitorDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLbV2MonitorConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2MonitorExists(t.Context(), "openstack_lb_monitor_v2.monitor_1", &monitor),
					resource.TestCheckResourceAttr(
						"openstack_lb_monitor_v2.monitor_1", "admin_state_up", "true"),
					resource.TestCheckResourceAttr(
						"openstack_lb_monitor_v2.monitor_1", "url_path", ""),
					resource.TestCheckResourceAttr(
						"openstack_lb_monitor_v2.monitor_1", "type", "PING"),
					resource.TestCheckResourceAttrPtr(
						"openstack_lb_monitor_v2.monitor_1", "id", &monitor.ID),
				),
			},
			{
				Config: testAccCheckLbV2MonitorConfigUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2MonitorExists(t.Context(), "openstack_lb_monitor_v2.monitor_1", &monitor),
					resource.TestCheckResourceAttr(
						"openstack_lb_monitor_v2.monitor_1", "admin_state_up", "false"),
					resource.TestCheckResourceAttr(
						"openstack_lb_monitor_v2.monitor_1", "url_path", "/host"),
					resource.TestCheckResourceAttr(
						"openstack_lb_monitor_v2.monitor_1", "type", "HTTP"),
					resource.TestCheckResourceAttr(
						"openstack_lb_monitor_v2.monitor_1", "http_method", "GET"),
					resource.TestCheckResourceAttr(
						"openstack_lb_monitor_v2.monitor_1", "expected_codes", "200-205"),
					resource.TestCheckResourceAttrPtr(
						"openstack_lb_monitor_v2.monitor_1", "id", &monitor.ID),
				),
			},
		},
	})
}

func testAccCheckLBV2MonitorDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		lbClient, err := config.LoadBalancerV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack load balancing client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_lb_monitor_v2" {
				continue
			}

			_, err := monitors.Get(ctx, lbClient, rs.Primary.ID).Extract()
			if err == nil {
				return fmt.Errorf("Monitor still exists: %s", rs.Primary.ID)
			}
		}

		return nil
	}
}

func testAccCheckLBV2MonitorExists(ctx context.Context, n string, monitor *monitors.Monitor) resource.TestCheckFunc {
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

		found, err := monitors.Get(ctx, lbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return errors.New("Monitor not found")
		}

		*monitor = *found

		return nil
	}
}

const testAccCheckLbV2MonitorConfig = `
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
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = openstack_lb_listener_v2.listener_1.id

  persistence {
    type        = "APP_COOKIE"
    cookie_name = "testCookie"
  }
}
`

func testAccCheckLbV2MonitorConfigBasic() string {
	return fmt.Sprintf(`
%s

resource "openstack_lb_monitor_v2" "monitor_1" {
  pool_id     = openstack_lb_pool_v2.pool_1.id
  type        = "PING"
  delay       = 20
  timeout     = 10
  max_retries = 5
}
`, testAccCheckLbV2MonitorConfig)
}

func testAccCheckLbV2MonitorConfigUpdate() string {
	return fmt.Sprintf(`
%s

resource "openstack_lb_monitor_v2" "monitor_1" {
  pool_id        = openstack_lb_pool_v2.pool_1.id
  admin_state_up = false
  type           = "HTTP"
  url_path       = "/host"
  http_method    = "GET"
  expected_codes = "200-205"
  delay          = 20
  timeout        = 10
  max_retries    = 5
}
`, testAccCheckLbV2MonitorConfig)
}
