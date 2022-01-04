package openstack

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas_v2/l7policies"
)

func TestAccLBV2L7Policy_basic(t *testing.T) {
	var l7Policy l7policies.L7Policy

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckLB(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV2L7PolicyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckLbV2L7PolicyConfigBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2L7PolicyExists("openstack_lb_l7policy_v2.l7policy_1", &l7Policy),
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
			{
				Config: testAccCheckLbV2L7PolicyConfigUpdate1(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2L7PolicyExists("openstack_lb_l7policy_v2.l7policy_1", &l7Policy),
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
			{
				Config: testAccCheckLbV2L7PolicyConfigUpdate2(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2L7PolicyExists("openstack_lb_l7policy_v2.l7policy_1", &l7Policy),
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
			{
				Config: testAccCheckLbV2L7PolicyConfigUpdate3(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV2L7PolicyExists("openstack_lb_l7policy_v2.l7policy_1", &l7Policy),
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

func testAccCheckLBV2L7PolicyDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	lbClient, err := chooseLBV2AccTestClient(config, osRegionName)
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

func testAccCheckLBV2L7PolicyExists(n string, l7Policy *l7policies.L7Policy) resource.TestCheckFunc {
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

const testAccCheckLbV2L7PolicyConfig = `
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
`

func testAccCheckLbV2L7PolicyConfigBasic() string {
	return fmt.Sprintf(`
%s

resource "openstack_lb_l7policy_v2" "l7policy_1" {
  name         = "test"
  action       = "REJECT"
  description  = "test description"
  position     = 1
  listener_id  = "${openstack_lb_listener_v2.listener_1.id}"
}
`, testAccCheckLbV2L7PolicyConfig)
}

func testAccCheckLbV2L7PolicyConfigUpdate1() string {
	return fmt.Sprintf(`
%s

resource "openstack_lb_l7policy_v2" "l7policy_1" {
  name         = "test"
  action       = "REDIRECT_TO_URL"
  description  = "test description"
  position     = 1
  listener_id  = "${openstack_lb_listener_v2.listener_1.id}"
  redirect_url = "http://www.example.com"
}
`, testAccCheckLbV2L7PolicyConfig)
}

func testAccCheckLbV2L7PolicyConfigUpdate2() string {
	return fmt.Sprintf(`
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
`, testAccCheckLbV2L7PolicyConfig)
}

func testAccCheckLbV2L7PolicyConfigUpdate3() string {
	return fmt.Sprintf(`
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
`, testAccCheckLbV2L7PolicyConfig)
}
