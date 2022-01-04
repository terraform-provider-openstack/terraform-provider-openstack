package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/lbaas/monitors"
)

func TestAccLBV1Monitor_basic(t *testing.T) {
	var monitor monitors.Monitor

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckDeprecated(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV1MonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbV1MonitorBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV1MonitorExists("openstack_lb_monitor_v1.monitor_1", &monitor),
				),
			},
			{
				Config: testAccLbV1MonitorUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("openstack_lb_monitor_v1.monitor_1", "delay", "20"),
				),
			},
		},
	})
}

func TestAccLBV1Monitor_timeout(t *testing.T) {
	var monitor monitors.Monitor

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheckDeprecated(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckLBV1MonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccLbV1MonitorTimeout,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckLBV1MonitorExists("openstack_lb_monitor_v1.monitor_1", &monitor),
				),
			},
		},
	})
}

func testAccCheckLBV1MonitorDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.NetworkingV2Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_lb_monitor_v1" {
			continue
		}

		_, err := monitors.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("LB monitor still exists")
		}
	}

	return nil
}

func testAccCheckLBV1MonitorExists(n string, monitor *monitors.Monitor) resource.TestCheckFunc {
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

		found, err := monitors.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Monitor not found")
		}

		*monitor = *found

		return nil
	}
}

const testAccLbV1MonitorBasic = `
resource "openstack_lb_monitor_v1" "monitor_1" {
  type = "PING"
  delay = 30
  timeout = 5
  max_retries = 3
  admin_state_up = "true"
}
`

const testAccLbV1MonitorUpdate = `
resource "openstack_lb_monitor_v1" "monitor_1" {
  type = "PING"
  delay = 20
  timeout = 5
  max_retries = 3
  admin_state_up = "true"
}
`

const testAccLbV1MonitorTimeout = `
resource "openstack_lb_monitor_v1" "monitor_1" {
  type = "PING"
  delay = 30
  timeout = 5
  max_retries = 3
  admin_state_up = "true"

  timeouts {
    create = "5m"
    delete = "5m"
  }
}
`
