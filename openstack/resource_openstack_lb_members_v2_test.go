package openstack

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/pools"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccLBV2Members_basic(t *testing.T) {
	var members []pools.Member

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheckLB(t)
			testAccPreCheckUseOctavia(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2MembersDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccLBV2MembersConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2MembersExists("openstack_lb_members_v2.members_1", &members),
					resource.TestCheckResourceAttr("openstack_lb_members_v2.members_1", "member.#", "2"),
					resource.TestCheckResourceAttr("openstack_lb_members_v2.members_1", "member.0.weight", "0"),
					resource.TestCheckResourceAttr("openstack_lb_members_v2.members_1", "member.1.weight", "1"),
					resource.TestCheckResourceAttrSet("openstack_lb_members_v2.members_1", "member.0.subnet_id"),
					resource.TestCheckResourceAttrSet("openstack_lb_members_v2.members_1", "member.1.subnet_id"),
				),
			},
			{
				Config: TestAccLBV2MembersConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2MembersExists("openstack_lb_members_v2.members_1", &members),
					resource.TestCheckResourceAttr("openstack_lb_members_v2.members_1", "member.#", "2"),
					resource.TestCheckResourceAttr("openstack_lb_members_v2.members_1", "member.0.weight", "10"),
					resource.TestCheckResourceAttr("openstack_lb_members_v2.members_1", "member.1.weight", "15"),
					resource.TestCheckResourceAttrSet("openstack_lb_members_v2.members_1", "member.0.subnet_id"),
					resource.TestCheckResourceAttrSet("openstack_lb_members_v2.members_1", "member.1.subnet_id"),
				),
			},
			{
				Config: TestAccLBV2MembersConfig_unset_subnet,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2MembersExists("openstack_lb_members_v2.members_1", &members),
					resource.TestCheckResourceAttr("openstack_lb_members_v2.members_1", "member.#", "2"),
					resource.TestCheckResourceAttr("openstack_lb_members_v2.members_1", "member.0.weight", "10"),
					resource.TestCheckResourceAttr("openstack_lb_members_v2.members_1", "member.1.weight", "15"),
					resource.TestCheckResourceAttr("openstack_lb_members_v2.members_1", "member.0.subnet_id", ""),
					resource.TestCheckResourceAttr("openstack_lb_members_v2.members_1", "member.1.subnet_id", ""),
				),
			},
			{
				Config: TestAccLBV2MembersConfig_delete_members,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2MembersExists("openstack_lb_members_v2.members_1", &members),
					resource.TestCheckResourceAttr("openstack_lb_members_v2.members_1", "member.#", "0"),
				),
			},
		},
	})
}

func testAccCheckLBV2MembersDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	lbClient, err := chooseLBV2AccTestClient(config, OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack load balancing client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_lb_members_v2" {
			continue
		}

		poolID := rs.Primary.Attributes["pool_id"]

		allPages, err := pools.ListMembers(lbClient, poolID, pools.ListMembersOpts{}).AllPages()
		if err != nil {
			if _, ok := err.(gophercloud.ErrDefault404); ok {
				return nil
			}
			return fmt.Errorf("Error getting openstack_lb_members_v2: %s", err)
		}

		members, err := pools.ExtractMembers(allPages)
		if err != nil {
			return fmt.Errorf("Unable to retrieve openstack_lb_members_v2: %s", err)
		}

		if len(members) > 0 {
			return fmt.Errorf("Members still exist: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV2MembersExists(n string, members *[]pools.Member) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		lbClient, err := chooseLBV2AccTestClient(config, OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack load balancing client: %s", err)
		}

		poolID := rs.Primary.Attributes["pool_id"]
		allPages, err := pools.ListMembers(lbClient, poolID, pools.ListMembersOpts{}).AllPages()
		if err != nil {
			return fmt.Errorf("Error getting openstack_lb_members_v2: %s", err)
		}

		found, err := pools.ExtractMembers(allPages)
		if err != nil {
			return fmt.Errorf("Unable to retrieve openstack_lb_members_v2: %s", err)
		}

		*members = found

		return nil
	}
}

const TestAccLBV2MembersConfig_basic = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  network_id = "${openstack_networking_network_v2.network_1.id}"
  cidr = "192.168.199.0/24"
  ip_version = 4
}

resource "openstack_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
  vip_address = "192.168.199.10"
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
}

resource "openstack_lb_members_v2" "members_1" {
  pool_id = "${openstack_lb_pool_v2.pool_1.id}"

  member {
    address = "192.168.199.110"
    protocol_port = 8080
    subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
    weight = 0
  }

  member {
    address = "192.168.199.111"
    protocol_port = 8080
    subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
  }

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const TestAccLBV2MembersConfig_update = `
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
}

resource "openstack_lb_members_v2" "members_1" {
  pool_id = "${openstack_lb_pool_v2.pool_1.id}"

  member {
    address = "192.168.199.110"
    protocol_port = 8080
    weight = 10
    admin_state_up = "true"
    subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
  }

  member {
    address = "192.168.199.111"
    protocol_port = 8080
    weight = 15
    admin_state_up = "true"
    subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
  }

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const TestAccLBV2MembersConfig_unset_subnet = `
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
}

resource "openstack_lb_members_v2" "members_1" {
  pool_id = "${openstack_lb_pool_v2.pool_1.id}"

  member {
    address = "192.168.199.110"
    protocol_port = 8080
    weight = 10
    admin_state_up = "true"
  }

  member {
    address = "192.168.199.111"
    protocol_port = 8080
    weight = 15
    admin_state_up = "true"
  }

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`

const TestAccLBV2MembersConfig_delete_members = `
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
}

resource "openstack_lb_members_v2" "members_1" {
  pool_id = "${openstack_lb_pool_v2.pool_1.id}"

  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
`
