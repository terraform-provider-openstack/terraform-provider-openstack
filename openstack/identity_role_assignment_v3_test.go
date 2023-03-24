package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitIdentityRoleAssignmentV3ID(t *testing.T) {
	domainID := "domain"
	projectID := "project"
	groupID := "group"
	userID := "user"
	roleID := "role"

	expected := "domain/project/group/user/role"
	actual := identityRoleAssignmentV3ID(domainID, projectID, groupID, userID, roleID)
	assert.Equal(t, expected, actual)
}

func TestUnitIdentityRoleAssignmentV3ParseID(t *testing.T) {
	id := "domain/project/group/user/role"

	expectedDomainID := "domain"
	expectedProjectID := "project"
	expectedGroupID := "group"
	expectedUserID := "user"
	expectedRoleID := "role"

	actualDomainID, actualProjectID, actualGroupID, actualUserID, actualRoleID, err := identityRoleAssignmentV3ParseID(id)
	assert.Equal(t, err, nil)
	assert.Equal(t, expectedDomainID, actualDomainID)
	assert.Equal(t, expectedProjectID, actualProjectID)
	assert.Equal(t, expectedGroupID, actualGroupID)
	assert.Equal(t, expectedUserID, actualUserID)
	assert.Equal(t, expectedRoleID, actualRoleID)
}
