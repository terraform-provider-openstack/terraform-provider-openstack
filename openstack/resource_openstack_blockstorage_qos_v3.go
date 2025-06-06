package openstack

import (
	"context"
	"fmt"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/qos"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceBlockStorageQosV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBlockStorageQosV3Create,
		ReadContext:   resourceBlockStorageQosV3Read,
		UpdateContext: resourceBlockStorageQosV3Update,
		DeleteContext: resourceBlockStorageQosV3Delete,
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

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"consumer": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"front-end", "back-end", "both",
				}, false),
			},

			"specs": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func resourceBlockStorageQosV3Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	blockStorageClient, err := config.BlockStorageV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	name := d.Get("name").(string)
	consumer := qos.QoSConsumer(d.Get("consumer").(string))
	specs := d.Get("specs").(map[string]any)
	createOpts := qos.CreateOpts{
		Name:     name,
		Consumer: consumer,
		Specs:    expandToMapStringString(specs),
	}

	log.Printf("[DEBUG] openstack_blockstorage_qos_v3 create options: %#v", createOpts)

	qosRes, err := qos.Create(ctx, blockStorageClient, &createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_blockstorage_qos_v3 %s: %s", name, err)
	}

	d.SetId(qosRes.ID)

	return resourceBlockStorageQosV3Read(ctx, d, meta)
}

func resourceBlockStorageQosV3Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	blockStorageClient, err := config.BlockStorageV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	qosRes, err := qos.Get(ctx, blockStorageClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_blockstorage_qos_v3"))
	}

	log.Printf("[DEBUG] Retrieved openstack_blockstorage_qos_v3 %s: %#v", d.Id(), qosRes)

	d.Set("region", GetRegion(d, config))
	d.Set("name", qosRes.Name)
	d.Set("consumer", qosRes.Consumer)

	if err := d.Set("specs", qosRes.Specs); err != nil {
		log.Printf("[WARN] Unable to set specs for openstack_blockstorage_qos_v3 %s: %s", d.Id(), err)
	}

	return nil
}

func resourceBlockStorageQosV3Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	blockStorageClient, err := config.BlockStorageV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	hasChange := false

	var updateOpts qos.UpdateOpts

	if d.HasChange("consumer") {
		hasChange = true
		consumer := qos.QoSConsumer(d.Get("consumer").(string))
		updateOpts.Consumer = consumer
	}

	if d.HasChange("specs") {
		oldSpecsRaw, newSpecsRaw := d.GetChange("specs")

		// Delete all old specs.
		var deleteKeys qos.DeleteKeysOpts
		for oldKey := range oldSpecsRaw.(map[string]any) {
			deleteKeys = append(deleteKeys, oldKey)
		}

		err = qos.DeleteKeys(ctx, blockStorageClient, d.Id(), deleteKeys).ExtractErr()
		if err != nil {
			return diag.Errorf("Error deleting specs for openstack_blockstorage_qos_v3 %s: %s", d.Id(), err)
		}

		// Add new specs to UpdateOpts
		newSpecs := expandToMapStringString(newSpecsRaw.(map[string]any))
		if len(newSpecs) > 0 {
			hasChange = true
			updateOpts.Specs = newSpecs
		}
	}

	if hasChange {
		_, err = qos.Update(ctx, blockStorageClient, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_blockstorage_qos_v3 %s: %s", d.Id(), err)
		}
	}

	return resourceBlockStorageQosV3Read(ctx, d, meta)
}

func resourceBlockStorageQosV3Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	blockStorageClient, err := config.BlockStorageV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	// remove all associations first
	err = qos.DisassociateAll(ctx, blockStorageClient, d.Id()).ExtractErr()
	if err != nil && CheckDeleted(d, err, "") != nil {
		return diag.FromErr(fmt.Errorf("Error deleting openstack_blockstorage_qos_v3 associations: %w", err))
	}

	// Delete the QoS itself
	err = qos.Delete(ctx, blockStorageClient, d.Id(), qos.DeleteOpts{}).ExtractErr()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_blockstorage_qos_v3"))
	}

	return nil
}
