package openstack

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/gophercloud/gophercloud/v2/openstack/loadbalancer/v2/l7policies"
)

func TestAccLBV2L7Rule_basic(t *testing.T) {
	var l7Policy l7policies.L7Policy
	var l7Rule l7policies.Rule

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV2L7RuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLbV2L7RuleConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2L7PolicyExists("openstack_lb_l7policy_v2.l7policy_1", &l7Policy),
					testAccCheckLBV2L7RuleExists("openstack_lb_l7rule_v2.l7rule_1", &l7Policy.ID, &l7Rule),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "admin_state_up", "true"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "type", "PATH"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "compare_type", "EQUAL_TO"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "value", "/api"),
					resource.TestCheckResourceAttrPtr(
						"openstack_lb_l7rule_v2.l7rule_1", "l7policy_id", &l7Policy.ID),
				),
			},
			{
				Config: testAccCheckLbV2L7RuleConfigUpdate1(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2L7PolicyExists("openstack_lb_l7policy_v2.l7policy_1", &l7Policy),
					testAccCheckLBV2L7RuleExists("openstack_lb_l7rule_v2.l7rule_1", &l7Policy.ID, &l7Rule),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "admin_state_up", "false"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "type", "HOST_NAME"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "compare_type", "STARTS_WITH"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "value", "example.com"),
					resource.TestCheckResourceAttrPtr(
						"openstack_lb_l7rule_v2.l7rule_1", "l7policy_id", &l7Policy.ID),
				),
			},
			{
				Config: testAccCheckLbV2L7RuleConfigUpdate2(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2L7PolicyExists("openstack_lb_l7policy_v2.l7policy_1", &l7Policy),
					testAccCheckLBV2L7RuleExists("openstack_lb_l7rule_v2.l7rule_1", &l7Policy.ID, &l7Rule),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "admin_state_up", "true"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "type", "COOKIE"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "compare_type", "ENDS_WITH"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "value", "cookie"),
					resource.TestCheckResourceAttrPtr(
						"openstack_lb_l7rule_v2.l7rule_1", "l7policy_id", &l7Policy.ID),
				),
			},
		},
	})
}

func testAccCheckLBV2L7RuleDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	lbClient, err := config.LoadBalancerV2Client(context.TODO(), osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack load balancing client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_lb_l7rule_v2" {
			continue
		}

		_, err := l7policies.Get(context.TODO(), lbClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("L7 Rule still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV2L7RuleExists(n string, l7PolicyID *string, l7Rule *l7policies.Rule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		lbClient, err := config.LoadBalancerV2Client(context.TODO(), osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack load balancing client: %s", err)
		}

		found, err := l7policies.GetRule(context.TODO(), lbClient, *l7PolicyID, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Rule not found")
		}

		*l7Rule = *found

		return nil
	}
}

const testAccCheckLbV2L7RuleConfig = `
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

resource "openstack_lb_l7policy_v2" "l7policy_1" {
  name         = "test"
  action       = "REDIRECT_TO_URL"
  description  = "test description"
  position     = 1
  listener_id  =  openstack_lb_listener_v2.listener_1.id
  redirect_url = "http://www.example.com"
}
`

func testAccCheckLbV2L7RuleConfigBasic() string {
	return fmt.Sprintf(`
%s

resource "openstack_lb_l7rule_v2" "l7rule_1" {
  admin_state_up = true

  l7policy_id  = openstack_lb_l7policy_v2.l7policy_1.id
  type         = "PATH"
  compare_type = "EQUAL_TO"
  value        = "/api"
}
`, testAccCheckLbV2L7RuleConfig)
}

func testAccCheckLbV2L7RuleConfigUpdate1() string {
	return fmt.Sprintf(`
%s

resource "openstack_lb_l7rule_v2" "l7rule_1" {
  admin_state_up = false

  l7policy_id  = openstack_lb_l7policy_v2.l7policy_1.id
  type         = "HOST_NAME"
  compare_type = "STARTS_WITH"
  value        = "example.com"
}
`, testAccCheckLbV2L7RuleConfig)
}

func testAccCheckLbV2L7RuleConfigUpdate2() string {
	return fmt.Sprintf(`
%s

resource "openstack_lb_l7rule_v2" "l7rule_1" {
  admin_state_up = true

  l7policy_id  = openstack_lb_l7policy_v2.l7policy_1.id
  type         = "COOKIE"
  key          = "test"
  compare_type = "ENDS_WITH"
  value        = "cookie"
}
`, testAccCheckLbV2L7RuleConfig)
}
