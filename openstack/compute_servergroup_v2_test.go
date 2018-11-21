package openstack

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/servergroups"
	th "github.com/gophercloud/gophercloud/testhelper"
	thclient "github.com/gophercloud/gophercloud/testhelper/client"
	"github.com/stretchr/testify/assert"
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

func TestExpandComputeServerGroupV2Policies(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	raw := []interface{}{
		"affinity",
	}
	client := thclient.ServiceClient()

	expectedPolicies := []string{
		"affinity",
	}
	expectedMicroversion := ""

	actualPolicies := expandComputeServerGroupV2Policies(client, raw)
	actualMicroversion := client.Microversion

	assert.Equal(t, expectedMicroversion, actualMicroversion)
	assert.Equal(t, expectedPolicies, actualPolicies)
}

func TestExpandComputeServerGroupV2PoliciesMicroversions(t *testing.T) {
	th.SetupHTTP()
	defer th.TeardownHTTP()
	raw := []interface{}{
		"affinity",
		"soft-anti-affinity",
		"soft-affinity",
	}
	client := thclient.ServiceClient()

	expectedPolicies := []string{
		"affinity",
		"soft-anti-affinity",
		"soft-affinity",
	}
	expectedMicroversion := "2.15"

	actualPolicies := expandComputeServerGroupV2Policies(client, raw)
	actualMicroversion := client.Microversion

	assert.Equal(t, expectedMicroversion, actualMicroversion)
	assert.Equal(t, expectedPolicies, actualPolicies)
}
