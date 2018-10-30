package openstack

import (
	"fmt"
	"log"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/containers"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/gophercloud/gophercloud/pagination"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceObjectStorageContainerV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceObjectStorageContainerV1Create,
		Read:   resourceObjectStorageContainerV1Read,
		Update: resourceObjectStorageContainerV1Update,
		Delete: resourceObjectStorageContainerV1Delete,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
			"container_read": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"container_sync_to": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"container_sync_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"container_write": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"content_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"versions_location": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},
			"metadata": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: false,
			},
			"force_destroy": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceObjectStorageContainerV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	objectStorageClient, err := config.objectStorageV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack object storage client: %s", err)
	}

	cn := d.Get("name").(string)

	createOpts := &containers.CreateOpts{
		ContainerRead:    d.Get("container_read").(string),
		ContainerSyncTo:  d.Get("container_sync_to").(string),
		ContainerSyncKey: d.Get("container_sync_key").(string),
		ContainerWrite:   d.Get("container_write").(string),
		ContentType:      d.Get("content_type").(string),
		VersionsLocation: d.Get("versions_location").(string),
		Metadata:         resourceContainerMetadataV2(d),
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	_, err = containers.Create(objectStorageClient, cn, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating OpenStack container: %s", err)
	}
	log.Printf("[INFO] Container ID: %s", cn)

	// Store the ID now
	d.SetId(cn)

	return resourceObjectStorageContainerV1Read(d, meta)
}

func resourceObjectStorageContainerV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	objectStorageClient, err := config.objectStorageV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack object storage client: %s", err)
	}

	result := containers.Get(objectStorageClient, d.Id(), nil)

	if result.Err != nil {
		return CheckDeleted(d, result.Err, "container")
	}

	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceObjectStorageContainerV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	objectStorageClient, err := config.objectStorageV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack object storage client: %s", err)
	}

	updateOpts := containers.UpdateOpts{
		ContainerRead:    d.Get("container_read").(string),
		ContainerSyncTo:  d.Get("container_sync_to").(string),
		ContainerSyncKey: d.Get("container_sync_key").(string),
		ContainerWrite:   d.Get("container_write").(string),
		ContentType:      d.Get("content_type").(string),
		VersionsLocation: d.Get("versions_location").(string),
	}

	if d.HasChange("versions_location") && d.Get("versions_location").(string) == "" {
		updateOpts.RemoveVersionsLocation = "true"
	}

	if d.HasChange("metadata") {
		updateOpts.Metadata = resourceContainerMetadataV2(d)
	}

	_, err = containers.Update(objectStorageClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating OpenStack container: %s", err)
	}

	return resourceObjectStorageContainerV1Read(d, meta)
}

func resourceObjectStorageContainerV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	objectStorageClient, err := config.objectStorageV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack object storage client: %s", err)
	}

	_, err = containers.Delete(objectStorageClient, d.Id()).Extract()
	if err != nil {
		gopherErr, ok := err.(gophercloud.ErrUnexpectedResponseCode)
		if ok && gopherErr.Actual == 409 && d.Get("force_destroy").(bool) {
			// Container may have things. Delete them.
			log.Printf("[DEBUG] Attempting to forceDestroy Openstack container %+v", err)

			container := d.Id()
			opts := &objects.ListOpts{
				Full: false,
			}
			// Retrieve a pager (i.e. a paginated collection)
			pager := objects.List(objectStorageClient, container, opts)
			// Define an anonymous function to be executed on each page's iteration
			err := pager.EachPage(func(page pagination.Page) (bool, error) {

				objectList, err := objects.ExtractNames(page)
				if err != nil {
					return false, fmt.Errorf("Error extracting names from objects from page %+v", err)
				}
				for _, object := range objectList {
					_, err = objects.Delete(objectStorageClient, container, object, objects.DeleteOpts{}).Extract()
					if err != nil {
						return false, fmt.Errorf("Error deleting object from container %+v", err)
					}
				}
				return true, nil
			})
			if err != nil {
				return err
			}
			return resourceObjectStorageContainerV1Delete(d, meta)
		}
		return fmt.Errorf("Error deleting OpenStack container: %s", err)
	}

	d.SetId("")
	return nil
}

func resourceContainerMetadataV2(d *schema.ResourceData) map[string]string {
	m := make(map[string]string)
	for key, val := range d.Get("metadata").(map[string]interface{}) {
		m[key] = val.(string)
	}
	return m
}
