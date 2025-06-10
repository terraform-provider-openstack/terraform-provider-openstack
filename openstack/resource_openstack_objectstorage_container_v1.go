package openstack

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/objectstorage/v1/containers"
	"github.com/gophercloud/gophercloud/v2/openstack/objectstorage/v1/objects"
	"github.com/gophercloud/gophercloud/v2/pagination"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceObjectStorageContainerV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceObjectStorageContainerV1Create,
		ReadContext:   resourceObjectStorageContainerV1Read,
		UpdateContext: resourceObjectStorageContainerV1Update,
		DeleteContext: resourceObjectStorageContainerV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		SchemaVersion: 1,
		StateUpgraders: []schema.StateUpgrader{
			{
				Type:    resourceObjectStorageContainerV1V0().CoreConfigSchema().ImpliedType(),
				Upgrade: resourceObjectStorageContainerStateUpgradeV0,
				Version: 0,
			},
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
				ForceNew: false,
			},
			"container_read": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"container_sync_to": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"container_sync_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"container_write": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"versioning": {
				Type:          schema.TypeBool,
				Optional:      true,
				Default:       false,
				ConflictsWith: []string{"versioning_legacy"},
			},
			"versioning_legacy": {
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"versions", "history",
							}, true),
						},
						"location": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				ConflictsWith: []string{"versioning"},
				// This method is not actually deprecated on Openstack layer. They
				// are strongly advising to use the method through `versioning`.
				// Deprecation notice is to drive users to use it as well but might
				// not be removed if it's not removed from upstream Openstack first
				Deprecated: "Use newer \"versioning\" implementation",
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: false,
			},
			"force_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"storage_policy": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"storage_class": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceObjectStorageContainerV1Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	objectStorageClient, err := config.ObjectStorageV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating OpenStack object storage client: %s", err)
	}

	cn := d.Get("name").(string)

	createOpts := &containerCreateOpts{
		CreateOpts: containers.CreateOpts{
			ContainerRead:    d.Get("container_read").(string),
			ContainerSyncTo:  d.Get("container_sync_to").(string),
			ContainerSyncKey: d.Get("container_sync_key").(string),
			ContainerWrite:   d.Get("container_write").(string),
			ContentType:      d.Get("content_type").(string),
			StoragePolicy:    d.Get("storage_policy").(string),
			VersionsEnabled:  d.Get("versioning").(bool),
			Metadata:         resourceContainerMetadataV2(d),
		},
		StorageClass: d.Get("storage_class").(string),
	}

	versioning := d.Get("versioning_legacy").(*schema.Set)
	if versioning.Len() > 0 {
		vParams := versioning.List()[0]
		if vRaw, ok := vParams.(map[string]any); ok {
			switch vRaw["type"].(string) {
			case "versions":
				createOpts.VersionsLocation = vRaw["location"].(string)
			case "history":
				createOpts.HistoryLocation = vRaw["location"].(string)
			}
		}
	}

	log.Printf("[DEBUG] Create Options for objectstorage_container_v1: %#v", createOpts)

	_, err = containers.Create(ctx, objectStorageClient, cn, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating objectstorage_container_v1: %s", err)
	}

	log.Printf("[INFO] objectstorage_container_v1 created with ID: %s", cn)

	// Store the ID now
	d.SetId(cn)

	return resourceObjectStorageContainerV1Read(ctx, d, meta)
}

func resourceObjectStorageContainerV1Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	objectStorageClient, err := config.ObjectStorageV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating OpenStack object storage client: %s", err)
	}

	result := containers.Get(ctx, objectStorageClient, d.Id(), nil)
	if result.Err != nil {
		return diag.FromErr(CheckDeleted(d, result.Err, "container"))
	}

	headers, err := result.Extract()
	if err != nil {
		return diag.Errorf("error extracting headers for objectstorage_container_v1 '%s': %s", d.Id(), err)
	}

	log.Printf("[DEBUG] Retrieved headers for objectstorage_container_v1 '%s': %#v", d.Id(), headers)

	metadata, err := result.ExtractMetadata()
	if err != nil {
		return diag.Errorf("error extracting metadata for objectstorage_container_v1 '%s': %s", d.Id(), err)
	}

	log.Printf("[DEBUG] Retrieved metadata for objectstorage_container_v1 '%s': %#v", d.Id(), metadata)

	d.Set("name", d.Id())

	if len(headers.Read) > 0 && headers.Read[0] != "" {
		d.Set("container_read", strings.Join(headers.Read, ","))
	}

	if len(headers.Write) > 0 && headers.Write[0] != "" {
		d.Set("container_write", strings.Join(headers.Write, ","))
	}

	if len(headers.StoragePolicy) > 0 {
		d.Set("storage_policy", headers.StoragePolicy)
	}

	versioningResource := resourceObjectStorageContainerV1().Schema["versioning_legacy"].Elem.(*schema.Resource)

	if headers.VersionsLocation != "" && headers.HistoryLocation != "" {
		return diag.Errorf("error reading versioning headers for objectstorage_container_v1 '%s': found location for both exclusive types, versions ('%s') and history ('%s')", d.Id(), headers.VersionsLocation, headers.HistoryLocation)
	}

	if headers.VersionsLocation != "" {
		versioning := map[string]any{
			"type":     "versions",
			"location": headers.VersionsLocation,
		}
		if err := d.Set("versioning_legacy", schema.NewSet(schema.HashResource(versioningResource), []any{versioning})); err != nil {
			return diag.Errorf("error setting 'versions' versioning for objectstorage_container_v1 '%s': %s", d.Id(), err)
		}
	}

	if headers.HistoryLocation != "" {
		versioning := map[string]any{
			"type":     "history",
			"location": headers.HistoryLocation,
		}
		if err := d.Set("versioning_legacy", schema.NewSet(schema.HashResource(versioningResource), []any{versioning})); err != nil {
			return diag.Errorf("error setting 'history' versioning for objectstorage_container_v1 '%s': %s", d.Id(), err)
		}
	}

	// Despite the create request "X-Object-Storage-Class" header, the
	// response header is "X-Storage-Class".
	d.Set("storage_class", result.Header.Get("X-Storage-Class"))

	d.Set("versioning", headers.VersionsEnabled)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceObjectStorageContainerV1Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	objectStorageClient, err := config.ObjectStorageV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating OpenStack object storage client: %s", err)
	}

	containerRead := d.Get("container_read").(string)
	containerSyncTo := d.Get("container_sync_to").(string)
	containerSyncKey := d.Get("container_sync_key").(string)
	containerWrite := d.Get("container_write").(string)
	contentType := d.Get("content_type").(string)

	updateOpts := containers.UpdateOpts{
		ContainerRead:    &containerRead,
		ContainerSyncTo:  &containerSyncTo,
		ContainerSyncKey: &containerSyncKey,
		ContainerWrite:   &containerWrite,
		ContentType:      &contentType,
	}

	if d.HasChange("versioning") {
		versioning := d.Get("versioning").(bool)
		updateOpts.VersionsEnabled = &versioning
	}

	if d.HasChange("versioning_legacy") {
		versioning := d.Get("versioning_legacy").(*schema.Set)
		if versioning.Len() == 0 {
			updateOpts.RemoveVersionsLocation = "true"
			updateOpts.RemoveHistoryLocation = "true"
		} else {
			vParams := versioning.List()[0]
			if vRaw, ok := vParams.(map[string]any); ok {
				if len(vRaw["location"].(string)) == 0 || len(vRaw["type"].(string)) == 0 {
					updateOpts.RemoveVersionsLocation = "true"
					updateOpts.RemoveHistoryLocation = "true"
				}

				switch vRaw["type"].(string) {
				case "versions":
					updateOpts.VersionsLocation = vRaw["location"].(string)
				case "history":
					updateOpts.HistoryLocation = vRaw["location"].(string)
				}
			}
		}
	}

	// remove legacy versioning first, before enabling the new versioning
	if updateOpts.VersionsEnabled != nil && *updateOpts.VersionsEnabled &&
		(updateOpts.RemoveVersionsLocation == "true" || updateOpts.RemoveHistoryLocation == "true") {
		opts := containers.UpdateOpts{
			RemoveVersionsLocation: "true",
			RemoveHistoryLocation:  "true",
		}

		_, err = containers.Update(ctx, objectStorageClient, d.Id(), opts).Extract()
		if err != nil {
			return diag.Errorf("error updating objectstorage_container_v1 '%s': %s", d.Id(), err)
		}
	}

	// remove new versioning first, before enabling the legacy versioning
	if (updateOpts.VersionsLocation != "" || updateOpts.HistoryLocation != "") &&
		updateOpts.VersionsEnabled != nil && !*updateOpts.VersionsEnabled {
		opts := containers.UpdateOpts{
			VersionsEnabled: updateOpts.VersionsEnabled,
		}

		_, err = containers.Update(ctx, objectStorageClient, d.Id(), opts).Extract()
		if err != nil {
			return diag.Errorf("error updating objectstorage_container_v1 '%s': %s", d.Id(), err)
		}
	}

	if d.HasChange("metadata") {
		updateOpts.Metadata = resourceContainerMetadataV2(d)
	}

	_, err = containers.Update(ctx, objectStorageClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("error updating objectstorage_container_v1 '%s': %s", d.Id(), err)
	}

	return resourceObjectStorageContainerV1Read(ctx, d, meta)
}

func resourceObjectStorageContainerV1Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	objectStorageClient, err := config.ObjectStorageV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating OpenStack object storage client: %s", err)
	}

	_, err = containers.Delete(ctx, objectStorageClient, d.Id()).Extract()
	if err != nil {
		if gophercloud.ResponseCodeIs(err, http.StatusConflict) && d.Get("force_destroy").(bool) {
			// Container may have things. Delete them.
			log.Printf("[DEBUG] Attempting to forceDestroy objectstorage_container_v1 '%s': %+v", d.Id(), err)

			container := d.Id()
			opts := &objects.ListOpts{
				Versions: true,
			}
			// Retrieve a pager (i.e. a paginated collection)
			pager := objects.List(objectStorageClient, container, opts)
			// Define an anonymous function to be executed on each page's iteration
			err := pager.EachPage(ctx, func(_ context.Context, page pagination.Page) (bool, error) {
				objectList, err := objects.ExtractInfo(page)
				if err != nil {
					return false, fmt.Errorf("error extracting names from objects from page for objectstorage_container_v1 '%s': %+w", container, err)
				}

				for _, object := range objectList {
					opts := objects.DeleteOpts{
						ObjectVersionID: object.VersionID,
					}

					_, err = objects.Delete(ctx, objectStorageClient, container, object.Name, opts).Extract()
					if err != nil {
						latest := "latest"
						if !object.IsLatest && object.VersionID != "" {
							latest = object.VersionID
						}

						return false, fmt.Errorf("error deleting object '%s@%s' from objectstorage_container_v1 '%s': %+w", object.Name, latest, container, err)
					}
				}

				return true, nil
			})
			if err != nil {
				return diag.FromErr(err)
			}

			return resourceObjectStorageContainerV1Delete(ctx, d, meta)
		}

		return diag.FromErr(CheckDeleted(d, err, fmt.Sprintf("error deleting objectstorage_container_v1 '%s'", d.Id())))
	}

	d.SetId("")

	return nil
}

func resourceContainerMetadataV2(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("metadata").(map[string]any) {
		m[key] = val.(string)
	}

	return m
}
