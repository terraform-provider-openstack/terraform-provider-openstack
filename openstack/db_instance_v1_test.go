package openstack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"

	"github.com/gophercloud/gophercloud/openstack/db/v1/databases"
	"github.com/gophercloud/gophercloud/openstack/db/v1/instances"
	"github.com/gophercloud/gophercloud/openstack/db/v1/users"
)

func TestExpandDatabaseInstanceV1Datastore(t *testing.T) {
	datastore := []interface{}{
		map[string]interface{}{
			"version": "foo",
			"type":    "bar",
		},
	}

	expected := instances.DatastoreOpts{
		Version: "foo",
		Type:    "bar",
	}

	actual := expandDatabaseInstanceV1Datastore(datastore)
	assert.Equal(t, expected, actual)
}

func TestExpandDatabaseInstanceV1Networks(t *testing.T) {
	network := []interface{}{
		map[string]interface{}{
			"uuid":        "foobar",
			"port":        "",
			"fixed_ip_v4": "",
			"fixed_ip_v6": "",
		},
	}

	expected := []instances.NetworkOpts{
		{
			UUID: "foobar",
		},
	}

	actual := expandDatabaseInstanceV1Networks(network)
	assert.Equal(t, expected, actual)
}

func TestExpandDatabaseInstanceV1Databases(t *testing.T) {
	dbs := []interface{}{
		map[string]interface{}{
			"name":    "testdb1",
			"charset": "utf8",
			"collate": "utf8_general_ci",
		},
		map[string]interface{}{
			"name":    "testdb2",
			"charset": "utf8",
			"collate": "utf8_general_ci",
		},
	}

	expected := databases.BatchCreateOpts{
		databases.CreateOpts{
			Name:    "testdb1",
			CharSet: "utf8",
			Collate: "utf8_general_ci",
		},
		databases.CreateOpts{
			Name:    "testdb2",
			CharSet: "utf8",
			Collate: "utf8_general_ci",
		},
	}

	actual := expandDatabaseInstanceV1Databases(dbs)
	assert.Equal(t, expected, actual)
}

func TestExpandDatabaseInstanceV1Users(t *testing.T) {
	userList := []interface{}{
		map[string]interface{}{
			"name":      "testuser",
			"password":  "testpassword",
			"databases": schema.NewSet(schema.HashString, []interface{}{"testdb1"}),
			"host":      "foobar",
		},
	}

	expected := users.BatchCreateOpts{
		users.CreateOpts{
			Name:     "testuser",
			Password: "testpassword",
			Databases: databases.BatchCreateOpts{
				databases.CreateOpts{
					Name: "testdb1",
				},
			},
			Host: "foobar",
		},
	}

	actual := expandDatabaseInstanceV1Users(userList)
	assert.Equal(t, expected, actual)
}
