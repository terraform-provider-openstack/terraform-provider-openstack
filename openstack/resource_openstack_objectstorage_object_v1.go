package openstack

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mitchellh/go-homedir"
)

func resourceObjectStorageObjectV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceObjectStorageObjectV1Put,
		Read:   resourceObjectStorageObjectV1Read,
		Update: resourceObjectStorageObjectV1Put,
		Delete: resourceObjectStorageObjectV1Delete,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"container_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"content_disposition": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"copy_from"},
			},

			"content_encoding": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"copy_from"},
			},

			"content_type": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"copy_from"},
			},

			"content": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"source", "copy_from", "object_manifest"},
			},

			"copy_from": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"content", "source", "object_manifest"},
			},

			"delete_after": {
				Type:     schema.TypeInt,
				Optional: true,
			},

			"delete_at": {
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				DiffSuppressFunc: suppressEquivilentTimeDiffs,
			},

			"detect_content_type": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			// this attribute is used to trigger resource updates
			// if the file content is changed
			"etag": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"metadata": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
			},

			"object_manifest": {
				Type:          schema.TypeString,
				Computed:      true,
				Optional:      true,
				ConflictsWith: []string{"copy_from", "source", "content"},
			},

			"source": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"content", "copy_from", "object_manifest"},
			},

			// Read Only
			"content_length": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"date": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"last_modified": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"trans_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceObjectStorageObjectV1Put(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	objectStorageClient, err := config.objectStorageV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack object storage client: %s", err)
	}

	name := d.Get("name").(string)
	cn := d.Get("container_name").(string)

	createOpts := &objects.CreateOpts{
		Metadata:         resourceObjectMetadataV1(d),
		TransferEncoding: "chunked",
	}

	var isValid bool
	if v, ok := d.GetOk("source"); ok {
		isValid = true
		source := v.(string)
		path, err := homedir.Expand(source)
		if err != nil {
			return fmt.Errorf("Error expanding homedir in source (%s): %s", source, err)
		}
		file, err := os.Open(path)
		defer file.Close()
		if err != nil {
			return fmt.Errorf("Error opening openstack swift object source (%s): %s", source, err)
		}
		fileinfo, err := file.Stat()
		if err != nil {
			return fmt.Errorf("Error opening openstack swift object source (%s): %s", source, err)
		}

		createOpts.Content = file
		createOpts.ContentLength = fileinfo.Size()
	}

	if v, ok := d.GetOk("content"); ok {
		isValid = true
		content := v.(string)
		createOpts.Content = bytes.NewReader([]byte(content))
		createOpts.ContentLength = int64(len(content))
	}

	if v, ok := d.GetOk("copy_from"); ok {
		isValid = true
		createOpts.CopyFrom = v.(string)
		createOpts.Content = bytes.NewReader([]byte(""))
	}

	if v, ok := d.GetOk("object_manifest"); ok {
		isValid = true
		createOpts.Content = bytes.NewReader([]byte(""))
		createOpts.ObjectManifest = v.(string)
	}

	if !isValid {
		return fmt.Errorf("Must specify \"source\", \"content\", \"copy_from\" or \"object_manifest\" field")
	}

	if v, ok := d.GetOk("content_disposition"); ok {
		createOpts.ContentDisposition = v.(string)
	}

	if v, ok := d.GetOk("content_encoding"); ok {
		createOpts.ContentEncoding = v.(string)
	}

	if v, ok := d.GetOk("content_type"); ok {
		createOpts.ContentType = v.(string)
	}

	if v, ok := d.GetOk("delete_after"); ok {
		createOpts.DeleteAfter = v.(int)
	}

	if v, ok := d.GetOk("delete_at"); ok && v != "" {
		t, err := time.Parse(time.RFC3339, fmt.Sprintf("%s", v))
		if err != nil {
			return fmt.Errorf("Error Parsing Swift Object Lifecycle Expiration Date: %s, %s", err.Error(), v)
		}

		createOpts.DeleteAt = int(t.Unix())
	}

	if v, ok := d.GetOk("detect_content_type"); ok && v.(bool) {
		createOpts.DetectContentType = "true"
	}

	// this attribute is used to trigger resource updates if the file content is changed
	if v, ok := d.GetOk("etag"); ok {
		createOpts.ETag = v.(string)
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	result, err := objects.Create(objectStorageClient, cn, name, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating OpenStack container object: %s", err)
	}
	log.Printf("[INFO] Object %s has been added to container : %s", name, cn)

	d.Set("etag", result.ETag)
	d.Set("content_length", result.ContentLength)
	d.Set("content_type", result.ContentType)
	if result.Date.Unix() > 0 {
		d.Set("date", result.Date.Format(time.RFC3339))
	}
	if result.LastModified.Unix() > 0 {
		d.Set("last_modified", result.LastModified.Format(time.RFC3339))
	}
	d.Set("trans_id", result.TransID)

	// Store the ID now
	d.SetId(fmt.Sprintf("%s/%s", cn, name))

	return resourceObjectStorageObjectV1Read(d, meta)
}

func resourceObjectStorageObjectV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	objectStorageClient, err := config.objectStorageV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack object storage client: %s", err)
	}

	name := d.Get("name").(string)
	cn := d.Get("container_name").(string)

	getOpts := &objects.GetOpts{}

	if v, ok := d.GetOk("tmp_url_sig"); ok {
		getOpts.Signature = v.(string)
	}
	if v, ok := d.GetOk("tmp_url_expires"); ok {
		getOpts.Expires = v.(string)
	}

	log.Printf("[DEBUG] Get Options: %#v", getOpts)
	result, err := objects.Get(objectStorageClient, cn, name, getOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error getting OpenStack container object: %s", err)
	}
	log.Printf("[INFO] Object %s has been added to container : %s", name, cn)

	d.Set("etag", result.ETag)
	d.Set("content_disposition", result.ContentDisposition)
	d.Set("content_encoding", result.ContentEncoding)
	d.Set("content_length", result.ContentLength)
	d.Set("content_type", result.ContentType)
	if result.Date.Unix() > 0 {
		d.Set("date", result.Date.Format(time.RFC3339))
	}
	if result.DeleteAt.Unix() > 0 {
		d.Set("delete_at", result.DeleteAt.Format(time.RFC3339))
	}
	if result.LastModified.Unix() > 0 {
		d.Set("last_modified", result.LastModified.Format(time.RFC3339))
	}
	d.Set("object_manifest", result.ObjectManifest)
	d.Set("trans_id", result.TransID)

	return nil
}

func resourceObjectStorageObjectV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	objectStorageClient, err := config.objectStorageV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack object storage client: %s", err)
	}

	name := d.Get("name").(string)
	cn := d.Get("container_name").(string)
	deleteOpts := &objects.DeleteOpts{}

	_, err = objects.Delete(objectStorageClient, cn, name, deleteOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error getting OpenStack container object: %s", err)
	}
	return nil
}

func resourceObjectMetadataV1(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("metadata").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}
