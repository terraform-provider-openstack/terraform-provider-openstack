package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitExpandIdentityUserV3MFARules(t *testing.T) {
	mfaRule1 := []any{"password", "totp"}
	mfaRule2 := []any{"password", "custom-auth-method"}

	mfaRules := []any{
		map[string]any{
			"rule": mfaRule1,
		},
		map[string]any{
			"rule": mfaRule2,
		},
	}

	expected := []any{mfaRule1, mfaRule2}

	actual := expandIdentityUserV3MFARules(mfaRules)
	assert.Equal(t, expected, actual)
}

func TestUnitFlattenIdentityUserV3MFARules(t *testing.T) {
	mfaRule1 := []any{"password", "totp"}
	mfaRule2 := []any{"password", "custom-auth-method"}

	mfaRules := []any{mfaRule1, mfaRule2}

	expected := []map[string]any{
		{
			"rule": mfaRule1,
		},
		{
			"rule": mfaRule2,
		},
	}

	actual := flattenIdentityUserV3MFARules(mfaRules)
	assert.Equal(t, expected, actual)
}
