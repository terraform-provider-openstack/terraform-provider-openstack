package openstack

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/rules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceNetworkingSecGroupRuleV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingSecGroupRuleV2Create,
		ReadContext:   resourceNetworkingSecGroupRuleV2Read,
		DeleteContext: resourceNetworkingSecGroupRuleV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Delete: schema.DefaultTimeout(10 * time.Minute),
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
				ForceNew: true,
			},

			"direction": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: resourceNetworkingSecGroupRuleV2Direction,
			},

			"ethertype": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: resourceNetworkingSecGroupRuleV2EtherType,
			},

			"port_range_min": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"protocol", "port_range_max"},
				ValidateFunc: validation.IntBetween(0, 65535),
			},

			"port_range_max": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"protocol", "port_range_min"},
				ValidateFunc: validation.IntBetween(0, 65535),
			},

			"protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: resourceNetworkingSecGroupRuleV2Protocol,
			},

			"remote_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"remote_ip_prefix", "remote_address_group_id"},
			},

			"remote_ip_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"remote_group_id", "remote_address_group_id"},
				StateFunc: func(v any) string {
					return strings.ToLower(v.(string))
				},
			},

			"remote_address_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"remote_group_id", "remote_ip_prefix"},
			},

			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceNetworkingSecGroupRuleV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	securityGroupID := d.Get("security_group_id").(string)
	config.Lock(securityGroupID)
	defer config.Unlock(securityGroupID)

	protocol := d.Get("protocol").(string)
	direction := d.Get("direction").(string)
	etherType := d.Get("ethertype").(string)
	opts := rules.CreateOpts{
		Direction:            rules.RuleDirection(direction),
		EtherType:            rules.RuleEtherType(etherType),
		Protocol:             rules.RuleProtocol(protocol),
		PortRangeMin:         d.Get("port_range_min").(int),
		PortRangeMax:         d.Get("port_range_max").(int),
		Description:          d.Get("description").(string),
		SecGroupID:           securityGroupID,
		RemoteGroupID:        d.Get("remote_group_id").(string),
		RemoteIPPrefix:       d.Get("remote_ip_prefix").(string),
		RemoteAddressGroupID: d.Get("remote_address_group_id").(string),
		ProjectID:            d.Get("tenant_id").(string),
	}

	log.Printf("[DEBUG] openstack_networking_secgroup_rule_v2 create options: %#v", opts)

	sgRule, err := rules.Create(ctx, networkingClient, opts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_networking_secgroup_rule_v2: %s", err)
	}

	d.SetId(sgRule.ID)

	log.Printf("[DEBUG] Created openstack_networking_secgroup_rule_v2 %s: %#v", sgRule.ID, sgRule)

	return resourceNetworkingSecGroupRuleV2Read(ctx, d, meta)
}

func resourceNetworkingSecGroupRuleV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	sgRule, err := rules.Get(ctx, networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting openstack_networking_secgroup_rule_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_networking_secgroup_rule_v2 %s: %#v", d.Id(), sgRule)

	d.Set("description", sgRule.Description)
	d.Set("direction", sgRule.Direction)
	d.Set("ethertype", sgRule.EtherType)
	d.Set("protocol", sgRule.Protocol)
	d.Set("port_range_min", sgRule.PortRangeMin)
	d.Set("port_range_max", sgRule.PortRangeMax)
	d.Set("remote_group_id", sgRule.RemoteGroupID)
	d.Set("remote_ip_prefix", sgRule.RemoteIPPrefix)
	d.Set("remote_address_group_id", sgRule.RemoteAddressGroupID)
	d.Set("security_group_id", sgRule.SecGroupID)
	d.Set("tenant_id", sgRule.TenantID)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceNetworkingSecGroupRuleV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	securityGroupID := d.Get("security_group_id").(string)
	config.Lock(securityGroupID)
	defer config.Unlock(securityGroupID)

	if err := rules.Delete(ctx, networkingClient, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_networking_secgroup_rule_v2"))
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    resourceNetworkingSecGroupRuleV2StateRefreshFunc(ctx, networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      0,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_networking_secgroup_rule_v2 %s to Delete:  %s", d.Id(), err)
	}

	d.SetId("")

	return nil
}
