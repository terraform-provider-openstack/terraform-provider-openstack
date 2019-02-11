package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gophercloud/gophercloud/openstack/db/v1/configurations"
)

func TestExpandDatabaseConfigurationV1Datastore(t *testing.T) {
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

func TestExpandDatabaseConfigurationV1Values(t *testing.T) {
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
	}

	expected := map[string]interface{}{
		"collation_server":   "latin1_swedish_ci",
		"collation_database": "latin1_swedish_ci",
		"max_connections":    200,
	}

	actual := expandDatabaseConfigurationV1Values(values)
	assert.Equal(t, expected, actual)
}
