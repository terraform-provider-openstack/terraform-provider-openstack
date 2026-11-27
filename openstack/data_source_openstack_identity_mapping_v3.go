package openstack

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/federation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceIdentityMappingV3() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceIdentityMappingV3Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"mapping_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"rules": identityMappingV3RulesSchema(false),
		},
	}
}

func dataSourceIdentityMappingV3Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	identityClient, err := config.IdentityV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	mappingID := d.Get("mapping_id").(string)
	mapping, err := federation.GetMapping(ctx, identityClient, mappingID).Extract()
	if err != nil {
		return diag.Errorf("Unable to query openstack_identity_mapping_v3: %s", err)
	}

	log.Printf("[DEBUG] Retrieved openstack_identity_mapping_v3 data source: %#v", mapping)

	rules, err := flattenIdentityMappingV3Rules(mapping.Rules)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(mapping.ID)
	d.Set("region", GetRegion(d, config))
	d.Set("mapping_id", mapping.ID)
	d.Set("rules", rules)

	return nil
}
