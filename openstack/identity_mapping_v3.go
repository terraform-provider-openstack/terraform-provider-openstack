package openstack

import (
	"encoding/json"
	"fmt"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/federation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"
)

func identityMappingV3RulesSchema(required bool) *schema.Schema {
	rulesSchema := &schema.Schema{
		Type: schema.TypeString,
	}

	if required {
		rulesSchema.Required = true
		rulesSchema.ValidateFunc = validateIdentityMappingV3Rules
		rulesSchema.DiffSuppressFunc = structure.SuppressJsonDiff
		rulesSchema.StateFunc = func(v any) string {
			normalized, _ := structure.NormalizeJsonString(v)

			return normalized
		}
	} else {
		rulesSchema.Computed = true
	}

	return rulesSchema
}

func validateIdentityMappingV3Rules(v any, k string) ([]string, []error) {
	var rules []federation.MappingRule

	if err := json.Unmarshal([]byte(v.(string)), &rules); err != nil {
		return nil, []error{fmt.Errorf("%q must be a JSON array of federation mapping rules: %w", k, err)}
	}
	if rules == nil {
		return nil, []error{fmt.Errorf("%q must be a JSON array of federation mapping rules, not null", k)}
	}

	return nil, nil
}

func expandIdentityMappingV3Rules(raw string) ([]federation.MappingRule, error) {
	var rules []federation.MappingRule

	if err := json.Unmarshal([]byte(raw), &rules); err != nil {
		return nil, fmt.Errorf("unable to parse federation mapping rules: %w", err)
	}

	return rules, nil
}

func flattenIdentityMappingV3Rules(rules []federation.MappingRule) (string, error) {
	raw, err := json.Marshal(rules)
	if err != nil {
		return "", fmt.Errorf("unable to serialize federation mapping rules: %w", err)
	}

	return string(raw), nil
}
