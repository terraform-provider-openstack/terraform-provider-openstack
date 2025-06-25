package openstack

import (
	"errors"
	"fmt"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/networks"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccOpenStackNetworkingSegmentV2DataSource_basic(t *testing.T) {
	var network networks.Network

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAdminOnly(t)
		},
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccOpenStackNetworkingSegmentV2DataSourceSegment,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkingV2NetworkExists(t.Context(), "openstack_networking_network_v2.network_1", &network),
					testAccCheckNetworkingSegmentV2DataSourceID("data.openstack_networking_segment_v2.segment_1"),
					resource.TestCheckResourceAttrPtr(
						"data.openstack_networking_segment_v2.segment_1", "network_id", &network.ID),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_segment_v2.segment_1", "name", ""),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_segment_v2.segment_1", "description", ""),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_segment_v2.segment_1", "network_type", "geneve"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_networking_segment_v2.segment_1", "segmentation_id"),
					// public network segment
					testAccCheckNetworkingSegmentV2DataSourceID("data.openstack_networking_segment_v2.segment_2"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_segment_v2.segment_2", "name", ""),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_segment_v2.segment_2", "description", ""),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_segment_v2.segment_2", "network_type", "flat"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_segment_v2.segment_2", "physical_network", "public"),
					resource.TestCheckResourceAttrSet(
						"data.openstack_networking_segment_v2.segment_2", "segmentation_id"),
					// check data source by segment_id
					testAccCheckNetworkingSegmentV2DataSourceID("data.openstack_networking_segment_v2.segment_3"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_segment_v2.segment_3", "name", "tf_test_segment"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_segment_v2.segment_3", "description", "my segment description"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_segment_v2.segment_3", "network_type", "local"),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_segment_v2.segment_3", "segmentation_id", "0"),
					resource.TestCheckResourceAttrPtr(
						"data.openstack_networking_segment_v2.segment_3", "network_id", &network.ID),
					resource.TestCheckResourceAttr(
						"data.openstack_networking_segment_v2.segment_3", "physical_network", ""),
				),
			},
		},
	})
}

func testAccCheckNetworkingSegmentV2DataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find segment data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("Segment data source ID not set")
		}

		return nil
	}
}

const testAccOpenStackNetworkingSegmentV2DataSourceSegment = `
resource "openstack_networking_network_v2" "network_1" {
  name = "tf_test_network"
  description = "my network description"
}

resource "openstack_networking_segment_v2" "segment_1" {
  network_id = openstack_networking_network_v2.network_1.id
  name = "tf_test_segment"
  description = "my segment description"
  network_type = "local"
}

data "openstack_networking_segment_v2" "segment_1" {
  network_id = openstack_networking_network_v2.network_1.id
  network_type = "geneve"
}

data "openstack_networking_segment_v2" "segment_2" {
  physical_network = "public"
}

data "openstack_networking_segment_v2" "segment_3" {
  segment_id = openstack_networking_segment_v2.segment_1.id
}
`
