package openstack

import (
	"context"
	"fmt"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/databases"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// databaseDatabaseV1StateRefreshFunc returns a retry.StateRefreshFunc
// that is used to watch a database.
func databaseDatabaseV1StateRefreshFunc(ctx context.Context, client *gophercloud.ServiceClient, instanceID string, dbName string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		pages, err := databases.List(client, instanceID).AllPages(ctx)
		if err != nil {
			return nil, "", fmt.Errorf("Unable to retrieve OpenStack databases: %w", err)
		}

		allDatabases, err := databases.ExtractDBs(pages)
		if err != nil {
			return nil, "", fmt.Errorf("Unable to extract OpenStack databases: %w", err)
		}

		for _, v := range allDatabases {
			if v.Name == dbName {
				return v, "ACTIVE", nil
			}
		}

		return nil, "BUILD", nil
	}
}

func databaseDatabaseV1Exists(ctx context.Context, client *gophercloud.ServiceClient, instanceID string, dbName string) (bool, error) {
	var exists bool

	var err error

	pages, err := databases.List(client, instanceID).AllPages(ctx)
	if err != nil {
		return exists, err
	}

	allDatabases, err := databases.ExtractDBs(pages)
	if err != nil {
		return exists, err
	}

	for _, v := range allDatabases {
		if v.Name == dbName {
			exists = true

			return exists, err
		}
	}

	return false, nil
}
