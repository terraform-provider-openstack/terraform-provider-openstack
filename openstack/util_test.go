package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitExpandToMapStringString(t *testing.T) {
	metadata := map[string]any{
		"contents": "junk",
	}

	expected := map[string]string{
		"contents": "junk",
	}

	actual := expandToMapStringString(metadata)
	assert.Equal(t, expected, actual)
}

func TestUnitExpandToStringSlice(t *testing.T) {
	data := []any{"foo", "bar"}

	expected := []string{"foo", "bar"}

	actual := expandToStringSlice(data)
	assert.Equal(t, expected, actual)
}

func TestUnitCompatibleMicroversion(t *testing.T) {
	actual, err := compatibleMicroversion("min", "2.1.0", "2.5")
	require.Error(t, err)
	assert.False(t, actual)

	actual, err = compatibleMicroversion("min", "2.1", "2.5.0")
	require.Error(t, err)
	assert.False(t, actual)

	actual, err = compatibleMicroversion("minn", "2.1", "2.5")
	require.Error(t, err)
	assert.False(t, actual)

	actual, err = compatibleMicroversion("min", "", "2.5")
	require.NoError(t, err)
	assert.False(t, actual)

	actual, err = compatibleMicroversion("min", "2.1", "")
	require.NoError(t, err)
	assert.False(t, actual)

	actual, err = compatibleMicroversion("min", "2.1", "2.5")
	require.NoError(t, err)
	assert.True(t, actual)

	actual, err = compatibleMicroversion("min", "2.1", "3.5")
	require.NoError(t, err)
	assert.False(t, actual)

	actual, err = compatibleMicroversion("min", "2.5", "2.1")
	require.NoError(t, err)
	assert.False(t, actual)

	actual, err = compatibleMicroversion("max", "2.5", "2.1")
	require.NoError(t, err)
	assert.True(t, actual)

	actual, err = compatibleMicroversion("min", "2.10", "2.17")
	require.NoError(t, err)
	assert.True(t, actual)
}

func TestUnitMapDiffWithNilValues(t *testing.T) {
	oldData := map[string]any{"a": "1", "b": "2"}
	newData := map[string]any{"a": "1", "c": "3"}

	result := mapDiffWithNilValues(oldData, newData)

	assert.Equal(t, "1", result["a"])
	assert.Nil(t, result["b"])
	assert.Equal(t, "3", result["c"])
	assert.Len(t, result, 3)
}

func TestUnitBuildRequestBoolType(t *testing.T) {
	v := SubnetCreateOpts{
		ValueSpecs: map[string]string{
			"key1": "value1",
			"key2": "true",
			"key3": "false",
		},
	}

	req, err := BuildRequest(v, "")
	require.NoError(t, err)

	expected := map[string]any{
		"": map[string]any{
			"key1":       "value1",
			"key2":       true,
			"key3":       false,
			"network_id": "",
		},
	}
	assert.Equal(t, expected, req)
}

func TestUnitparsePairedIDs(t *testing.T) {
	id := "foo/bar"
	expectedParentID := "foo"
	expectedChildID := "bar"

	actualParentID, actualChildID, err := parsePairedIDs(id, "")
	require.NoError(t, err)
	assert.Equal(t, expectedParentID, actualParentID)
	assert.Equal(t, expectedChildID, actualChildID)
}
