package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas/members"
)

func TestAccLBV1Member_basic(t *testing.T) {
	var member members.Member

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckDeprecated(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV1MemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbV1MemberBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV1MemberExists("openstack_lb_member_v1.member_1", &member),
				),
			},
			{
				Config: testAccLbV1MemberUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openstack_lb_member_v1.member_1", "admin_state_up", "false"),
				),
			},
		},
	})
}

func TestAccLBV1Member_timeout(t *testing.T) {
	var member members.Member

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckDeprecated(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV1MemberDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbV1MemberTimeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV1MemberExists("openstack_lb_member_v1.member_1", &member),
				),
			},
		},
	})
}

func testAccCheckLBV1MemberDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.NetworkingV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_lb_member_v1" {
			continue
		}

		_, err := members.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("LB Member still exists")
		}
	}

	return nil
}

func testAccCheckLBV1MemberExists(n string, member *members.Member) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.NetworkingV2Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %s", err)
		}

		found, err := members.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Member not found")
		}

		*member = *found

		return nil
	}
}

const testAccLbV1MemberBasic = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_lb_pool_v1" "pool_1" {
  name = "pool_1"
  protocol = "HTTP"
  lb_method = "ROUND_ROBIN"
  subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
}

resource "openstack_lb_member_v1" "member_1" {
  address = "192.168.199.10"
  port = 80
  admin_state_up = true
  pool_id = "${openstack_lb_pool_v1.pool_1.id}"
}
`

const testAccLbV1MemberUpdate = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_lb_pool_v1" "pool_1" {
  name = "pool_1"
  protocol = "HTTP"
  lb_method = "ROUND_ROBIN"
  subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
}

resource "openstack_lb_member_v1" "member_1" {
  address = "192.168.199.10"
  port = 80
  admin_state_up = false
  pool_id = "${openstack_lb_pool_v1.pool_1.id}"
}
`

const testAccLbV1MemberTimeout = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  admin_state_up = "true"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  cidr = "192.168.199.0/24"
  ip_version = 4
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_lb_pool_v1" "pool_1" {
  name = "pool_1"
  protocol = "HTTP"
  lb_method = "ROUND_ROBIN"
  subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"
}

resource "openstack_lb_member_v1" "member_1" {
  address = "192.168.199.10"
  port = 80
  admin_state_up = true
  pool_id = "${openstack_lb_pool_v1.pool_1.id}"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`
