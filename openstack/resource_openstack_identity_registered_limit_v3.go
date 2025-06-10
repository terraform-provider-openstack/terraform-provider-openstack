package openstack

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/registeredlimits"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceIdentityRegisteredLimitV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityRegisteredLimitV3Create,
		ReadContext:   resourceIdentityRegisteredLimitV3Read,
		UpdateContext: resourceIdentityRegisteredLimitV3Update,
		DeleteContext: resourceIdentityRegisteredLimitV3Delete,
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

			"service_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"default_limit": {
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

func resourceIdentityRegisteredLimitV3Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	identityClient, err := config.IdentityV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	createOpts := registeredlimits.BatchCreateOpts{
		registeredlimits.CreateOpts{
			RegionID:     GetRegion(d, config),
			ServiceID:    d.Get("service_id").(string),
			ResourceName: d.Get("resource_name").(string),
			DefaultLimit: d.Get("default_limit").(int),
			Description:  d.Get("description").(string),
		},
	}

	log.Printf("[DEBUG] openstack_identity_registered_limit_v3 create options: %#v", createOpts)

	registeredlimit, err := registeredlimits.BatchCreate(ctx, identityClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_identity_registered_limit_v3: %s", err)
	}

	d.SetId(registeredlimit[0].ID)

	return resourceIdentityRegisteredLimitV3Read(ctx, d, meta)
}

func resourceIdentityRegisteredLimitV3Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	identityClient, err := config.IdentityV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	registeredlimit, err := registeredlimits.Get(ctx, identityClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_identity_registered_limit_v3"))
	}

	log.Printf("[DEBUG] Retrieved openstack_identity_registered_limit_v3: %#v", registeredlimit)

	d.Set("region", GetRegion(d, config))
	d.Set("service_id", registeredlimit.ServiceID)
	d.Set("resource_name", registeredlimit.ResourceName)
	d.Set("default_limit", registeredlimit.DefaultLimit)
	d.Set("description", registeredlimit.Description)

	return nil
}

func resourceIdentityRegisteredLimitV3Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	identityClient, err := config.IdentityV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	var hasChange bool

	var updateOpts registeredlimits.UpdateOpts

	if d.HasChange("default_limit") {
		hasChange = true
		defaultLimit := d.Get("default_limit").(int)
		updateOpts.DefaultLimit = &defaultLimit
	}

	if d.HasChange("service_id") {
		hasChange = true
		updateOpts.ServiceID = d.Get("service_id").(string)
	}

	if d.HasChange("resource_name") {
		hasChange = true
		updateOpts.ResourceName = d.Get("resource_name").(string)
	}

	if d.HasChange("description") {
		hasChange = true
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}

	if hasChange {
		_, err := registeredlimits.Update(ctx, identityClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_identity_registered_limit_v3 %s: %s", d.Id(), err)
		}
	}

	return resourceIdentityRegisteredLimitV3Read(ctx, d, meta)
}

func resourceIdentityRegisteredLimitV3Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	identityClient, err := config.IdentityV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	err = registeredlimits.Delete(ctx, identityClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_identity_registered_limit_v3"))
	}

	return nil
}
