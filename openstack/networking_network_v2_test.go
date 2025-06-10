package openstack

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
)

func TestUnitExpandNetworkingNetworkSegmentsV2(t *testing.T) {
	r := resourceNetworkingNetworkV2()
	d := r.TestResourceData()
	d.SetId("1")

	segments1 := map[string]any{
		"physical_network": "aaa",
		"network_type":     "type11",
		"segmentation_id":  11,
	}
	segments2 := map[string]any{
		"physical_network": "bbb",
		"network_type":     "type12",
		"segmentation_id":  12,
	}
	segments := []map[string]any{segments1, segments2}
	d.Set("segments", segments)

	expectedSegments := []provider.Segment{
		{
			PhysicalNetwork: "aaa",
			NetworkType:     "type11",
			SegmentationID:  11,
		},
		{
			PhysicalNetwork: "bbb",
			NetworkType:     "type12",
			SegmentationID:  12,
		},
	}

	actualSegments := expandNetworkingNetworkSegmentsV2(d.Get("segments").(*schema.Set))

	assert.ElementsMatch(t, expectedSegments, actualSegments)
}
