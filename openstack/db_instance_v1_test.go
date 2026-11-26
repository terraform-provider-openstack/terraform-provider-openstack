package openstack

import (
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/databases"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/instances"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/users"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitExpandDatabaseInstanceV1Datastore(t *testing.T) {
	datastore := []any{
		map[string]any{
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

func TestUnitExpandDatabaseInstanceV1Networks(t *testing.T) {
	network := []any{
		map[string]any{
			"uuid":        "legacy-network-id",
			"port":        "legacy-port-id",
			"fixed_ip_v4": "192.0.2.10",
			"fixed_ip_v6": "2001:db8::10",
			"network_id":  "",
			"subnet_id":   "",
			"ip_address":  "",
		},
	}

	expected := []databaseInstanceV1NetworkOpts{
		{
			UUID:      "legacy-network-id",
			Port:      "legacy-port-id",
			V4FixedIP: "192.0.2.10",
			V6FixedIP: "2001:db8::10",
		},
	}

	actual, err := expandDatabaseInstanceV1Networks(network)
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestUnitExpandDatabaseInstanceV1NetworksModern(t *testing.T) {
	network := []any{
		map[string]any{
			"uuid":        "",
			"port":        "",
			"fixed_ip_v4": "",
			"fixed_ip_v6": "",
			"network_id":  "modern-network-id",
			"subnet_id":   "modern-subnet-id",
			"ip_address":  "192.0.2.20",
		},
	}

	expected := []databaseInstanceV1NetworkOpts{
		{
			NetworkID: "modern-network-id",
			SubnetID:  "modern-subnet-id",
			IPAddress: "192.0.2.20",
		},
	}

	actual, err := expandDatabaseInstanceV1Networks(network)
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestUnitExpandDatabaseInstanceV1NetworksMixed(t *testing.T) {
	network := []any{
		map[string]any{
			"uuid":        "legacy-network-id",
			"port":        "",
			"fixed_ip_v4": "",
			"fixed_ip_v6": "",
			"network_id":  "modern-network-id",
			"subnet_id":   "",
			"ip_address":  "",
		},
	}

	actual, err := expandDatabaseInstanceV1Networks(network)
	assert.Nil(t, actual)
	assert.EqualError(t, err, "network.0 cannot mix legacy fields uuid, port, fixed_ip_v4, or fixed_ip_v6 with modern fields network_id, subnet_id, or ip_address")
}

func TestUnitDatabaseInstanceV1NetworkOptsToMapLegacy(t *testing.T) {
	network := databaseInstanceV1NetworkOpts{
		UUID:      "legacy-network-id",
		Port:      "legacy-port-id",
		V4FixedIP: "192.0.2.10",
		V6FixedIP: "2001:db8::10",
	}

	expected := map[string]any{
		"net-id":      "legacy-network-id",
		"port-id":     "legacy-port-id",
		"v4-fixed-ip": "192.0.2.10",
		"v6-fixed-ip": "2001:db8::10",
	}

	assert.Equal(t, expected, network.ToMap())
}

func TestUnitDatabaseInstanceV1NetworkOptsToMapModern(t *testing.T) {
	network := databaseInstanceV1NetworkOpts{
		NetworkID: "modern-network-id",
		SubnetID:  "modern-subnet-id",
		IPAddress: "192.0.2.20",
	}

	expected := map[string]any{
		"network_id": "modern-network-id",
		"subnet_id":  "modern-subnet-id",
		"ip_address": "192.0.2.20",
	}

	assert.Equal(t, expected, network.ToMap())
}

func TestUnitDatabaseInstanceV1CreateOptsToInstanceCreateMapModernNetworks(t *testing.T) {
	createOpts := databaseInstanceV1CreateOpts{
		FlavorRef: "flavor-id",
		Name:      "db-instance",
		Size:      10,
		Datastore: &instances.DatastoreOpts{
			Version: "mysql-8.0",
			Type:    "mysql",
		},
		Networks: []databaseInstanceV1NetworkOpts{
			{
				NetworkID: "modern-network-id",
				SubnetID:  "modern-subnet-id",
				IPAddress: "192.0.2.20",
			},
		},
	}

	expected := map[string]any{
		"instance": map[string]any{
			"flavorRef": "flavor-id",
			"name":      "db-instance",
			"datastore": map[string]any{
				"version": "mysql-8.0",
				"type":    "mysql",
			},
			"nics": []map[string]any{
				{
					"network_id": "modern-network-id",
					"subnet_id":  "modern-subnet-id",
					"ip_address": "192.0.2.20",
				},
			},
			"volume": map[string]any{
				"size": 10,
			},
		},
	}

	actual, err := createOpts.ToInstanceCreateMap()
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestUnitDatabaseInstanceV1CreateOptsToInstanceCreateMapWithVolumeType(t *testing.T) {
	createOpts := databaseInstanceV1CreateOpts{
		FlavorRef:  "flavor-id",
		Name:       "db-instance",
		Size:       10,
		VolumeType: "ssd",
	}

	expected := map[string]any{
		"instance": map[string]any{
			"flavorRef": "flavor-id",
			"name":      "db-instance",
			"volume": map[string]any{
				"size": 10,
				"type": "ssd",
			},
		},
	}

	actual, err := createOpts.ToInstanceCreateMap()
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestUnitDatabaseInstanceV1CreateOptsToInstanceCreateMapWithoutVolumeType(t *testing.T) {
	createOpts := databaseInstanceV1CreateOpts{
		FlavorRef: "flavor-id",
		Name:      "db-instance",
		Size:      10,
	}

	expected := map[string]any{
		"instance": map[string]any{
			"flavorRef": "flavor-id",
			"name":      "db-instance",
			"volume": map[string]any{
				"size": 10,
			},
		},
	}

	actual, err := createOpts.ToInstanceCreateMap()
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestUnitExpandDatabaseInstanceV1Databases(t *testing.T) {
	dbs := []any{
		map[string]any{
			"name":    "testdb1",
			"charset": "utf8",
			"collate": "utf8_general_ci",
		},
		map[string]any{
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

func TestUnitExpandDatabaseInstanceV1Users(t *testing.T) {
	userList := []any{
		map[string]any{
			"name":      "testuser",
			"password":  "testpassword",
			"databases": schema.NewSet(schema.HashString, []any{"testdb1"}),
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
