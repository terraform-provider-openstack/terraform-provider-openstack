package openstack

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/fwaas_v2/rules"
)

func resourceFWRuleV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFWRuleV2Create,
		ReadContext:   resourceFWRuleV2Read,
		UpdateContext: resourceFWRuleV2Update,
		DeleteContext: resourceFWRuleV2Delete,
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
			},

			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"any", "icmp", "tcp", "udp",
				}, true),
				Default: "any",
			},

			"action": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"allow", "deny", "reject",
				}, true),
				Default: "deny",
			},

			"ip_version": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateFunc: validation.IntInSlice([]int{
					4, 6,
				}),
				Default: 4,
			},

			"source_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"destination_ip_address": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"source_port": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"protocol"},
			},

			"destination_port": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"protocol"},
			},

			"shared": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceFWRuleV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	shared := d.Get("shared").(bool)
	enabled := d.Get("enabled").(bool)
	ruleCreateOpts := rules.CreateOpts{
		Name:                 d.Get("name").(string),
		Description:          d.Get("description").(string),
		Protocol:             rules.Protocol((d.Get("protocol").(string))),
		Action:               rules.Action((d.Get("action").(string))),
		IPVersion:            gophercloud.IPVersion(d.Get("ip_version").(int)),
		SourceIPAddress:      d.Get("source_ip_address").(string),
		DestinationIPAddress: d.Get("destination_ip_address").(string),
		SourcePort:           d.Get("source_port").(string),
		DestinationPort:      d.Get("destination_port").(string),
		Shared:               &shared,
		Enabled:              &enabled,
		TenantID:             d.Get("tenant_id").(string),
	}

	log.Printf("[DEBUG] openstack_fw_rule_v2 create options: %#v", ruleCreateOpts)

	rule, err := rules.Create(networkingClient, ruleCreateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_fw_rule_v2: %s", err)
	}

	log.Printf("[DEBUG] Created openstack_fw_rule_v2 %s: %#v", rule.ID, rule)

	d.SetId(rule.ID)

	return resourceFWRuleV2Read(ctx, d, meta)
}

func resourceFWRuleV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	rule, err := rules.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_fw_rule_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_fw_rule_v2 %s: %#v", d.Id(), rule)

	d.Set("name", rule.Name)
	d.Set("description", rule.Description)
	d.Set("action", rule.Action)
	d.Set("ip_version", rule.IPVersion)
	d.Set("source_ip_address", rule.SourceIPAddress)
	d.Set("destination_ip_address", rule.DestinationIPAddress)
	d.Set("source_port", rule.SourcePort)
	d.Set("destination_port", rule.DestinationPort)
	d.Set("shared", rule.Shared)
	d.Set("enabled", rule.Enabled)

	if rule.Protocol == "" {
		d.Set("protocol", "any")
	} else {
		d.Set("protocol", rule.Protocol)
	}

	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceFWRuleV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var updateOpts rules.UpdateOpts
	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}

	if d.HasChange("protocol") {
		protocol := rules.Protocol(d.Get("protocol").(string))
		updateOpts.Protocol = &protocol
	}

	if d.HasChange("action") {
		action := rules.Action(d.Get("action").(string))
		updateOpts.Action = &action
	}

	if d.HasChange("ip_version") {
		ipVersion := gophercloud.IPVersion(d.Get("ip_version").(int))
		updateOpts.IPVersion = &ipVersion
	}

	if d.HasChange("source_ip_address") {
		sourceIPAddress := d.Get("source_ip_address").(string)
		updateOpts.SourceIPAddress = &sourceIPAddress

		// Also include the ip_version.
		ipVersion := gophercloud.IPVersion(d.Get("ip_version").(int))
		updateOpts.IPVersion = &ipVersion
	}

	if d.HasChange("source_port") {
		sourcePort := d.Get("source_port").(string)
		if sourcePort == "" {
			sourcePort = "0"
		}
		updateOpts.SourcePort = &sourcePort

		// Also include the protocol.
		protocol := rules.Protocol(d.Get("protocol").(string))
		updateOpts.Protocol = &protocol
	}

	if d.HasChange("destination_ip_address") {
		destinationIPAddress := d.Get("destination_ip_address").(string)
		updateOpts.DestinationIPAddress = &destinationIPAddress

		// Also include the ip_version.
		ipVersion := gophercloud.IPVersion(d.Get("ip_version").(int))
		updateOpts.IPVersion = &ipVersion
	}

	if d.HasChange("destination_port") {
		destinationPort := d.Get("destination_port").(string)
		if destinationPort == "" {
			destinationPort = "0"
		}

		updateOpts.DestinationPort = &destinationPort

		// Also include the protocol.
		protocol := rules.Protocol(d.Get("protocol").(string))
		updateOpts.Protocol = &protocol
	}

	if d.HasChange("enabled") {
		enabled := d.Get("enabled").(bool)
		updateOpts.Enabled = &enabled
	}

	if d.HasChange("shared") {
		shared := d.Get("shared").(bool)
		updateOpts.Enabled = &shared
	}

	log.Printf("[DEBUG] openstack_fw_rule_v2 %s update options: %#v", d.Id(), updateOpts)
	err = rules.Update(networkingClient, d.Id(), updateOpts).Err
	if err != nil {
		return diag.Errorf("Error updating openstack_fw_rule_v2 %s: %s", d.Id(), err)
	}

	return resourceFWRuleV2Read(ctx, d, meta)
}

func resourceFWRuleV2Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	err = rules.Delete(networkingClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.Errorf("Error deleting openstack_fw_rule_v2 %s: %s", d.Id(), err)
	}

	return nil
}
