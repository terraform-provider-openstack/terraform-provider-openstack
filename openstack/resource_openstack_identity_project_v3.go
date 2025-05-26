package openstack

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/v2/openstack/identity/v3/projects"
)

func resourceIdentityProjectV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityProjectV3Create,
		ReadContext:   resourceIdentityProjectV3Read,
		UpdateContext: resourceIdentityProjectV3Update,
		DeleteContext: resourceIdentityProjectV3Delete,
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

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"is_domain": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"parent_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"extra": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func resourceIdentityProjectV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	enabled := d.Get("enabled").(bool)
	isDomain := d.Get("is_domain").(bool)
	createOpts := projects.CreateOpts{
		Description: d.Get("description").(string),
		DomainID:    d.Get("domain_id").(string),
		Enabled:     &enabled,
		IsDomain:    &isDomain,
		Name:        d.Get("name").(string),
		Extra:       d.Get("extra").(map[string]interface{}),
		ParentID:    d.Get("parent_id").(string),
	}

	if v, ok := d.GetOk("tags"); ok {
		tags := v.(*schema.Set).List()
		createOpts.Tags = expandToStringSlice(tags)
	}

	log.Printf("[DEBUG] openstack_identity_project_v3 create options: %#v", createOpts)
	project, err := projects.Create(ctx, identityClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_identity_project_v3: %s", err)
	}

	d.SetId(project.ID)

	return resourceIdentityProjectV3Read(ctx, d, meta)
}

func resourceIdentityProjectV3Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	project, err := projects.Get(ctx, identityClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_identity_project_v3"))
	}

	log.Printf("[DEBUG] Retrieved openstack_identity_project_v3 %s: %#v", d.Id(), project)

	d.Set("description", project.Description)
	d.Set("domain_id", project.DomainID)
	d.Set("enabled", project.Enabled)
	d.Set("is_domain", project.IsDomain)
	d.Set("name", project.Name)
	d.Set("extra", expandToMapStringString(project.Extra))
	d.Set("parent_id", project.ParentID)
	d.Set("region", GetRegion(d, config))
	d.Set("tags", project.Tags)

	return nil
}

func resourceIdentityProjectV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	var hasChange bool
	var updateOpts projects.UpdateOpts

	if d.HasChange("domain_id") {
		hasChange = true
		updateOpts.DomainID = d.Get("domain_id").(string)
	}

	if d.HasChange("enabled") {
		hasChange = true
		enabled := d.Get("enabled").(bool)
		updateOpts.Enabled = &enabled
	}

	if d.HasChange("is_domain") {
		hasChange = true
		isDomain := d.Get("is_domain").(bool)
		updateOpts.IsDomain = &isDomain
	}

	if d.HasChange("name") {
		hasChange = true
		updateOpts.Name = d.Get("name").(string)
	}

	if d.HasChange("parent_id") {
		hasChange = true
		updateOpts.ParentID = d.Get("parent_id").(string)
	}

	if d.HasChange("description") {
		hasChange = true
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}

	if d.HasChange("extra") {
		hasChange = true
		updateOpts.Extra = resourceIdentityProjectV3ExtraChange(d)
	}

	if d.HasChange("tags") {
		hasChange = true
		if v, ok := d.GetOk("tags"); ok {
			tags := v.(*schema.Set).List()
			tagsToUpdate := expandToStringSlice(tags)
			updateOpts.Tags = &tagsToUpdate
		} else {
			updateOpts.Tags = &[]string{}
		}
	}

	if hasChange {
		_, err := projects.Update(ctx, identityClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_identity_project_v3 %s: %s", d.Id(), err)
		}
	}

	return resourceIdentityProjectV3Read(ctx, d, meta)
}

func resourceIdentityProjectV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	identityClient, err := config.IdentityV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack identity client: %s", err)
	}

	if d.Get("is_domain").(bool) {
		log.Printf("[DEBUG] openstack_identity_project_v3 %s is domain, disabling", d.Id())
		updateOpts := projects.UpdateOpts{
			Enabled: new(bool),
		}
		_, err := projects.Update(ctx, identityClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error disabling domain openstack_identity_project_v3 %s: %s", d.Id(), err)
		}
	}

	err = projects.Delete(ctx, identityClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_identity_project_v3"))
	}

	return nil
}

func resourceIdentityProjectV3ExtraChange(d *schema.ResourceData) map[string]interface{} {
	o, n := d.GetChange("extra")
	oldExtra := o.(map[string]interface{})
	newExtra := n.(map[string]interface{})
	extra := newExtra

	for oldKey := range oldExtra {
		// unset old keys
		if _, ok := newExtra[oldKey]; !ok {
			extra[oldKey] = nil
		}
	}

	return extra
}
