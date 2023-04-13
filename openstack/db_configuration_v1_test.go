package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gophercloud/gophercloud/openstack/db/v1/configurations"
)

func TestUnitExpandDatabaseConfigurationV1Datastore(t *testing.T) {
	datastore := []interface{}{
		map[string]interface{}{
			"version": "foo",
			"type":    "bar",
		},
	}

	expected := configurations.DatastoreOpts{
		Version: "foo",
		Type:    "bar",
	}

	actual := expandDatabaseConfigurationV1Datastore(datastore)
	assert.Equal(t, expected, actual)
}

func TestUnitExpandDatabaseConfigurationV1Values(t *testing.T) {
	values := []interface{}{
		map[string]interface{}{
			"name":  "collation_server",
			"value": "latin1_swedish_ci",
		},
		map[string]interface{}{
			"name":  "collation_database",
			"value": "latin1_swedish_ci",
		},
		map[string]interface{}{
			"name":  "max_connections",
			"value": "200",
		},
		map[string]interface{}{
			"name":        "collation_connection",
			"value":       "47",
			"string_type": false,
		},
		map[string]interface{}{
			"name":        "connect_timeout",
			"value":       "3",
			"string_type": true,
		},
		map[string]interface{}{
			"name":  "autocommit",
			"value": "true",
		},
		map[string]interface{}{
			"name":        "sync_binlog",
			"value":       "true",
			"string_type": true,
		},
	}

	expected := map[string]interface{}{
		"collation_server":     "latin1_swedish_ci",
		"collation_database":   "latin1_swedish_ci",
		"max_connections":      200,
		"collation_connection": 47,
		"connect_timeout":      "3",
		"autocommit":           true,
		"sync_binlog":          "true",
	}

	actual := expandDatabaseConfigurationV1Values(values)
	assert.Equal(t, expected, actual)
}
