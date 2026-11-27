package openstack

import (
	"fmt"
	"net/http"
	"strconv"
	"sync/atomic"
	"testing"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/databases"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/instances"
	"github.com/gophercloud/gophercloud/v2/openstack/db/v1/users"
	th "github.com/gophercloud/gophercloud/v2/testhelper"
	thclient "github.com/gophercloud/gophercloud/v2/testhelper/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnitDatabaseInstanceV1SizePlan(t *testing.T) {
	tests := []struct {
		name        string
		size        int
		unknown     bool
		create      bool
		change      map[string]any
		wantError   bool
		wantReplace bool
	}{
		{name: "increase", size: 11},
		{name: "unchanged", size: 10},
		{name: "decrease", size: 9, wantError: true},
		{name: "unknown", unknown: true},
		{name: "create", size: 9, create: true},
		{name: "replacement", size: 9, change: map[string]any{"name": "replacement"}, wantReplace: true},
		{name: "nested replacement", size: 9, change: map[string]any{
			"datastore": []any{map[string]any{"type": "mysql", "version": "8.4"}},
		}, wantReplace: true},
		{name: "configuration update", size: 9, change: map[string]any{"configuration_id": "new-config"}, wantError: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resource := resourceDatabaseInstanceV1()
			config := map[string]any{
				"name":      "basic",
				"flavor_id": "flavor-1",
				"size":      10,
				"datastore": []any{map[string]any{"type": "mysql", "version": "8.0"}},
			}
			data := schema.TestResourceDataRaw(t, resource.Schema, config)
			data.SetId("instance-1")
			state := data.State()
			if tt.create {
				state = nil
			}

			config["size"] = tt.size
			for key, value := range tt.change {
				config[key] = value
			}

			if tt.unknown {
				// Terraform SDK's unknown-value sentinel for legacy ResourceConfig.
				config["size"] = "74D93920-ED26-11E3-AC10-0800200C9A66"
			}

			planConfig := terraform.NewResourceConfigRaw(config)

			diff, err := resource.Diff(t.Context(), state, planConfig, nil)
			if tt.wantError {
				require.ErrorContains(t, err, "size cannot decrease from 10 to 9 GB")

				return
			}

			require.NoError(t, err)

			if !tt.create {
				assert.Equal(t, tt.wantReplace, diff != nil && diff.RequiresNew())
			}
		})
	}
}

func TestUnitDatabaseInstanceV1SizeSupportsUpdate(t *testing.T) {
	resource := resourceDatabaseInstanceV1()

	assert.False(t, resource.Schema["size"].ForceNew)
	assert.NotNil(t, resource.Timeouts.Update)
}

func TestUnitDatabaseInstanceV1SizeValidation(t *testing.T) {
	validate := resourceDatabaseInstanceV1().Schema["size"].ValidateFunc

	for _, size := range []int{-1, 0, 1, 10} {
		t.Run(strconv.Itoa(size), func(t *testing.T) {
			_, errs := validate(size, "size")
			if size < 1 {
				assert.NotEmpty(t, errs)
			} else {
				assert.Empty(t, errs)
			}
		})
	}
}

func TestUnitDatabaseInstanceV1UpdateState(t *testing.T) {
	tests := []struct {
		name              string
		resizeCode        int
		instanceStatus    string
		configurationCode int
		failAttach        bool
		wantError         bool
	}{
		{name: "resize rejected", resizeCode: http.StatusBadRequest, wantError: true},
		{name: "resize failed", instanceStatus: "ERROR", wantError: true},
		{name: "detach failed after resize", configurationCode: http.StatusBadRequest, wantError: true},
		{name: "attach failed after resize", failAttach: true, wantError: true},
		{name: "success"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeServer := th.SetupHTTP()
			t.Cleanup(fakeServer.Teardown)

			var configurationCalls atomic.Int32

			fakeServer.Mux.HandleFunc("/instances/instance-1/action", func(w http.ResponseWriter, r *http.Request) {
				th.TestMethod(t, r, http.MethodPost)
				th.TestJSONRequest(t, r, `{"resize":{"volume":{"size":11}}}`)

				if tt.resizeCode != 0 {
					w.WriteHeader(tt.resizeCode)

					return
				}

				w.WriteHeader(http.StatusAccepted)
			})
			fakeServer.Mux.HandleFunc("/instances/instance-1", func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodPut {
					call := configurationCalls.Add(1)
					if call%2 == 1 {
						th.TestJSONRequest(t, r, `{"instance":{}}`)
					} else {
						th.TestJSONRequest(t, r, `{"instance":{"configuration":"new-config"}}`)
					}

					if tt.configurationCode != 0 {
						w.WriteHeader(tt.configurationCode)
					} else if tt.failAttach && call%2 == 0 {
						w.WriteHeader(http.StatusBadRequest)
					} else {
						w.WriteHeader(http.StatusAccepted)
					}

					return
				}

				th.TestMethod(t, r, http.MethodGet)

				status := tt.instanceStatus
				if status == "" {
					status = "ACTIVE"
				}

				w.Header().Set("Content-Type", "application/json")
				_, err := fmt.Fprintf(w, `{"instance":{"id":"instance-1","name":"basic","status":%q,"volume":{"size":11}}}`, status)
				assert.NoError(t, err)
			})

			config := &Config{}
			config.OsClient = thclient.ServiceClient(fakeServer).ProviderClient
			config.OsClient.EndpointLocator = func(_ gophercloud.EndpointOpts) (string, error) {
				return fakeServer.Endpoint(), nil
			}

			resource := resourceDatabaseInstanceV1()
			data := schema.TestResourceDataRaw(t, resource.Schema, map[string]any{
				"size":             10,
				"configuration_id": "old-config",
			})
			data.SetId("instance-1")

			diff := &terraform.InstanceDiff{Attributes: map[string]*terraform.ResourceAttrDiff{
				"size":             {Old: "10", New: "11"},
				"configuration_id": {Old: "old-config", New: "new-config"},
			}}

			state, diags := resource.Apply(t.Context(), data.State(), diff, config)
			require.NotNil(t, state)
			assert.Equal(t, tt.wantError, diags.HasError(), "%v", diags)

			if !tt.wantError {
				assert.Equal(t, "11", state.Attributes["size"])
				assert.Equal(t, "new-config", state.Attributes["configuration_id"])
				assert.EqualValues(t, 2, configurationCalls.Load())

				return
			}

			assert.Equal(t, "10", state.Attributes["size"])
			assert.Equal(t, "old-config", state.Attributes["configuration_id"])

			if tt.resizeCode != 0 || tt.instanceStatus == "ERROR" {
				assert.Zero(t, configurationCalls.Load())

				return
			}

			// After a successful resize and a failed configuration update, refresh
			// recovers the actual size while retaining the configuration change.
			refreshed := resource.Data(state)
			diags = resourceDatabaseInstanceV1Read(t.Context(), refreshed, config)
			require.False(t, diags.HasError(), "%v", diags)
			assert.Equal(t, 11, refreshed.Get("size"))
			assert.Equal(t, "old-config", refreshed.Get("configuration_id"))
		})
	}
}

func TestUnitDatabaseInstanceV1VolumeResizeStateRefresh(t *testing.T) {
	tests := []struct {
		name       string
		status     string
		size       int
		code       int
		wantState  string
		wantError  bool
		wantAbsent bool
	}{
		{name: "active with old size", status: "ACTIVE", size: 10, wantState: "RESIZE"},
		{name: "healthy with old size", status: "HEALTHY", size: 10, wantState: "RESIZE"},
		{name: "resizing with new size", status: "RESIZE", size: 11, wantState: "RESIZE"},
		{name: "active with new size", status: "ACTIVE", size: 11, wantState: "ACTIVE"},
		{name: "healthy with new size", status: "HEALTHY", size: 11, wantState: "HEALTHY"},
		{name: "error", status: "ERROR", size: 10, wantState: "ERROR", wantError: true},
		{name: "lowercase error", status: "error", size: 10, wantState: "error", wantError: true},
		{name: "deleted", code: http.StatusNotFound, wantState: "DELETED", wantAbsent: true},
		{name: "API error", code: http.StatusInternalServerError, wantError: true, wantAbsent: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeServer := th.SetupHTTP()
			t.Cleanup(fakeServer.Teardown)
			fakeServer.Mux.HandleFunc("/instances/instance-1", func(w http.ResponseWriter, r *http.Request) {
				th.TestMethod(t, r, http.MethodGet)

				if tt.code != 0 {
					w.WriteHeader(tt.code)

					return
				}

				w.Header().Set("Content-Type", "application/json")
				_, err := fmt.Fprintf(w, `{"instance":{"id":"instance-1","status":%q,"volume":{"size":%d}}}`, tt.status, tt.size)
				assert.NoError(t, err)
			})

			refresh := databaseInstanceV1VolumeResizeStateRefreshFunc(t.Context(), thclient.ServiceClient(fakeServer), "instance-1", 11)
			result, state, err := refresh()
			if tt.wantError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.wantState, state)

			if tt.wantAbsent {
				assert.Nil(t, result)
			} else {
				instance, ok := result.(*instances.Instance)
				require.True(t, ok)
				assert.Equal(t, tt.size, instance.Volume.Size)
			}
		})
	}
}

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
