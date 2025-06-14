package openstack

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/transfer/accept"
	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/transfer/request"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDNSTransferAcceptV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSTransferAcceptV2Create,
		ReadContext:   resourceDNSTransferAcceptV2Read,
		DeleteContext: resourceDNSTransferAcceptV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDNSTransferAcceptV2Import,
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

			"zone_transfer_request_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"key": {
				Type:     schema.TypeString,
				Required: true,
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
				ForceNew: true,
			},
		},
	}
}

func resourceDNSTransferAcceptV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	dnsClient, err := config.DNSV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack DNS client: %s", err)
	}

	createOpts := TransferAcceptCreateOpts{
		accept.CreateOpts{
			ZoneTransferRequestID: d.Get("zone_transfer_request_id").(string),
			Key:                   d.Get("key").(string),
		},
		MapValueSpecs(d),
	}

	log.Printf("[DEBUG] openstack_dns_transfer_accept_v2 create options: %#v", createOpts)

	n, err := accept.Create(ctx, dnsClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("Error creating openstack_transfer_accept_zone_v2: %s", err)
	}

	// Key is returned only once
	d.Set("key", n.Key)

	if d.Get("disable_status_check").(bool) {
		d.SetId(n.ID)

		log.Printf("[DEBUG] Created OpenStack Zone Transfer accept %s: %#v", n.ID, n)

		return resourceDNSTransferAcceptV2Read(ctx, d, meta)
	}

	stateConf := &retry.StateChangeConf{
		Target:     []string{"COMPLETE"},
		Pending:    []string{"PENDING"},
		Refresh:    dnsTransferAcceptV2RefreshFunc(ctx, dnsClient, n.ID),
		Timeout:    d.Timeout(schema.TimeoutCreate),
		Delay:      0,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"Error waiting for openstack_dns_transfer_accept_v2 %s to become active: %s", d.Id(), err)
	}

	d.SetId(n.ID)

	log.Printf("[DEBUG] Created OpenStack Zone Transfer accept %s: %#v", n.ID, n)

	return resourceDNSTransferAcceptV2Read(ctx, d, meta)
}

func resourceDNSTransferAcceptV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	dnsClient, err := config.DNSV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack DNS client: %s", err)
	}

	n, err := accept.Get(ctx, dnsClient, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error retrieving openstack_dns_transfer_accept_v2"))
	}

	log.Printf("[DEBUG] Retrieved openstack_dns_transfer_accept_v2 %s: %#v", d.Id(), n)

	d.Set("region", GetRegion(d, config))
	d.Set("zone_transfer_request_id", n.ZoneTransferRequestID)

	return nil
}

func resourceDNSTransferAcceptV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	dnsClient, err := config.DNSV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("Error creating OpenStack DNS client: %s", err)
	}

	if d.Get("disable_status_check").(bool) {
		return nil
	}

	err = request.Delete(ctx, dnsClient, d.Get("zone_transfer_request_id").(string)).ExtractErr()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "Error deleting openstack_dns_transfer_request_v2"))
	}

	stateConf := &retry.StateChangeConf{
		Target:     []string{"DELETED"},
		Pending:    []string{"ACTIVE"},
		Refresh:    dnsTransferAcceptV2RefreshFunc(ctx, dnsClient, d.Id()),
		Timeout:    d.Timeout(schema.TimeoutDelete),
		Delay:      0,
		MinTimeout: 3 * time.Second,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf(
			"Error waiting for openstack_dns_transfer_accept_v2 %s to become deleted: %s", d.Id(), err)
	}

	return nil
}

func resourceDNSTransferAcceptV2Import(ctx context.Context, d *schema.ResourceData, meta any) ([]*schema.ResourceData, error) {
	config := meta.(*Config)

	dnsClient, err := config.DNSV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return nil, fmt.Errorf("Error creating OpenStack DNS client: %w", err)
	}

	n, err := accept.Get(ctx, dnsClient, d.Id()).Extract()
	if err != nil {
		return nil, fmt.Errorf("Error retrieving openstack_dns_transfer_accept_v2 %s: %w", d.Id(), err)
	}

	d.Set("key", n.Key)

	return []*schema.ResourceData{d}, nil
}
