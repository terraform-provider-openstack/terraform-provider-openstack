package openstack

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/transfer/request"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDNSTransferRequestV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSTransferRequestV2Create,
		ReadContext:   resourceDNSTransferRequestV2Read,
		UpdateContext: resourceDNSTransferRequestV2Update,
		DeleteContext: resourceDNSTransferRequestV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDNSTransferRequestV2Import,
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

			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"target_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"key": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: false,
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
		},
	}
}

func resourceDNSTransferRequestV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	dnsClient, err := config.DNSV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack DNS client: %s", err)
	}

	createOpts := TransferRequestCreateOpts{
		request.CreateOpts{
			TargetProjectID: d.Get("target_project_id").(string),
			Description:     d.Get("description").(string),
		},
		MapValueSpecs(d),
	}

	log.Printf("[DEBUG] openstack_dns_transfer_request_v2 create options: %#v", createOpts)

	zoneID := d.Get("zone_id").(string)

	n, err := request.Create(ctx, dnsClient, zoneID, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_transfer_request_zone_v2: %s", err)
	}

	if d.Get("disable_status_check").(bool) {
		d.SetId(n.ID)

		log.Printf("[DEBUG] Created OpenStack Zone Transfer request %s: %#v", n.ID, n)

		return resourceDNSTransferRequestV2Read(ctx, d, meta)
	}

	stateConf := &retry.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Pending:    []string{"PENDING"},
		Refresh:    dnsTransferRequestV2RefreshFunc(ctx, dnsClient, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"Error waiting for openstack_dns_transfer_request_v2 %s to become active: %s", d.Id(), err)
	}

	d.SetId(n.ID)

	log.Printf("[DEBUG] Created OpenStack Zone Transfer request %s: %#v", n.ID, n)

	return resourceDNSTransferRequestV2Read(ctx, d, meta)
}

func resourceDNSTransferRequestV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	dnsClient, err := config.DNSV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack DNS client: %s", err)
	}

	n, err := request.Get(ctx, dnsClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_dns_transfer_request_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_dns_transfer_request_v2 %s: %#v", d.Id(), n)

	d.Set("region", GetRegion(d, config))
	d.Set("zone_id", n.ZoneID)
	d.Set("target_project_id", n.TargetProjectID)
	d.Set("description", n.Description)
	d.Set("key", n.Key)

	return nil
}

func resourceDNSTransferRequestV2Update(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	dnsClient, err := config.DNSV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack DNS client: %s", err)
	}

	var updateOpts request.UpdateOpts

	changed := false

	if d.HasChange("target_project_id") {
		updateOpts.TargetProjectID = d.Get("target_project_id").(string)
		changed = true
	}

	if d.HasChange("description") {
		updateOpts.Description = d.Get("description").(string)
		changed = true
	}

	if !changed {
		return resourceDNSTransferRequestV2Read(ctx, d, meta)
	}

	log.Printf("[DEBUG] Updating openstack_dns_transfer_request_v2 %s with options: %#v", d.Id(), updateOpts)

	_, err = request.Update(ctx, dnsClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("Error updating openstack_dns_transfer_request_v2 %s: %s", d.Id(), err)
	}

	if d.Get("disable_status_check").(bool) {
		return resourceDNSTransferRequestV2Read(ctx, d, meta)
	}

	stateConf := &retry.StateChangeConf{
		Target:     []string{"ACTIVE"},
		Pending:    []string{"PENDING"},
		Refresh:    dnsTransferRequestV2RefreshFunc(ctx, dnsClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutUpdate),
		Delay:      0,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"Error waiting for openstack_dns_transfer_request_v2 %s to become active: %s", d.Id(), err)
	}

	return resourceDNSTransferRequestV2Read(ctx, d, meta)
}

func resourceDNSTransferRequestV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	dnsClient, err := config.DNSV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack DNS client: %s", err)
	}

	err = request.Delete(ctx, dnsClient, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_dns_transfer_request_v2"))
	}

	if d.Get("disable_status_check").(bool) {
		return nil
	}

	stateConf := &retry.StateChangeConf{
		Target:     []string{"DELETED"},
		Pending:    []string{"ACTIVE", "PENDING"},
		Refresh:    dnsTransferRequestV2RefreshFunc(ctx, dnsClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      0,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"Error waiting for openstack_dns_transfer_request_v2 %s to become deleted: %s", d.Id(), err)
	}

	return nil
}

func resourceDNSTransferRequestV2Import(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	dnsClient, err := config.DNSV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return nil, fmt.Errorf("Error creating OpenStack DNS client: %w", err)
	}

	n, err := request.Get(ctx, dnsClient, d.Id()).Extract()
	if err != nil {
		return nil, fmt.Errorf("Error retrieving openstack_dns_transfer_request_v2 %s: %w", d.Id(), err)
	}

	d.Set("zone_id", n.ZoneID)

	return []*schema.ResourceData{d}, nil
}
