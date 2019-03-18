package openstack

import (
	"fmt"
	"log"
	"time"

	"github.com/gophercloud/gophercloud/openstack/keymanager/v1/secrets"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceKeymanagerSecretV1() *schema.Resource {
	return &schema.Resource{
		Create: resourceKeymanagerSecretV1Create,
		Read:   resourceKeymanagerSecretV1Read,
		Update: resourceKeymanagerSecretV1Update,
		Delete: resourceKeymanagerSecretV1Delete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
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
				ForceNew: true,
			},

			"bit_length": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"creator_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"mode": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"secret_ref": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"secret_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"symmetric", "public", "private", "passphrase", "certificate", "opaque",
				}, true),
				ForceNew: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"payload": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				ForceNew:  false,
			},

			"payload_content_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"payload_content_encoding": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"expiration": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"content_types": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: false,
			},

			"all_metadata": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func resourceKeymanagerSecretV1Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	kmClient, err := config.keymanagerV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack keymanager client: %s", err)
	}

	secretType := keymanagerSecretV1SecretType(d.Get("secret_type").(string))

	createOpts := secrets.CreateOpts{
		Name:                   d.Get("name").(string),
		Algorithm:              d.Get("algorithm").(string),
		BitLength:              d.Get("bit_length").(int),
		Mode:                   d.Get("mode").(string),
		PayloadContentType:     d.Get("payload_content_type").(string),
		PayloadContentEncoding: d.Get("payload_content_encoding").(string),
		SecretType:             secretType,
	}

	log.Printf("[DEBUG] Create Options for resource_keymanager_secret_v1: %#v", createOpts)

	//Add payload here so it does not get printed in the above log
	createOpts.Payload = d.Get("payload").(string)

	var secret *secrets.Secret
	secret, err = secrets.Create(kmClient, createOpts).Extract()

	if err != nil {
		return fmt.Errorf("Error creating openstack_keymanager_secret_v1: %s", err)
	}

	uuid := keymanagerSecretV1GetUUIDfromSecretRef(secret.SecretRef)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"NOT_CREATED"},
		Target:     []string{"ACTIVE"},
		Refresh:    keymanagerSecretV1WaitForSecretCreation(kmClient, uuid),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	_, err = stateConf.WaitForState()

	if err != nil {
		return CheckDeleted(d, err, "Error creating openstack_keymanager_secret_v1")
	}

	d.SetId(uuid)
	d.Partial(true)

	var metadataCreateOpts secrets.MetadataOpts
	metadataCreateOpts = flattenKeyManagerSecretMetadataV1(d)

	log.Printf("[DEBUG] Metadata Create Options for resource_keymanager_secret_metadata_v1 %s: %#v", uuid, metadataCreateOpts)

	if len(metadataCreateOpts) > 0 {
		_, err = secrets.CreateMetadata(kmClient, uuid, metadataCreateOpts).Extract()

		if err != nil {
			return fmt.Errorf("Error creating metadata for openstack_keymanager_secret_v1 with ID %v", uuid)
		}

		stateConf = &resource.StateChangeConf{
			Pending:    []string{"NOT_CREATED"},
			Target:     []string{"ACTIVE"},
			Refresh:    keymanagerSecretMetadataV1WaitForSecretMetadataCreation(kmClient, uuid),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      0,
			MinTimeout: 2 * time.Second,
		}

		_, err = stateConf.WaitForState()

		if err != nil {
			return fmt.Errorf("Error creating metadata for openstack_keymanager_secret_v1: %s: %s", uuid, err)
		}
	}

	d.Partial(false)

	return resourceKeymanagerSecretV1Read(d, meta)
}

func resourceKeymanagerSecretV1Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	kmClient, err := config.keymanagerV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack barbican client: %s", err)
	}

	secret, err := secrets.Get(kmClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Error retrieving openstack_keymanager_secret_v1")
	}

	log.Printf("[DEBUG] Retrieved openstack_keymanager_secret_v1 %s: %#v", d.Id(), secret)

	d.Set("name", secret.Name)

	d.Set("bit_length", secret.BitLength)
	d.Set("algorithm", secret.Algorithm)
	d.Set("creator_id", secret.CreatorID)
	d.Set("mode", secret.Mode)
	d.Set("secret_ref", secret.SecretRef)
	d.Set("secret_type", secret.SecretType)
	d.Set("status", secret.Status)
	d.Set("created_at", secret.Created.Format(time.RFC3339))
	d.Set("updated_at", secret.Updated.Format(time.RFC3339))
	d.Set("expiration", secret.Expiration.Format(time.RFC3339))
	d.Set("content_types", secret.ContentTypes)

	metadataMap, err := secrets.GetMetadata(kmClient, d.Id()).Extract()
	if err != nil {
		log.Printf("[DEBUG] Unable to set metadata: %s", err)
	}
	d.Set("all_metadata", metadataMap)

	// Set the region
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceKeymanagerSecretV1Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	kmClient, err := config.keymanagerV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack barbican client: %s", err)
	}

	var hasChange = false
	var updateOpts secrets.UpdateOpts
	if d.HasChange("payload_content_type") {
		hasChange = true
	}
	// This is not optional so we have to set it regardless
	updateOpts.ContentType = d.Get("payload_content_type").(string)

	if d.HasChange("payload_content_encoding") {
		hasChange = true
		updateOpts.ContentEncoding = d.Get("content_encoding").(string)
	}

	// Print the update options before we set the payload
	log.Printf("[DEBUG] Update Options for resource_keymanager_secret_v1: %#v", updateOpts)

	if d.HasChange("payload") {
		hasChange = true
		updateOpts.Payload = d.Get("payload").(string)
	}

	if hasChange {
		err := secrets.Update(kmClient, d.Id(), updateOpts).Err
		if err != nil {
			return err
		}
	}

	if d.HasChange("metadata") {
		var metadataToDelete []string
		var metadataToAdd []string
		var metadataToUpdate []string

		o, n := d.GetChange("metadata")
		oldMetadata := o.(map[string]interface{})
		newMetadata := n.(map[string]interface{})

		// Determine if any metadata keys were removed from the configuration.
		// Then request those keys to be deleted.
		for oldKey := range oldMetadata {
			if _, ok := newMetadata[oldKey]; !ok {
				metadataToDelete = append(metadataToDelete, oldKey)
			}
		}

		log.Printf("[DEBUG] Deleting the following items from metadata for openstack_keymanager_secret_v1 %s: %v", d.Id(), metadataToDelete)

		for _, key := range metadataToDelete {
			err := secrets.DeleteMetadatum(kmClient, d.Id(), key).ExtractErr()
			if err != nil {
				return fmt.Errorf("Error deleting openstack_keymanager_secret_v1 %s metadata %s: %s", d.Id(), key, err)
			}
		}

		// Determine if any metadata keys were updated or added in the configuration.
		// Then request those keys to be updated or added.
		for newKey, newValue := range newMetadata {
			if oldValue, ok := oldMetadata[newKey]; ok {
				if newValue != oldValue {
					metadataToUpdate = append(metadataToUpdate, newKey)
				}
			} else {
				metadataToAdd = append(metadataToAdd, newKey)
			}
		}

		log.Printf("[DEBUG] Updating the following items in metadata for openstack_keymanager_secret_v1 %s: %v", d.Id(), metadataToUpdate)

		for _, key := range metadataToUpdate {
			var metadatumOpts secrets.MetadatumOpts
			metadatumOpts.Key = key
			metadatumOpts.Value = newMetadata[key].(string)
			_, err := secrets.UpdateMetadatum(kmClient, d.Id(), metadatumOpts).Extract()
			if err != nil {
				return fmt.Errorf("Error updating openstack_keymanager_secret_v1 %s metadata %s: %s", d.Id(), key, err)
			}
		}

		log.Printf("[DEBUG] Adding the following items to metadata for openstack_keymanager_secret_v1 %s: %v", d.Id(), metadataToAdd)

		for _, key := range metadataToAdd {
			var metadatumOpts secrets.MetadatumOpts
			metadatumOpts.Key = key
			metadatumOpts.Value = newMetadata[key].(string)
			err := secrets.CreateMetadatum(kmClient, d.Id(), metadatumOpts).Err
			if err != nil {
				return fmt.Errorf("Error adding openstack_keymanager_secret_v1 %s metadata %s: %s", d.Id(), key, err)
			}
		}
	}

	return resourceKeymanagerSecretV1Read(d, meta)
}

func resourceKeymanagerSecretV1Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	kmClient, err := config.keymanagerV1Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack barbican client: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"ACTIVE"},
		Target:     []string{"DELETED"},
		Refresh:    keymanagerSecretV1WaitForSecretDeletion(kmClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 2 * time.Second,
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return err
	}

	return nil
}
