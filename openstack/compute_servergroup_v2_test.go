package openstack

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/servergroups"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	thclient "github.com/gophercloud/gophercloud/v2/testhelper/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitComputeServerGroupV2CreateOpts(t *testing.T) {
	createOpts := ComputeServerGroupV2CreateOpts{
		servergroups.CreateOpts{
			Name:     "foo",
			Policies: []string{"affinity"},
		},
		map[string]string{
			"foo": "bar",
		},
	}

	expected := map[string]any{
		"server_group": map[string]any{
			"name":     "foo",
			"policies": []any{"affinity"},
			"foo":      "bar",
		},
	}

	actual, err := createOpts.ToServerGroupCreateMap()

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestUnitExpandComputeServerGroupV2PoliciesMicroversions(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	raw := []any{
		"affinity",
		"soft-anti-affinity",
		"soft-affinity",
		"custom-policy",
	}
	client := thclient.ServiceClient(fakeServer)

	expectedPolicies := []string{
		"affinity",
		"soft-anti-affinity",
		"soft-affinity",
		"custom-policy",
	}
	expectedMicroversion := "2.15"

	actualPolicies := expandComputeServerGroupV2Policies(client, raw)
	actualMicroversion := client.Microversion

	assert.Equal(t, expectedMicroversion, actualMicroversion)
	assert.Equal(t, expectedPolicies, actualPolicies)
}

func TestUnitExpandComputeServerGroupV2PoliciesMicroversionsLegacy(t *testing.T) {
	fakeServer := th.SetupHTTP()
	defer fakeServer.Teardown()

	raw := []any{
		"anti-affinity",
		"affinity",
	}
	client := thclient.ServiceClient(fakeServer)

	expectedPolicies := []string{
		"anti-affinity",
		"affinity",
	}
	expectedMicroversion := ""

	actualPolicies := expandComputeServerGroupV2Policies(client, raw)
	actualMicroversion := client.Microversion

	assert.Equal(t, expectedMicroversion, actualMicroversion)
	assert.Equal(t, expectedPolicies, actualPolicies)
}
