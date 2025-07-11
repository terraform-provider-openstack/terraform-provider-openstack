package openstack

import (
	"context"
	"log"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/attributestags"
	"github.com/gophercloud/gophercloud/v2/openstack/networking/v2/extensions/trunks"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNetworkingTrunkV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNetworkingTrunkV2Create,
		ReadContext:   resourceNetworkingTrunkV2Read,
		UpdateContext: resourceNetworkingTrunkV2Update,
		DeleteContext: resourceNetworkingTrunkV2Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
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
				ForceNew: false,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"port_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"sub_port": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port_id": {
							Type:     schema.TypeString,
							Required: true,
						},
						"segmentation_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"segmentation_id": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},

			"tags": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: false,
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

func resourceNetworkingTrunkV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	client, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	createOpts := trunks.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		PortID:      d.Get("port_id").(string),
		TenantID:    d.Get("tenant_id").(string),
		Subports:    expandNetworkingTrunkV2Subports(d.Get("sub_port").(*schema.Set)),
	}

	if v, ok := getOkExists(d, "admin_state_up"); ok {
		asu := v.(bool)
		createOpts.AdminStateUp = &asu
	}

	log.Printf("[DEBUG] openstack_networking_trunk_v2 create options: %#v", createOpts)

	trunk, err := trunks.Create(ctx, client, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_networking_trunk_v2: %s", err)
	}

	log.Printf("[DEBUG] Waiting for openstack_networking_trunk_v2 %s to become available.", trunk.ID)

	stateConf := &retry.StateChangeConf{
		Target:     []string{"ACTIVE", "DOWN"},
		Refresh:    networkingTrunkV2StateRefreshFunc(ctx, client, trunk.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_networking_trunk_v2 %s to become available: %s", trunk.ID, err)
	}

	d.SetId(trunk.ID)

	tags := networkingV2AttributesTags(d)
	if len(tags) > 0 {
		tagOpts := attributestags.ReplaceAllOpts{Tags: tags}

		tags, err := attributestags.ReplaceAll(ctx, client, "trunks", trunk.ID, tagOpts).Extract()
		if err != nil {
			return diag.Errorf("Error setting tags on openstack_networking_trunk_v2 %s: %s", trunk.ID, err)
		}

		log.Printf("[DEBUG] Set tags %s on openstack_networking_trunk_v2 %s", tags, trunk.ID)
	}

	log.Printf("[DEBUG] Created openstack_networking_trunk_v2 %s: %#v", trunk.ID, trunk)

	return resourceNetworkingTrunkV2Read(ctx, d, meta)
}

func resourceNetworkingTrunkV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	client, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	trunk, err := trunks.Get(ctx, client, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting openstack_networking_trunk_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_networking_trunk_v2 %s: %#v", d.Id(), trunk)

	d.Set("region", GetRegion(d, config))
	d.Set("name", trunk.Name)
	d.Set("description", trunk.Description)
	d.Set("port_id", trunk.PortID)
	d.Set("admin_state_up", trunk.AdminStateUp)
	d.Set("tenant_id", trunk.TenantID)

	networkingV2ReadAttributesTags(d, trunk.Tags)

	err = d.Set("sub_port", flattenNetworkingTrunkV2Subports(trunk.Subports))
	if err != nil {
		log.Printf("[DEBUG] Unable to set openstack_networking_trunk_v2 %s sub_port: %s", d.Id(), err)
	}

	return nil
}

func resourceNetworkingTrunkV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	client, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	// Work with basic trunk update options.
	var (
		updateTrunk bool
		updateOpts  trunks.UpdateOpts
	)

	if d.HasChange("name") {
		updateTrunk = true
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}

	if d.HasChange("description") {
		updateTrunk = true
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}

	if d.HasChange("admin_state_up") {
		updateTrunk = true
		asu := d.Get("admin_state_up").(bool)
		updateOpts.AdminStateUp = &asu
	}

	if updateTrunk {
		log.Printf("[DEBUG] openstack_networking_trunk_v2 %s update options: %#v", d.Id(), updateOpts)

		_, err = trunks.Update(ctx, client, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_networking_trunk_v2 %s: %s", d.Id(), err)
		}
	}

	// Update subports  if needed.
	if d.HasChange("sub_port") {
		o, n := d.GetChange("sub_port")
		oldSubport := o.(*schema.Set)
		newSubport := n.(*schema.Set)

		// Delete all old subports, regardless of if they still exist.
		// If they do still exist, they will be re-added below.
		if oldSubport.Len() != 0 {
			removeSubports := expandNetworkingTrunkV2SubportsRemove(oldSubport)
			removeSubportsOpts := trunks.RemoveSubportsOpts{
				Subports: removeSubports,
			}

			log.Printf("[DEBUG] Deleting old subports for openstack_networking_trunk_v2 %s: %#v", d.Id(), removeSubportsOpts)

			_, err := trunks.RemoveSubports(ctx, client, d.Id(), removeSubportsOpts).Extract()
			if err != nil {
				return diag.Errorf("Error removing subports for openstack_networking_trunk_v2 %s: %s", d.Id(), err)
			}
		}

		// Add any new sub_port and re-add previously set subports.
		if newSubport.Len() != 0 {
			addSubports := expandNetworkingTrunkV2Subports(newSubport)
			addSubportsOpts := trunks.AddSubportsOpts{
				Subports: addSubports,
			}

			log.Printf("[DEBUG] openstack_networking_trunk_v2 %s subports update options: %#v", d.Id(), addSubports)

			_, err := trunks.AddSubports(ctx, client, d.Id(), addSubportsOpts).Extract()
			if err != nil {
				return diag.Errorf("Error updating openstack_networking_trunk_v2 %s subports: %s", d.Id(), err)
			}
		}
	}

	if d.HasChange("tags") {
		tags := networkingV2UpdateAttributesTags(d)
		tagOpts := attributestags.ReplaceAllOpts{Tags: tags}

		tags, err := attributestags.ReplaceAll(ctx, client, "trunks", d.Id(), tagOpts).Extract()
		if err != nil {
			return diag.Errorf("Error setting tags on openstack_networking_trunk_v2 %s: %s", d.Id(), err)
		}

		log.Printf("[DEBUG] Set tags %s on openstack_networking_trunk_v2 %s", tags, d.Id())
	}

	return resourceNetworkingTrunkV2Read(ctx, d, meta)
}

func resourceNetworkingTrunkV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	client, err := config.NetworkingV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack networking client: %s", err)
	}

	if err := trunks.Delete(ctx, client, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_networking_trunk_v2"))
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"ACTIVE", "DOWN"},
		Target:     []string{"DELETED"},
		Refresh:    networkingTrunkV2StateRefreshFunc(ctx, client, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      0,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_networking_trunk_v2 %s to Delete:  %s", d.Id(), err)
	}

	return nil
}
