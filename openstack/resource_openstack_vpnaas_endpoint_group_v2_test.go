package openstack

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/vpnaas/endpointgroups"
)

func TestAccGroupVPNaaSV2_basic(t *testing.T) {
	var group endpointgroups.EndpointGroup
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckVPN(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckEndpointGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointGroupV2Basic,
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

func TestAccGroupVPNaaSV2_update(t *testing.T) {
	var group endpointgroups.EndpointGroup
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckVPN(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckEndpointGroupV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEndpointGroupV2Basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEndpointGroupV2Exists(
						"openstack_vpnaas_endpoint_group_v2.group_1", &group),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_endpoint_group_v2.group_1", "name", &group.Name),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_endpoint_group_v2.group_1", "type", &group.Type),
					testAccCheckEndpoints("openstack_vpnaas_endpoint_group_v2.group_1", &group.Endpoints),
				),
			},
			{
				Config: testAccEndpointGroupV2Update,
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
	networkingClient, err := config.NetworkingV2Client(context.TODO(), osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_vpnaas_group" {
			continue
		}
		_, err = endpointgroups.Get(context.TODO(), networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("EndpointGroup (%s) still exists", rs.Primary.ID)
		}
		if !gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
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
		networkingClient, err := config.NetworkingV2Client(context.TODO(), osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %s", err)
		}

		var found *endpointgroups.EndpointGroup

		found, err = endpointgroups.Get(context.TODO(), networkingClient, rs.Primary.ID).Extract()
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
		var endpointsList []string
		// Find all "endpoints.<number>" keys and collect the values.
		// The <number> values are seemingly random and very large.
		for k, v := range rs.Primary.Attributes {
			println("[DEBUG] key:", k, "value:", v)
			if strings.HasPrefix(k, "endpoints.") && k[10] >= '0' && k[10] <= '9' {
				endpointsList = append(endpointsList, v)
			}
		}

		if len(*actual) != len(endpointsList) {
			return fmt.Errorf("The number of endpoints did not match. Expected: %v but got %v", len(*actual), len(endpointsList))
		}

		sort.Strings(endpointsList)
		sort.Strings(*actual)

		for i, endpoint := range endpointsList {
			if endpoint != (*actual)[i] {
				return fmt.Errorf("The endpoints did not match. Expected: '%v' but got '%v'", endpoint, (*actual)[i])
			}
		}
		return nil
	}
}

const testAccEndpointGroupV2Basic = `
	resource "openstack_vpnaas_endpoint_group_v2" "group_1" {
		name = "Group 1"
		type = "cidr"
		endpoints = ["10.3.0.0/24",
			"10.2.0.0/24",]
	}
`

const testAccEndpointGroupV2Update = `
	resource "openstack_vpnaas_endpoint_group_v2" "group_1" {
		name = "Updated Group 1"
		type = "cidr"
		endpoints = ["10.2.0.0/24",
			"10.3.0.0/24",]
	}
`
