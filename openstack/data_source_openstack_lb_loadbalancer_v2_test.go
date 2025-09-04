package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceLBV2LoadBalancer_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLbV2LoadBalancerConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.openstack_lb_loadbalancer_v2.lb_ds", "name", "loadbalancer_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_loadbalancer_v2.lb_ds", "description", "loadbalancer_1 description"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_loadbalancer_v2.lb_ds", "tags.#", "1"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_loadbalancer_v2.lb_ds", "tags.0", "tag1"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_lb_loadbalancer_v2.lb_ds", "vip_subnet_id"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_lb_loadbalancer_v2.lb_ds", "vip_port_id"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_loadbalancer_v2.lb_ds", "pools.#", "1"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_lb_loadbalancer_v2.lb_ds", "pools.0.id"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_loadbalancer_v2.lb_ds", "listeners.#", "1"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_lb_loadbalancer_v2.lb_ds", "listeners.0.id"),
				),
			},
		},
	})
}

func TestAccDataSourceLBV2LoadBalancer_secGroup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLbV2LoadBalancerConfigSecGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.openstack_lb_loadbalancer_v2.lb_ds", "name", "loadbalancer_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_loadbalancer_v2.lb_ds", "security_group_ids.#", "1"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_lb_loadbalancer_v2.lb_ds", "vip_port_id"),
				),
			},
		},
	})
}

func TestAccDataSourceLBV2LoadBalancer_vipNetwork(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLbV2LoadBalancerConfigVIPNetwork,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.openstack_lb_loadbalancer_v2.lb_ds", "name", "loadbalancer_1"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_lb_loadbalancer_v2.lb_ds", "vip_network_id"),
				),
			},
		},
	})
}

func TestAccDataSourceLBV2LoadBalancer_vipPortID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLbV2LoadBalancerConfigVIPPortID,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.openstack_lb_loadbalancer_v2.lb_ds", "name", "loadbalancer_1"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_lb_loadbalancer_v2.lb_ds", "vip_port_id"),
				),
			},
		},
	})
}

const testAccDataSourceLbV2LoadBalancerConfigBasic = `
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
  tags = ["tag1"]

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}

resource "openstack_lb_listener_v2" "listener_1" {
  name            = "listener_1"
  protocol        = "HTTP"
  protocol_port   = 80
  loadbalancer_id = openstack_lb_loadbalancer_v2.loadbalancer_1.id
}

resource "openstack_lb_pool_v2" "pool_1" {
  name        = "pool_1"
  protocol    = "HTTP"
  lb_method   = "ROUND_ROBIN"
  listener_id = openstack_lb_listener_v2.listener_1.id
}

data "openstack_lb_loadbalancer_v2" "lb_ds" {
  name = openstack_lb_loadbalancer_v2.loadbalancer_1.name
  description = openstack_lb_loadbalancer_v2.loadbalancer_1.description
  vip_address = openstack_lb_loadbalancer_v2.loadbalancer_1.vip_address
}
`

const testAccDataSourceLbV2LoadBalancerConfigSecGroup = `
resource "openstack_networking_secgroup_v2" "secgroup_1" {
  name = "secgroup_1"
  description = "secgroup_1"
}

resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  network_id = openstack_networking_network_v2.network_1.id
  cidr = "192.168.199.0/24"
}

resource "openstack_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  description = "loadbalancer_1 description"
  loadbalancer_provider = "octavia"
  vip_subnet_id = openstack_networking_subnet_v2.subnet_1.id
  security_group_ids = [
    openstack_networking_secgroup_v2.secgroup_1.id
  ]

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}

data "openstack_lb_loadbalancer_v2" "lb_ds" {
  loadbalancer_id = openstack_lb_loadbalancer_v2.loadbalancer_1.id
  description = openstack_lb_loadbalancer_v2.loadbalancer_1.description
  vip_address = openstack_lb_loadbalancer_v2.loadbalancer_1.vip_address
}
`

const testAccDataSourceLbV2LoadBalancerConfigVIPNetwork = `
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
  vip_network_id = openstack_networking_network_v2.network_1.id
  depends_on = ["openstack_networking_subnet_v2.subnet_1"]

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}

data "openstack_lb_loadbalancer_v2" "lb_ds" {
  loadbalancer_id = openstack_lb_loadbalancer_v2.loadbalancer_1.id
  vip_address = openstack_lb_loadbalancer_v2.loadbalancer_1.vip_address
}
`

const testAccDataSourceLbV2LoadBalancerConfigVIPPortID = `
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

resource "openstack_networking_port_v2" "port_1" {
  name           = "port_1"
  network_id     = openstack_networking_network_v2.network_1.id
  admin_state_up = "true"
  depends_on = ["openstack_networking_subnet_v2.subnet_1"]
}

resource "openstack_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  description = "loadbalancer_1 description"
  loadbalancer_provider = "octavia"
  vip_port_id = openstack_networking_port_v2.port_1.id
  depends_on = ["openstack_networking_port_v2.port_1"]

  timeouts {
    create = "15m"
    update = "15m"
    delete = "15m"
  }
}

data "openstack_lb_loadbalancer_v2" "lb_ds" {
  loadbalancer_id = openstack_lb_loadbalancer_v2.loadbalancer_1.id
  description = openstack_lb_loadbalancer_v2.loadbalancer_1.description
}
`
