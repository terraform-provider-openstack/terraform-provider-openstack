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
