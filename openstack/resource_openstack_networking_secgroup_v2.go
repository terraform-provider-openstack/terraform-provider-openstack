package openstack

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/extensions/security/rules"
)

func resourceNetworkingSecGroupV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingSecGroupV2Create,
		ReadContext:   resourceNetworkingSecGroupV2Read,
		UpdateContext: resourceNetworkingSecGroupV2Update,
		DeleteContext: resourceNetworkingSecGroupV2Delete,
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

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"delete_default_rules": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"all_tags": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"rule": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Set:      networkingSecgroupV2RuleHash,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"description": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},

						"direction": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: false,
						},

						"ethertype": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: false,
						},

						"protocol": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},

						"port_range_min": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: false,
						},

						"port_range_max": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: false,
						},

						"remote_ip_prefix": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
							StateFunc: func(v interface{}) string {
								return strings.ToLower(v.(string))
							},
						},

						"remote_group_id": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: false,
						},

						"self": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
							ForceNew: false,
						},
					},
				},
			},
		},
	}
}

func resourceNetworkingSecGroupV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	// Before creating the security group, make sure all rules are valid.
	if err := networkingSecgroupV2RulesCheckForErrors(d); err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	opts := groups.CreateOpts{
		Name:        name,
		Description: d.Get("description").(string),
		TenantID:    d.Get("tenant_id").(string),
	}

	log.Printf("[DEBUG] openstack_networking_secgroup_v2 create options: %#v", opts)
	sg, err := groups.Create(networkingClient, opts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_networking_secgroup_v2: %s", err)
	}

	// Delete the default security group rules if it has been requested.
	deleteDefaultRules := d.Get("delete_default_rules").(bool)
	if deleteDefaultRules {
		sgID := sg.ID
		sg, err := groups.Get(networkingClient, sgID).Extract()
		if err != nil {
			return diag.Errorf("Error retrieving the created openstack_networking_secgroup_v2 %s: %s", sgID, err)
		}

		for _, rule := range sg.Rules {
			if err := rules.Delete(networkingClient, rule.ID).ExtractErr(); err != nil {
				return diag.Errorf("Error deleting a default rule for openstack_networking_secgroup_v2 %s: %s", sgID, err)
			}
		}
	}

	d.SetId(sg.ID)

	tags := networkingV2AttributesTags(d)
	if len(tags) > 0 {
		tagOpts := attributestags.ReplaceAllOpts{Tags: tags}
		tags, err := attributestags.ReplaceAll(networkingClient, "security-groups", sg.ID, tagOpts).Extract()
		if err != nil {
			return diag.Errorf("Error setting tags on openstack_networking_secgroup_v2 %s: %s", sg.ID, err)
		}
		log.Printf("[DEBUG] Set tags %s on openstack_networking_secgroup_v2 %s", tags, sg.ID)
	}

	// Now that the security group has been created, iterate through each rule and create it
	createRuleOptsList, err := expandNetworkingSecgroupV2CreateRules(d)
	if err != nil {
		return diag.Errorf("Error creating openstack_networking_secgroup_v2 %s rules:", err)
	}

	for _, createRuleOpts := range createRuleOptsList {
		_, err := rules.Create(networkingClient, createRuleOpts).Extract()
		if err != nil {
			return diag.Errorf("Error creating openstack_networking_secgroup_v2 %s rule: %s", name, err)
		}
	}

	log.Printf("[DEBUG] Created openstack_networking_secgroup_v2: %#v", sg)

	return resourceNetworkingSecGroupV2Read(ctx, d, meta)
}

func resourceNetworkingSecGroupV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	sg, err := groups.Get(networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_networking_secgroup_v2"))
	}

	d.Set("description", sg.Description)
	d.Set("tenant_id", sg.TenantID)
	d.Set("name", sg.Name)
	d.Set("region", GetRegion(d, config))

	networkingV2ReadAttributesTags(d, sg.Tags)

	sgrPager, err := rules.List(networkingClient, rules.ListOpts{SecGroupID: d.Id()}).AllPages()
	if err != nil {
		return diag.Errorf("Error retrieving openstack_networking_secgroup_v2: %s", err)
	}

	sgRules, err := rules.ExtractRules(sgrPager)
	if err != nil {
		return diag.Errorf("Error retrieving openstack_networking_secgroup_v2: %s", err)
	}

	rules, err := flattenNetworkingSecgroupV2Rules(networkingClient, d, sgRules)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Retrieved openstack_networking_secgroup_v2 %s rules: %#v", d.Id(), rules)

	if err := d.Set("rule", rules); err != nil {
		return diag.Errorf("Unable to set openstack_networking_secgroup_v2 %s rules: %s", d.Id(), err)
	}

	return nil
}

func resourceNetworkingSecGroupV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	var (
		updated    bool
		updateOpts groups.UpdateOpts
	)

	if d.HasChange("name") {
		updated = true
		updateOpts.Name = d.Get("name").(string)
	}

	if d.HasChange("description") {
		updated = true
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}

	if updated {
		log.Printf("[DEBUG] Updating openstack_networking_secgroup_v2 %s with options: %#v", d.Id(), updateOpts)
		_, err = groups.Update(networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_networking_secgroup_v2: %s", err)
		}
	}

	if d.HasChange("tags") {
		tags := networkingV2UpdateAttributesTags(d)
		tagOpts := attributestags.ReplaceAllOpts{Tags: tags}
		tags, err := attributestags.ReplaceAll(networkingClient, "security-groups", d.Id(), tagOpts).Extract()
		if err != nil {
			return diag.Errorf("Error setting tags on openstack_networking_secgroup_v2 %s: %s", d.Id(), err)
		}
		log.Printf("[DEBUG] Set tags %s on openstack_networking_secgroup_v2 %s", tags, d.Id())
	}

	if d.HasChange("rule") {
		oldSGRaw, newSGRaw := d.GetChange("rule")
		oldSGRSet, newSGRSet := oldSGRaw.(*schema.Set), newSGRaw.(*schema.Set)
		secgrouprulesToAdd := newSGRSet.Difference(oldSGRSet)
		secgrouprulesToRemove := oldSGRSet.Difference(newSGRSet)

		log.Printf("[DEBUG] openstack_networking_secgroup_v2 %s rules to add: %v", d.Id(), secgrouprulesToAdd)
		log.Printf("[DEBUG] openstack_networking_secgroup_v2 %s rules to remove: %v", d.Id(), secgrouprulesToRemove)

		for _, rawRule := range secgrouprulesToAdd.List() {
			createRuleOpts, err := expandNetworkingSecgroupV2CreateRule(d, rawRule)
			if err != nil {
				return diag.Errorf("Error adding rule to openstack_networking_secgroup_v2 %s: %s", d.Id(), err)
			}

			_, erro := rules.Create(networkingClient, createRuleOpts).Extract()
			if erro != nil {
				return diag.Errorf("Error adding rule to openstack_networking_secgroup_v2 %s: %s", d.Id(), err)
			}
		}

		for _, r := range secgrouprulesToRemove.List() {
			rule := expandNetworkingSecgroupV2Rule(d, r)

			log.Printf("[DEBUG] openstack_networking_secgroup_v2 %s removing rule %v", d.Id(), rule.ID)
			err := rules.Delete(networkingClient, rule.ID).ExtractErr()
			if err != nil {
				if _, ok := err.(gophercloud.ErrDefault404); ok {
					continue
				}

				return diag.Errorf("Error removing rule %s from openstack_networking_secgroup_v2 %s: %s", rule.ID, d.Id(), err)
			}
		}
	}

	return resourceNetworkingSecGroupV2Read(ctx, d, meta)
}

func resourceNetworkingSecGroupV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	networkingClient, err := config.NetworkingV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    networkingSecgroupV2StateRefreshFuncDelete(networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error deleting openstack_networking_secgroup_v2: %s", err)
	}

	return diag.FromErr(err)
}
