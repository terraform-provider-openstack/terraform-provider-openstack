package openstack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandIdentityUserV3MFARules(t *testing.T) {
	mfaRule_1 := []interface{}{"password", "totp"}
	mfaRule_2 := []interface{}{"password", "custom-auth-method"}

	mfaRules := []interface{}{
		map[string]interface{}{
			"rule": mfaRule_1,
		},
		map[string]interface{}{
			"rule": mfaRule_2,
		},
	}

	expected := []interface{}{mfaRule_1, mfaRule_2}

	actual := expandIdentityUserV3MFARules(mfaRules)
	assert.Equal(t, expected, actual)
}

func TestFlattenIdentityUserV3MFARules(t *testing.T) {
	mfaRule_1 := []interface{}{"password", "totp"}
	mfaRule_2 := []interface{}{"password", "custom-auth-method"}

	mfaRules := []interface{}{mfaRule_1, mfaRule_2}

	expected := []map[string]interface{}{
		{
			"rule": mfaRule_1,
		},
		{
			"rule": mfaRule_2,
		},
	}

	actual := flattenIdentityUserV3MFARules(mfaRules)
	assert.Equal(t, expected, actual)
}
