package openstack

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/applicationcredentials"
	"github.com/stretchr/testify/assert"
)

func TestUnitFlattenIdentityApplicationCredentialRolesV3(t *testing.T) {
	role1 := applicationcredentials.Role{
		ID:   "123",
		Name: "foo",
	}
	role2 := applicationcredentials.Role{
		ID:   "321",
		Name: "bar",
	}

	roles := []applicationcredentials.Role{role1, role2}

	expected := []string{"foo", "bar"}

	actual := flattenIdentityApplicationCredentialRolesV3(roles)
	assert.Equal(t, expected, actual)
}

func TestUnitExpandIdentityApplicationCredentialRolesV3(t *testing.T) {
	role1 := applicationcredentials.Role{
		Name: "foo",
	}
	role2 := applicationcredentials.Role{
		Name: "bar",
	}

	roles := []any{role1.Name, role2.Name}

	expected := []applicationcredentials.Role{role1, role2}

	actual := expandIdentityApplicationCredentialRolesV3(roles)
	assert.Equal(t, expected, actual)
}
