package openstack

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/agents"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOpenStackNetworkingAgentV2DataSource_agentID(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccOpenStackNetworkingAgentV2DataSourceAgentID(osDRAgentID + "1"),
				ExpectError: regexp.MustCompile(`No openstack_networking_agent_v2 found`),
			},
			{
				Config: testAccOpenStackNetworkingAgentV2DataSourceAgentID(osDRAgentID),
				Check:  testAccCheckNetworkingV2AgentChecks(t),
			},
		},
	})
}

func TestAccOpenStackNetworkingAgentV2DataSource_agentType(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccOpenStackNetworkingAgentV2DataSourceAgentType("test"),
				ExpectError: regexp.MustCompile(`No openstack_networking_agent_v2 found`),
			},
			{
				Config: testAccOpenStackNetworkingAgentV2DataSourceAgentType("BGP dynamic routing agent"),
				Check:  testAccCheckNetworkingV2AgentChecks(t),
			},
		},
	})
}

func TestAccOpenStackNetworkingAgentV2DataSource_alive(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccOpenStackNetworkingAgentV2DataSourceAlive("false"),
				ExpectError: regexp.MustCompile(`No openstack_networking_agent_v2 found`),
			},
			{
				Config: testAccOpenStackNetworkingAgentV2DataSourceAlive("true"),
				Check:  testAccCheckNetworkingV2AgentChecks(t),
			},
		},
	})
}

func TestAccOpenStackNetworkingAgentV2DataSource_az(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccOpenStackNetworkingAgentV2DataSourceAZ("test"),
				ExpectError: regexp.MustCompile(`No openstack_networking_agent_v2 found`),
			},
			{
				Config: testAccOpenStackNetworkingAgentV2DataSourceAZ(""),
				Check:  testAccCheckNetworkingV2AgentChecks(t),
			},
		},
	})
}

func TestAccOpenStackNetworkingAgentV2DataSource_binary(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccOpenStackNetworkingAgentV2DataSourceBinary("test"),
				ExpectError: regexp.MustCompile(`No openstack_networking_agent_v2 found`),
			},
			{
				Config: testAccOpenStackNetworkingAgentV2DataSourceBinary("neutron-bgp-dragent"),
				Check:  testAccCheckNetworkingV2AgentChecks(t),
			},
		},
	})
}

func TestAccOpenStackNetworkingAgentV2DataSource_description(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccOpenStackNetworkingAgentV2DataSourceDesc("test"),
				ExpectError: regexp.MustCompile(`No openstack_networking_agent_v2 found`),
			},
			{
				Config: testAccOpenStackNetworkingAgentV2DataSourceDesc(""),
				Check:  testAccCheckNetworkingV2AgentChecks(t),
			},
		},
	})
}

func TestAccOpenStackNetworkingAgentV2DataSource_host(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccOpenStackNetworkingAgentV2DataSourceHost("test"),
				ExpectError: regexp.MustCompile(`No openstack_networking_agent_v2 found`),
			},
			{
				Config: testAccOpenStackNetworkingAgentV2DataSourceHost(osDRAgentHost),
				Check:  testAccCheckNetworkingV2AgentChecks(t),
			},
		},
	})
}

func TestAccOpenStackNetworkingAgentV2DataSource_topic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config:      testAccOpenStackNetworkingAgentV2DataSourceTopic("test"),
				ExpectError: regexp.MustCompile(`No openstack_networking_agent_v2 found`),
			},
			{
				Config: testAccOpenStackNetworkingAgentV2DataSourceTopic("bgp_dragent"),
				Check:  testAccCheckNetworkingV2AgentChecks(t),
			},
		},
	})
}

func testAccCheckNetworkingV2AgentChecks(t *testing.T) resource.TestCheckFunc {
	var agent agents.Agent

	return resource.ComposeTestCheckFunc(
		testAccCheckNetworkingV2AgentExists(t.Context(), "data.openstack_networking_agent_v2.agent", &agent),
		resource.TestCheckResourceAttrPtr("data.openstack_networking_agent_v2.agent", "id", &agent.ID),
		resource.TestCheckResourceAttr("data.openstack_networking_agent_v2.agent", "agent_type", "BGP dynamic routing agent"),
		resource.TestCheckResourceAttr("data.openstack_networking_agent_v2.agent", "alive", "true"),
		resource.TestCheckResourceAttr("data.openstack_networking_agent_v2.agent", "binary", "neutron-bgp-dragent"),
		resource.TestCheckResourceAttrPtr("data.openstack_networking_agent_v2.agent", "host", &agent.Host),
		resource.TestCheckResourceAttr("data.openstack_networking_agent_v2.agent", "topic", "bgp_dragent"),
		resource.TestCheckResourceAttr("data.openstack_networking_agent_v2.agent", "admin_state_up", "true"),
		resource.TestCheckResourceAttr("data.openstack_networking_agent_v2.agent", "resources_synced", "false"),
		resource.TestCheckResourceAttr("data.openstack_networking_agent_v2.agent", "configurations.advertise_routes", "0"),
		resource.TestCheckResourceAttr("data.openstack_networking_agent_v2.agent", "configurations.bgp_peers", "0"),
		resource.TestCheckResourceAttr("data.openstack_networking_agent_v2.agent", "configurations.bgp_speakers", "0"),
		resource.TestCheckResourceAttrSet("data.openstack_networking_agent_v2.agent", "created_at"),
		resource.TestCheckResourceAttrSet("data.openstack_networking_agent_v2.agent", "started_at"),
		resource.TestCheckResourceAttrSet("data.openstack_networking_agent_v2.agent", "heartbeat_timestamp"),
	)
}

func testAccCheckNetworkingV2AgentExists(ctx context.Context, n string, agent *agents.Agent) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)

		networkingClient, err := config.NetworkingV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %w", err)
		}

		found, err := agents.Get(ctx, networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return errors.New("Agent not found")
		}

		*agent = *found

		return nil
	}
}

func testAccOpenStackNetworkingAgentV2DataSourceAgentID(id string) string {
	return fmt.Sprintf(`
data "openstack_networking_agent_v2" "agent" {
  agent_id = "%s"
}
	`, id)
}

func testAccOpenStackNetworkingAgentV2DataSourceAgentType(agentType string) string {
	return fmt.Sprintf(`
data "openstack_networking_agent_v2" "agent" {
  agent_type = "%s"
}`, agentType)
}

func testAccOpenStackNetworkingAgentV2DataSourceAlive(alive string) string {
	return fmt.Sprintf(`
data "openstack_networking_agent_v2" "agent" {
  agent_type = "BGP dynamic routing agent"
  alive = %s
}`, alive)
}

func testAccOpenStackNetworkingAgentV2DataSourceAZ(az string) string {
	return fmt.Sprintf(`
data "openstack_networking_agent_v2" "agent" {
  agent_type = "BGP dynamic routing agent"
  availability_zone = "%s"
}`, az)
}

func testAccOpenStackNetworkingAgentV2DataSourceBinary(binary string) string {
	return fmt.Sprintf(`
data "openstack_networking_agent_v2" "agent" {
  binary = "%s"
}`, binary)
}

func testAccOpenStackNetworkingAgentV2DataSourceDesc(desc string) string {
	return fmt.Sprintf(`
data "openstack_networking_agent_v2" "agent" {
  agent_type = "BGP dynamic routing agent"
  description = "%s"
}`, desc)
}

func testAccOpenStackNetworkingAgentV2DataSourceHost(host string) string {
	return fmt.Sprintf(`
data "openstack_networking_agent_v2" "agent" {
  agent_type = "BGP dynamic routing agent"
  host = "%s"
}`, host)
}

func testAccOpenStackNetworkingAgentV2DataSourceTopic(topic string) string {
	return fmt.Sprintf(`
data "openstack_networking_agent_v2" "agent" {
  topic = "%s"
}`, topic)
}
