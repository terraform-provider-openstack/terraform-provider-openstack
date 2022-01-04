package openstack

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/l7policies"
)

func TestAccLBV2L7Rule_basic(t *testing.T) {
	var l7rule l7policies.Rule

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
					testAccCheckLBV2L7RuleExists("openstack_lb_l7rule_v2.l7rule_1", &l7rule),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "type", "PATH"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "compare_type", "EQUAL_TO"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "value", "/api"),
					resource.TestMatchResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "listener_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
					resource.TestMatchResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "l7policy_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
				),
			},
			{
				Config: testAccCheckLbV2L7RuleConfigUpdate1(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2L7RuleExists("openstack_lb_l7rule_v2.l7rule_1", &l7rule),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "type", "HOST_NAME"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "compare_type", "EQUAL_TO"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "value", "www.example.com"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "invert", "true"),
					resource.TestMatchResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "listener_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
					resource.TestMatchResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "l7policy_id",
						regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")),
				),
			},
			{
				Config: testAccCheckLbV2L7RuleConfigUpdate2(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2L7RuleExists("openstack_lb_l7rule_v2.l7rule_1", &l7rule),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "type", "HOST_NAME"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "compare_type", "EQUAL_TO"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "value", "www.example.com"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "invert", "true"),
				),
			},
			{
				Config: testAccCheckLbV2L7RuleConfigUpdate3(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2L7RuleExists("openstack_lb_l7rule_v2.l7rule_1", &l7rule),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "type", "HEADER"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "compare_type", "EQUAL_TO"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "key", "Host"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "value", "www.example.com"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "invert", "false"),
				),
			},
			{
				Config: testAccCheckLbV2L7RuleConfigUpdate4(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2L7RuleExists("openstack_lb_l7rule_v2.l7rule_1", &l7rule),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "type", "HOST_NAME"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "compare_type", "EQUAL_TO"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "key", ""),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "value", "www.example.com"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "invert", "false"),
				),
			},
			{
				Config: testAccCheckLbV2L7RuleConfigUpdate5(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2L7RuleExists("openstack_lb_l7rule_v2.l7rule_1", &l7rule),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "type", "COOKIE"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "compare_type", "EQUAL_TO"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "key", "X-Ref"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "value", "foo"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "invert", "false"),
				),
			},
			{
				Config: testAccCheckLbV2L7RuleConfigUpdate6(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2L7RuleExists("openstack_lb_l7rule_v2.l7rule_1", &l7rule),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "type", "PATH"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "compare_type", "STARTS_WITH"),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "key", ""),
					resource.TestCheckResourceAttr(
						"openstack_lb_l7rule_v2.l7rule_1", "value", "/images"),
				),
			},
		},
	})
}

func testAccCheckLBV2L7RuleDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	lbClient, err := chooseLBV2AccTestClient(config, osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack load balancing client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_lb_l7rule_v2" {
			continue
		}

		l7policyID := ""
		for k, v := range rs.Primary.Attributes {
			if k == "l7policy_id" {
				l7policyID = v
				break
			}
		}

		if l7policyID == "" {
			return fmt.Errorf("Unable to find l7policy_id")
		}

		_, err := l7policies.GetRule(lbClient, l7policyID, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("L7 Rule still exists: %s", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckLBV2L7RuleExists(n string, l7rule *l7policies.Rule) resource.TestCheckFunc {
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

		l7policyID := ""
		for k, v := range rs.Primary.Attributes {
			if k == "l7policy_id" {
				l7policyID = v
				break
			}
		}

		if l7policyID == "" {
			return fmt.Errorf("Unable to find l7policy_id")
		}

		found, err := l7policies.GetRule(lbClient, l7policyID, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Policy not found")
		}

		*l7rule = *found

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
  network_id = "${openstack_networking_network_v2.network_1.id}"
}

resource "openstack_lb_loadbalancer_v2" "loadbalancer_1" {
  name = "loadbalancer_1"
  vip_subnet_id = "${openstack_networking_subnet_v2.subnet_1.id}"

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
  loadbalancer_id = "${openstack_lb_loadbalancer_v2.loadbalancer_1.id}"
}

resource "openstack_lb_l7policy_v2" "l7policy_1" {
  name         = "test"
  action       = "REDIRECT_TO_URL"
  description  = "test description"
  position     = 1
  listener_id  = "${openstack_lb_listener_v2.listener_1.id}"
  redirect_url = "http://www.example.com"
}
`

func testAccCheckLbV2L7RuleConfigBasic() string {
	return fmt.Sprintf(`
%s

resource "openstack_lb_l7rule_v2" "l7rule_1" {
  l7policy_id  = "${openstack_lb_l7policy_v2.l7policy_1.id}"
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
  l7policy_id  = "${openstack_lb_l7policy_v2.l7policy_1.id}"
  type         = "HOST_NAME"
  compare_type = "EQUAL_TO"
  value        = "www.example.com"
  invert       = true
}
`, testAccCheckLbV2L7RuleConfig)
}

func testAccCheckLbV2L7RuleConfigUpdate2() string {
	return fmt.Sprintf(`
%s

resource "openstack_lb_l7rule_v2" "l7rule_1" {
  l7policy_id  = "${openstack_lb_l7policy_v2.l7policy_1.id}"
  type         = "HOST_NAME"
  compare_type = "EQUAL_TO"
  value        = "www.example.com"
  invert       = true
}
`, testAccCheckLbV2L7RuleConfig)
}

func testAccCheckLbV2L7RuleConfigUpdate3() string {
	return fmt.Sprintf(`
%s

resource "openstack_lb_l7rule_v2" "l7rule_1" {
  l7policy_id  = "${openstack_lb_l7policy_v2.l7policy_1.id}"
  type         = "HEADER"
  compare_type = "EQUAL_TO"
  key          = "Host"
  value        = "www.example.com"
}
`, testAccCheckLbV2L7RuleConfig)
}

func testAccCheckLbV2L7RuleConfigUpdate4() string {
	return fmt.Sprintf(`
%s

resource "openstack_lb_l7rule_v2" "l7rule_1" {
  l7policy_id  = "${openstack_lb_l7policy_v2.l7policy_1.id}"
  type         = "HOST_NAME"
  compare_type = "EQUAL_TO"
  value        = "www.example.com"
}
`, testAccCheckLbV2L7RuleConfig)
}

func testAccCheckLbV2L7RuleConfigUpdate5() string {
	return fmt.Sprintf(`
%s

resource "openstack_lb_l7rule_v2" "l7rule_1" {
  l7policy_id  = "${openstack_lb_l7policy_v2.l7policy_1.id}"
  type         = "COOKIE"
  compare_type = "EQUAL_TO"
  key          = "X-Ref"
  value        = "foo"
}
`, testAccCheckLbV2L7RuleConfig)
}

func testAccCheckLbV2L7RuleConfigUpdate6() string {
	return fmt.Sprintf(`
%s

resource "openstack_lb_l7rule_v2" "l7rule_1" {
  l7policy_id  = "${openstack_lb_l7policy_v2.l7policy_1.id}"
  type         = "PATH"
  compare_type = "STARTS_WITH"
  value        = "/images"
}
`, testAccCheckLbV2L7RuleConfig)
}
