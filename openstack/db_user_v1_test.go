package openstack

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/databases"
	"github.com/stretchr/testify/assert"
)

func TestUnitExpandDatabaseUserV1Databases(t *testing.T) {
	dbs := []any{"db1", "db2"}

	expected := databases.BatchCreateOpts{
		databases.CreateOpts{
			Name: "db1",
		},
		databases.CreateOpts{
			Name: "db2",
		},
	}

	actual := expandDatabaseUserV1Databases(dbs)
	assert.Equal(t, expected, actual)
}

func TestUnitFlattenDatabaseUserV1Databases(t *testing.T) {
	dbs := []databases.Database{
		{
			Name: "db1",
		},
		{
			Name: "db2",
		},
	}

	expected := []string{"db1", "db2"}

	actual := flattenDatabaseUserV1Databases(dbs)
	assert.Equal(t, expected, actual)
}
