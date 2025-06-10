package openstack

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/containerinfra/v1/clustertemplates"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitExpandContainerInfraV1LabelsMap(t *testing.T) {
	labels := map[string]any{
		"foo": "bar",
		"bar": "baz",
	}

	expectedLabels := map[string]string{
		"foo": "bar",
		"bar": "baz",
	}

	actualLabels, err := expandContainerInfraV1LabelsMap(labels)
	require.NoError(t, err)
	assert.Equal(t, expectedLabels, actualLabels)
}

func TestUnitExpandContainerInfraV1LabelsString(t *testing.T) {
	labels := map[string]any{
		"foo": "bar",
		"bar": "baz",
	}

	expectedLabels1 := "{'foo':'bar','bar':'baz'}"
	expectedLabels2 := "{'bar':'baz','foo':'bar'}"

	actualLabels, err := expandContainerInfraV1LabelsString(labels)
	require.NoError(t, err)

	if actualLabels != expectedLabels1 && actualLabels != expectedLabels2 {
		t.Fatalf("Unexpected labels. Got %s, expected %s or %s",
			actualLabels, expectedLabels1, expectedLabels2)
	}
}

func TestUnitContainerInfraClusterTemplateV1AppendUpdateOpts(t *testing.T) {
	actualUpdateOpts := []clustertemplates.UpdateOptsBuilder{}

	expectedUpdateOpts := []clustertemplates.UpdateOptsBuilder{
		clustertemplates.UpdateOpts{
			Op:    clustertemplates.ReplaceOp,
			Path:  "/master_lb_enabled",
			Value: "True",
		},
		clustertemplates.UpdateOpts{
			Op:    clustertemplates.ReplaceOp,
			Path:  "/registry_enabled",
			Value: "True",
		},
	}

	actualUpdateOpts = containerInfraClusterTemplateV1AppendUpdateOpts(
		actualUpdateOpts, "master_lb_enabled", "True")

	actualUpdateOpts = containerInfraClusterTemplateV1AppendUpdateOpts(
		actualUpdateOpts, "registry_enabled", "True")

	assert.Equal(t, expectedUpdateOpts, actualUpdateOpts)
}
