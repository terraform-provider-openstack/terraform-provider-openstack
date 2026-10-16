package openstack

import (
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// getWriteOnlyString returns the value of a write-only string attribute.
// Write-only attributes are never stored in the state or in the ResourceData,
// so d.Get/d.GetOk always return an empty value for them: the only way to read
// them is through the raw configuration.
func getWriteOnlyString(d *schema.ResourceData, key string) (string, bool) {
	v, diags := d.GetRawConfigAt(cty.GetAttrPath(key))
	if diags.HasError() || v.IsNull() || !v.IsKnown() || v.Type() != cty.String {
		return "", false
	}

	s := v.AsString()

	return s, s != ""
}
