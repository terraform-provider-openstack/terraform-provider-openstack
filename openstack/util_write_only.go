package openstack

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// getWriteOnlyString returns the value of a write-only string attribute.
// Write-only attributes are never stored in the state nor in the ResourceData,
// so d.Get/d.GetOk always return an empty value for them: the only way to read
// them is through the raw configuration.
func getWriteOnlyString(d *schema.ResourceData, key string) (string, bool) {
	cfg := d.GetRawConfig()
	if cfg.IsNull() || !cfg.IsKnown() || !cfg.Type().IsObjectType() || !cfg.Type().HasAttribute(key) {
		return "", false
	}

	v := cfg.GetAttr(key)
	if v.IsNull() || !v.IsKnown() || !v.Type().IsPrimitiveType() || v.Type().FriendlyName() != "string" {
		return "", false
	}

	s := v.AsString()

	return s, s != ""
}
