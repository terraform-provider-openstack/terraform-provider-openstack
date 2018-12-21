package openstack

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/l7policies"
)

func TestAccLBV2L7policy_basic(t *testing.T) {
	var l7Policy l7policies.L7Policy

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckLB(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckLBV2L7policyDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckLBV2L7policyConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2L7policyExists("openstack_lb_l7policy_v2.l7policy_1", &l7Policy),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7policy_v2.l7policy_1", "name", "test"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7policy_v2.l7policy_1", "description", "test description"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7policy_v2.l7policy_1", "action", "REJECT"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7policy_v2.l7policy_1", "position", "1"),
					resource.TestMatchResourceAttr(
						"openstack_lb_l7policy_v2.l7policy_1", "listener_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
				),
			},
			resource.TestStep{
				Config: testAccCheckLBV2L7policyConfig_update1,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2L7policyExists("openstack_lb_l7policy_v2.l7policy_1", &l7Policy),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7policy_v2.l7policy_1", "name", "test"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7policy_v2.l7policy_1", "description", "test description"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7policy_v2.l7policy_1", "action", "REDIRECT_TO_URL"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7policy_v2.l7policy_1", "position", "1"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7policy_v2.l7policy_1", "redirect_url", "http://www.example.com"),
					resource.TestMatchResourceAttr(
						"openstack_lb_l7policy_v2.l7policy_1", "listener_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
				),
			},
			resource.TestStep{
				Config: testAccCheckLBV2L7policyConfig_update2,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2L7policyExists("openstack_lb_l7policy_v2.l7policy_1", &l7Policy),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7policy_v2.l7policy_1", "name", "test_updated"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7policy_v2.l7policy_1", "description", ""),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7policy_v2.l7policy_1", "action", "REDIRECT_TO_POOL"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7policy_v2.l7policy_1", "position", "1"),
					resource.TestMatchResourceAttr(
						"openstack_lb_l7policy_v2.l7policy_1", "listener_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
					resource.TestMatchResourceAttr(
						"openstack_lb_l7policy_v2.l7policy_1", "redirect_pool_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
				),
			},
			resource.TestStep{
				Config: testAccCheckLBV2L7policyConfig_update3,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2L7policyExists("openstack_lb_l7policy_v2.l7policy_1", &l7Policy),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7policy_v2.l7policy_1", "name", "test_updated"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7policy_v2.l7policy_1", "description", ""),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7policy_v2.l7policy_1", "action", "REJECT"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7policy_v2.l7policy_1", "position", "1"),
					resource.TestMatchResourceAttr(
						"openstack_lb_l7policy_v2.l7policy_1", "listener_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
				),
			},
		},
	})
}

func testAccCheckLBV2L7policyDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	lbClient, err := chooseLBV2AccTestClient(config, OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack load balancing client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_lb_l7policy_v2" {
			continue
		}

		_, err := l7policies.Get(lbClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("L7 Policy still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV2L7policyExists(n string, l7Policy *l7policies.L7Policy) resource.TestCheckFunc {
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

		found, err := l7policies.Get(lbClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Policy not found")
		}

		*l7Policy = *found

		return nil
	}
}

const testAccCheckLBV2L7policyConfig = `
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
`

var testAccCheckLBV2L7policyConfig_basic = fmt.Sprintf(`
%s

resource "openstack_lb_l7policy_v2" "l7policy_1" {
  name         = "test"
  action       = "REJECT"
  description  = "test description"
  position     = 1
  listener_id  = "${openstack_lb_listener_v2.listener_1.id}"
}
`, testAccCheckLBV2L7policyConfig)

var testAccCheckLBV2L7policyConfig_update1 = fmt.Sprintf(`
%s

resource "openstack_lb_l7policy_v2" "l7policy_1" {
  name         = "test"
  action       = "REDIRECT_TO_URL"
  description  = "test description"
  position     = 1
  listener_id  = "${openstack_lb_listener_v2.listener_1.id}"
  redirect_url = "http://www.example.com"
}
`, testAccCheckLBV2L7policyConfig)

var testAccCheckLBV2L7policyConfig_update2 = fmt.Sprintf(`
%s

resource "openstack_lb_pool_v2" "pool_1" {
  name            = "pool_1"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = "${openstack_lb_loadbalancer_v2.loadbalancer_1.id}"
}

resource "openstack_lb_l7policy_v2" "l7policy_1" {
  name             = "test_updated"
  action           = "REDIRECT_TO_POOL"
  position         = 1
  listener_id      = "${openstack_lb_listener_v2.listener_1.id}"
  redirect_pool_id = "${openstack_lb_pool_v2.pool_1.id}"
}
`, testAccCheckLBV2L7policyConfig)

var testAccCheckLBV2L7policyConfig_update3 = fmt.Sprintf(`
%s

resource "openstack_lb_pool_v2" "pool_1" {
  name            = "pool_1"
  protocol        = "HTTP"
  lb_method       = "ROUND_ROBIN"
  loadbalancer_id = "${openstack_lb_loadbalancer_v2.loadbalancer_1.id}"
}

resource "openstack_lb_l7policy_v2" "l7policy_1" {
  name             = "test_updated"
  action           = "REJECT"
  position         = 1
  listener_id      = "${openstack_lb_listener_v2.listener_1.id}"
}
`, testAccCheckLBV2L7policyConfig)
