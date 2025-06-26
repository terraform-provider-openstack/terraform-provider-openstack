package openstack

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/vpnaas/services"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccServiceVPNaaSV2_basic(t *testing.T) {
	var service services.Service

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckNonAdminOnly(t)
			testAccPreCheckVPN(t)
			t.Skip("Currently failing in GH-A")
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckServiceV2Destroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccServiceV2Basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServiceV2Exists(t.Context(),
						"openstack_vpnaas_service_v2.service_1", &service),
					resource.TestCheckResourceAttrPtr("openstack_vpnaas_service_v2.service_1", "router_id", &service.RouterID),
					resource.TestCheckResourceAttr("openstack_vpnaas_service_v2.service_1", "admin_state_up", strconv.FormatBool(service.AdminStateUp)),
				),
			},
		},
	})
}

func testAccCheckServiceV2Destroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		networkingClient, err := config.NetworkingV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_vpnaas_service" {
				continue
			}

			_, err = services.Get(ctx, networkingClient, rs.Primary.ID).Extract()
			if err == nil {
				return fmt.Errorf("Service (%s) still exists", rs.Primary.ID)
			}

			if !gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return err
			}
		}

		return nil
	}
}

func testAccCheckServiceV2Exists(ctx context.Context, n string, serv *services.Service) resource.TestCheckFunc {
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

		var found *services.Service

		found, err = services.Get(ctx, networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		*serv = *found

		return nil
	}
}

func testAccServiceV2Basic() string {
	return fmt.Sprintf(`
	resource "openstack_networking_router_v2" "router_1" {
	  name = "router_1"
	  admin_state_up = "true"
	  external_network_id = "%s"
	}

	resource "openstack_vpnaas_service_v2" "service_1" {
		router_id = openstack_networking_router_v2.router_1.id
		admin_state_up = "false"
	}
	`, osExtGwID)
}
