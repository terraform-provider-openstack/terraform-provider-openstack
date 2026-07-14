package openstack

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/sharetypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSharedFilesystemShareTypeV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSharedFilesystemShareTypeV2Create,
		ReadContext:   resourceSharedFilesystemShareTypeV2Read,
		UpdateContext: resourceSharedFilesystemShareTypeV2Update,
		DeleteContext: resourceSharedFilesystemShareTypeV2Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
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

			// Setting/reading the description requires Manila microversion
			// 2.50 or later.
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"is_public": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			// The "driver_handles_share_servers" key is mandatory, as
			// required by the Shared File System API. Other well-known
			// keys, such as "snapshot_support", as well as arbitrary
			// backend-specific keys, are also accepted here.
			"extra_specs": {
				Type:     schema.TypeMap,
				Required: true,
			},

			"is_default": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

// resourceSharedFilesystemShareTypeV2ExtraSpecs splits the "extra_specs" map
// into the strongly-typed options accepted by the share type Create call
// and the remaining, arbitrary extra specs that must be applied separately
// via the extra_specs API once the share type exists.
func resourceSharedFilesystemShareTypeV2ExtraSpecs(d *schema.ResourceData) (sharetypes.ExtraSpecsOpts, map[string]any, error) {
	raw := d.Get("extra_specs").(map[string]any)

	var opts sharetypes.ExtraSpecsOpts

	dhss, ok := raw["driver_handles_share_servers"]
	if !ok {
		return opts, nil, fmt.Errorf("extra_specs must contain a driver_handles_share_servers key")
	}

	dhssBool, err := strconv.ParseBool(fmt.Sprintf("%v", dhss))
	if err != nil {
		return opts, nil, fmt.Errorf("driver_handles_share_servers must be a boolean value: %s", err)
	}

	// There is an open issue in gophercloud v2 where creating a share type will fail if
	// driver_handles_share_servers is false. Work around this by setting it to true, then
	// updating it to the desired value with SetExtraSpecs.
	opts.DriverHandlesShareServers = true

	if ss, ok := raw["snapshot_support"]; ok {
		ssBool, err := strconv.ParseBool(fmt.Sprintf("%v", ss))
		if err != nil {
			return opts, nil, fmt.Errorf("snapshot_support must be a boolean value: %s", err)
		}

		opts.SnapshotSupport = &ssBool
	}

	remaining := make(map[string]any)

	remaining["driver_handles_share_servers"] = dhssBool

	for k, v := range raw {
		if k == "driver_handles_share_servers" || k == "snapshot_support" {
			continue
		}

		remaining[k] = v
	}

	return opts, remaining, nil
}

func resourceSharedFilesystemShareTypeV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	sfsClient, err := config.SharedfilesystemV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	name := d.Get("name").(string)
	isPublic := d.Get("is_public").(bool)

	extraSpecsOpts, remainingExtraSpecs, err := resourceSharedFilesystemShareTypeV2ExtraSpecs(d)
	if err != nil {
		return diag.FromErr(err)
	}

	createOpts := sharetypes.CreateOpts{
		Name:       name,
		IsPublic:   isPublic,
		ExtraSpecs: extraSpecsOpts,
	}

	log.Printf("[DEBUG] openstack_sharedfilesystem_sharetype_v2 create options: %#v", createOpts)

	st, err := sharetypes.Create(ctx, sfsClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_sharedfilesystem_sharetype_v2 %s: %s", name, err)
	}

	d.SetId(st.ID)

	// The Create call only accepts the driver_handles_share_servers and
	// snapshot_support extra specs. Any other extra specs supplied by the
	// user must be applied afterwards.
	if len(remainingExtraSpecs) > 0 {
		setOpts := sharetypes.SetExtraSpecsOpts{ExtraSpecs: remainingExtraSpecs}

		if _, err := sharetypes.SetExtraSpecs(ctx, sfsClient, st.ID, setOpts).Extract(); err != nil {
			return diag.Errorf("Error setting extra_specs for openstack_sharedfilesystem_sharetype_v2 %s: %s", st.ID, err)
		}
	}

	// The Create call does not accept a description. If one was supplied,
	// apply it with a follow-up Update call.
	if v, ok := d.GetOk("description"); ok {
		sfsClient.Microversion = sharedFilesystemV2ShareTypeMinMicroversion

		description := v.(string)
		updateOpts := sharetypes.UpdateOpts{Description: &description}

		if _, err := sharetypes.Update(ctx, sfsClient, st.ID, updateOpts).Extract(); err != nil {
			return diag.Errorf("Error setting description for openstack_sharedfilesystem_sharetype_v2 %s: %s", st.ID, err)
		}
	}

	return resourceSharedFilesystemShareTypeV2Read(ctx, d, meta)
}

func resourceSharedFilesystemShareTypeV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	sfsClient, err := config.SharedfilesystemV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	// Show share type details requires microversion 2.50 or later.
	sfsClient.Microversion = sharedFilesystemV2ShareTypeMinMicroversion

	st, err := sharetypes.Get(ctx, sfsClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_sharedfilesystem_sharetype_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_sharedfilesystem_sharetype_v2 %s: %#v", d.Id(), st)

	d.Set("name", st.Name)
	d.Set("is_public", st.IsPublic)
	d.Set("is_default", st.IsDefault)
	d.Set("region", GetRegion(d, config))

	if st.Description != nil {
		d.Set("description", *st.Description)
	} else {
		d.Set("description", "")
	}

	es, err := sharetypes.GetExtraSpecs(ctx, sfsClient, d.Id()).Extract()
	if err != nil {
		return diag.Errorf("Error reading extra_specs for openstack_sharedfilesystem_sharetype_v2 %s: %s", d.Id(), err)
	}

	if err := d.Set("extra_specs", es); err != nil {
		log.Printf("[WARN] Unable to set extra_specs for openstack_sharedfilesystem_sharetype_v2 %s: %s", d.Id(), err)
	}

	return nil
}

func resourceSharedFilesystemShareTypeV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	sfsClient, err := config.SharedfilesystemV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	hasChange := false

	var updateOpts sharetypes.UpdateOpts

	if d.HasChange("name") {
		hasChange = true
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}

	if d.HasChange("description") {
		hasChange = true
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}

	if d.HasChange("is_public") {
		hasChange = true
		isPublic := d.Get("is_public").(bool)
		updateOpts.IsPublic = &isPublic
	}

	if hasChange {
		// Update (and Show) share type support requires microversion 2.50
		// or later.
		sfsClient.Microversion = sharedFilesystemV2ShareTypeMinMicroversion

		log.Printf("[DEBUG] openstack_sharedfilesystem_sharetype_v2 %s update options: %#v", d.Id(), updateOpts)

		if _, err := sharetypes.Update(ctx, sfsClient, d.Id(), updateOpts).Extract(); err != nil {
			return diag.Errorf("Error updating openstack_sharedfilesystem_sharetype_v2 %s: %s", d.Id(), err)
		}
	}

	if d.HasChange("extra_specs") {
		oldRaw, newRaw := d.GetChange("extra_specs")
		oldES := oldRaw.(map[string]any)
		newES := newRaw.(map[string]any)

		// Unset extra specs that were removed.
		for oldKey := range oldES {
			if _, ok := newES[oldKey]; !ok {
				if err := sharetypes.UnsetExtraSpecs(ctx, sfsClient, d.Id(), oldKey).ExtractErr(); err != nil {
					return diag.Errorf("Error removing extra_spec %q from openstack_sharedfilesystem_sharetype_v2 %s: %s", oldKey, d.Id(), err)
				}
			}
		}

		// Set extra specs that were added or changed.
		toSet := map[string]any{}

		for newKey, newValue := range newES {
			if oldValue, ok := oldES[newKey]; !ok || oldValue != newValue {
				toSet[newKey] = newValue
			}
		}

		if len(toSet) > 0 {
			setOpts := sharetypes.SetExtraSpecsOpts{ExtraSpecs: toSet}

			if _, err := sharetypes.SetExtraSpecs(ctx, sfsClient, d.Id(), setOpts).Extract(); err != nil {
				return diag.Errorf("Error updating extra_specs for openstack_sharedfilesystem_sharetype_v2 %s: %s", d.Id(), err)
			}
		}
	}

	return resourceSharedFilesystemShareTypeV2Read(ctx, d, meta)
}

func resourceSharedFilesystemShareTypeV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	sfsClient, err := config.SharedfilesystemV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	if err := sharetypes.Delete(ctx, sfsClient, d.Id()).ExtractErr(); err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_sharedfilesystem_sharetype_v2"))
	}

	return nil
}
