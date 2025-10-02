package openstack

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNetworkingPortV2V0() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"all_fixed_ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func convertFixedIPs(list []any) []map[string]any {
	newList := make([]map[string]any, len(list))
	for i, ipRaw := range list {
		newList[i] = map[string]any{
			"ip_address": ipRaw.(string),
			"subnet_id":  "",
		}
	}

	return newList
}

func upgradeNetworkingPortV2StateV0toV1(_ context.Context, rawState map[string]any, _ any) (map[string]any, error) {
	if v, ok := rawState["all_fixed_ips"]; ok {
		if list, ok := v.([]any); ok {
			rawState["all_fixed_ips"] = convertFixedIPs(list)
		}
	}

	return rawState, nil
}
