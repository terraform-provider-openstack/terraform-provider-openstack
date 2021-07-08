package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/servergroups"
	th "github.com/gophercloud/gophercloud/testhelper"
	thclient "github.com/gophercloud/gophercloud/testhelper/client"
)

func TestComputeServerGroupV2CreateOpts(t *testing.T) {
	createOpts := ComputeServerGroupV2CreateOpts{
		servergroups.CreateOpts{
			Name:     "foo",
			Policies: []string{"affinity"},
		},
		map[string]string{
			"foo": "bar",
		},
	}

	expected := map[string]interface{}{
		"server_group": map[string]interface{}{
			"name":     "foo",
			"policies": []interface{}{"affinity"},
			"foo":      "bar",
		},
	}

	actual, err := createOpts.ToServerGroupCreateMap()

	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestExpandComputeServerGroupV2PoliciesMicroversions(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	raw := []interface{}{
		"affinity",
		"soft-anti-affinity",
		"soft-affinity",
		"custom-policy",
	}
	client := thclient.ServiceClient()

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

func TestExpandComputeServerGroupV2PoliciesMicroversionsLegacy(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	raw := []interface{}{
		"anti-affinity",
		"affinity",
	}
	client := thclient.ServiceClient()

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
