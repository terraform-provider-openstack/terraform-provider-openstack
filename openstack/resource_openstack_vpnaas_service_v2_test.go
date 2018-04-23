package openstack

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/services"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"regexp"
)

func TestAccServiceV2_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckVPN(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServiceV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccServiceV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceV2Exists(
						"openstack_vpnaas_service_v2.service_1", "", ""),
				),
			},
		},
	})
}

func TestAccServiceV2_update(t *testing.T) {
	errorRegExp, err := regexp.Compile("openstack_vpnaas_service_v2.service_1: 1 error")
	if err != nil {
		t.Error("Couldn't compile regular expression")
	}
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckVPN(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckServiceV2Destroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccServiceV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceV2Exists(
						"openstack_vpnaas_service_v2.service_1", "", ""),
				),
			},
			// We expect the update to fail because a service cannot be updated while it
			// is not being used in an IPSec site connection.
			resource.TestStep{
				Config:      testAccServiceV2_update,
				ExpectError: errorRegExp,
			},
		},
	})
}

func testAccCheckServiceV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_vpnaas_service" {
			continue
		}
		_, err = services.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Service (%s) still exists.", rs.Primary.ID)
		}
		if _, ok := err.(gophercloud.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckServiceV2Exists(n, name, description string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %s", err)
		}

		var found *services.Service

		found, err = services.Get(networkingClient, rs.Primary.ID).Extract()

		switch {
		case name != found.Name:
			err = fmt.Errorf("Expected name <%s>, but found <%s>", name, found.Name)
		case description != found.Description:
			err = fmt.Errorf("Expected description <%s>, but found <%s>", description, found.Description)
		}

		if err != nil {
			return err
		}

		return nil
	}
}

var testAccServiceV2_basic = fmt.Sprintf(`
	resource "openstack_networking_router_v2" "router_1" {
	  name = "router_1"
	  admin_state_up = "true"
	  external_network_id = "%s"
	}
	
	resource "openstack_vpnaas_service_v2" "service_1" {
		router_id = "${openstack_networking_router_v2.router_1.id}",
		admin_state_up = "false"
	}
	`, OS_EXTGW_ID)

var testAccServiceV2_update = fmt.Sprintf(`
	resource "openstack_networking_router_v2" "router_1" {
	  name = "router_1"
	  admin_state_up = "true"
	  external_network_id = "%s"
	}
	
	resource "openstack_vpnaas_service_v2" "service_1" {
		router_id = "${openstack_networking_router_v2.router_1.id}",
		admin_state_up = "true"
		description = "An updated VPN service"
	}
	`, OS_EXTGW_ID)
