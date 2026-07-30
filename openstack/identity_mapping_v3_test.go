package openstack

import (
	"reflect"
	"testing"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/federation"
)

func TestIdentityMappingV3RulesRoundTrip(t *testing.T) {
	t.Parallel()

	input := `[
		{
			"local": [
				{"user": {"name": "{0}", "type": "ephemeral"}},
				{"group": {"name": "federated", "domain": {"id": "default"}}},
				{"projects": [{"name": "project-{0}", "roles": [{"name": "member"}]}]}
			],
			"remote": [
				{"type": "OIDC-preferred_username"},
				{"type": "OIDC-groups", "any_one_of": ["engineering"], "regex": true}
			]
		}
	]`

	rules, err := expandIdentityMappingV3Rules(input)
	if err != nil {
		t.Fatalf("unexpected expand error: %s", err)
	}

	expectedType := federation.UserTypeEphemeral
	expected := []federation.MappingRule{
		{
			Local: []federation.RuleLocal{
				{
					User: &federation.RuleUser{
						Name: "{0}",
						Type: &expectedType,
					},
				},
				{
					Group: &federation.Group{
						Name: "federated",
						Domain: &federation.Domain{
							ID: "default",
						},
					},
				},
				{
					Projects: []federation.RuleProject{
						{
							Name: "project-{0}",
							Roles: []federation.RuleProjectRole{
								{Name: "member"},
							},
						},
					},
				},
			},
			Remote: []federation.RuleRemote{
				{Type: "OIDC-preferred_username"},
				{
					Type:     "OIDC-groups",
					Regex:    boolPtr(true),
					AnyOneOf: []string{"engineering"},
				},
			},
		},
	}

	if !reflect.DeepEqual(rules, expected) {
		t.Fatalf("unexpected expanded rules:\nwant: %#v\ngot:  %#v", expected, rules)
	}

	flattened, err := flattenIdentityMappingV3Rules(rules)
	if err != nil {
		t.Fatalf("unexpected flatten error: %s", err)
	}

	roundTripped, err := expandIdentityMappingV3Rules(flattened)
	if err != nil {
		t.Fatalf("unexpected second expand error: %s", err)
	}
	if !reflect.DeepEqual(roundTripped, expected) {
		t.Fatalf("rules changed after round trip:\nwant: %#v\ngot:  %#v", expected, roundTripped)
	}
}

func TestValidateIdentityMappingV3Rules(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:  "valid",
			value: `[{"local":[{"user":{"name":"{0}"}}],"remote":[{"type":"username"}]}]`,
		},
		{
			name:    "object",
			value:   `{"local":[],"remote":[]}`,
			wantErr: true,
		},
		{
			name:    "null",
			value:   `null`,
			wantErr: true,
		},
		{
			name:    "malformed",
			value:   `[`,
			wantErr: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			_, errors := validateIdentityMappingV3Rules(test.value, "rules")
			if test.wantErr && len(errors) == 0 {
				t.Fatal("expected validation error")
			}
			if !test.wantErr && len(errors) != 0 {
				t.Fatalf("unexpected validation error: %s", errors[0])
			}
		})
	}
}

func TestIdentityMappingV3Registration(t *testing.T) {
	t.Parallel()

	provider := Provider()

	if _, ok := provider.ResourcesMap["openstack_identity_mapping_v3"]; !ok {
		t.Fatal("openstack_identity_mapping_v3 resource is not registered")
	}
	if _, ok := provider.DataSourcesMap["openstack_identity_mapping_v3"]; !ok {
		t.Fatal("openstack_identity_mapping_v3 data source is not registered")
	}
}

func boolPtr(value bool) *bool {
	return &value
}
