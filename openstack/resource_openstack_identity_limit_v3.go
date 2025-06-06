package openstack

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/limits"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIdentityLimitV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityLimitV3Create,
		ReadContext:   resourceIdentityLimitV3Read,
		UpdateContext: resourceIdentityLimitV3Update,
		DeleteContext: resourceIdentityLimitV3Delete,
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

			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"service_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_limit": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceIdentityLimitV3Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	identityClient, err := config.IdentityV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	createOpts := limits.BatchCreateOpts{
		limits.CreateOpts{
			RegionID:      GetRegion(d, config),
			DomainID:      d.Get("domain_id").(string),
			ProjectID:     d.Get("project_id").(string),
			ServiceID:     d.Get("service_id").(string),
			ResourceName:  d.Get("resource_name").(string),
			ResourceLimit: d.Get("resource_limit").(int),
			Description:   d.Get("description").(string),
		},
	}

	log.Printf("[DEBUG] openstack_identity_limit_v3 create options: %#v", createOpts)

	limit, err := limits.BatchCreate(ctx, identityClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_identity_limit_v3: %s", err)
	}

	d.SetId(limit[0].ID)

	return resourceIdentityLimitV3Read(ctx, d, meta)
}

func resourceIdentityLimitV3Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	identityClient, err := config.IdentityV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	limit, err := limits.Get(ctx, identityClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_identity_limit_v3"))
	}

	log.Printf("[DEBUG] Retrieved openstack_identity_limit_v3: %#v", limit)

	d.Set("region", GetRegion(d, config))
	d.Set("domain_id", limit.DomainID)
	d.Set("project_id", limit.ProjectID)
	d.Set("service_id", limit.ServiceID)
	d.Set("resource_name", limit.ResourceName)
	d.Set("resource_limit", limit.ResourceLimit)
	d.Set("description", limit.Description)

	return nil
}

func resourceIdentityLimitV3Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	identityClient, err := config.IdentityV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	var hasChange bool

	var updateOpts limits.UpdateOpts

	if d.HasChange("resource_limit") {
		hasChange = true
		resourceLimit := d.Get("resource_limit").(int)
		updateOpts.ResourceLimit = &resourceLimit
	}

	if d.HasChange("description") {
		hasChange = true
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}

	if hasChange {
		_, err := limits.Update(ctx, identityClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_identity_limit_v3 %s: %s", d.Id(), err)
		}
	}

	return resourceIdentityLimitV3Read(ctx, d, meta)
}

func resourceIdentityLimitV3Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	identityClient, err := config.IdentityV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	err = limits.Delete(ctx, identityClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_identity_limit_v3"))
	}

	return nil
}
