package openstack

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/containerinfra/v1/clustertemplates"

	"github.com/stretchr/testify/assert"
)

func TestExpandContainerInfraV1LabelsMap(t *testing.T) {
	labels := map[string]interface{}{
		"foo": "bar",
		"bar": "baz",
	}

	expectedLabels := map[string]string{
		"foo": "bar",
		"bar": "baz",
	}

	actualLabels, err := expandContainerInfraV1LabelsMap(labels)
	assert.Equal(t, err, nil)
	assert.Equal(t, expectedLabels, actualLabels)
}

func TestExpandContainerInfraV1LabelsString(t *testing.T) {
	labels := map[string]interface{}{
		"foo": "bar",
		"bar": "baz",
	}

	expectedLabels_1 := "foo=bar,bar=baz"
	expectedLabels_2 := "bar=baz,foo=bar"

	actualLabels, err := expandContainerInfraV1LabelsString(labels)
	assert.Equal(t, err, nil)

	if actualLabels != expectedLabels_1 && actualLabels != expectedLabels_2 {
		t.Fatalf("Unexpected labels. Got %s, expected %s or %s",
			actualLabels, expectedLabels_1, expectedLabels_2)
	}

}

func TestContainerInfraClusterTemplateV1AppendUpdateOpts(t *testing.T) {
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
