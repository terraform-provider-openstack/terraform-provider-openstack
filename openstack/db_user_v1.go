package openstack

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/databases"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/users"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func expandDatabaseUserV1Databases(rawDatabases []any) databases.BatchCreateOpts {
	var dbs databases.BatchCreateOpts

	for _, db := range rawDatabases {
		dbs = append(dbs, databases.CreateOpts{
			Name: db.(string),
		})
	}

	return dbs
}

func flattenDatabaseUserV1Databases(dbs []databases.Database) []string {
	databases := make([]string, 0, len(dbs))
	for _, db := range dbs {
		databases = append(databases, db.Name)
	}

	return databases
}

// databaseUserV1StateRefreshFunc returns a retry.StateRefreshFunc that is used to watch db user.
func databaseUserV1StateRefreshFunc(ctx context.Context, client *gophercloud.ServiceClient, instanceID string, userName string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		pages, err := users.List(client, instanceID).AllPages(ctx)
		if err != nil {
			return nil, "", fmt.Errorf("Unable to retrieve OpenStack database users: %w", err)
		}

		allUsers, err := users.ExtractUsers(pages)
		if err != nil {
			return nil, "", fmt.Errorf("Unable to extract OpenStack database users: %w", err)
		}

		for _, v := range allUsers {
			if v.Name == userName {
				return v, "ACTIVE", nil
			}
		}

		return nil, "BUILD", nil
	}
}

// databaseUserV1Exists is used to check whether user exists on particular database instance.
func databaseUserV1Exists(ctx context.Context, client *gophercloud.ServiceClient, instanceID string, userName string) (bool, users.User, error) {
	var exists bool

	var err error

	var userObj users.User

	pages, err := users.List(client, instanceID).AllPages(ctx)
	if err != nil {
		return exists, userObj, err
	}

	allUsers, err := users.ExtractUsers(pages)
	if err != nil {
		return exists, userObj, err
	}

	for _, v := range allUsers {
		if v.Name == userName {
			exists = true

			return exists, v, nil
		}
	}

	return false, userObj, err
}
