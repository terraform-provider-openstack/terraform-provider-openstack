package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceLBV2Monitor_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLbV2MonitorConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.openstack_lb_monitor_v2.monitor_ds_1", "name", "monitor_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_monitor_v2.monitor_ds_1", "pools.#", "1"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_lb_monitor_v2.monitor_ds_1", "pools.0.id"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_monitor_v2.monitor_ds_2", "type", "HTTP"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_monitor_v2.monitor_ds_2", "http_method", "GET"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_monitor_v2.monitor_ds_2", "expected_codes", "200-205"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_monitor_v2.monitor_ds_2", "name", "monitor_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_monitor_v2.monitor_ds_2", "pools.#", "1"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_lb_monitor_v2.monitor_ds_2", "pools.0.id"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_monitor_v2.monitor_ds_2", "url_path", "/host"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_monitor_v2.monitor_ds_2", "type", "HTTP"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_monitor_v2.monitor_ds_2", "http_method", "GET"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_monitor_v2.monitor_ds_2", "expected_codes", "200-205"),
				),
			},
		},
	})
}

const testAccDataSourceLbV2MonitorConfigBasic = `
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
  description = "loadbalancer_1 description"
  loadbalancer_provider = "octavia"
  vip_subnet_id = openstack_networking_subnet_v2.subnet_1.id

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}

resource "openstack_lb_listener_v2" "listener_1" {
  name            = "listener_1"
  description     = "listener_1 description"
  protocol        = "HTTP"
  protocol_port   = 80
  loadbalancer_id = openstack_lb_loadbalancer_v2.loadbalancer_1.id
  tags = [
    "tag1",
    "tag2",
    "tag3",
  ]
}

resource "openstack_lb_pool_v2" "pool_1" {
  name       = "pool_1"
  protocol   = "HTTP"
  lb_method  = "ROUND_ROBIN"
  listener_id = openstack_lb_listener_v2.listener_1.id
  tags = [
	  "tag1",
    "tag2",
    "tag3",
  ]
}

resource "openstack_lb_member_v2" "member_1" {
  pool_id       = openstack_lb_pool_v2.pool_1.id
  name          = "member_1"
  address       = "192.168.199.13"
  protocol_port = 8080
}

resource "openstack_lb_monitor_v2" "monitor_1" {
  pool_id          = openstack_lb_pool_v2.pool_1.id
  name             = "monitor_1"
  delay            = 20
  timeout          = 10
  max_retries      = 3
  max_retries_down = 7
  type             = "HTTP"
  url_path         = "/host"
  http_method      = "GET"
  expected_codes   = "200-205"
}

data "openstack_lb_monitor_v2" "monitor_ds_1" {
  monitor_id = openstack_lb_monitor_v2.monitor_1.id
}

data "openstack_lb_monitor_v2" "monitor_ds_2" {
  name           = "monitor_1"
  pool_id        = openstack_lb_pool_v2.pool_1.id
  url_path       = "/host"
  type           = "HTTP"
  http_method    = "GET"
  expected_codes = "200-205"
  depends_on = [
    openstack_lb_monitor_v2.monitor_1,
  ]
}
`
