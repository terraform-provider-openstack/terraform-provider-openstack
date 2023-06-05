package openstack

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/groups"
)

func dataSourceFWGroupV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFWGroupV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"shared": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  true,
			},

			"ingress_firewall_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"egress_firewall_policy_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"ports": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceFWGroupV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	listOpts := groups.ListOpts{}

	if v, ok := d.GetOk("group_id"); ok {
		listOpts.ID = v.(string)
	}

	if v, ok := d.GetOk("name"); ok {
		listOpts.Name = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		listOpts.Description = v.(string)
	}

	if v, ok := d.GetOk("project_id"); ok {
		listOpts.ProjectID = v.(string)
	}

	if v, ok := d.GetOk("tenant_id"); ok {
		listOpts.ProjectID = v.(string)
	}

	if v, ok := d.GetOk("shared"); ok {
		listOpts.Shared = v.(*bool)
	}

	if v, ok := d.GetOk("ingress_firewall_policy_id"); ok {
		listOpts.IngressFirewallPolicyID = v.(string)
	}

	if v, ok := d.GetOk("egress_firewall_policy_id"); ok {
		listOpts.EgressFirewallPolicyID = v.(string)
	}

	pages, err := groups.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("Unable to list Groups: %s", err)
	}

	allGroups, err := groups.ExtractGroups(pages)
	if err != nil {
		return diag.Errorf("Unable to retrieve Groups: %s", err)
	}

	if len(allGroups) < 1 {
		return diag.Errorf("No Group found")
	}

	if len(allGroups) > 1 {
		return diag.Errorf("More than one Group found")
	}

	group := allGroups[0]

	log.Printf("[DEBUG] Retrieved Group %s: %+v", group.ID, group)

	d.SetId(group.ID)

	d.Set("name", group.Name)
	d.Set("description", group.Description)
	d.Set("project_id", group.ProjectID)
	d.Set("tenant_id", group.TenantID)
	d.Set("shared", group.Shared)
	d.Set("admin_state_up", group.AdminStateUp)
	d.Set("status", group.Status)
	d.Set("ingress_firewall_policy_id", group.IngressFirewallPolicyID)
	d.Set("egress_firewall_policy_id", group.EgressFirewallPolicyID)
	d.Set("ports", group.Ports)
	d.Set("region", GetRegion(d, config))

	return nil
}
