package openstack

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/quotasets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBlockStorageQuotasetV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBlockStorageQuotasetV3Create,
		ReadContext:   resourceBlockStorageQuotasetV3Read,
		UpdateContext: resourceBlockStorageQuotasetV3Update,
		Delete:        schema.RemoveFromState,
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

			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"volumes": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"snapshots": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"gigabytes": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"per_volume_gigabytes": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"backups": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"backup_gigabytes": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"groups": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"volume_type_quota": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func resourceBlockStorageQuotasetV3Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)
	region := GetRegion(d, config)

	blockStorageClient, err := config.BlockStorageV3Client(ctx, region)
	if err != nil {
		return diag.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	updateOpts := quotasets.UpdateOpts{}
	projectID := d.Get("project_id").(string)

	if v, ok := getOkExists(d, "volumes"); ok {
		value := v.(int)
		updateOpts.Volumes = &value
	}

	if v, ok := getOkExists(d, "snapshots"); ok {
		value := v.(int)
		updateOpts.Snapshots = &value
	}

	if v, ok := getOkExists(d, "gigabytes"); ok {
		value := v.(int)
		updateOpts.Gigabytes = &value
	}

	if v, ok := getOkExists(d, "per_volume_gigabytes"); ok {
		value := v.(int)
		updateOpts.PerVolumeGigabytes = &value
	}

	if v, ok := getOkExists(d, "backups"); ok {
		value := v.(int)
		updateOpts.Backups = &value
	}

	if v, ok := getOkExists(d, "backup_gigabytes"); ok {
		value := v.(int)
		updateOpts.BackupGigabytes = &value
	}

	if v, ok := getOkExists(d, "groups"); ok {
		value := v.(int)
		updateOpts.Groups = &value
	}

	volumeTypeQuotaRaw := d.Get("volume_type_quota").(map[string]any)

	volumeTypeQuota, err := blockStorageQuotasetVolTypeQuotaToInt(volumeTypeQuotaRaw)
	if err != nil {
		return diag.Errorf("Error parsing volume_type_quota in openstack_blockstorage_quotaset_v3: %s", err)
	}

	updateOpts.Extra = volumeTypeQuota

	q, err := quotasets.Update(ctx, blockStorageClient, projectID, updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_blockstorage_quotaset_v3: %s", err)
	}

	id := fmt.Sprintf("%s/%s", projectID, region)
	d.SetId(id)

	log.Printf("[DEBUG] Created openstack_blockstorage_quotaset_v3 %#v", q)

	return resourceBlockStorageQuotasetV3Read(ctx, d, meta)
}

func resourceBlockStorageQuotasetV3Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)
	region := GetRegion(d, config)

	blockStorageClient, err := config.BlockStorageV3Client(ctx, region)
	if err != nil {
		return diag.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	// Depending on the provider version the resource was created, the resource id
	// can be either <project_id> or <project_id>/<region>. This parses the project_id
	// in both cases
	projectID := strings.Split(d.Id(), "/")[0]

	q, err := quotasets.Get(ctx, blockStorageClient, projectID).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_blockstorage_quotaset_v3"))
	}

	log.Printf("[DEBUG] Retrieved openstack_blockstorage_quotaset_v3 %s: %#v", d.Id(), q)

	d.Set("project_id", projectID)
	d.Set("region", region)
	d.Set("volumes", q.Volumes)
	d.Set("snapshots", q.Snapshots)
	d.Set("gigabytes", q.Gigabytes)
	d.Set("per_volume_gigabytes", q.PerVolumeGigabytes)
	d.Set("backups", q.Backups)
	d.Set("backup_gigabytes", q.BackupGigabytes)
	d.Set("groups", q.Groups)

	// We only set volume_type_quota when user is defining them
	volumeTypeQuotaProvided := d.Get("volume_type_quota").(map[string]any)
	if len(volumeTypeQuotaProvided) > 0 {
		volumeTypeQuota, err := blockStorageQuotasetVolTypeQuotaToStr(q.Extra)
		if err != nil {
			log.Printf("[WARN] Unable to read openstack_blockstorage_quotaset_v3 %s volume_type_quotas: %s", d.Id(), err)
		}

		if err := d.Set("volume_type_quota", volumeTypeQuota); err != nil {
			log.Printf("[WARN] Unable to set openstack_blockstorage_quotaset_v3 %s volume_type_quotas: %s", d.Id(), err)
		}
	}

	return nil
}

func resourceBlockStorageQuotasetV3Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	blockStorageClient, err := config.BlockStorageV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	var (
		hasChange  bool
		updateOpts quotasets.UpdateOpts
	)

	if d.HasChange("volumes") {
		hasChange = true
		volumes := d.Get("volumes").(int)
		updateOpts.Volumes = &volumes
	}

	if d.HasChange("snapshots") {
		hasChange = true
		snapshots := d.Get("snapshots").(int)
		updateOpts.Snapshots = &snapshots
	}

	if d.HasChange("gigabytes") {
		hasChange = true
		gigabytes := d.Get("gigabytes").(int)
		updateOpts.Gigabytes = &gigabytes
	}

	if d.HasChange("per_volume_gigabytes") {
		hasChange = true
		perVolumeGigabytes := d.Get("per_volume_gigabytes").(int)
		updateOpts.PerVolumeGigabytes = &perVolumeGigabytes
	}

	if d.HasChange("backups") {
		hasChange = true
		backups := d.Get("backups").(int)
		updateOpts.Backups = &backups
	}

	if d.HasChange("backup_gigabytes") {
		hasChange = true
		backupGigabytes := d.Get("backup_gigabytes").(int)
		updateOpts.BackupGigabytes = &backupGigabytes
	}

	if d.HasChange("groups") {
		hasChange = true
		groups := d.Get("groups").(int)
		updateOpts.Groups = &groups
	}

	if d.HasChange("volume_type_quota") {
		volumeTypeQuotaRaw := d.Get("volume_type_quota").(map[string]any)

		// if len(volumeTypeQuotaRaw) == 0 it can lead to error when trying to do an update with
		// zero attributes. Not updating when a user removes all attributes is acceptable
		// as this attributes are not removed anyways.
		if len(volumeTypeQuotaRaw) > 0 {
			volumeTypeQuota, err := blockStorageQuotasetVolTypeQuotaToInt(volumeTypeQuotaRaw)
			if err != nil {
				return diag.Errorf("Error parsing volume_type_quota in openstack_blockstorage_quotaset_v3: %s", err)
			}

			updateOpts.Extra = volumeTypeQuota
			hasChange = true
		}
	}

	if hasChange {
		log.Printf("[DEBUG] openstack_blockstorage_quotaset_v3 %s update options: %#v", d.Id(), updateOpts)
		projectID := d.Get("project_id").(string)

		_, err = quotasets.Update(ctx, blockStorageClient, projectID, updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_blockstorage_quotaset_v3: %s", err)
		}
	}

	return resourceBlockStorageQuotasetV3Read(ctx, d, meta)
}
