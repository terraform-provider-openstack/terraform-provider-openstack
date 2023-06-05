package openstack

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/policies"
)

func resourceFWPolicyV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFWPolicyV2Create,
		ReadContext:   resourceFWPolicyV2Read,
		UpdateContext: resourceFWPolicyV2Update,
		DeleteContext: resourceFWPolicyV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"audited": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"shared": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"rules": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceFWPolicyV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	audited := d.Get("audited").(bool)
	opts := policies.CreateOpts{
		Name:          d.Get("name").(string),
		Description:   d.Get("description").(string),
		Audited:       &audited,
		TenantID:      d.Get("tenant_id").(string),
		FirewallRules: expandToStringSlice(d.Get("rules").([]interface{})),
	}

	if r, ok := d.GetOk("shared"); ok {
		shared := r.(bool)
		opts.Shared = &shared
	}

	log.Printf("[DEBUG] openstack_fw_policy_v2 create options: %#v", opts)

	policy, err := policies.Create(networkingClient, opts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_fw_policy_v2: %s", err)
	}

	log.Printf("[DEBUG] openstack_fw_policy_v2 %s created: %#v", policy.ID, policy)

	d.SetId(policy.ID)

	return resourceFWPolicyV2Read(ctx, d, meta)
}

func resourceFWPolicyV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	policy, err := policies.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_fw_policy_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_fw_policy_v2 %s: %#v", d.Id(), policy)

	d.Set("name", policy.Name)
	d.Set("description", policy.Description)
	d.Set("shared", policy.Shared)
	d.Set("audited", policy.Audited)
	d.Set("tenant_id", policy.TenantID)
	d.Set("rules", policy.Rules)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceFWPolicyV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var updateOpts policies.UpdateOpts

	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}

	if d.HasChange("rules") {
		rules := expandToStringSlice(d.Get("rules").([]interface{}))
		updateOpts.FirewallRules = &rules
	}

	log.Printf("[DEBUG] openstack_fw_policy_v2 %s update options: %#v", d.Id(), updateOpts)

	err = policies.Update(networkingClient, d.Id(), updateOpts).Err
	if err != nil {
		return diag.Errorf("Error updating openstack_fw_policy_v2 %s: %s", d.Id(), err)
	}

	return resourceFWPolicyV2Read(ctx, d, meta)
}

func resourceFWPolicyV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	_, err = policies.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_fw_policy_v2"))
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    fwPolicyV2DeleteFunc(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("Error waiting for openstack_fw_policy_v2 %s to be deleted: %s", d.Id(), err)
	}

	return nil
}
