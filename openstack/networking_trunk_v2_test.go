package openstack

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/trunks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestUnitFlattenNetworkingTrunkSubportsV2(t *testing.T) {
	subports := []trunks.Subport{
		{
			PortID:           "port_id_1",
			SegmentationType: "type_1",
			SegmentationID:   111,
		},
		{
			PortID:           "port_id_2",
			SegmentationType: "type_2",
			SegmentationID:   222,
		},
	}

	expectedSubports := []map[string]any{
		{
			"port_id":           "port_id_1",
			"segmentation_id":   111,
			"segmentation_type": "type_1",
		},
		{
			"port_id":           "port_id_2",
			"segmentation_id":   222,
			"segmentation_type": "type_2",
		},
	}

	actualSubports := flattenNetworkingTrunkV2Subports(subports)

	assert.ElementsMatch(t, expectedSubports, actualSubports)
}

func TestUnitExpandNetworkingTrunkSubportsV2(t *testing.T) {
	r := resourceNetworkingTrunkV2()
	d := r.TestResourceData()
	d.SetId("1")

	subport1 := map[string]any{
		"port_id":           "port_id_1",
		"segmentation_id":   111,
		"segmentation_type": "type_1",
	}
	subport2 := map[string]any{
		"port_id":           "port_id_2",
		"segmentation_id":   222,
		"segmentation_type": "type_2",
	}
	subports := []map[string]any{subport1, subport2}
	d.Set("sub_port", subports)

	expectedSubports := []trunks.Subport{
		{
			PortID:           "port_id_1",
			SegmentationType: "type_1",
			SegmentationID:   111,
		},
		{
			PortID:           "port_id_2",
			SegmentationType: "type_2",
			SegmentationID:   222,
		},
	}

	actualSubports := expandNetworkingTrunkV2Subports(d.Get("sub_port").(*schema.Set))

	assert.ElementsMatch(t, expectedSubports, actualSubports)
}

func TestUnitExpandNetworkingTrunkSubportsRemoveV2(t *testing.T) {
	r := resourceNetworkingTrunkV2()
	d := r.TestResourceData()
	d.SetId("1")

	subport1 := map[string]any{
		"port_id":           "port_id_3",
		"segmentation_id":   333,
		"segmentation_type": "type_3",
	}
	subport2 := map[string]any{
		"port_id":           "port_id_4",
		"segmentation_id":   444,
		"segmentation_type": "type_4",
	}
	subports := []map[string]any{subport1, subport2}
	d.Set("sub_port", subports)

	expectedRemoveSubports := []trunks.RemoveSubport{
		{
			PortID: "port_id_3",
		},
		{
			PortID: "port_id_4",
		},
	}

	actualRemoveSubports := expandNetworkingTrunkV2SubportsRemove(d.Get("sub_port").(*schema.Set))

	assert.ElementsMatch(t, expectedRemoveSubports, actualRemoveSubports)
}
