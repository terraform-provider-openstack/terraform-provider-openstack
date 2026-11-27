package openstack

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/databases"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/instances"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/users"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func databaseInstanceV1CustomizeDiff(resourceSchema map[string]*schema.Schema) schema.CustomizeDiffFunc {
	return func(_ context.Context, d *schema.ResourceDiff, _ any) error {
		if d.Id() == "" || !d.HasChange("size") || !d.NewValueKnown("size") {
			return nil
		}

		// A replacement instance may have a smaller volume than the old one.
		for key, field := range resourceSchema {
			if field.ForceNew && d.HasChange(key) {
				return nil
			}
		}

		oldSize, newSize := d.GetChange("size")
		if newSize.(int) < oldSize.(int) {
			return fmt.Errorf("decreasing openstack_db_instance_v1 volume size is not supported: size cannot decrease from %d to %d GB", oldSize, newSize)
		}

		return nil
	}
}

func expandDatabaseInstanceV1Datastore(rawDatastore []any) instances.DatastoreOpts {
	v := rawDatastore[0].(map[string]any)
	datastore := instances.DatastoreOpts{
		Version: v["version"].(string),
		Type:    v["type"].(string),
	}

	return datastore
}

func expandDatabaseInstanceV1Networks(rawNetworks []any) []instances.NetworkOpts {
	networks := make([]instances.NetworkOpts, 0, len(rawNetworks))

	for _, v := range rawNetworks {
		network := v.(map[string]any)
		networks = append(networks, instances.NetworkOpts{
			UUID:      network["uuid"].(string),
			Port:      network["port"].(string),
			V4FixedIP: network["fixed_ip_v4"].(string),
			V6FixedIP: network["fixed_ip_v6"].(string),
		})
	}

	return networks
}

func expandDatabaseInstanceV1Databases(rawDatabases []any) databases.BatchCreateOpts {
	var dbs databases.BatchCreateOpts

	for _, v := range rawDatabases {
		db := v.(map[string]any)
		dbs = append(dbs, databases.CreateOpts{
			Name:    db["name"].(string),
			CharSet: db["charset"].(string),
			Collate: db["collate"].(string),
		})
	}

	return dbs
}

func expandDatabaseInstanceV1Users(rawUsers []any) users.BatchCreateOpts {
	var userList users.BatchCreateOpts

	for _, v := range rawUsers {
		user := v.(map[string]any)
		userList = append(userList, users.CreateOpts{
			Name:      user["name"].(string),
			Password:  user["password"].(string),
			Databases: expandInstanceV1UserDatabases(user["databases"].(*schema.Set).List()),
			Host:      user["host"].(string),
		})
	}

	return userList
}

// databaseInstanceV1StateRefreshFunc returns a retry.StateRefreshFunc
// that is used to watch a database instance.
func databaseInstanceV1StateRefreshFunc(ctx context.Context, client *gophercloud.ServiceClient, instanceID string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		i, err := instances.Get(ctx, client, instanceID).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return i, "DELETED", nil
			}

			return nil, "", err
		}

		if i.Status == "error" {
			return i, i.Status, errors.New("There was an error creating the database instance")
		}

		return i, i.Status, nil
	}
}

// databaseInstanceV1VolumeResizeStateRefreshFunc waits until the database
// instance is ready and reports the requested volume size.
func databaseInstanceV1VolumeResizeStateRefreshFunc(
	ctx context.Context,
	client *gophercloud.ServiceClient,
	instanceID string,
	expectedSize int,
) retry.StateRefreshFunc {
	return func() (any, string, error) {
		i, err := instances.Get(ctx, client, instanceID).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return i, "DELETED", nil
			}

			return nil, "", err
		}

		if strings.EqualFold(i.Status, "error") {
			return i, i.Status, errors.New("database instance volume resize failed")
		}

		if (i.Status != "ACTIVE" && i.Status != "HEALTHY") || i.Volume.Size != expectedSize {
			return i, "RESIZE", nil
		}

		return i, i.Status, nil
	}
}

func expandInstanceV1UserDatabases(v []any) databases.BatchCreateOpts {
	var dbs databases.BatchCreateOpts

	for _, db := range v {
		dbs = append(dbs, databases.CreateOpts{
			Name: db.(string),
		})
	}

	return dbs
}
