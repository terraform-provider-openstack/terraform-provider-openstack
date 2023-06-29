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

			"shared": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
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
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
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

	if v, ok := d.GetOk("admin_state_up"); ok {
		AdminStateUp := v.(bool)
		listOpts.AdminStateUp = &AdminStateUp
	}

	if v, ok := d.GetOk("ingress_firewall_policy_id"); ok {
		listOpts.IngressFirewallPolicyID = v.(string)
	}

	if v, ok := d.GetOk("egress_firewall_policy_id"); ok {
		listOpts.EgressFirewallPolicyID = v.(string)
	}

	if v, ok := d.GetOk("status"); ok {
		listOpts.Status = v.(string)
	}

	pages, err := groups.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("Unable to list openstack_fw_group_v2 groups: %s", err)
	}

	allGroups, err := groups.ExtractGroups(pages)
	if err != nil {
		return diag.Errorf("Unable to retrieve openstack_fw_group_v2: %s", err)
	}

	if len(allGroups) < 1 {
		return diag.Errorf("Your openstack_fw_group_v2 query returned no results")
	}

	if len(allGroups) > 1 {
		return diag.Errorf("Your openstack_fw_group_v2 query returned more than one result")
	}

	group := allGroups[0]

	log.Printf("[DEBUG] Retrieved openstack_fw_policy_v2 %s: %+v", group.ID, group)

	d.SetId(group.ID)

	d.Set("name", group.Name)
	d.Set("description", group.Description)
	d.Set("tenant_id", group.TenantID)
	d.Set("project_id", group.ProjectID)
	d.Set("shared", group.Shared)
	d.Set("admin_state_up", group.AdminStateUp)
	d.Set("ingress_firewall_policy_id", group.IngressFirewallPolicyID)
	d.Set("egress_firewall_policy_id", group.EgressFirewallPolicyID)
	d.Set("ports", group.Ports)
	d.Set("status", group.Status)
	d.Set("region", GetRegion(d, config))

	return nil
}
