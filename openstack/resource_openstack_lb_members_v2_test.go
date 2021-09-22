package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/pools"
)

func testAccCheckLBV2MembersComputeHash(members *[]pools.Member, weight int, address string, idx *int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		membersResource := resourceMembersV2().Schema["member"].Elem.(*schema.Resource)
		f := schema.HashResource(membersResource)

		for _, m := range flattenLBMembersV2(*members) {
			if m["address"] == address && m["weight"] == weight {
				*idx = f(m)
				break
			}
		}

		return nil
	}
}

func testCheckResourceAttrWithIndexesAddr(name, format string, idx *int, value string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return resource.TestCheckResourceAttr(name, fmt.Sprintf(format, *idx), value)(s)
	}
}

func testCheckResourceAttrSetWithIndexesAddr(name, format string, idx *int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return resource.TestCheckResourceAttrSet(name, fmt.Sprintf(format, *idx))(s)
	}
}

func TestAccLBV2Members_basic(t *testing.T) {
	var members []pools.Member
	var idx1 int
	var idx2 int

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
			testAccPreCheckUseOctavia(t)
			testAccPreCheckOctaviaBatchMembersEnv(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV2MembersDestroy,
		Steps: []resource.TestStep{
			{
				Config: TestAccLbV2MembersConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2MembersExists("openstack_lb_members_v2.members_1", &members),
					resource.TestCheckResourceAttr("openstack_lb_members_v2.members_1", "member.#", "2"),
					testAccCheckLBV2MembersComputeHash(&members, 0, "192.168.199.110", &idx1),
					testAccCheckLBV2MembersComputeHash(&members, 1, "192.168.199.111", &idx2),
					testCheckResourceAttrWithIndexesAddr("openstack_lb_members_v2.members_1", "member.%d.weight", &idx1, "0"),
					testCheckResourceAttrWithIndexesAddr("openstack_lb_members_v2.members_1", "member.%d.weight", &idx2, "1"),
					testCheckResourceAttrWithIndexesAddr("openstack_lb_members_v2.members_1", "member.%d.backup", &idx1, "false"),
					testCheckResourceAttrWithIndexesAddr("openstack_lb_members_v2.members_1", "member.%d.backup", &idx2, "true"),
					testCheckResourceAttrSetWithIndexesAddr("openstack_lb_members_v2.members_1", "member.%d.subnet_id", &idx1),
					testCheckResourceAttrSetWithIndexesAddr("openstack_lb_members_v2.members_1", "member.%d.subnet_id", &idx2),
				),
			},
			{
				Config: TestAccLbV2MembersConfigUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2MembersExists("openstack_lb_members_v2.members_1", &members),
					resource.TestCheckResourceAttr("openstack_lb_members_v2.members_1", "member.#", "2"),
					testAccCheckLBV2MembersComputeHash(&members, 10, "192.168.199.110", &idx1),
					testAccCheckLBV2MembersComputeHash(&members, 15, "192.168.199.111", &idx2),
					testCheckResourceAttrWithIndexesAddr("openstack_lb_members_v2.members_1", "member.%d.weight", &idx1, "10"),
					testCheckResourceAttrWithIndexesAddr("openstack_lb_members_v2.members_1", "member.%d.weight", &idx2, "15"),
					testCheckResourceAttrWithIndexesAddr("openstack_lb_members_v2.members_1", "member.%d.backup", &idx1, "true"),
					testCheckResourceAttrWithIndexesAddr("openstack_lb_members_v2.members_1", "member.%d.backup", &idx2, "false"),
					testCheckResourceAttrSetWithIndexesAddr("openstack_lb_members_v2.members_1", "member.%d.subnet_id", &idx1),
					testCheckResourceAttrSetWithIndexesAddr("openstack_lb_members_v2.members_1", "member.%d.subnet_id", &idx2),
				),
			},
			{
				Config: TestAccLbV2MembersConfigUnsetSubnet,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2MembersExists("openstack_lb_members_v2.members_1", &members),
					resource.TestCheckResourceAttr("openstack_lb_members_v2.members_1", "member.#", "2"),
					testAccCheckLBV2MembersComputeHash(&members, 10, "192.168.199.110", &idx1),
					testAccCheckLBV2MembersComputeHash(&members, 15, "192.168.199.111", &idx2),
					testCheckResourceAttrWithIndexesAddr("openstack_lb_members_v2.members_1", "member.%d.weight", &idx1, "10"),
					testCheckResourceAttrWithIndexesAddr("openstack_lb_members_v2.members_1", "member.%d.weight", &idx2, "15"),
					testCheckResourceAttrWithIndexesAddr("openstack_lb_members_v2.members_1", "member.%d.subnet_id", &idx1, ""),
					testCheckResourceAttrWithIndexesAddr("openstack_lb_members_v2.members_1", "member.%d.subnet_id", &idx2, ""),
				),
			},
			{
				Config: TestAccLbV2MembersConfigDeleteMembers,
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
	lbClient, err := chooseLBV2AccTestClient(config, osRegionName)
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
		lbClient, err := chooseLBV2AccTestClient(config, osRegionName)
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

const TestAccLbV2MembersConfigBasic = `
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
	backup = true
  }

  timeouts {
    create = "10m"
    update = "10m"
    delete = "10m"
  }
}
`

const TestAccLbV2MembersConfigUpdate = `
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
	backup = true
}

  member {
    address = "192.168.199.111"
    protocol_port = 8080
    weight = 15
    admin_state_up = "true"
	subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
	backup = false
  }

  timeouts {
    create = "10m"
    update = "10m"
    delete = "10m"
  }
}
`

const TestAccLbV2MembersConfigUnsetSubnet = `
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
    create = "10m"
    update = "10m"
    delete = "10m"
  }
}
`

const TestAccLbV2MembersConfigDeleteMembers = `
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
    create = "10m"
    update = "10m"
    delete = "10m"
  }
}
`
