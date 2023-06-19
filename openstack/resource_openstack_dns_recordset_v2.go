package openstack

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/openstack/dns/v2/recordsets"
)

func resourceDNSRecordSetV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSRecordSetV2Create,
		ReadContext:   resourceDNSRecordSetV2Read,
		UpdateContext: resourceDNSRecordSetV2Update,
		DeleteContext: resourceDNSRecordSetV2Delete,
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
				ForceNew: true,
				Computed: true,
			},

			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
			},

			"records": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: false,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: false,
			},

			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"value_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},

			"disable_status_check": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceDNSRecordSetV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	dnsClient, err := config.DNSV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack DNS client: %s", err)
	}

	if err := dnsClientSetAuthHeader(d, dnsClient); err != nil {
		return diag.Errorf("Error setting dns client auth headers: %s", err)
	}

	records := []string{}
	if v, ok := d.GetOk("records"); ok {
		for _, vv := range v.(*schema.Set).List() {
			records = append(records, vv.(string))
		}
	}

	createOpts := RecordSetCreateOpts{
		recordsets.CreateOpts{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
			Records:     records,
			TTL:         d.Get("ttl").(int),
			Type:        d.Get("type").(string),
		},
		MapValueSpecs(d),
	}

	log.Printf("[DEBUG] openstack_dns_recordset_v2 create options: %#v", createOpts)

	zoneID := d.Get("zone_id").(string)
	n, err := recordsets.Create(dnsClient, zoneID, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_dns_recordset_v2: %s", err)
	}

	if !d.Get("disable_status_check").(bool) {
		stateConf := &resource.StateChangeConf{
			Target:     []string{"ACTIVE"},
			Pending:    []string{"PENDING"},
			Refresh:    dnsRecordSetV2RefreshFunc(dnsClient, zoneID, n.ID),
			Timeout:    d.Timeout(schema.TimeoutCreate),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf(
				"Error waiting for openstack_dns_recordset_v2 %s to become active: %s", d.Id(), err)
		}
	}
	id := fmt.Sprintf("%s/%s", zoneID, n.ID)
	d.SetId(id)

	log.Printf("[DEBUG] Created openstack_dns_recordset_v2 %s: %#v", n.ID, n)
	return resourceDNSRecordSetV2Read(ctx, d, meta)
}

func resourceDNSRecordSetV2Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	dnsClient, err := config.DNSV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack DNS client: %s", err)
	}

	if err := dnsClientSetAuthHeader(d, dnsClient); err != nil {
		return diag.Errorf("Error setting dns client auth headers: %s", err)
	}

	// Obtain relevant info from parsing the ID
	zoneID, recordsetID, err := dnsRecordSetV2ParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	n, err := recordsets.Get(dnsClient, zoneID, recordsetID).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_dns_recordset_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_dns_recordset_v2 %s: %#v", recordsetID, n)

	d.Set("name", n.Name)
	d.Set("description", n.Description)
	records := []interface{}{}
	for _, c := range n.Records {
		records = append(records, c)
	}
	d.Set("records", schema.NewSet(schema.HashString, records))
	d.Set("ttl", n.TTL)
	d.Set("type", n.Type)
	d.Set("zone_id", zoneID)
	d.Set("region", GetRegion(d, config))
	d.Set("project_id", n.ProjectID)

	return nil
}

func resourceDNSRecordSetV2Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	dnsClient, err := config.DNSV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack DNS client: %s", err)
	}

	if err := dnsClientSetAuthHeader(d, dnsClient); err != nil {
		return diag.Errorf("Error setting dns client auth headers: %s", err)
	}

	changed := false
	var updateOpts recordsets.UpdateOpts
	if d.HasChange("ttl") {
		ttl := d.Get("ttl").(int)
		updateOpts.TTL = &ttl
		changed = true
	}

	if d.HasChange("records") {
		records := []string{}
		if v, ok := d.GetOk("records"); ok {
			for _, vv := range v.(*schema.Set).List() {
				records = append(records, vv.(string))
			}
		}
		updateOpts.Records = records
		changed = true
	}

	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
		changed = true
	}

	if !changed {
		// Nothing in OpenStack fields really changed, so just return zone from OpenStack
		return resourceDNSRecordSetV2Read(ctx, d, meta)
	}

	// Obtain relevant info from parsing the ID
	zoneID, recordsetID, err := dnsRecordSetV2ParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Updating openstack_dns_recordset_v2 %s with options: %#v", recordsetID, updateOpts)

	_, err = recordsets.Update(dnsClient, zoneID, recordsetID, updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error updating openstack_dns_recordset_v2 %s: %s", d.Id(), err)
	}

	if !d.Get("disable_status_check").(bool) {
		stateConf := &resource.StateChangeConf{
			Target:     []string{"ACTIVE"},
			Pending:    []string{"PENDING"},
			Refresh:    dnsRecordSetV2RefreshFunc(dnsClient, zoneID, recordsetID),
			Timeout:    d.Timeout(schema.TimeoutUpdate),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf(
				"Error waiting for openstack_dns_recordset_v2 %s to become active: %s", d.Id(), err)
		}
	}

	return resourceDNSRecordSetV2Read(ctx, d, meta)
}

func resourceDNSRecordSetV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	dnsClient, err := config.DNSV2Client(GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack DNS client: %s", err)
	}

	if err := dnsClientSetAuthHeader(d, dnsClient); err != nil {
		return diag.Errorf("Error setting dns client auth headers: %s", err)
	}

	// Obtain relevant info from parsing the ID
	zoneID, recordsetID, err := dnsRecordSetV2ParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	err = recordsets.Delete(dnsClient, zoneID, recordsetID).ExtractErr()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_dns_recordset_v2"))
	}

	if !d.Get("disable_status_check").(bool) {
		stateConf := &resource.StateChangeConf{
			Target:     []string{"DELETED"},
			Pending:    []string{"ACTIVE", "PENDING"},
			Refresh:    dnsRecordSetV2RefreshFunc(dnsClient, zoneID, recordsetID),
			Timeout:    d.Timeout(schema.TimeoutDelete),
			Delay:      5 * time.Second,
			MinTimeout: 3 * time.Second,
		}

		_, err = stateConf.WaitForStateContext(ctx)
		if err != nil {
			return diag.Errorf(
				"Error waiting for openstack_dns_recordset_v2 %s to become deleted: %s", d.Id(), err)
		}
	}

	return nil
}
