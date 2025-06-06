package openstack

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/textproto"
	"strings"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/objectstorage/v1/accounts"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceObjectStorageAccountV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceObjectStorageAccountV1Create,
		ReadContext:   resourceObjectStorageAccountV1Read,
		UpdateContext: resourceObjectStorageAccountV1Update,
		DeleteContext: resourceObjectStorageAccountV1Delete,
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
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"metadata": {
				Type:             schema.TypeMap,
				Optional:         true,
				DiffSuppressFunc: resourceAccountMetadataDiffSuppressFunc,
			},

			// computed
			"headers": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"bytes_used": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"quota_bytes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"container_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"object_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceObjectStorageAccountV1Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	objectStorageClient, err := config.ObjectStorageV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating OpenStack object storage client: %s", err)
	}

	projectID, err := modifyClientEndpoint(ctx, objectStorageClient, d.Get("project_id").(string))
	if err != nil {
		return diag.FromErr(err)
	}

	opts := &accounts.UpdateOpts{
		Metadata: resourceAccountExpandMetadata(d),
	}

	log.Printf("[DEBUG] Create Options for objectstorage_account_v1: %#v", opts)

	_, err = accounts.Update(ctx, objectStorageClient, opts).Extract()
	if err != nil {
		return diag.Errorf("error creating objectstorage_account_v1: %s", err)
	}

	d.SetId(projectID)

	return resourceObjectStorageAccountV1Read(ctx, d, meta)
}

func resourceObjectStorageAccountV1Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	objectStorageClient, err := config.ObjectStorageV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating OpenStack object storage client: %s", err)
	}

	projectID, err := modifyClientEndpoint(ctx, objectStorageClient, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	result := accounts.Get(ctx, objectStorageClient, nil)
	if result.Err != nil {
		return diag.FromErr(CheckDeleted(d, result.Err, "account"))
	}

	h := make(map[string]string)

	err = result.ExtractInto(&h)
	if err != nil {
		return diag.Errorf("error extracting headers for objectstorage_account_v1 '%s': %s", d.Id(), err)
	}

	headers, err := result.Extract()
	if err != nil {
		return diag.Errorf("error extracting headers for objectstorage_account_v1 '%s': %s", d.Id(), err)
	}

	log.Printf("[DEBUG] Retrieved headers for objectstorage_account_v1 '%s': %#v", d.Id(), headers)

	metadata, err := result.ExtractMetadata()
	if err != nil {
		return diag.Errorf("error extracting metadata for objectstorage_account_v1 '%s': %s", d.Id(), err)
	}

	log.Printf("[DEBUG] Retrieved metadata for objectstorage_account_v1 '%s': %#v", d.Id(), metadata)

	d.Set("project_id", projectID)
	d.Set("bytes_used", headers.BytesUsed)

	if headers.QuotaBytes != nil {
		d.Set("quota_bytes", *headers.QuotaBytes)
	} else {
		d.Set("quota_bytes", 0)
	}

	d.Set("container_count", headers.ContainerCount)
	d.Set("object_count", headers.ObjectCount)
	d.Set("metadata", metadata)
	d.Set("headers", h)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceObjectStorageAccountV1Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	objectStorageClient, err := config.ObjectStorageV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating OpenStack object storage client: %s", err)
	}

	_, err = modifyClientEndpoint(ctx, objectStorageClient, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	updateOpts := accounts.UpdateOpts{}
	if d.HasChange("metadata") {
		updateOpts.Metadata, updateOpts.RemoveMetadata = resourceAccountMetadataChange(d)
	}

	if len(updateOpts.Metadata) > 0 || len(updateOpts.RemoveMetadata) > 0 {
		_, err = accounts.Update(ctx, objectStorageClient, updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating objectstorage_account_v1 '%s': %s", d.Id(), err)
		}
	}

	return resourceObjectStorageAccountV1Read(ctx, d, meta)
}

func resourceObjectStorageAccountV1Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	objectStorageClient, err := config.ObjectStorageV1Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating OpenStack object storage client: %s", err)
	}

	_, err = modifyClientEndpoint(ctx, objectStorageClient, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	metadata := resourceAccountExpandMetadata(d)
	opts := &accounts.UpdateOpts{
		RemoveMetadata: make([]string, 0, len(metadata)),
	}

	for k := range metadata {
		if k == "Quota-Bytes" {
			continue
		}

		opts.RemoveMetadata = append(opts.RemoveMetadata, k)
	}

	if len(opts.RemoveMetadata) == 0 {
		return nil
	}

	_, err = accounts.Update(ctx, objectStorageClient, opts).Extract()
	if err != nil {
		return diag.Errorf("error deleting objectstorage_account_v1 '%s': %s", d.Id(), err)
	}

	return nil
}

func modifyClientEndpoint(ctx context.Context, client *gophercloud.ServiceClient, projectID string) (string, error) {
	if projectID == "" {
		v, err := getTokenInfo(ctx, client)
		if err != nil {
			return "", fmt.Errorf("failed to obtain token info: %w", err)
		}

		if v.projectID == "" {
			return "", errors.New("project_id must be provided, when a token has no project scope")
		}

		projectID = v.projectID
		log.Printf("[DEBUG] detected the %s project_id from the token", projectID)
	}

	v := strings.SplitN(client.Endpoint, "/AUTH_", 2)
	if len(v) != 2 {
		return "", errors.New("could not extract project_id from the endpoint")
	}

	if projectID == v[1] {
		log.Printf("[DEBUG] project_id is the same as the one extracted from the endpoint")
	} else {
		ep := v[0] + "/AUTH_" + projectID
		log.Printf("[DEBUG] modifying the endpoint according to the %s project_id: %s", projectID, ep)
		client.Endpoint = ep
	}

	return projectID, nil
}

func resourceAccountExpandMetadata(d *schema.ResourceData) map[string]string {
	m := d.Get("metadata").(map[string]any)
	metadata := make(map[string]string, len(m))

	for k, v := range m {
		metadata[k] = v.(string)
	}

	return metadata
}

func resourceAccountMetadataChange(d *schema.ResourceData) (map[string]string, []string) {
	o, n := d.GetChange("metadata")
	oldMetadata := o.(map[string]any)
	newMetadata := n.(map[string]any)

	removeKeys := make([]string, 0, len(oldMetadata))

	for oldKey := range oldMetadata {
		// we cannot remove the Quota-Bytes key
		if oldKey == "Quota-Bytes" {
			continue
		}

		if _, ok := newMetadata[oldKey]; !ok {
			removeKeys = append(removeKeys, oldKey)
		}
	}

	metadata := make(map[string]string, len(newMetadata))

	for k, v := range newMetadata {
		if v, ok := oldMetadata[k]; ok && v == newMetadata[k] {
			continue
		}

		metadata[k] = v.(string)
	}

	return metadata, removeKeys
}

func resourceAccountMetadataToCanonicalHeader(m map[string]any) map[string]string {
	c := make(map[string]string, len(m))
	for k, v := range m {
		c[textproto.CanonicalMIMEHeaderKey(k)] = v.(string)
	}

	return c
}

func resourceAccountMetadataDiffSuppressFunc(_, _, _ string, d *schema.ResourceData) bool {
	o, n := d.GetChange("metadata")
	metadataOld := resourceAccountMetadataToCanonicalHeader(o.(map[string]any))
	metadataNew := resourceAccountMetadataToCanonicalHeader(n.(map[string]any))

	for k, v := range metadataNew {
		o := metadataOld[k]
		if v == o {
			continue
		}

		return false
	}

	for k, v := range metadataOld {
		n := metadataNew[k]
		if v == n || k == "Quota-Bytes" && n == "" {
			continue
		}

		return false
	}

	return true
}
