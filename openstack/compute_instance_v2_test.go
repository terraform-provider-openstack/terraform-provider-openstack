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
		// 0
		{
			old:      ":foo:123",
			new:      "az1",
			suppress: true,
		},
		// 1
		{
			old:      "az:bar321:",
			new:      "az",
			suppress: true,
		},
		// 2
		{
			old:      "az",
			new:      "az:baz:456",
			suppress: true,
		},
		// 3
		{
			old:      "az1:baz:654",
			new:      "az2:baz:456",
			suppress: false,
		},
		// 4
		{
			old:      ":baz:654",
			new:      "az2:baz:456",
			suppress: true,
		},
		// 5
		{
			old:      "az:baz:654",
			new:      ":baz:456",
			suppress: false,
		},
		// 6
		{
			old:      "az:qux:",
			new:      "az1",
			suppress: false,
		},
		// 7
		{
			old:      "az1",
			new:      "az2",
			suppress: false,
		},
		// 8
		{
			old:      "::",
			new:      "az",
			suppress: false,
		},
		// 9
		{
			old:      ":foo:",
			new:      "az",
			suppress: true,
		},
		// 10
		{
			old:      ":::::",
			new:      "az",
			suppress: false,
		},
	}

	for i, testCase := range testCases {
		assert.Equal(t, testCase.suppress, suppressAvailabilityZoneDetailDiffs("", testCase.old, testCase.new, nil), fmt.Sprintf("Test case: %d", i))
	}
}
