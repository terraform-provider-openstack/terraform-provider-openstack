package openstack

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/federation"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIdentityMappingV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityMappingV3Create,
		ReadContext:   resourceIdentityMappingV3Read,
		UpdateContext: resourceIdentityMappingV3Update,
		DeleteContext: resourceIdentityMappingV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"mapping_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"rules": identityMappingV3RulesSchema(true),
		},
	}
}

func resourceIdentityMappingV3Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	identityClient, err := config.IdentityV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	rules, err := expandIdentityMappingV3Rules(d.Get("rules").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	createOpts := federation.CreateMappingOpts{
		Rules: rules,
	}
	mappingID := d.Get("mapping_id").(string)

	log.Printf("[DEBUG] openstack_identity_mapping_v3 create options: %#v", createOpts)

	mapping, err := federation.CreateMapping(ctx, identityClient, mappingID, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_identity_mapping_v3: %s", err)
	}

	d.SetId(mapping.ID)

	return resourceIdentityMappingV3Read(ctx, d, meta)
}

func resourceIdentityMappingV3Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	identityClient, err := config.IdentityV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	mapping, err := federation.GetMapping(ctx, identityClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_identity_mapping_v3"))
	}

	log.Printf("[DEBUG] Retrieved openstack_identity_mapping_v3: %#v", mapping)

	rules, err := flattenIdentityMappingV3Rules(mapping.Rules)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("region", GetRegion(d, config))
	d.Set("mapping_id", mapping.ID)
	d.Set("rules", rules)

	return nil
}

func resourceIdentityMappingV3Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	identityClient, err := config.IdentityV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	if d.HasChange("rules") {
		rules, err := expandIdentityMappingV3Rules(d.Get("rules").(string))
		if err != nil {
			return diag.FromErr(err)
		}

		updateOpts := federation.UpdateMappingOpts{
			Rules: rules,
		}

		log.Printf("[DEBUG] openstack_identity_mapping_v3 update options: %#v", updateOpts)

		_, err = federation.UpdateMapping(ctx, identityClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_identity_mapping_v3 %s: %s", d.Id(), err)
		}
	}

	return resourceIdentityMappingV3Read(ctx, d, meta)
}

func resourceIdentityMappingV3Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	identityClient, err := config.IdentityV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	err = federation.DeleteMapping(ctx, identityClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_identity_mapping_v3"))
	}

	return nil
}
