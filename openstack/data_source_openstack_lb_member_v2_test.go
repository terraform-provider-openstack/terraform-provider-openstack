package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSourceLBV2Member_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceLbV2MemberConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.openstack_lb_member_v2.mb_ds_1", "name", "member_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_member_v2.mb_ds_2", "name", "member_2"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_lb_member_v2.mb_ds_2", "member_id"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_member_v2.mb_ds_1", "tags.#", "3"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_member_v2.mb_ds_1", "tags.0", "tag1"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_member_v2.mb_ds_1", "monitor_address", "192.168.199.73"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_member_v2.mb_ds_1", "monitor_port", "8181"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_member_v2.mb_ds_2", "monitor_address", ""),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_member_v2.mb_ds_2", "monitor_port", "0"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_member_v2.mb_ds_1", "protocol_port", "8080"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_member_v2.mb_ds_2", "protocol_port", "9090"),
				),
			},
			{
				Config: testAccDataSourceLBV2MemberManifestUpdate1(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.openstack_lb_member_v2.mb_ds_3", "name", "member_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_member_v2.mb_ds_4", "name", "member_1"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_member_v2.mb_ds_3", "tags.#", "2"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_member_v2.mb_ds_3", "tags.0", "tag11"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_member_v2.mb_ds_4", "tags.#", "2"),
					resource.TestCheckResourceAttr(
						"data.openstack_lb_member_v2.mb_ds_4", "tags.0", "tag13"),
				),
			},
		},
	})
}

const testAccDataSourceLbV2MemberConfigBasic = `
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
}

resource "openstack_lb_member_v2" "member_1" {
  pool_id         = openstack_lb_pool_v2.pool_1.id
  name            = "member_1"
  address         = "192.168.199.13"
  protocol_port   = 8080
  monitor_address = "192.168.199.73"
  monitor_port    = 8181
  tags = [
	  "tag1",
    "tag2",
    "tag3",
  ]
}

resource "openstack_lb_member_v2" "member_2" {
  pool_id       = openstack_lb_pool_v2.pool_1.id
  name          = "member_2"
  address       = "192.168.199.23"
  protocol_port = 9090
  subnet_id     = openstack_networking_subnet_v2.subnet_1.id
}

resource "openstack_lb_monitor_v2" "monitor_1" {
  pool_id     = openstack_lb_pool_v2.pool_1.id
  name        = "monitor_1"
  type        = "PING"
  delay       = 20
  timeout     = 10
  max_retries = 5
}

data "openstack_lb_member_v2" "mb_ds_1" {
  name = openstack_lb_member_v2.member_1.name
  pool_id = openstack_lb_pool_v2.pool_1.id
  tags = [
	  "tag1",
    "tag2",
    "tag3",
  ]
}

data "openstack_lb_member_v2" "mb_ds_2" {
  member_id = openstack_lb_member_v2.member_2.id
  pool_id = openstack_lb_pool_v2.pool_1.id
}
`

func testAccDataSourceLBV2MemberManifestUpdate1() string {
	return fmt.Sprintf(`
%s
resource "openstack_lb_member_v2" "member_3" {
  pool_id         = openstack_lb_pool_v2.pool_1.id
  name            = "member_1"
  address         = "192.168.199.33"
  protocol_port   = 8080
  monitor_address = "192.168.199.93"
  monitor_port    = 8181
  tags = [
	  "tag11",
    "tag12",
  ]
}

resource "openstack_lb_member_v2" "member_4" {
  pool_id       = openstack_lb_pool_v2.pool_1.id
  name          = "member_1"
  address       = "192.168.199.43"
  protocol_port = 9090
  subnet_id     = openstack_networking_subnet_v2.subnet_1.id
  tags = [
	  "tag13",
    "tag14",
  ]
}

data "openstack_lb_member_v2" "mb_ds_3" {
  name = openstack_lb_member_v2.member_3.name
  pool_id = openstack_lb_pool_v2.pool_1.id
  tags = [
	  "tag11",
  ]
}

data "openstack_lb_member_v2" "mb_ds_4" {
  member_id = openstack_lb_member_v2.member_4.id
  pool_id = openstack_lb_pool_v2.pool_1.id
  tags = [
	  "tag13",
    "tag14",
  ]
}
`, testAccDataSourceLbV2MemberConfigBasic)
}
