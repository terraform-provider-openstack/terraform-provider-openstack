package openstack

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestLBV2FlavorCreateOpts(t *testing.T) {
	opts := flavorsCreateOpts{
		Name:            "name",
		Description:     "description",
		FlavorProfileID: "id",
		Enabled:         new(bool),
	}

	b, err := opts.ToFlavorCreateMap()
	if err != nil {
		t.Fatalf("Error creating flavor create map: %v", err)
	}

	expected := map[string]any{
		"flavor": map[string]any{
			"name":              "name",
			"description":       "description",
			"flavor_profile_id": "id",
			"enabled":           false,
		},
	}
	if diff := cmp.Diff(expected, b); diff != "" {
		t.Fatalf("Values are not the same:\n%s", diff)
	}
}
