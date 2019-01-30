package openstack

import (
	"fmt"
	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/secrets"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"time"
)

func resourceKeymanagerSecretMetadataV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeymanagerSecretMetadataV1Create,
		Read:   resourceKeymanagerSecretMetadataV1Read,
		Update: resourceKeymanagerSecretMetadataV1Update,
		Delete: resourceKeymanagerSecretMetadataV1Delete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"secret_ref": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Required: true,
			},
		},
	}
}

func resourceKeymanagerSecretMetadataV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	kmClient, err := config.keymanagerV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack keymanager client: %s", err)
	}

	data := d.Get("metadata").(map[string]interface{})

	metadata := make(map[string]string)
	for k, v := range data {
		metadata[k] = v.(string)
	}

	var createOpts secrets.MetadataOpts
	createOpts = metadata

	log.Printf("[DEBUG] Create Options for resource_keymanager_secret_metadata_v1: %#v", createOpts)

	secret_ref := d.Get("secret_ref").(string)
	uuid := keymanagerSecretV1GetUUIDfromSecretRef(secret_ref)
	_, err = secrets.CreateMetadata(kmClient, uuid, createOpts).Extract()

	if err != nil {
		return fmt.Errorf("Error creating metadata for secret with ID %v", uuid)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"NOT_CREATED"},
		Target:     []string{"ACTIVE"},
		Refresh:    keymanagerSecretMetadataV1WaitForSecretMetadataCreation(kmClient, uuid),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	_, err = stateConf.WaitForState()

	if err != nil {
		return fmt.Errorf("Error creating OpenStack barbican secret metadata: %s", err)
	}

	d.SetId(uuid)

	return resourceKeymanagerSecretMetadataV1Read(d, meta)
}

func resourceKeymanagerSecretMetadataV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	kmClient, err := config.keymanagerV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack barbican client: %s", err)
	}

	uuid := d.Id()
	metadata, err := secrets.GetMetadata(kmClient, uuid).Extract()
	if err != nil {
		return CheckDeleted(d, err, "secret")
	}

	log.Printf("[DEBUG] Retrieved secret metadata %s: %+v", uuid, metadata)

	for key, value := range metadata {
		d.Set(key, value)
	}
	return nil
}

func resourceKeymanagerSecretMetadataV1Update(d *schema.ResourceData, meta interface{}) error {
	var hasChange = false
	if d.HasChange("metadata") {
		hasChange = true
	}
	if hasChange {
		// for metadata, creating automatically deletes the old metadata and
		// replaces it with the new metadata.
		resourceKeymanagerSecretMetadataV1Create(d, meta)
	}

	return resourceKeymanagerSecretV1Read(d, meta)
}

func resourceKeymanagerSecretMetadataV1Delete(d *schema.ResourceData, meta interface{}) error {

	d.Set("metadata", "")

	return resourceKeymanagerSecretMetadataV1Create(d, meta)
}
