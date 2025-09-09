package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceLBV2Listener_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLbV2ListenerConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.openstack_lb_listener_v2.ls_ds", "name", "listener_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_listener_v2.ls_ds", "description", "listener_1 description"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_listener_v2.ls_ds", "tags.#", "3"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_listener_v2.ls_ds", "tags.0", "tag1"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_lb_listener_v2.ls_ds", "provisioning_status"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_lb_listener_v2.ls_ds", "operating_status"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_listener_v2.ls_ds", "protocol", "HTTP"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_listener_v2.ls_ds", "protocol_port", "80"),
					resource.TestCheckResourceAttrPair(
						"data.openstack_lb_listener_v2.ls_ds", "loadbalancers.0.id",
						"openstack_lb_loadbalancer_v2.loadbalancer_1", "id"),
				),
			},
		},
	})
}

const testAccDataSourceLbV2ListenerConfigBasic = `
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
  name        = "pool_1"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  loadbalancer_id = openstack_lb_loadbalancer_v2.loadbalancer_1.id
}

data "openstack_lb_listener_v2" "ls_ds" {
  listener_id = openstack_lb_listener_v2.listener_1.id
  tags = [
    "tag1",
    "tag3",
  ]
}
`
