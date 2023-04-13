package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitResourceNetworkingQoSRuleV2BuildID(t *testing.T) {
	tests := map[string][]string{
		"2f22cc5f-3807-4433-9230-0558f9539d94/398b1a1c-8254-4ce9-8fc3-9d1994437e9d": {"2f22cc5f-3807-4433-9230-0558f9539d94", "398b1a1c-8254-4ce9-8fc3-9d1994437e9d"},
		"ecefc62f-b81b-4745-8711-07fc1095baf0/5010886f-d108-47e0-9f82-d43f27891067": {"ecefc62f-b81b-4745-8711-07fc1095baf0", "5010886f-d108-47e0-9f82-d43f27891067"},
		"a2b9d7c2-c49f-4da8-89dd-b8728ce88ffd/aaa7165b-a738-4314-a36c-a4a5e7b3f201": {"a2b9d7c2-c49f-4da8-89dd-b8728ce88ffd", "aaa7165b-a738-4314-a36c-a4a5e7b3f201"},
	}

	for expected, params := range tests {
		actual := resourceNetworkingQoSRuleV2BuildID(params[0], params[1])
		assert.Equal(t, expected, actual)
	}
}

func TestUnitResourceNetworkingQoSRuleV2ParseID(t *testing.T) {
	tests := map[string][]string{
		"2f22cc5f-3807-4433-9230-0558f9539d94/398b1a1c-8254-4ce9-8fc3-9d1994437e9d": {"2f22cc5f-3807-4433-9230-0558f9539d94", "398b1a1c-8254-4ce9-8fc3-9d1994437e9d"},
		"ecefc62f-b81b-4745-8711-07fc1095baf0/5010886f-d108-47e0-9f82-d43f27891067": {"ecefc62f-b81b-4745-8711-07fc1095baf0", "5010886f-d108-47e0-9f82-d43f27891067"},
		"a2b9d7c2-c49f-4da8-89dd-b8728ce88ffd/aaa7165b-a738-4314-a36c-a4a5e7b3f201": {"a2b9d7c2-c49f-4da8-89dd-b8728ce88ffd", "aaa7165b-a738-4314-a36c-a4a5e7b3f201"},
	}

	for params, expected := range tests {
		actualQoSPolicy, actualQoSRule, err := resourceNetworkingQoSRuleV2ParseID(params)
		assert.NoError(t, err)
		assert.Equal(t, expected[0], actualQoSPolicy)
		assert.Equal(t, expected[1], actualQoSRule)
	}
}
