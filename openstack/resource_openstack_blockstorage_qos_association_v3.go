package openstack

import (
	"context"
	"fmt"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/blockstorage/v3/qos"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceBlockStorageQosAssociationV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBlockStorageQosAssociationV3Create,
		ReadContext:   resourceBlockStorageQosAssociationV3Read,
		DeleteContext: resourceBlockStorageQosAssociationV3Delete,
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

			"qos_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"volume_type_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceBlockStorageQosAssociationV3Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	blockStorageClient, err := config.BlockStorageV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	qosID := d.Get("qos_id").(string)
	vtID := d.Get("volume_type_id").(string)
	associateOpts := qos.AssociateOpts{
		VolumeTypeID: vtID,
	}

	id := fmt.Sprintf("%s/%s", qosID, vtID)

	log.Printf("[DEBUG] openstack_blockstorage_qos_association_v3 create options: %#v", associateOpts)

	err = qos.Associate(ctx, blockStorageClient, qosID, associateOpts).ExtractErr()
	if err != nil {
		return diag.Errorf("Error creating openstack_blockstorage_qos_association_v3 %s: %s", id, err)
	}

	d.SetId(id)

	return resourceBlockStorageQosAssociationV3Read(ctx, d, meta)
}

func resourceBlockStorageQosAssociationV3Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	blockStorageClient, err := config.BlockStorageV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	qosID, vtID, err := parsePairedIDs(d.Id(), "openstack_blockstorage_qos_association_v3")
	if err != nil {
		return diag.FromErr(err)
	}

	allPages, err := qos.ListAssociations(blockStorageClient, qosID).AllPages(ctx)
	if err != nil {
		return diag.Errorf("Error retrieving associations openstack_blockstorage_qos_association_v3 for qos: %s", qosID)
	}

	allAssociations, err := qos.ExtractAssociations(allPages)
	if err != nil {
		return diag.Errorf("Error extracting associations openstack_blockstorage_qos_association_v3 for qos: %s", qosID)
	}

	found := false

	for _, association := range allAssociations {
		if association.ID == vtID {
			found = true

			break
		}
	}

	if !found {
		return diag.Errorf("Error getting qos association openstack_blockstorage_qos_association_v3 for qos: %s and vt: %s", qosID, vtID)
	}

	d.Set("region", GetRegion(d, config))
	d.Set("qos_id", qosID)
	d.Set("volume_type_id", vtID)

	return nil
}

func resourceBlockStorageQosAssociationV3Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	blockStorageClient, err := config.BlockStorageV3Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	qosID, vtID, err := parsePairedIDs(d.Id(), "openstack_blockstorage_qos_association_v3")
	if err != nil {
		return diag.FromErr(err)
	}

	disassociateOpts := qos.DisassociateOpts{
		VolumeTypeID: vtID,
	}

	err = qos.Disassociate(ctx, blockStorageClient, qosID, disassociateOpts).ExtractErr()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_blockstorage_qos_association_v3"))
	}

	return nil
}
