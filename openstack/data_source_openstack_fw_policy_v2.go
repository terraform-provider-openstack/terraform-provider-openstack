package openstack

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/policies"
)

func dataSourceFWPolicyV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFWPolicyV2Read,

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

			"policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"tenant_id": {
				Type:          schema.TypeString,
				Computed:      true,
				Optional:      true,
				ConflictsWith: []string{"project_id"},
			},

			"project_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"tenant_id"},
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"shared": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"audited": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceFWPolicyV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	listOpts := policies.ListOpts{}

	if v, ok := d.GetOk("name"); ok {
		listOpts.Name = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		listOpts.Description = v.(string)
	}

	if v, ok := d.GetOk("policy_id"); ok {
		listOpts.ID = v.(string)
	}

	if v, ok := d.GetOk("tenant_id"); ok {
		listOpts.TenantID = v.(string)
	}

	if v, ok := d.GetOk("project_id"); ok {
		listOpts.ProjectID = v.(string)
	}

	if v, ok := d.GetOk("shared"); ok {
		shared := v.(bool)
		listOpts.Shared = &shared
	}

	if v, ok := d.GetOk("audited"); ok {
		audited := v.(bool)
		listOpts.Audited = &audited
	}

	pages, err := policies.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("Unable to list openstack_fw_policy_v2 policies: %s", err)
	}

	allFWPolicies, err := policies.ExtractPolicies(pages)
	if err != nil {
		return diag.Errorf("Unable to retrieve openstack_fw_policy_v2: %s", err)
	}

	if len(allFWPolicies) < 1 {
		return diag.Errorf("Your openstack_fw_policy_v2 query returned no results")
	}

	if len(allFWPolicies) > 1 {
		return diag.Errorf("Your openstack_fw_policy_v2 query returned more than one result")
	}

	policy := allFWPolicies[0]

	log.Printf("[DEBUG] Retrieved openstack_fw_policy_v2 %s: %#v", policy.ID, policy)

	d.SetId(policy.ID)

	d.Set("name", policy.Name)
	d.Set("tenant_id", policy.TenantID)
	d.Set("project_id", policy.ProjectID)
	d.Set("description", policy.Description)
	d.Set("shared", policy.Shared)
	d.Set("audited", policy.Audited)
	d.Set("rules", policy.Rules)
	d.Set("region", GetRegion(d, config))

	return nil
}
