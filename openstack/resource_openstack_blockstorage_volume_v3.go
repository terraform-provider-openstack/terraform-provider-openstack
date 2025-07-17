package openstack

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/volumeattach"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceBlockStorageVolumeV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBlockStorageVolumeV3Create,
		ReadContext:   resourceBlockStorageVolumeV3Read,
		UpdateContext: resourceBlockStorageVolumeV3Update,
		DeleteContext: resourceBlockStorageVolumeV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceBlockStorageVolumeV3Import,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"size": {
				Type:     schema.TypeInt,
				Required: true,
			},

			"enable_online_resize": {
				Type:     schema.TypeBool,
				Optional: true,
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

			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: false,
				Computed: true,
			},

			"snapshot_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"source_vol_id", "image_id", "backup_id"},
			},

			"source_vol_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"snapshot_id", "image_id", "backup_id"},
			},

			"image_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"snapshot_id", "source_vol_id", "backup_id"},
			},

			"backup_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"snapshot_id", "source_vol_id", "image_id"},
			},

			"volume_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"volume_retype_policy": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "never",
				ValidateFunc: validation.StringInSlice([]string{
					"never", "on-demand",
				}, true),
			},

			"consistency_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"source_replica": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"attachment": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"device": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Set: blockStorageVolumeV3AttachmentHash,
			},

			"scheduler_hints": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"different_host": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"same_host": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"query": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"local_to_instance": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"additional_properties": {
							Type:     schema.TypeMap,
							Optional: true,
							ForceNew: true,
						},
					},
				},
				Set: blockStorageVolumeV3SchedulerHintsHash,
			},
		},
	}
}

func resourceBlockStorageVolumeV3Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	blockStorageClient, err := config.BlockStorageV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	metadata := d.Get("metadata").(map[string]any)
	createOpts := &volumes.CreateOpts{
		AvailabilityZone:   d.Get("availability_zone").(string),
		ConsistencyGroupID: d.Get("consistency_group_id").(string),
		Description:        d.Get("description").(string),
		ImageID:            d.Get("image_id").(string),
		Metadata:           expandToMapStringString(metadata),
		Name:               d.Get("name").(string),
		Size:               d.Get("size").(int),
		SnapshotID:         d.Get("snapshot_id").(string),
		SourceReplica:      d.Get("source_replica").(string),
		SourceVolID:        d.Get("source_vol_id").(string),
		VolumeType:         d.Get("volume_type").(string),
	}

	var schedulerHints volumes.SchedulerHintOpts

	schedulerHintsRaw := d.Get("scheduler_hints").(*schema.Set).List()
	if len(schedulerHintsRaw) > 0 {
		log.Printf("[DEBUG] openstack_blockstorage_volume_v3 scheduler hints: %+v", schedulerHintsRaw[0])
		schedulerHints = resourceBlockStorageVolumeV3SchedulerHints(schedulerHintsRaw[0].(map[string]any))
	}

	if v := d.Get("backup_id").(string); v != "" {
		blockStorageClient.Microversion = blockstorageV3VolumeFromBackupMicroversion
		createOpts.BackupID = v
	}

	log.Printf("[DEBUG] openstack_blockstorage_volume_v3 create options: %#v", createOpts)

	v, err := volumes.Create(ctx, blockStorageClient, createOpts, schedulerHints).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_blockstorage_volume_v3: %s", err)
	}

	d.SetId(v.ID)

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"downloading", "creating"},
		Target:     []string{"available"},
		Refresh:    blockStorageVolumeV3StateRefreshFunc(ctx, blockStorageClient, v.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"Error waiting for openstack_blockstorage_volume_v3 %s to become ready: %s", v.ID, err)
	}

	return resourceBlockStorageVolumeV3Read(ctx, d, meta)
}

func resourceBlockStorageVolumeV3Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	blockStorageClient, err := config.BlockStorageV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	v, err := volumes.Get(ctx, blockStorageClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_blockstorage_volume_v3"))
	}

	log.Printf("[DEBUG] Retrieved openstack_blockstorage_volume_v3 %s: %#v", d.Id(), v)

	d.Set("size", v.Size)
	d.Set("description", v.Description)
	d.Set("availability_zone", v.AvailabilityZone)
	d.Set("name", v.Name)
	d.Set("snapshot_id", v.SnapshotID)
	d.Set("backup_id", v.BackupID)
	d.Set("source_vol_id", v.SourceVolID)
	d.Set("volume_type", v.VolumeType)
	d.Set("metadata", v.Metadata)
	d.Set("region", GetRegion(d, config))

	if _, exists := d.GetOk("volume_retype_policy"); !exists {
		d.Set("volume_retype_policy", "never")
	}

	attachments := flattenBlockStorageVolumeV3Attachments(v.Attachments)
	log.Printf("[DEBUG] openstack_blockstorage_volume_v3 %s attachments: %#v", d.Id(), attachments)

	if err := d.Set("attachment", attachments); err != nil {
		log.Printf(
			"[DEBUG] unable to set openstack_blockstorage_volume_v3 %s attachments: %s", d.Id(), err)
	}

	return nil
}

func resourceBlockStorageVolumeV3Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	blockStorageClient, err := config.BlockStorageV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	updateOpts := volumes.UpdateOpts{
		Name:        &name,
		Description: &description,
	}

	if d.HasChange("metadata") {
		metadata := d.Get("metadata").(map[string]any)
		updateOpts.Metadata = expandToMapStringString(metadata)
	}

	var v *volumes.Volume
	if d.HasChange("size") {
		v, err = volumes.Get(ctx, blockStorageClient, d.Id()).Extract()
		if err != nil {
			return diag.Errorf("Error extending openstack_blockstorage_volume_v3 %s: %s", d.Id(), err)
		}

		if v.Status == "in-use" {
			if v, ok := d.Get("enable_online_resize").(bool); ok && !v {
				return diag.Errorf(
					`Error extending openstack_blockstorage_volume_v3 %s,
					volume is attached to the instance and
					resizing online is disabled,
					see enable_online_resize option`, d.Id())
			}

			blockStorageClient.Microversion = blockstorageV3ResizeOnlineInUse
		}

		extendOpts := volumes.ExtendSizeOpts{
			NewSize: d.Get("size").(int),
		}

		err = volumes.ExtendSize(ctx, blockStorageClient, d.Id(), extendOpts).ExtractErr()
		if err != nil {
			return diag.Errorf("Error extending openstack_blockstorage_volume_v3 %s size: %s", d.Id(), err)
		}

		stateConf := &retry.StateChangeConf{
			Pending:    []string{"extending"},
			Target:     []string{"available", "in-use"},
			Refresh:    blockStorageVolumeV3StateRefreshFunc(ctx, blockStorageClient, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      0,
			MinTimeout: 3 * time.Second,
		}

		_, err := stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf(
				"Error waiting for openstack_blockstorage_volume_v3 %s to become ready: %s", d.Id(), err)
		}
	}

	if d.HasChange("volume_type") {
		_, err = volumes.Get(ctx, blockStorageClient, d.Id()).Extract()
		if err != nil {
			return diag.Errorf("Error changing volume type openstack_blockstorage_volume_v3 %s: %s", d.Id(), err)
		}

		retypeOptions := &volumes.ChangeTypeOpts{
			NewType:         d.Get("volume_type").(string),
			MigrationPolicy: volumes.MigrationPolicy(d.Get("volume_retype_policy").(string)),
		}

		err := volumes.ChangeType(ctx, blockStorageClient, d.Id(), retypeOptions).ExtractErr()
		if err != nil {
			return diag.Errorf("Error changing volume %s type: %s", d.Id(), err)
		}

		stateConf := &retry.StateChangeConf{
			Pending:    []string{"retyping"},
			Target:     []string{"available", "in-use"},
			Refresh:    blockStorageVolumeV3StateRefreshFunc(ctx, blockStorageClient, d.Id()),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      0,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf(
				"Error waiting for openstack_blockstorage_volume_v3 %s to become ready: %s", d.Id(), err)
		}
	}

	_, err = volumes.Update(ctx, blockStorageClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error updating openstack_blockstorage_volume_v3 %s: %s", d.Id(), err)
	}

	return resourceBlockStorageVolumeV3Read(ctx, d, meta)
}

func resourceBlockStorageVolumeV3Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	blockStorageClient, err := config.BlockStorageV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	v, err := volumes.Get(ctx, blockStorageClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_blockstorage_volume_v3"))
	}

	// make sure this volume is detached from all instances before deleting
	if len(v.Attachments) > 0 {
		computeClient, err := config.ComputeV2Client(ctx, GetRegion(d, config))
		if err != nil {
			return diag.Errorf("Error creating OpenStack compute client: %s", err)
		}

		for _, volumeAttachment := range v.Attachments {
			log.Printf("[DEBUG] openstack_blockstorage_volume_v3 %s attachment: %#v", d.Id(), volumeAttachment)

			serverID := volumeAttachment.ServerID
			attachmentID := volumeAttachment.ID

			if err := volumeattach.Delete(ctx, computeClient, serverID, attachmentID).ExtractErr(); err != nil {
				// It's possible the volume was already detached by
				// openstack_compute_volume_attach_v2, so consider
				// a 404 acceptable and continue.
				if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
					continue
				}

				// A 409 is also acceptable because there's another
				// concurrent action happening.
				if gophercloud.ResponseCodeIs(err, http.StatusConflict) {
					continue
				}

				return diag.Errorf(
					"Error detaching openstack_blockstorage_volume_v3 %s from %s: %s", d.Id(), serverID, err)
			}
		}

		stateConf := &retry.StateChangeConf{
			Pending:    []string{"in-use", "attaching", "detaching"},
			Target:     []string{"available", "deleted"},
			Refresh:    blockStorageVolumeV3StateRefreshFunc(ctx, blockStorageClient, d.Id()),
			Timeout:    10 * time.Minute,
			Delay:      0,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf(
				"Error waiting for openstack_blockstorage_volume_v3 %s to become available: %s", d.Id(), err)
		}
	}

	// It's possible that this volume was used as a boot device and is currently
	// in a "deleting" state from when the instance was terminated.
	// If this is true, just move on. It'll eventually delete.
	if v.Status != "deleting" {
		if err := volumes.Delete(ctx, blockStorageClient, d.Id(), nil).ExtractErr(); err != nil {
			return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_blockstorage_volume_v3"))
		}
	}

	stateConf := &retry.StateChangeConf{
		Pending:    []string{"deleting", "downloading", "available"},
		Target:     []string{"deleted"},
		Refresh:    blockStorageVolumeV3StateRefreshFunc(ctx, blockStorageClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      0,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("Error waiting for openstack_blockstorage_volume_v3 %s to Delete:  %s", d.Id(), err)
	}

	return nil
}

func resourceBlockStorageVolumeV3Import(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	blockStorageClient, err := config.BlockStorageV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return nil, fmt.Errorf("error creating OpenStack block storage client: %w", err)
	}

	v, err := volumes.Get(ctx, blockStorageClient, d.Id()).Extract()
	if err != nil {
		return nil, fmt.Errorf("error retrieving openstack_blockstorage_volume_v3 %s: %w", d.Id(), err)
	}

	log.Printf("[DEBUG] Retrieved openstack_blockstorage_volume_v3 %s: %#v", d.Id(), v)

	// As we can have a snapshot_id and a image_id when the VolumeImageMetadata is stored in the snapshot for example,
	// we want to exclude image_id when one of the others ids is used
	_, snapshotExists := d.GetOk("snapshot_id")
	_, backupExists := d.GetOk("backup_id")
	_, sourcevolExists := d.GetOk("source_vol_id")

	if v.VolumeImageMetadata != nil && !snapshotExists && !backupExists && !sourcevolExists {
		d.Set("image_id", v.VolumeImageMetadata["image_id"])
	}

	return []*schema.ResourceData{d}, nil
}
