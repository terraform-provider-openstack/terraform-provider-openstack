package openstack

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/databases"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/instances"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/users"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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

// databaseInstanceV1CustomizeDiff returns a schema.CustomizeDiffFunc that
// refuses to shrink the volume of an existing instance. Trove can only grow a
// volume, and replacing the instance to get a smaller one would destroy its
// data. A smaller size is accepted when the instance is replaced anyway.
func databaseInstanceV1CustomizeDiff(resourceSchema map[string]*schema.Schema) schema.CustomizeDiffFunc {
	return func(_ context.Context, d *schema.ResourceDiff, _ any) error {
		if d.Id() == "" || !d.HasChange("size") || !d.NewValueKnown("size") {
			return nil
		}

		o, n := d.GetChange("size")
		oldSize, newSize := o.(int), n.(int)

		if newSize >= oldSize {
			return nil
		}

		// d.HasChange reports nested sets, such as user.databases, as
		// changed even when they are not, so look at the diff itself.
		for _, k := range d.GetChangedKeysPrefix("") {
			if v, ok := resourceSchema[strings.SplitN(k, ".", 2)[0]]; ok && v.ForceNew {
				return nil
			}
		}

		return fmt.Errorf("the size of openstack_db_instance_v1 %s cannot be decreased from %d to %d GB: "+
			"the volume of a database instance can only grow. To replace the instance with a smaller volume, "+
			"taint it first (terraform taint); its data will not be kept", d.Id(), oldSize, newSize)
	}
}

// databaseInstanceV1WaitForResize waits until a resize action on a database
// instance ends, and checks with applied that the resize took effect.
func databaseInstanceV1WaitForResize(ctx context.Context, client *gophercloud.ServiceClient, instanceID string,
	timeout time.Duration, applied func(*instances.Instance) bool,
) (*instances.Instance, error) {
	stateConf := &retry.StateChangeConf{
		Pending:    []string{"RESIZE"},
		Target:     []string{"ACTIVE", "HEALTHY"},
		Refresh:    databaseInstanceV1ResizeStateRefreshFunc(ctx, client, instanceID, applied),
		Timeout:    timeout,
		Delay:      0,
		MinTimeout: 3 * time.Second,
	}

	v, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return nil, err
	}

	return v.(*instances.Instance), nil
}

// databaseInstanceV1ResizeStateRefreshFunc returns a retry.StateRefreshFunc
// that is used to watch a resize of a database instance. Trove reports the
// instance as RESIZE from the moment it accepts the action, and as ACTIVE or
// HEALTHY again when the task ends, also when it failed and was rolled back.
// So ACTIVE or HEALTHY is only reported when applied returns true.
func databaseInstanceV1ResizeStateRefreshFunc(ctx context.Context, client *gophercloud.ServiceClient, instanceID string,
	applied func(*instances.Instance) bool,
) retry.StateRefreshFunc {
	var resizing bool

	return func() (any, string, error) {
		i, err := instances.Get(ctx, client, instanceID).Extract()
		if err != nil {
			return nil, "", err
		}

		switch strings.ToUpper(i.Status) {
		case "RESIZE":
			resizing = true

			return i, "RESIZE", nil
		case "ACTIVE", "HEALTHY":
			if applied(i) {
				return i, strings.ToUpper(i.Status), nil
			}

			if resizing {
				return i, i.Status, fmt.Errorf("the resize was not applied and the instance is %s again, "+
					"it may have been rolled back", i.Status)
			}

			// The resize has not started yet.
			return i, "RESIZE", nil
		case "ERROR":
			if i.Fault != nil && i.Fault.Message != "" {
				return i, i.Status, fmt.Errorf("the instance is in %s state: %s", i.Status, i.Fault.Message)
			}

			return i, i.Status, fmt.Errorf("the instance is in %s state", i.Status)
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
