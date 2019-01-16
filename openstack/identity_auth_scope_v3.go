package openstack

import (
	"github.com/gophercloud/gophercloud/openstack/identity/v3/tokens"
)

func flattenIdentityAuthScopeV3Roles(roles []tokens.Role) []map[string]string {
	var allRoles []map[string]string

	for _, r := range roles {
		allRoles = append(allRoles, map[string]string{
			"role_name": r.Name,
			"role_id":   r.ID,
		})

	}

	return allRoles
}
