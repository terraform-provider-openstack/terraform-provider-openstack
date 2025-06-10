package openstack

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/qos/rules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNetworkingQoSDSCPMarkingRuleV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNetworkingQoSDSCPMarkingRuleV2Read,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"qos_policy_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"dscp_mark": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
				ForceNew: false,
			},
		},
	}
}

func dataSourceNetworkingQoSDSCPMarkingRuleV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	listOpts := rules.DSCPMarkingRulesListOpts{}

	if v, ok := d.GetOk("dscp_mark"); ok {
		listOpts.DSCPMark = v.(int)
	}

	qosPolicyID := d.Get("qos_policy_id").(string)

	pages, err := rules.ListDSCPMarkingRules(networkingClient, qosPolicyID, listOpts).AllPages(ctx)
	if err != nil {
		return diag.Errorf("Unable to retrieve openstack_networking_qos_dscp_marking_rule_v2: %s", err)
	}

	allRules, err := rules.ExtractDSCPMarkingRules(pages)
	if err != nil {
		return diag.Errorf("Unable to extract openstack_networking_qos_dscp_marking_rule_v2: %s", err)
	}

	if len(allRules) < 1 {
		return diag.Errorf("Your query returned no openstack_networking_qos_dscp_marking_rule_v2. " +
			"Please change your search criteria and try again.")
	}

	if len(allRules) > 1 {
		return diag.Errorf("Your query returned more than one openstack_networking_qos_dscp_marking_rule_v2." +
			" Please try a more specific search criteria")
	}

	rule := allRules[0]
	id := resourceNetworkingQoSRuleV2BuildID(qosPolicyID, rule.ID)

	log.Printf("[DEBUG] Retrieved openstack_networking_qos_dscp_marking_rule_v2 %s: %+v", id, rule)
	d.SetId(id)

	d.Set("qos_policy_id", qosPolicyID)
	d.Set("dscp_mark", rule.DSCPMark)
	d.Set("region", GetRegion(d, config))

	return nil
}
