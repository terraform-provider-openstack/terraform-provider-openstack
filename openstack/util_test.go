package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitExpandToMapStringString(t *testing.T) {
	metadata := map[string]interface{}{
		"contents": "junk",
	}

	expected := map[string]string{
		"contents": "junk",
	}

	actual := expandToMapStringString(metadata)
	assert.Equal(t, expected, actual)
}

func TestUnitExpandToStringSlice(t *testing.T) {
	data := []interface{}{"foo", "bar"}

	expected := []string{"foo", "bar"}

	actual := expandToStringSlice(data)
	assert.Equal(t, expected, actual)
}

func TestUnitCompatibleMicroversion(t *testing.T) {
	actual, err := compatibleMicroversion("min", "2.1.0", "2.5")
	assert.NotNil(t, err)
	assert.False(t, actual)

	actual, err = compatibleMicroversion("min", "2.1", "2.5.0")
	assert.NotNil(t, err)
	assert.False(t, actual)

	actual, err = compatibleMicroversion("minn", "2.1", "2.5")
	assert.NotNil(t, err)
	assert.False(t, actual)

	actual, err = compatibleMicroversion("min", "", "2.5")
	assert.Nil(t, err)
	assert.False(t, actual)

	actual, err = compatibleMicroversion("min", "2.1", "")
	assert.Nil(t, err)
	assert.False(t, actual)

	actual, err = compatibleMicroversion("min", "2.1", "2.5")
	assert.Nil(t, err)
	assert.True(t, actual)

	actual, err = compatibleMicroversion("min", "2.1", "3.5")
	assert.Nil(t, err)
	assert.False(t, actual)

	actual, err = compatibleMicroversion("min", "2.5", "2.1")
	assert.Nil(t, err)
	assert.False(t, actual)

	actual, err = compatibleMicroversion("max", "2.5", "2.1")
	assert.Nil(t, err)
	assert.True(t, actual)

	actual, err = compatibleMicroversion("min", "2.10", "2.17")
	assert.Nil(t, err)
	assert.True(t, actual)
}

func TestUnitMapDiffWithNilValues(t *testing.T) {
	oldData := map[string]interface{}{"a": "1", "b": "2"}
	newData := map[string]interface{}{"a": "1", "c": "3"}

	result := mapDiffWithNilValues(oldData, newData)

	assert.Equal(t, result["a"], "1")
	assert.Equal(t, result["b"], nil)
	assert.Equal(t, result["c"], "3")
	assert.Equal(t, len(result), 3)
}
