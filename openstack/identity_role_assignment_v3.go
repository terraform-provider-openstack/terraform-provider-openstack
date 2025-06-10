package openstack

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/roles"
	"github.com/gophercloud/gophercloud/v2/pagination"
)

// Role assignments have no ID in OpenStack.
// Build an ID out of the IDs that make up the role assignment.
func identityRoleAssignmentV3ID(domainID, projectID, groupID, userID, roleID string) string {
	return fmt.Sprintf("%s/%s/%s/%s/%s", domainID, projectID, groupID, userID, roleID)
}

func identityRoleAssignmentV3ParseID(roleAssignmentID string) (string, string, string, string, string, error) {
	split := strings.Split(roleAssignmentID, "/")

	if len(split) != 5 {
		return "", "", "", "", "", fmt.Errorf("Malformed ID: %s", roleAssignmentID)
	}

	return split[0], split[1], split[2], split[3], split[4], nil
}

func identityRoleAssignmentV3FindAssignment(ctx context.Context, identityClient *gophercloud.ServiceClient, id string) (roles.RoleAssignment, error) {
	var assignment roles.RoleAssignment

	domainID, projectID, groupID, userID, roleID, err := identityRoleAssignmentV3ParseID(id)
	if err != nil {
		return assignment, err
	}

	opts := roles.ListAssignmentsOnResourceOpts{
		GroupID:   groupID,
		DomainID:  domainID,
		ProjectID: projectID,
		UserID:    userID,
	}

	pager := roles.ListAssignmentsOnResource(identityClient, opts)

	found := false
	err = pager.EachPage(ctx, func(_ context.Context, page pagination.Page) (bool, error) {
		assignmentList, err := roles.ExtractRoles(page)
		if err != nil {
			return false, err
		}

		for _, a := range assignmentList {
			if a.ID == roleID {
				found = true
				assignment = roles.RoleAssignment{
					Role: roles.AssignedRole{
						ID: a.ID,
					},
					Scope: roles.Scope{
						Domain: roles.Domain{
							ID: domainID,
						},
						Project: roles.Project{
							ID: projectID,
						},
					},
					User: roles.User{
						ID: userID,
					},
					Group: roles.Group{
						ID: groupID,
					},
				}

				return false, nil
			}
		}

		return true, nil
	})

	if !found {
		return assignment, gophercloud.ErrUnexpectedResponseCode{Actual: http.StatusNotFound}
	}

	return assignment, err
}
