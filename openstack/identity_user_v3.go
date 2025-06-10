package openstack

import (
	"fmt"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/users"
)

func getUserOptions() [4]users.Option {
	return [4]users.Option{
		users.IgnoreChangePasswordUponFirstUse,
		users.IgnorePasswordExpiry,
		users.IgnoreLockoutFailureAttempts,
		users.MultiFactorAuthEnabled,
	}
}

func expandIdentityUserV3MFARules(rules []any) []any {
	mfaRules := make([]any, 0, len(rules))

	for _, rule := range rules {
		ruleMap := rule.(map[string]any)
		ruleList := ruleMap["rule"].([]any)
		mfaRules = append(mfaRules, ruleList)
	}

	return mfaRules
}

func flattenIdentityUserV3MFARules(v []any) []map[string]any {
	mfaRules := []map[string]any{}

	for _, rawRule := range v {
		mfaRule := map[string]any{
			"rule": rawRule,
		}
		mfaRules = append(mfaRules, mfaRule)
	}

	return mfaRules
}

// Ensure that password_expires_at query matches format explained in
// https://developer.openstack.org/api-ref/identity/v3/#list-users
func validatePasswordExpiresAtQuery(v any, k string) (ws []string, errors []error) {
	value := v.(string)
	values := strings.SplitN(value, ":", 2)

	if len(values) != 2 {
		err := fmt.Errorf("%s '%s' does not match expected format: {operator}:{timestamp}", k, value)
		errors = append(errors, err)
	}

	operator, timestamp := values[0], values[1]

	validOperators := map[string]bool{
		"lt":  true,
		"lte": true,
		"gt":  true,
		"gte": true,
		"eq":  true,
		"neq": true,
	}
	if !validOperators[operator] {
		err := fmt.Errorf("'%s' is not a valid operator for %s. Choose one of 'lt', 'lte', 'gt', 'gte', 'eq', 'neq'", operator, k)
		errors = append(errors, err)
	}

	_, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		err = fmt.Errorf("'%s' is not a valid timestamp for %s. It should be in the form 'YYYY-MM-DDTHH:mm:ssZ'", timestamp, k)
		errors = append(errors, err)
	}

	return ws, errors
}
