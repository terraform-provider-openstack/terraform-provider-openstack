package openstack

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/aggregates"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceComputeAggregateV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceComputeAggregateV2Create,
		ReadContext:   resourceComputeAggregateV2Read,
		UpdateContext: resourceComputeAggregateV2Update,
		DeleteContext: resourceComputeAggregateV2Delete,
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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"zone": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"metadata": {
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				DefaultFunc: func() (any, error) { return map[string]any{}, nil },
			},

			"hosts": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				DefaultFunc: func() (any, error) { return []string{}, nil },
			},
		},
	}
}

func resourceComputeAggregateV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	computeClient, err := config.ComputeV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	aggregate, err := aggregates.Create(ctx, computeClient, aggregates.CreateOpts{
		Name:             d.Get("name").(string),
		AvailabilityZone: d.Get("zone").(string),
	}).Extract()
	if err != nil {
		return diag.Errorf("Error creating OpenStack aggregate: %s", err)
	}

	idStr := strconv.Itoa(aggregate.ID)
	d.SetId(idStr)

	h, ok := d.GetOk("hosts")
	if ok {
		hosts := h.(*schema.Set)
		for _, host := range hosts.List() {
			_, err = aggregates.AddHost(ctx, computeClient, aggregate.ID, aggregates.AddHostOpts{Host: host.(string)}).Extract()
			if err != nil {
				return diag.Errorf("Error adding host %s to OpenStack aggregate: %s", host.(string), err)
			}
		}
	}

	_, err = aggregates.SetMetadata(ctx, computeClient, aggregate.ID, aggregates.SetMetadataOpts{Metadata: d.Get("metadata").(map[string]any)}).Extract()
	if err != nil {
		return diag.Errorf("Error setting metadata: %s", err)
	}

	return nil
}

func resourceComputeAggregateV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	computeClient, err := config.ComputeV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Can't convert ID to integer: %s", err)
	}

	aggregate, err := aggregates.Get(ctx, computeClient, id).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error getting host aggregate"))
	}

	// Metadata is redundant with Availability Zone
	metadata := aggregate.Metadata
	_, ok := metadata["availability_zone"]

	if ok {
		delete(metadata, "availability_zone")
	}

	d.Set("name", aggregate.Name)
	d.Set("zone", aggregate.AvailabilityZone)
	d.Set("hosts", aggregate.Hosts)
	d.Set("metadata", metadata)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceComputeAggregateV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	computeClient, err := config.ComputeV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Can't convert ID to integer: %s", err)
	}

	var updateOpts aggregates.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}

	if d.HasChange("zone") {
		updateOpts.AvailabilityZone = d.Get("zone").(string)
	}

	if updateOpts != (aggregates.UpdateOpts{}) {
		_, err = aggregates.Update(ctx, computeClient, id, updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating OpenStack aggregate: %s", err)
		}
	}

	if d.HasChange("hosts") {
		o, n := d.GetChange("hosts")
		oldHosts, newHosts := o.(*schema.Set), n.(*schema.Set)
		hostsToDelete := oldHosts.Difference(newHosts)
		hostsToAdd := newHosts.Difference(oldHosts)

		for _, h := range hostsToDelete.List() {
			host := h.(string)
			log.Printf("[DEBUG] Removing host '%s' from aggregate '%s'", host, d.Get("name"))

			_, err = aggregates.RemoveHost(ctx, computeClient, id, aggregates.RemoveHostOpts{Host: host}).Extract()
			if err != nil {
				return diag.Errorf("Error deleting host %s from OpenStack aggregate: %s", host, err)
			}
		}

		for _, h := range hostsToAdd.List() {
			host := h.(string)
			log.Printf("[DEBUG] Adding host '%s' to aggregate '%s'", host, d.Get("name"))

			_, err = aggregates.AddHost(ctx, computeClient, id, aggregates.AddHostOpts{Host: host}).Extract()
			if err != nil {
				return diag.Errorf("Error adding host %s to OpenStack aggregate: %s", host, err)
			}
		}
	}

	if d.HasChange("metadata") {
		oldMetadata, newMetadata := d.GetChange("metadata")
		metadata := mapDiffWithNilValues(oldMetadata.(map[string]any), newMetadata.(map[string]any))

		_, err = aggregates.SetMetadata(ctx, computeClient, id, aggregates.SetMetadataOpts{Metadata: metadata}).Extract()
		if err != nil {
			return diag.Errorf("Error setting metadata: %s", err)
		}
	}

	return nil
}

func resourceComputeAggregateV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	computeClient, err := config.ComputeV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack compute client: %s", err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("Can't convert ID to integer: %s", err)
	}

	// OpenStack do not delete the host aggregate if it's not empty
	hostsToDelete := d.Get("hosts").(*schema.Set)
	for _, h := range hostsToDelete.List() {
		host := h.(string)
		log.Printf("[DEBUG] Removing host '%s' from aggregate '%s'", host, d.Get("name"))

		_, err = aggregates.RemoveHost(ctx, computeClient, id, aggregates.RemoveHostOpts{Host: host}).Extract()
		if err != nil {
			return diag.Errorf("Error deleting host %s from OpenStack aggregate: %s", host, err)
		}
	}

	err = aggregates.Delete(ctx, computeClient, id).ExtractErr()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting OpenStack aggregate"))
	}

	return nil
}
