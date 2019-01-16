package openstack

import (
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/roles"
	"github.com/gophercloud/gophercloud/pagination"
)

// Role assignments have no ID in OpenStack.
// Build an ID out of the IDs that make up the role assignment
func identityRoleAssignmentV3ID(domainID, projectID, groupID, userID, roleID string) string {
	return fmt.Sprintf("%s/%s/%s/%s/%s", domainID, projectID, groupID, userID, roleID)
}

func identityRoleAssignmentV3ParseID(roleAssignmentID string) (string, string, string, string, string) {
	split := strings.Split(roleAssignmentID, "/")
	return split[0], split[1], split[2], split[3], split[4]
}

func identityRoleAssignmentV3FindAssignment(identityClient *gophercloud.ServiceClient, id string) (roles.RoleAssignment, error) {
	domainID, projectID, groupID, userID, roleID := identityRoleAssignmentV3ParseID(id)

	var opts roles.ListAssignmentsOpts
	opts = roles.ListAssignmentsOpts{
		GroupID:        groupID,
		ScopeDomainID:  domainID,
		ScopeProjectID: projectID,
		UserID:         userID,
	}

	pager := roles.ListAssignments(identityClient, opts)
	var assignment roles.RoleAssignment

	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		assignmentList, err := roles.ExtractRoleAssignments(page)
		if err != nil {
			return false, err
		}

		for _, a := range assignmentList {
			if a.Role.ID == roleID {
				assignment = a
				return false, nil
			}
		}

		return true, nil
	})

	return assignment, err
}
