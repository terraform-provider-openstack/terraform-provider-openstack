package openstack

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/rules"
)

func dataSourceFWRuleV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFWRuleV2Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"rule_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"name": {
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

			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"icmp", "tcp", "udp",
				}, true),
			},

			"action": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"allow", "deny", "reject",
				}, true),
			},

			"ip_version": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateFunc: validation.IntInSlice([]int{
					4, 6,
				}),
			},

			"source_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"source_port": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"destination_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"destination_port": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"shared": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"firewall_policy_id": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceFWRuleV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	listOpts := rules.ListOpts{}

	if v, ok := d.GetOk("rule_id"); ok {
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

	if v, ok := d.GetOk("protocol"); ok {
		listOpts.Protocol = rules.Protocol(v.(string))
	}

	if v, ok := d.GetOk("action"); ok {
		listOpts.Action = rules.Action(v.(string))
	}

	if v, ok := d.GetOk("ip_version"); ok {
		listOpts.IPVersion = v.(int)
	}

	if v, ok := d.GetOk("source_ip_address"); ok {
		listOpts.SourceIPAddress = v.(string)
	}

	if v, ok := d.GetOk("source_port"); ok {
		listOpts.SourcePort = v.(string)
	}

	if v, ok := d.GetOk("destination_ip_address"); ok {
		listOpts.DestinationIPAddress = v.(string)
	}

	if v, ok := d.GetOk("destination_port"); ok {
		listOpts.DestinationPort = v.(string)
	}

	if v, ok := d.GetOk("shared"); ok {
		shared := v.(bool)
		listOpts.Shared = &shared
	}

	if v, ok := d.GetOk("enabled"); ok {
		enabled := v.(bool)
		listOpts.Enabled = &enabled
	}

	pages, err := rules.List(networkingClient, listOpts).AllPages()
	if err != nil {
		return diag.Errorf("Unable to list openstack_fw_rule_v2 rules: %s", err)
	}

	allFWRules, err := rules.ExtractRules(pages)
	if err != nil {
		return diag.Errorf("Unable to retrieve openstack_fw_rule_v2: %s", err)
	}

	if len(allFWRules) < 1 {
		return diag.Errorf("Your openstack_fw_rule_v2 query returned no results")
	}

	if len(allFWRules) > 1 {
		return diag.Errorf("Your openstack_fw_rule_v2 query returned more than one result")
	}

	rule := allFWRules[0]

	log.Printf("[DEBUG] Retrieved openstack_fw_rule_v2 %s: %#v", rule.ID, rule)

	d.SetId(rule.ID)

	d.Set("name", rule.Name)
	d.Set("description", rule.Description)
	d.Set("tenant_id", rule.TenantID)
	d.Set("project_id", rule.ProjectID)
	d.Set("action", rule.Action)
	d.Set("ip_version", rule.IPVersion)
	d.Set("source_ip_address", rule.SourceIPAddress)
	d.Set("source_port", rule.SourcePort)
	d.Set("destination_ip_address", rule.DestinationIPAddress)
	d.Set("destination_port", rule.DestinationPort)
	d.Set("shared", rule.Shared)
	d.Set("enabled", rule.Enabled)
	d.Set("firewall_policy_id", rule.FirewallPolicyID)
	d.Set("region", GetRegion(d, config))

	if rule.Protocol == "" {
		d.Set("protocol", "any")
	} else {
		d.Set("protocol", rule.Protocol)
	}

	return nil
}
