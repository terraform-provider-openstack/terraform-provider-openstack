package openstack

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/quotas"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDNSQuotaV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSQuotaV2Create,
		ReadContext:   resourceDNSQuotaV2Read,
		UpdateContext: resourceDNSQuotaV2Update,
		Delete:        schema.RemoveFromState,
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

			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"api_export_size": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"recordset_records": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"zone_records": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"zone_recordsets": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"zones": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceDNSQuotaV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)
	region := GetRegion(d, config)

	dnsClient, err := config.DNSV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack DNS client: %s", err)
	}

	if err := dnsClientSetAuthHeader(ctx, d, dnsClient); err != nil {
		return diag.Errorf("Error setting dns client auth headers: %s", err)
	}

	updateOpts := quotas.UpdateOpts{}
	projectID := d.Get("project_id").(string)

	if v, ok := getOkExists(d, "api_export_size"); ok {
		v := v.(int)
		updateOpts.APIExporterSize = &v
	}

	if v, ok := getOkExists(d, "recordset_records"); ok {
		v := v.(int)
		updateOpts.RecordsetRecords = &v
	}

	if v, ok := getOkExists(d, "zone_records"); ok {
		v := v.(int)
		updateOpts.ZoneRecords = &v
	}

	if v, ok := getOkExists(d, "zone_recordsets"); ok {
		v := v.(int)
		updateOpts.ZoneRecordsets = &v
	}

	if v, ok := getOkExists(d, "zones"); ok {
		v := v.(int)
		updateOpts.Zones = &v
	}

	q, err := quotas.Update(ctx, dnsClient, projectID, updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_dns_quota_v2: %s", err)
	}

	id := fmt.Sprintf("%s/%s", projectID, region)
	d.SetId(id)

	log.Printf("[DEBUG] Created openstack_dns_quota_v2 %#v", q)

	return resourceDNSQuotaV2Read(ctx, d, meta)
}

func resourceDNSQuotaV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)
	region := GetRegion(d, config)

	dnsClient, err := config.DNSV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack DNS client: %s", err)
	}

	// Depending on the provider version the resource was created, the resource id
	// can be either <project_id> or <project_id>/<region>. This parses the project_id
	// in both cases
	projectID := strings.Split(d.Id(), "/")[0]
	d.Set("project_id", projectID)

	if err := dnsClientSetAuthHeader(ctx, d, dnsClient); err != nil {
		return diag.Errorf("Error setting dns client auth headers: %s", err)
	}

	q, err := quotas.Get(ctx, dnsClient, projectID).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_dns_quota_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_dns_quota_v2 %s: %#v", d.Id(), q)

	d.Set("region", region)
	d.Set("api_export_size", q.APIExporterSize)
	d.Set("recordset_records", q.RecordsetRecords)
	d.Set("zone_records", q.ZoneRecords)
	d.Set("zone_recordsets", q.ZoneRecordsets)
	d.Set("zones", q.Zones)

	return nil
}

func resourceDNSQuotaV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	dnsClient, err := config.DNSV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack DNS client: %s", err)
	}

	if err := dnsClientSetAuthHeader(ctx, d, dnsClient); err != nil {
		return diag.Errorf("Error setting dns client auth headers: %s", err)
	}

	var (
		hasChange  bool
		updateOpts quotas.UpdateOpts
	)

	if d.HasChange("api_export_size") {
		hasChange = true
		v := d.Get("api_export_size").(int)
		updateOpts.APIExporterSize = &v
	}

	if d.HasChange("recordset_records") {
		hasChange = true
		v := d.Get("recordset_records").(int)
		updateOpts.RecordsetRecords = &v
	}

	if d.HasChange("zone_records") {
		hasChange = true
		v := d.Get("zone_records").(int)
		updateOpts.ZoneRecords = &v
	}

	if d.HasChange("zone_recordsets") {
		hasChange = true
		v := d.Get("zone_recordsets").(int)
		updateOpts.ZoneRecordsets = &v
	}

	if d.HasChange("zones") {
		hasChange = true
		v := d.Get("zones").(int)
		updateOpts.Zones = &v
	}

	if hasChange {
		log.Printf("[DEBUG] openstack_dns_quota_v2 %s update options: %#v", d.Id(), updateOpts)
		projectID := d.Get("project_id").(string)

		_, err := quotas.Update(ctx, dnsClient, projectID, updateOpts).Extract()
		if err != nil {
			return diag.Errorf("Error updating openstack_dns_quota_v2: %s", err)
		}
	}

	return resourceDNSQuotaV2Read(ctx, d, meta)
}
