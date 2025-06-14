package openstack

import (
	"context"
	"log"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/groups"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/security/rules"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

			"stateful": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
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
		},
	}
}

func resourceNetworkingSecGroupV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	opts := groups.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		TenantID:    d.Get("tenant_id").(string),
	}

	if v, ok := getOkExists(d, "stateful"); ok {
		v := v.(bool)
		opts.Stateful = &v
	}

	log.Printf("[DEBUG] openstack_networking_secgroup_v2 create options: %#v", opts)

	sg, err := groups.Create(ctx, networkingClient, opts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_networking_secgroup_v2: %s", err)
	}

	// Delete the default security group rules if it has been requested.
	deleteDefaultRules := d.Get("delete_default_rules").(bool)
	if deleteDefaultRules {
		sgID := sg.ID
		sg, err := groups.Get(ctx, networkingClient, sgID).Extract()
		if err != nil {
			return diag.Errorf("Error retrieving the created openstack_networking_secgroup_v2 %s: %s", sgID, err)
		}

		for _, rule := range sg.Rules {
			if err := rules.Delete(ctx, networkingClient, rule.ID).ExtractErr(); err != nil {
				return diag.Errorf("Error deleting a default rule for openstack_networking_secgroup_v2 %s: %s", sgID, err)
			}
		}
	}

	d.SetId(sg.ID)

	tags := networkingV2AttributesTags(d)
	if len(tags) > 0 {
		tagOpts := attributestags.ReplaceAllOpts{Tags: tags}

		tags, err := attributestags.ReplaceAll(ctx, networkingClient, "security-groups", sg.ID, tagOpts).Extract()
		if err != nil {
			return diag.Errorf("Error setting tags on openstack_networking_secgroup_v2 %s: %s", sg.ID, err)
		}

		log.Printf("[DEBUG] Set tags %s on openstack_networking_secgroup_v2 %s", tags, sg.ID)
	}

	log.Printf("[DEBUG] Created openstack_networking_secgroup_v2: %#v", sg)

	return resourceNetworkingSecGroupV2Read(ctx, d, meta)
}

func resourceNetworkingSecGroupV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	sg, err := groups.Get(ctx, networkingClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_networking_secgroup_v2"))
	}

	log.Printf("[DEBUG] Created openstack_networking_secgroup_v2: %#v", sg)

	d.Set("description", sg.Description)
	d.Set("tenant_id", sg.TenantID)
	d.Set("name", sg.Name)
	d.Set("stateful", sg.Stateful)
	d.Set("region", GetRegion(d, config))

	networkingV2ReadAttributesTags(d, sg.Tags)

	return nil
}

func resourceNetworkingSecGroupV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
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

	if d.HasChange("stateful") {
		updated = true
		stateful := d.Get("stateful").(bool)
		updateOpts.Stateful = &stateful
	}

	if updated {
		log.Printf("[DEBUG] Updating openstack_networking_secgroup_v2 %s with options: %#v", d.Id(), updateOpts)

		_, err = groups.Update(ctx, networkingClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_networking_secgroup_v2: %s", err)
		}
	}

	if d.HasChange("tags") {
		tags := networkingV2UpdateAttributesTags(d)
		tagOpts := attributestags.ReplaceAllOpts{Tags: tags}

		tags, err := attributestags.ReplaceAll(ctx, networkingClient, "security-groups", d.Id(), tagOpts).Extract()
		if err != nil {
			return diag.Errorf("Error setting tags on openstack_networking_secgroup_v2 %s: %s", d.Id(), err)
		}

		log.Printf("[DEBUG] Set tags %s on openstack_networking_secgroup_v2 %s", tags, d.Id())
	}

	return resourceNetworkingSecGroupV2Read(ctx, d, meta)
}

func resourceNetworkingSecGroupV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	networkingClient, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    networkingSecgroupV2StateRefreshFuncDelete(ctx, networkingClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      0,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error deleting openstack_networking_secgroup_v2: %s", err)
	}

	return diag.FromErr(err)
}
