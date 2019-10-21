package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandToMapStringString(t *testing.T) {
	metadata := map[string]interface{}{
		"contents": "junk",
	}

	expected := map[string]string{
		"contents": "junk",
	}

	actual := expandToMapStringString(metadata)
	assert.Equal(t, expected, actual)
}

func TestExpandToStringSlice(t *testing.T) {
	data := []interface{}{"foo", "bar"}

	expected := []string{"foo", "bar"}

	actual := expandToStringSlice(data)
	assert.Equal(t, expected, actual)
}

func TestCompatibleMicroversion(t *testing.T) {
	actual, err := compatibleMicroversion("min", "2.1.0", "2.5")
	assert.NotNil(t, err)

	actual, err = compatibleMicroversion("min", "2.1", "2.5.0")
	assert.NotNil(t, err)

	actual, err = compatibleMicroversion("minn", "2.1", "2.5")
	assert.NotNil(t, err)

	actual, err = compatibleMicroversion("min", "", "2.5")
	assert.Nil(t, err)
	assert.Equal(t, false, actual)

	actual, err = compatibleMicroversion("min", "2.1", "")
	assert.Nil(t, err)
	assert.Equal(t, false, actual)

	actual, err = compatibleMicroversion("min", "2.1", "2.5")
	assert.Nil(t, err)
	assert.Equal(t, true, actual)

	actual, err = compatibleMicroversion("min", "2.1", "3.5")
	assert.Nil(t, err)
	assert.Equal(t, false, actual)

	actual, err = compatibleMicroversion("min", "2.5", "2.1")
	assert.Nil(t, err)
	assert.Equal(t, false, actual)

	actual, err = compatibleMicroversion("max", "2.5", "2.1")
	assert.Nil(t, err)
	assert.Equal(t, true, actual)

	actual, err = compatibleMicroversion("min", "2.10", "2.17")
	assert.Nil(t, err)
	assert.Equal(t, true, actual)
}
