package openstack

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/databases"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/instances"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/users"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type databaseInstanceV1CreateOpts struct {
	FlavorRef  string
	Size       int
	VolumeType string
	Name       string
	Databases  databases.CreateOptsBuilder
	Users      users.CreateOptsBuilder
	Datastore  *instances.DatastoreOpts
	Networks   []databaseInstanceV1NetworkOpts
}

func (opts databaseInstanceV1CreateOpts) ToInstanceCreateMap() (map[string]any, error) {
	if opts.Size > 300 || opts.Size < 1 {
		err := gophercloud.ErrInvalidInput{}
		err.Argument = "instances.CreateOpts.Size"
		err.Value = opts.Size
		err.Info = "Size (GB) must be between 1-300"

		return nil, err
	}

	if opts.FlavorRef == "" {
		return nil, gophercloud.ErrMissingInput{Argument: "instances.CreateOpts.FlavorRef"}
	}

	instance := map[string]any{
		"flavorRef": opts.FlavorRef,
	}

	if opts.Name != "" {
		instance["name"] = opts.Name
	}

	if opts.Databases != nil {
		dbs, err := opts.Databases.ToDBCreateMap()
		if err != nil {
			return nil, err
		}

		instance["databases"] = dbs["databases"]
	}

	if opts.Users != nil {
		userList, err := opts.Users.ToUserCreateMap()
		if err != nil {
			return nil, err
		}

		instance["users"] = userList["users"]
	}

	if opts.Datastore != nil {
		datastore, err := opts.Datastore.ToMap()
		if err != nil {
			return nil, err
		}

		instance["datastore"] = datastore
	}

	if len(opts.Networks) > 0 {
		networks := make([]map[string]any, len(opts.Networks))
		for i, network := range opts.Networks {
			networks[i] = network.ToMap()
		}

		instance["nics"] = networks
	}

	volume := map[string]any{
		"size": opts.Size,
	}

	if opts.VolumeType != "" {
		volume["type"] = opts.VolumeType
	}

	instance["volume"] = volume

	return map[string]any{"instance": instance}, nil
}

type databaseInstanceV1NetworkOpts struct {
	UUID      string
	Port      string
	V4FixedIP string
	V6FixedIP string
	NetworkID string
	SubnetID  string
	IPAddress string
}

func (opts databaseInstanceV1NetworkOpts) ToMap() map[string]any {
	if opts.usesModernKeys() {
		return opts.toModernMap()
	}

	return opts.toLegacyMap()
}

func (opts databaseInstanceV1NetworkOpts) usesModernKeys() bool {
	return opts.NetworkID != "" || opts.SubnetID != "" || opts.IPAddress != ""
}

func (opts databaseInstanceV1NetworkOpts) usesLegacyKeys() bool {
	return opts.UUID != "" || opts.Port != "" || opts.V4FixedIP != "" || opts.V6FixedIP != ""
}

func (opts databaseInstanceV1NetworkOpts) toModernMap() map[string]any {
	network := make(map[string]any)

	if opts.NetworkID != "" {
		network["network_id"] = opts.NetworkID
	}

	if opts.SubnetID != "" {
		network["subnet_id"] = opts.SubnetID
	}

	if opts.IPAddress != "" {
		network["ip_address"] = opts.IPAddress
	}

	return network
}

func (opts databaseInstanceV1NetworkOpts) toLegacyMap() map[string]any {
	network := make(map[string]any)

	if opts.UUID != "" {
		network["net-id"] = opts.UUID
	}

	if opts.Port != "" {
		network["port-id"] = opts.Port
	}

	if opts.V4FixedIP != "" {
		network["v4-fixed-ip"] = opts.V4FixedIP
	}

	if opts.V6FixedIP != "" {
		network["v6-fixed-ip"] = opts.V6FixedIP
	}

	return network
}

func expandDatabaseInstanceV1Datastore(rawDatastore []any) instances.DatastoreOpts {
	v := rawDatastore[0].(map[string]any)
	datastore := instances.DatastoreOpts{
		Version: v["version"].(string),
		Type:    v["type"].(string),
	}

	return datastore
}

func expandDatabaseInstanceV1Networks(rawNetworks []any) ([]databaseInstanceV1NetworkOpts, error) {
	networks := make([]databaseInstanceV1NetworkOpts, 0, len(rawNetworks))

	for i, v := range rawNetworks {
		rawNetwork := v.(map[string]any)
		network := databaseInstanceV1NetworkOpts{
			UUID:      rawNetwork["uuid"].(string),
			Port:      rawNetwork["port"].(string),
			V4FixedIP: rawNetwork["fixed_ip_v4"].(string),
			V6FixedIP: rawNetwork["fixed_ip_v6"].(string),
			NetworkID: rawNetwork["network_id"].(string),
			SubnetID:  rawNetwork["subnet_id"].(string),
			IPAddress: rawNetwork["ip_address"].(string),
		}

		if network.usesLegacyKeys() && network.usesModernKeys() {
			return nil, fmt.Errorf("network.%d cannot mix legacy fields uuid, port, fixed_ip_v4, or fixed_ip_v6 with modern fields network_id, subnet_id, or ip_address", i)
		}

		networks = append(networks, network)
	}

	return networks, nil
}

func expandDatabaseInstanceV1Databases(rawDatabases []any) databases.BatchCreateOpts {
	dbs := make(databases.BatchCreateOpts, 0, len(rawDatabases))

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
	userList := make(users.BatchCreateOpts, 0, len(rawUsers))

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

func expandInstanceV1UserDatabases(v []any) databases.BatchCreateOpts {
	dbs := make(databases.BatchCreateOpts, 0, len(v))

	for _, db := range v {
		dbs = append(dbs, databases.CreateOpts{
			Name: db.(string),
		})
	}

	return dbs
}
