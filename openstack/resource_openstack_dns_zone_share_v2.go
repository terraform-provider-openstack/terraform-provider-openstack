package openstack

import (
	"context"
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/zones"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDNSZoneShareV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSZoneShareV2Create,
		ReadContext:   resourceDNSZoneShareV2Read,
		DeleteContext: resourceDNSZoneShareV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDNSZoneShareV2Import,
		},
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the DNS zone to share",
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The owner project ID. If omitted, it is derived from the zone details.",
			},
			"target_project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The target project ID with which to share the DNS zone",
			},
		},
	}
}

func resourceDNSZoneShareV2Create(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	dnsClient, err := config.DNSV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating OpenStack DNS client: %s", err)
	}

	if err := dnsClientSetAuthHeader(ctx, d, dnsClient); err != nil {
		return diag.Errorf("Error setting dns client auth headers: %s", err)
	}

	dnsClient.Microversion = "2.1"

	zoneID := d.Get("zone_id").(string)
	targetProjectID := d.Get("target_project_id").(string)
	shareOpts := zones.ShareZoneOpts{
		TargetProjectID: targetProjectID,
	}

	share, err := zones.Share(ctx, dnsClient, zoneID, shareOpts).Extract()
	if err != nil {
		return diag.Errorf("error sharing DNS zone %s with project %s: %s", zoneID, targetProjectID, err)
	}

	d.SetId(share.ID)

	return resourceDNSZoneShareV2Read(ctx, d, meta)
}

func resourceDNSZoneShareV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	dnsClient, err := config.DNSV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating OpenStack DNS client: %s", err)
	}

	if err := dnsClientSetAuthHeader(ctx, d, dnsClient); err != nil {
		return diag.Errorf("Error setting dns client auth headers: %s", err)
	}

	dnsClient.Microversion = "2.1"

	zoneID := d.Get("zone_id").(string)
	share, err := zones.GetShare(ctx, dnsClient, zoneID, d.Id()).Extract()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "error retrieving share details for DNS zone: "+zoneID))
	}

	d.Set("zone_id", share.ZoneID)
	d.Set("project_id", share.ProjectID)
	d.Set("target_project_id", share.TargetProjectID)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceDNSZoneShareV2Delete(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	dnsClient, err := config.DNSV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating OpenStack DNS client: %s", err)
	}

	if err := dnsClientSetAuthHeader(ctx, d, dnsClient); err != nil {
		return diag.Errorf("Error setting dns client auth headers: %s", err)
	}

	dnsClient.Microversion = "2.1"

	zoneID := d.Get("zone_id").(string)

	err = zones.Unshare(ctx, dnsClient, zoneID, d.Id()).ExtractErr()
	if err != nil {
		return diag.FromErr(CheckDeleted(d, err, "error unsharing DNS zone"))
	}

	return nil
}

// resourceDNSZoneShareV2Import supports both full-format
// (<zone_id>/<share_id>/<project_id>) and simplified (<zone_id>/<share_id>)
// imports.
func resourceDNSZoneShareV2Import(_ context.Context, d *schema.ResourceData, _ any) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 && len(parts) != 3 {
		return nil, fmt.Errorf("unexpected ID format (%s). Expected either <zone_id>/<share_id> or <zone_id>/<share_id>/<project_id>", d.Id())
	}

	zoneID := parts[0]
	shareID := parts[1]

	d.Set("zone_id", zoneID)
	d.SetId(shareID)

	if len(parts) == 3 {
		projectID := parts[2]
		d.Set("project_id", projectID)
	}

	return []*schema.ResourceData{d}, nil
}
