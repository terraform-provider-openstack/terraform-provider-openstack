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
				ForceNew: false,
			},
			"bit_length": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"algorithm": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"creator_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
				Computed: true,
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
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"payload": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"payload_content_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"payload_content_encoding": {
				Type:     schema.TypeString,
				Optional: true,
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

	var createOpts secrets.CreateOptsBuilder

	secretType := keymanagerSecretV1SecretType(d.Get("secret_type").(string))

	createOpts = &secrets.CreateOpts{
		Name:                   d.Get("name").(string),
		Algorithm:              d.Get("algorithm").(string),
		BitLength:              d.Get("bit_length").(int),
		Mode:                   d.Get("mode").(string),
		Payload:                d.Get("payload").(string),
		PayloadContentType:     d.Get("payload_content_type").(string),
		PayloadContentEncoding: d.Get("payload_content_encoding").(string),
		SecretType:             secretType,
	}

	log.Printf("[DEBUG] Create Options for resource_keymanager_secret_v1: %#v", createOpts)

	var secret *secrets.Secret
	secret, err = secrets.Create(kmClient, createOpts).Extract()

	if err != nil {
		return fmt.Errorf("Error creating OpenStack barbican secret: %s", err)
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

	var metadataCreateOpts secrets.MetadataOpts
	metadataCreateOpts = keymanagerSecretMetadataV1(d)

	log.Printf("[DEBUG] Create Options for resource_keymanager_secret_metadata_v1: %#v", createOpts)

	_, err = secrets.CreateMetadata(kmClient, uuid, metadataCreateOpts).Extract()

	if err != nil {
		return fmt.Errorf("Error creating metadata for secret with ID %v", uuid)
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
		return fmt.Errorf("Error creating OpenStack barbican secret metadata: %s", err)
	}

	d.SetId(uuid)

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
		return CheckDeleted(d, err, "secret")
	}

	log.Printf("[DEBUG] Retrieved openstack_keymanager_secret_v1 with id %s: %+v", d.Id(), secret)

	d.Set("name", secret.Name)

	d.Set("bit_length", secret.BitLength)
	d.Set("algorithm", secret.Algorithm)
	d.Set("creator_id", secret.CreatorID)
	d.Set("mode", secret.Mode)
	d.Set("secret_ref", secret.SecretRef)
	d.Set("secret_type", secret.SecretType)
	d.Set("status", secret.Status)
	d.Set("created", secret.Created.Format(time.RFC3339))
	d.Set("updated", secret.Updated.Format(time.RFC3339))
	d.Set("expiration", secret.Expiration.Format(time.RFC3339))
	d.Set("content_types", secret.ContentTypes)

	metadataMap, err := secrets.GetMetadata(kmClient, d.Id()).Extract()
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
		oldMetadata, newMetadata := d.GetChange("metadata")
		var metadataToDelete []string

		// Determine if any metadata keys were removed from the configuration.
		// Then request those keys to be deleted.
		for oldKey := range oldMetadata.(map[string]interface{}) {
			var found bool
			for newKey := range newMetadata.(map[string]interface{}) {
				if oldKey == newKey {
					found = true
				}
			}

			if !found {
				metadataToDelete = append(metadataToDelete, oldKey)
			}
		}

		for _, key := range metadataToDelete {
			err := secrets.DeleteMetadatum(kmClient, d.Id(), key).ExtractErr()
			if err != nil {
				return fmt.Errorf("Error deleting metadata (%s) from secret (%s): %s", key, d.Id(), err)
			}
		}

		oldMetadata, newMetadata = d.GetChange("metadata")
		var metadataToAdd []string
		var metadataToUpdate []string

		// Determine if any metadata keys were removed from the configuration.
		// Then request those keys to be deleted.
		for newKey, newValue := range newMetadata.(map[string]interface{}) {
			var found bool
			for oldKey, oldValue := range oldMetadata.(map[string]interface{}) {
				found = true
				if oldKey == newKey {
					if newValue != oldValue {
						metadataToUpdate = append(metadataToUpdate, newKey)
					}
				}
			}

			if !found {
				metadataToAdd = append(metadataToAdd, newKey)
			}
		}

		for _, key := range metadataToUpdate {
			var metadatumOpts secrets.MetadatumOpts
			metadatumOpts.Key = key
			metadatumOpts.Value = newMetadata.(map[string]interface{})[key].(string)
			_, err := secrets.UpdateMetadatum(kmClient, d.Id(), metadatumOpts).Extract()
			if err != nil {
				return fmt.Errorf("Error updating OpenStack secret (%s) metadata: %s", d.Id(), err)
			}
		}

		for _, key := range metadataToAdd {
			var metadatumOpts secrets.MetadatumOpts
			metadatumOpts.Key = key
			metadatumOpts.Value = newMetadata.(map[string]interface{})[key].(string)
			err := secrets.CreateMetadatum(kmClient, d.Id(), metadatumOpts).Err
			if err != nil {
				return fmt.Errorf("Error updating OpenStack secret (%s) metadata: %s", d.Id(), err)
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
