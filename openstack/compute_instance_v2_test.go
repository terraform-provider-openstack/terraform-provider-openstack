package openstack

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuppressAvailabilityZoneDetailDiffs(t *testing.T) {
	testCases := []struct {
		old      string
		new      string
		suppress bool
	}{
		{
			old:      ":foo:123",
			new:      "az1",
			suppress: true,
		},
		{
			old:      "az:bar321:",
			new:      "az",
			suppress: true,
		},
		{
			old:      "az",
			new:      "az:baz:456",
			suppress: true,
		},
		{
			old:      "az1:baz:654",
			new:      "az2:baz:456",
			suppress: false,
		},
		{
			old:      ":baz:654",
			new:      "az2:baz:456",
			suppress: true,
		},
		{
			old:      "az:baz:654",
			new:      ":baz:456",
			suppress: false,
		},
		{
			old:      "az:qux:",
			new:      "az1",
			suppress: false,
		},
		{
			old:      "az1",
			new:      "az2",
			suppress: false,
		},
	}

	for i, testCase := range testCases {
		assert.Equal(t, testCase.suppress, suppressAvailabilityZoneDetailDiffs("", testCase.old, testCase.new, nil), fmt.Sprintf("Test case: %d", i))
	}
}
