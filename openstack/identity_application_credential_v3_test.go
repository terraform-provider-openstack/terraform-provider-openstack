package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/applicationcredentials"
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

	roles := []interface{}{role1.Name, role2.Name}

	expected := []applicationcredentials.Role{role1, role2}

	actual := expandIdentityApplicationCredentialRolesV3(roles)
	assert.Equal(t, expected, actual)
}
