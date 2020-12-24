package openstack

import (
	"fmt"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/aggregates"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"strconv"
	"time"
)

func resourceComputeAggregateV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeAggregateV2Create,
		Read:   resourceComputeAggregateV2Read,
		Update: resourceComputeAggregateV2Update,
		Delete: resourceComputeAggregateV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"zone": {
				Type:     schema.TypeString,
				Required: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
			"hosts": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}

func resourceComputeAggregateV2Create(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}

	aggregate, err := aggregates.Create(computeClient, aggregates.CreateOpts{
		Name:             d.Get("name").(string),
		AvailabilityZone: d.Get("zone").(string),
	}).Extract()
	if err != nil {
		return fmt.Errorf("Error creating OpenStack aggregate: %s", err)
	}
	id_str := strconv.Itoa(aggregate.ID)
	d.SetId(id_str)

	hosts, ok := d.GetOk("hosts")
	if ok {
		for _, host := range hosts.([]string) {
			_, err = aggregates.AddHost(computeClient, aggregate.ID, aggregates.AddHostOpts{Host: host}).Extract()
			if err != nil {
				return fmt.Errorf("Error adding host %s to Openstack aggregate: %s", host, err)
			}
		}
	}

	_, err = aggregates.SetMetadata(computeClient, aggregate.ID, aggregates.SetMetadataOpts{Metadata: d.Get("metadata").(map[string]interface{})}).Extract()
	if err != nil {
		return fmt.Errorf("Error setting metadata: %s", err)
	}

	return nil

}

func resourceComputeAggregateV2Read(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Can't convert ID to integer: %s", err)
	}

	aggregate, err := aggregates.Get(computeClient, id).Extract()
	if err != nil {
		return fmt.Errorf("Error getting host aggregate: %s", err)
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

	return nil
}

func resourceComputeAggregateV2Update(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Can't convert ID to integer: %s", err)
	}

	var updateOpts aggregates.UpdateOpts
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("zone") {
		updateOpts.AvailabilityZone = d.Get("zone").(string)
	}

	if updateOpts != (aggregates.UpdateOpts{}) {
		_, err = aggregates.Update(computeClient, id, updateOpts).Extract()
		if err != nil {
			return fmt.Errorf("Error updating OpenStack aggregate: %s", err)
		}
	}

	if d.HasChange("hosts") {
		old_hosts, new_hosts := d.GetChange("hosts")
		hosts_to_delete := arrayDifference(old_hosts, new_hosts)
		hosts_to_add := arrayDifference(new_hosts, old_hosts)
		for _, host := range hosts_to_delete {
			_, err = aggregates.RemoveHost(computeClient, id, aggregates.RemoveHostOpts{Host: host}).Extract()
			if err != nil {
				return fmt.Errorf("Error adding host %s to Openstack aggregate: %s", host, err)
			}
		}
		for _, host := range hosts_to_add {
			_, err = aggregates.AddHost(computeClient, id, aggregates.AddHostOpts{Host: host}).Extract()
			if err != nil {
				return fmt.Errorf("Error adding host %s to Openstack aggregate: %s", host, err)
			}
		}
	}

	if d.HasChange("metadata") {
		_, err = aggregates.SetMetadata(computeClient, id, aggregates.SetMetadataOpts{Metadata: d.Get("metadata").(map[string]interface{})}).Extract()
		if err != nil {
			return fmt.Errorf("Error setting metadata: %s", err)
		}
	}

	return nil

}

func resourceComputeAggregateV2Delete(d *schema.ResourceData, meta interface{}) error {

	config := meta.(*Config)
	computeClient, err := config.ComputeV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack compute client: %s", err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("Can't convert ID to integer: %s", err)
	}

	err = aggregates.Delete(computeClient, id).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting Openstack aggregate: %s", err)
	}

	return nil

}

func arrayDifference(a, b interface{}) (diff []string) {
	m := make(map[string]bool)

	for _, item := range b.([]string) {
		m[item] = true
	}
	for _, item := range a.([]string) {
		_, ok := m[item]
		if !ok {
			diff = append(diff, item)
		}
	}
	return
}
