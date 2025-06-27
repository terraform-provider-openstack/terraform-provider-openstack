package openstack

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/segments"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccNetworkingV2Segment_basic(t *testing.T) {
	var segment segments.Segment

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckNetworkingV2SegmentDestroy(t.Context()),
		Steps: []resource.TestStep{
			{
				Config: testAccNetworkingV2SegmentBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2SegmentExists(t.Context(), "openstack_networking_segment_v2.segment_1", &segment),
					resource.TestCheckResourceAttr(
						"openstack_networking_segment_v2.segment_1", "name", "segment_1"),
					resource.TestCheckResourceAttr(
						"openstack_networking_segment_v2.segment_1", "description", "my segment description"),
				),
			},
			{
				Config: testAccNetworkingV2SegmentUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"openstack_networking_segment_v2.segment_1", "name", ""),
					resource.TestCheckResourceAttr(
						"openstack_networking_segment_v2.segment_1", "description", ""),
					resource.TestCheckResourceAttrPtr(
						"openstack_networking_subnet_v2.subnet_1", "segment_id", &segment.ID),
				),
			},
		},
	})
}

func testAccCheckNetworkingV2SegmentDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		config := testAccProvider.Meta().(*Config)

		networkingClient, err := config.NetworkingV2Client(ctx, osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack networking client: %w", err)
		}

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "openstack_networking_network_v2" {
				continue
			}

			_, err := segments.Get(ctx, networkingClient, rs.Primary.ID).Extract()
			if err == nil {
				return errors.New("Segment still exists")
			}
		}

		return nil
	}
}

func testAccCheckNetworkingV2SegmentExists(ctx context.Context, n string, segment *segments.Segment) resource.TestCheckFunc {
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

		found, err := segments.Get(ctx, networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return errors.New("Segment not found")
		}

		*segment = *found

		return nil
	}
}

const testAccNetworkingV2SegmentBasic = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_1"
  description = "my segment description"
  admin_state_up = "true"
}

resource "openstack_networking_segment_v2" "segment_1" {
  name = "segment_1"
  description = "my segment description"
  network_id = openstack_networking_network_v2.network_1.id
  network_type = "geneve"
}
`

const testAccNetworkingV2SegmentUpdate = `
resource "openstack_networking_network_v2" "network_1" {
  name = "network_2"
  admin_state_up = "true"
}

resource "openstack_networking_segment_v2" "segment_1" {
  network_id = openstack_networking_network_v2.network_1.id
  network_type = "geneve"
}

resource "openstack_networking_subnet_v2" "subnet_1" {
  name = "subnet_1"
  cidr = "192.168.10.0/24"
  ip_version = 4
  network_id = openstack_networking_network_v2.network_1.id
  segment_id = openstack_networking_segment_v2.segment_1.id
}
`
