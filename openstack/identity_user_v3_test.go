package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnitExpandIdentityUserV3MFARules(t *testing.T) {
	mfaRule1 := []interface{}{"password", "totp"}
	mfaRule2 := []interface{}{"password", "custom-auth-method"}

	mfaRules := []interface{}{
		map[string]interface{}{
			"rule": mfaRule1,
		},
		map[string]interface{}{
			"rule": mfaRule2,
		},
	}

	expected := []interface{}{mfaRule1, mfaRule2}

	actual := expandIdentityUserV3MFARules(mfaRules)
	assert.Equal(t, expected, actual)
}

func TestUnitFlattenIdentityUserV3MFARules(t *testing.T) {
	mfaRule1 := []interface{}{"password", "totp"}
	mfaRule2 := []interface{}{"password", "custom-auth-method"}

	mfaRules := []interface{}{mfaRule1, mfaRule2}

	expected := []map[string]interface{}{
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
