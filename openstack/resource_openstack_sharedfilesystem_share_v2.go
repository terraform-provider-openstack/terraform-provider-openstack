package openstack

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/errors"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/messages"
	"github.com/gophercloud/gophercloud/v2/openstack/sharedfilesystems/v2/shares"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

const (
	// Major share functionality appeared in 2.14.
	minManilaShareMicroversion = "2.14"
)

func resourceSharedFilesystemShareV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSharedFilesystemShareV2Create,
		ReadContext:   resourceSharedFilesystemShareV2Read,
		UpdateContext: resourceSharedFilesystemShareV2Update,
		DeleteContext: resourceSharedFilesystemShareV2Delete,
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
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"share_proto": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"NFS", "CIFS", "CEPHFS", "GLUSTERFS", "HDFS", "MAPRFS",
				}, true),
			},

			"size": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntAtLeast(1),
			},

			"share_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"snapshot_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"is_public": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
			},

			"share_network_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"export_locations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"path": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"preferred": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"has_replicas": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"host": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"replication_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"share_server_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"all_metadata": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func resourceSharedFilesystemShareV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	sfsClient, err := config.SharedfilesystemV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	sfsClient.Microversion = minManilaShareMicroversion

	isPublic := d.Get("is_public").(bool)

	metadataRaw := d.Get("metadata").(map[string]any)
	metadata := make(map[string]string, len(metadataRaw))

	for k, v := range metadataRaw {
		if stringVal, ok := v.(string); ok {
			metadata[k] = stringVal
		}
	}

	createOpts := shares.CreateOpts{
		Name:             d.Get("name").(string),
		Description:      d.Get("description").(string),
		ShareProto:       d.Get("share_proto").(string),
		Size:             d.Get("size").(int),
		SnapshotID:       d.Get("snapshot_id").(string),
		IsPublic:         &isPublic,
		Metadata:         metadata,
		ShareNetworkID:   d.Get("share_network_id").(string),
		AvailabilityZone: d.Get("availability_zone").(string),
	}

	if v, ok := getOkExists(d, "share_type"); ok {
		createOpts.ShareType = v.(string)
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	timeout := d.Timeout(schema.TimeoutCreate)

	log.Printf("[DEBUG] Attempting to create share")

	var share *shares.Share

	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		share, err = shares.Create(ctx, sfsClient, createOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}

		return nil
	})
	if err != nil {
		detailedErr := errors.ErrorDetails{}

		e := errors.ExtractErrorInto(err, &detailedErr)
		if e != nil {
			return diag.Errorf("Error creating share: %s: %s", err, e)
		}

		for k, msg := range detailedErr {
			return diag.Errorf("Error creating share: %s (%d): %s", k, msg.Code, msg.Message)
		}
	}

	d.SetId(share.ID)

	// Wait for share to become active before continuing
	err = waitForSFV2Share(ctx, sfsClient, share.ID, "available", []string{"creating", "manage_starting"}, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceSharedFilesystemShareV2Read(ctx, d, meta)
}

func resourceSharedFilesystemShareV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	sfsClient, err := config.SharedfilesystemV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	sfsClient.Microversion = minManilaShareMicroversion

	share, err := shares.Get(ctx, sfsClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "share"))
	}

	log.Printf("[DEBUG] Retrieved share %s: %#v", d.Id(), share)

	exportLocationsRaw, err := shares.ListExportLocations(ctx, sfsClient, d.Id()).Extract()
	if err != nil {
		return diag.Errorf("Failed to retrieve share's export_locations %s: %s", d.Id(), err)
	}

	log.Printf("[DEBUG] Retrieved share's export_locations %s: %#v", d.Id(), exportLocationsRaw)

	exportLocations := make([]map[string]string, 0, len(exportLocationsRaw))
	for _, v := range exportLocationsRaw {
		exportLocations = append(exportLocations, map[string]string{
			"path":      v.Path,
			"preferred": strconv.FormatBool(v.Preferred),
		})
	}

	if err = d.Set("export_locations", exportLocations); err != nil {
		log.Printf("[DEBUG] Unable to set export_locations: %s", err)
	}

	d.Set("name", share.Name)
	d.Set("description", share.Description)
	d.Set("share_proto", share.ShareProto)
	d.Set("size", share.Size)
	d.Set("share_type", share.ShareTypeName)
	d.Set("snapshot_id", share.SnapshotID)
	d.Set("is_public", share.IsPublic)
	d.Set("all_metadata", share.Metadata)
	d.Set("share_network_id", share.ShareNetworkID)
	d.Set("availability_zone", share.AvailabilityZone)
	// Computed
	d.Set("region", GetRegion(d, config))
	d.Set("project_id", share.ProjectID)
	d.Set("has_replicas", share.HasReplicas)
	d.Set("host", share.Host)
	d.Set("replication_type", share.ReplicationType)
	d.Set("share_server_id", share.ShareServerID)

	return nil
}

func resourceSharedFilesystemShareV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	sfsClient, err := config.SharedfilesystemV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	sfsClient.Microversion = minManilaShareMicroversion

	timeout := d.Timeout(schema.TimeoutUpdate)

	var updateOpts shares.UpdateOpts

	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateOpts.DisplayName = &name
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.DisplayDescription = &description
	}

	if d.HasChange("is_public") {
		isPublic := d.Get("is_public").(bool)
		updateOpts.IsPublic = &isPublic
	}

	if updateOpts != (shares.UpdateOpts{}) {
		// Wait for share to become active before continuing
		err = waitForSFV2Share(ctx, sfsClient, d.Id(), "available", []string{"creating", "manage_starting", "extending", "shrinking"}, timeout)
		if err != nil {
			return diag.FromErr(err)
		}

		log.Printf("[DEBUG] Attempting to update share")

		err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
			_, err := shares.Update(ctx, sfsClient, d.Id(), updateOpts).Extract()
			if err != nil {
				return checkForRetryableError(err)
			}

			return nil
		})
		if err != nil {
			detailedErr := errors.ErrorDetails{}

			e := errors.ExtractErrorInto(err, &detailedErr)
			if e != nil {
				return diag.Errorf("Error updating %s share: %s: %s", d.Id(), err, e)
			}

			for k, msg := range detailedErr {
				return diag.Errorf("Error updating %s share: %s (%d): %s", d.Id(), k, msg.Code, msg.Message)
			}
		}

		// Wait for share to become active before continuing
		err = waitForSFV2Share(ctx, sfsClient, d.Id(), "available", []string{"creating", "manage_starting", "extending", "shrinking"}, timeout)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("size") {
		var pending []string

		oldSize, newSize := d.GetChange("size")

		if newSize.(int) > oldSize.(int) {
			pending = append(pending, "extending")
			resizeOpts := shares.ExtendOpts{NewSize: newSize.(int)}
			log.Printf("[DEBUG] Resizing share %s with options: %#v", d.Id(), resizeOpts)

			err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
				err := shares.Extend(ctx, sfsClient, d.Id(), resizeOpts).Err
				log.Printf("[DEBUG] Resizing share %s with options: %#v", d.Id(), resizeOpts)

				if err != nil {
					return checkForRetryableError(err)
				}

				return nil
			})
		} else if newSize.(int) < oldSize.(int) {
			pending = append(pending, "shrinking")
			resizeOpts := shares.ShrinkOpts{NewSize: newSize.(int)}
			log.Printf("[DEBUG] Resizing share %s with options: %#v", d.Id(), resizeOpts)

			err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
				err := shares.Shrink(ctx, sfsClient, d.Id(), resizeOpts).Err
				log.Printf("[DEBUG] Resizing share %s with options: %#v", d.Id(), resizeOpts)

				if err != nil {
					return checkForRetryableError(err)
				}

				return nil
			})
		}

		if err != nil {
			detailedErr := errors.ErrorDetails{}

			e := errors.ExtractErrorInto(err, &detailedErr)
			if e != nil {
				return diag.Errorf("Unable to resize %s share: %s: %s", d.Id(), err, e)
			}

			for k, msg := range detailedErr {
				return diag.Errorf("Unable to resize %s share: %s (%d): %s", d.Id(), k, msg.Code, msg.Message)
			}
		}

		// Wait for share to become active before continuing
		err = waitForSFV2Share(ctx, sfsClient, d.Id(), "available", pending, timeout)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("metadata") {
		metadataToDelete := make(map[string]string)
		metadataToUpdate := make(map[string]string)

		o, n := d.GetChange("metadata")
		oldMetadata := o.(map[string]any)
		newMetadata := n.(map[string]any)
		existingMetadata := d.Get("all_metadata").(map[string]any)

		// Determine if any metadata keys were removed from the configuration.
		// Then request those keys to be deleted.
		for oldKey, oldValue := range oldMetadata {
			if _, ok := newMetadata[oldKey]; !ok {
				metadataToDelete[oldKey] = oldValue.(string)
			}
		}

		log.Printf("[DEBUG] Deleting the following items from metadata for openstack_sharedfilesystem_share_v2 %s: %v", d.Id(), metadataToDelete)

		for oldKey := range metadataToDelete {
			err := shares.DeleteMetadatum(ctx, sfsClient, d.Id(), oldKey).ExtractErr()
			if err != nil && CheckDeleted(d, err, "") != nil {
				return diag.Errorf("Error deleting openstack_sharedfilesystem_share_v2 %s metadata %s: %s", d.Id(), oldKey, err)
			}
		}

		for newKey, newValue := range newMetadata {
			metadataToUpdate[newKey] = newValue.(string)
		}

		for newKey, newValue := range existingMetadata {
			metadataToUpdate[newKey] = newValue.(string)
		}

		// Remove already removed metadata from the update list
		for oldKey := range metadataToDelete {
			delete(metadataToUpdate, oldKey)
		}

		log.Printf("[DEBUG] Updating the following items in metadata for openstack_sharedfilesystem_share_v2 %s: %v", d.Id(), metadataToUpdate)

		_, err := shares.UpdateMetadata(ctx, sfsClient, d.Id(), shares.UpdateMetadataOpts{Metadata: metadataToUpdate}).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_sharedfilesystem_share_v2 %s metadata: %s", d.Id(), err)
		}
	}

	return resourceSharedFilesystemShareV2Read(ctx, d, meta)
}

func resourceSharedFilesystemShareV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	sfsClient, err := config.SharedfilesystemV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack sharedfilesystem client: %s", err)
	}

	timeout := d.Timeout(schema.TimeoutDelete)

	log.Printf("[DEBUG] Attempting to delete share %s", d.Id())

	err = retry.RetryContext(ctx, timeout, func() *retry.RetryError {
		err = shares.Delete(ctx, sfsClient, d.Id()).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}

		return nil
	})
	if err != nil {
		e := CheckDeleted(d, err, "")
		if e == nil {
			return nil
		}

		detailedErr := errors.ErrorDetails{}

		e = errors.ExtractErrorInto(err, &detailedErr)
		if e != nil {
			return diag.Errorf("Unable to delete %s share: %s: %s", d.Id(), err, e)
		}

		for k, msg := range detailedErr {
			return diag.Errorf("Unable to delete %s share: %s (%d): %s", d.Id(), k, msg.Code, msg.Message)
		}
	}

	// Wait for share to become deleted before continuing
	pending := []string{"", "deleting", "available"}

	err = waitForSFV2Share(ctx, sfsClient, d.Id(), "deleted", pending, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// Full list of the share statuses: https://developer.openstack.org/api-ref/shared-file-system/#shares
func waitForSFV2Share(ctx context.Context, sfsClient *gophercloud.ServiceClient, id string, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for share %s to become %s.", id, target)

	stateConf := &retry.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceSFV2ShareRefreshFunc(ctx, sfsClient, id),
		Timeout:    timeout,
		Delay:      0,
		MinTimeout: 1 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		if gophercloud.ResponseCodeIs(err, http.StatusNotFound) {
			switch target {
			case "deleted":
				return nil
			default:
				return fmt.Errorf("Error: share %s not found: %w", id, err)
			}
		}

		errorMessage := fmt.Sprintf("Error waiting for share %s to become %s", id, target)
		msg := resourceSFSV2ShareManilaMessage(ctx, sfsClient, id)

		if msg == nil {
			return fmt.Errorf("%s: %w", errorMessage, err)
		}

		return fmt.Errorf("%s: %w: the latest manila message (%s): %s", errorMessage, err, msg.CreatedAt, msg.UserMessage)
	}

	return nil
}

func resourceSFV2ShareRefreshFunc(ctx context.Context, sfsClient *gophercloud.ServiceClient, id string) retry.StateRefreshFunc {
	return func() (any, string, error) {
		share, err := shares.Get(ctx, sfsClient, id).Extract()
		if err != nil {
			return nil, "", err
		}

		return share, share.Status, nil
	}
}

func resourceSFSV2ShareManilaMessage(ctx context.Context, sfsClient *gophercloud.ServiceClient, id string) *messages.Message {
	// we can simply set this, because this function is called after the error occurred
	sfsClient.Microversion = "2.37"

	listOpts := messages.ListOpts{
		ResourceID: id,
		SortKey:    "created_at",
		SortDir:    "desc",
		Limit:      1,
	}

	allPages, err := messages.List(sfsClient, listOpts).AllPages(ctx)
	if err != nil {
		log.Printf("[DEBUG] Unable to retrieve messages: %v", err)

		return nil
	}

	allMessages, err := messages.ExtractMessages(allPages)
	if err != nil {
		log.Printf("[DEBUG] Unable to extract messages: %v", err)

		return nil
	}

	if len(allMessages) == 0 {
		log.Printf("[DEBUG] No messages found")

		return nil
	}

	return &allMessages[0]
}
