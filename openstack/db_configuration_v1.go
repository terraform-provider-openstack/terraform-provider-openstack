package openstack

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/configurations"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

func expandDatabaseConfigurationV1Datastore(rawDatastore []any) configurations.DatastoreOpts {
	v := rawDatastore[0].(map[string]any)
	datastore := configurations.DatastoreOpts{
		Version: v["version"].(string),
		Type:    v["type"].(string),
	}

	return datastore
}

func expandDatabaseConfigurationV1Values(rawValues []any) map[string]any {
	values := make(map[string]any)

	for _, rawValue := range rawValues {
		v := rawValue.(map[string]any)
		name := v["name"].(string)
		value := v["value"]

		if isStringType, ok := v["string_type"].(bool); !ok || !isStringType {
			// check if value can be converted into int
			if valueInt, err := strconv.Atoi(value.(string)); err == nil {
				value = valueInt
				// check if value can be converted into bool
			} else if valueBool, err := strconv.ParseBool(value.(string)); err == nil {
				value = valueBool
			}
		}

		values[name] = value
	}

	return values
}

// databaseConfigurationV1StateRefreshFunc returns a retry.StateRefreshFunc that is used to watch
// an cloud database instance.
func databaseConfigurationV1StateRefreshFunc(ctx context.Context, client *gophercloud.ServiceClient, cgroupID string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		i, err := configurations.Get(ctx, client, cgroupID).Extract()
		if err != nil {
			if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
				return i, "DELETED", nil
			}

			return nil, "", err
		}

		return i, "ACTIVE", nil
	}
}
