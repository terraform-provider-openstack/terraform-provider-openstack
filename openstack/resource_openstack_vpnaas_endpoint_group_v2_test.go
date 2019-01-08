package openstack

import (
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/vpnaas/endpointgroups"
	"github.com/hashicorp/terraform/flatmap"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccGroupV2_basic(t *testing.T) {
	var group endpointgroups.EndpointGroup
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckVPN(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEndpointGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointGroupV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEndpointGroupV2Exists(
						"openstack_vpnaas_endpoint_group_v2.group_1", &group),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_endpoint_group_v2.group_1", "name", &group.Name),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_endpoint_group_v2.group_1", "type", &group.Type),
					testAccCheckEndpoints("openstack_vpnaas_endpoint_group_v2.group_1", &group.Endpoints),
				),
			},
		},
	})
}

func TestAccGroupV2_update(t *testing.T) {
	var group endpointgroups.EndpointGroup
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckVPN(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEndpointGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointGroupV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEndpointGroupV2Exists(
						"openstack_vpnaas_endpoint_group_v2.group_1", &group),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_endpoint_group_v2.group_1", "name", &group.Name),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_endpoint_group_v2.group_1", "type", &group.Type),
					testAccCheckEndpoints("openstack_vpnaas_endpoint_group_v2.group_1", &group.Endpoints),
				),
			},
			{
				Config: testAccEndpointGroupV2_update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEndpointGroupV2Exists(
						"openstack_vpnaas_endpoint_group_v2.group_1", &group),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_endpoint_group_v2.group_1", "name", &group.Name),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_endpoint_group_v2.group_1", "type", &group.Type),
					testAccCheckEndpoints("openstack_vpnaas_endpoint_group_v2.group_1", &group.Endpoints),
				),
			},
		},
	})
}

func testAccCheckEndpointGroupV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkingClient, err := config.networkingV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_vpnaas_group" {
			continue
		}
		_, err = endpointgroups.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("EndpointGroup (%s) still exists.", rs.Primary.ID)
		}
		if _, ok := err.(gophercloud.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckEndpointGroupV2Exists(n string, group *endpointgroups.EndpointGroup) resource.TestCheckFunc {
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

		var found *endpointgroups.EndpointGroup

		found, err = endpointgroups.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*group = *found

		return nil
	}
}

func testAccCheckEndpoints(n string, actual *[]string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		endpoints := flatmap.Expand(rs.Primary.Attributes, "endpoints")
		i := 0
		for _, raw := range endpoints.([]interface{}) {
			endpoint := raw.(string)
			if endpoint != (*actual)[i] {
				return fmt.Errorf("The endpoints did not match. Expected: %v but got %v", endpoint, (*actual)[i])
			}
			i++
		}
		return nil
	}
}

var testAccEndpointGroupV2_basic = `
	resource "openstack_vpnaas_endpoint_group_v2" "group_1" {
		name = "Group 1"
		type = "cidr"
		endpoints = ["10.3.0.0/24",
			"10.2.0.0/24",]
	}
`

var testAccEndpointGroupV2_update = `
	resource "openstack_vpnaas_endpoint_group_v2" "group_1" {
		name = "Updated Group 1"
		type = "cidr"
		endpoints = ["10.2.0.0/24",
			"10.3.0.0/24",]
	}
`
