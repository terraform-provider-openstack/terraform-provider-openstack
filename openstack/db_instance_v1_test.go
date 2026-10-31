package openstack

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
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

func TestUnitDatabaseInstanceV1CustomizeDiffSize(t *testing.T) {
	testCases := []struct {
		name          string
		create        bool
		change        map[string]any
		expectError   bool
		expectReplace bool
	}{
		{name: "grow", change: map[string]any{"size": 20}},
		{name: "resize flavor", change: map[string]any{"flavor_id": "flavor-2"}},
		{name: "resize flavor and grow", change: map[string]any{"flavor_id": "flavor-2", "size": 20}},
		{name: "shrink", change: map[string]any{"size": 5}, expectError: true},
		{name: "shrink and resize flavor", change: map[string]any{"flavor_id": "flavor-2", "size": 5}, expectError: true},
		{name: "shrink and replace", change: map[string]any{"name": "other", "size": 5}, expectReplace: true},
		{
			name: "shrink and replace on a nested change",
			change: map[string]any{
				"size":      5,
				"datastore": []any{map[string]any{"type": "mysql", "version": "8.4"}},
			},
			expectReplace: true,
		},
		{
			name: "shrink and replace on a user change",
			change: map[string]any{
				"size": 5,
				"user": []any{map[string]any{
					"name":      "testuser",
					"password":  "testpassword",
					"host":      "%",
					"databases": []any{"testdb", "otherdb"},
				}},
			},
			expectReplace: true,
		},
		{name: "create", create: true, change: map[string]any{"size": 5}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := resourceDatabaseInstanceV1()

			config := map[string]any{
				"name":      "basic",
				"flavor_id": "flavor-1",
				"size":      10,
				"datastore": []any{map[string]any{"type": "mysql", "version": "8.0"}},
				"database":  []any{map[string]any{"name": "testdb"}},
				"user": []any{map[string]any{
					"name":      "testuser",
					"password":  "testpassword",
					"host":      "%",
					"databases": []any{"testdb"},
				}},
			}

			var state *terraform.InstanceState

			if !tc.create {
				d := schema.TestResourceDataRaw(t, r.Schema, config)
				d.SetId("instance-1")
				state = d.State()
			}

			for k, v := range tc.change {
				config[k] = v
			}

			diff, err := r.Diff(t.Context(), state, terraform.NewResourceConfigRaw(config), nil)
			if tc.expectError {
				require.ErrorContains(t, err, "cannot be decreased from 10 to 5 GB")

				return
			}

			require.NoError(t, err)
			require.NotNil(t, diff)

			if !tc.create {
				assert.Equal(t, tc.expectReplace, diff.RequiresNew())
			}
		})
	}
}

func TestUnitDatabaseInstanceV1ResizeStateRefreshFunc(t *testing.T) {
	type response struct {
		status string
		flavor string
		fault  string
	}

	testCases := []struct {
		name          string
		responses     []response
		expectStates  []string
		expectErrorOn int
		expectError   string
	}{
		{
			name:         "resized",
			responses:    []response{{"RESIZE", "flavor-1", ""}, {"RESIZE", "flavor-2", ""}, {"ACTIVE", "flavor-2", ""}},
			expectStates: []string{"RESIZE", "RESIZE", "ACTIVE"},
		},
		{
			name:         "resized and healthy",
			responses:    []response{{"RESIZE", "flavor-1", ""}, {"HEALTHY", "flavor-2", ""}},
			expectStates: []string{"RESIZE", "HEALTHY"},
		},
		{
			name:         "not started yet",
			responses:    []response{{"ACTIVE", "flavor-1", ""}, {"RESIZE", "flavor-1", ""}, {"ACTIVE", "flavor-2", ""}},
			expectStates: []string{"RESIZE", "RESIZE", "ACTIVE"},
		},
		{
			name:          "rolled back",
			responses:     []response{{"RESIZE", "flavor-1", ""}, {"ACTIVE", "flavor-1", ""}},
			expectStates:  []string{"RESIZE", "ACTIVE"},
			expectErrorOn: 2,
			expectError:   "the resize was not applied",
		},
		{
			name:          "error",
			responses:     []response{{"RESIZE", "flavor-1", ""}, {"ERROR", "flavor-1", "No valid host was found"}},
			expectStates:  []string{"RESIZE", "ERROR"},
			expectErrorOn: 2,
			expectError:   "the instance is in ERROR state: No valid host was found",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakeServer := th.SetupHTTP()
			defer fakeServer.Teardown()

			var calls int

			fakeServer.Mux.HandleFunc("/instances/instance-1", func(w http.ResponseWriter, r *http.Request) {
				th.TestMethod(t, r, http.MethodGet)

				resp := tc.responses[calls]
				calls++

				fault := ""
				if resp.fault != "" {
					fault = fmt.Sprintf(`, "fault": {"message": %q}`, resp.fault)
				}

				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, `{"instance": {"id": "instance-1", "status": %q, "flavor": {"id": %q}, "volume": {"size": 10}%s}}`,
					resp.status, resp.flavor, fault)
			})

			refresh := databaseInstanceV1ResizeStateRefreshFunc(t.Context(), thclient.ServiceClient(fakeServer), "instance-1",
				func(i *instances.Instance) bool {
					return i.Flavor.ID == "flavor-2"
				})

			for n, expectState := range tc.expectStates {
				_, state, err := refresh()
				assert.Equal(t, expectState, state)

				if n+1 == tc.expectErrorOn {
					require.ErrorContains(t, err, tc.expectError)
				} else {
					require.NoError(t, err)
				}
			}
		})
	}
}

func TestUnitDatabaseInstanceV1UpdateResize(t *testing.T) {
	testCases := []struct {
		name           string
		change         map[string]string
		growWithFlavor bool
		refuseFlavor   bool
		expectActions  []string
		expectError    string
	}{
		{
			name:          "flavor",
			change:        map[string]string{"flavor_id": "flavor-2"},
			expectActions: []string{`{"resize":{"flavorRef":"flavor-2"}}`},
		},
		{
			name:          "volume",
			change:        map[string]string{"size": "20"},
			expectActions: []string{`{"resize":{"volume":{"size":20}}}`},
		},
		{
			name:   "flavor before volume",
			change: map[string]string{"flavor_id": "flavor-2", "size": "20"},
			expectActions: []string{
				`{"resize":{"flavorRef":"flavor-2"}}`,
				`{"resize":{"volume":{"size":20}}}`,
			},
		},
		{
			name:           "volume grown with the flavor",
			change:         map[string]string{"flavor_id": "flavor-2", "size": "20"},
			growWithFlavor: true,
			expectActions:  []string{`{"resize":{"flavorRef":"flavor-2"}}`},
		},
		{
			name:          "flavor resize refused",
			change:        map[string]string{"flavor_id": "flavor-2", "size": "20"},
			refuseFlavor:  true,
			expectActions: []string{`{"resize":{"flavorRef":"flavor-2"}}`},
			expectError:   "Error resizing openstack_db_instance_v1 instance-1 flavor",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fakeServer := th.SetupHTTP()
			defer fakeServer.Teardown()

			flavor, size := "flavor-1", 10

			var actions []string

			// The fake Trove applies a resize at once and reports the
			// instance as ACTIVE.
			fakeServer.Mux.HandleFunc("/instances/instance-1/action", func(w http.ResponseWriter, r *http.Request) {
				th.TestMethod(t, r, http.MethodPost)

				body, err := io.ReadAll(r.Body)
				assert.NoError(t, err)

				actions = append(actions, string(body))

				var action struct {
					Resize struct {
						FlavorRef string `json:"flavorRef"`
						Volume    struct {
							Size int `json:"size"`
						} `json:"volume"`
					} `json:"resize"`
				}

				assert.NoError(t, json.Unmarshal(body, &action))

				if action.Resize.FlavorRef != "" && tc.refuseFlavor {
					w.WriteHeader(http.StatusBadRequest)

					return
				}

				if action.Resize.FlavorRef != "" {
					flavor = action.Resize.FlavorRef

					if tc.growWithFlavor {
						size = 20
					}
				}

				if action.Resize.Volume.Size != 0 {
					size = action.Resize.Volume.Size
				}

				w.WriteHeader(http.StatusAccepted)
			})
			fakeServer.Mux.HandleFunc("/instances/instance-1", func(w http.ResponseWriter, r *http.Request) {
				th.TestMethod(t, r, http.MethodGet)

				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, `{"instance": {"id": "instance-1", "name": "basic", "status": "ACTIVE", "flavor": {"id": %q}, "volume": {"size": %d}}}`,
					flavor, size)
			})

			config := &Config{}
			config.OsClient = thclient.ServiceClient(fakeServer).ProviderClient
			config.OsClient.EndpointLocator = func(gophercloud.EndpointOpts) (string, error) {
				return fakeServer.Endpoint(), nil
			}

			r := resourceDatabaseInstanceV1()
			d := schema.TestResourceDataRaw(t, r.Schema, map[string]any{
				"name":      "basic",
				"flavor_id": "flavor-1",
				"size":      10,
				"datastore": []any{map[string]any{"type": "mysql", "version": "8.0"}},
			})
			d.SetId("instance-1")

			diff := &terraform.InstanceDiff{Attributes: map[string]*terraform.ResourceAttrDiff{}}
			for k, v := range tc.change {
				diff.Attributes[k] = &terraform.ResourceAttrDiff{Old: d.State().Attributes[k], New: v}
			}

			state, diags := r.Apply(t.Context(), d.State(), diff, config)
			assert.Equal(t, tc.expectActions, actions)

			if tc.expectError != "" {
				require.True(t, diags.HasError())
				assert.Contains(t, diags[0].Summary, tc.expectError)
			} else {
				require.False(t, diags.HasError(), "%v", diags)
			}

			// The state matches what Trove reports, also after a failure.
			assert.Equal(t, flavor, state.Attributes["flavor_id"])
			assert.Equal(t, strconv.Itoa(size), state.Attributes["size"])
		})
	}
}
