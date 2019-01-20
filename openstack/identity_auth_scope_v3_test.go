package openstack

import (
	"testing"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
	"github.com/stretchr/testify/assert"
)

func TestFlattenIdentityAuthScopeV3Roles(t *testing.T) {
	roles := []tokens.Role{
		{
			ID:   "1",
			Name: "foo",
		},
		{
			ID:   "2",
			Name: "bar",
		},
	}

	expected := []map[string]string{
		{
			"role_id":   "1",
			"role_name": "foo",
		},
		{
			"role_name": "bar",
			"role_id":   "2",
		},
	}

	actual := flattenIdentityAuthScopeV3Roles(roles)
	assert.Equal(t, expected, actual)
}
