package openstack

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/zones"
)

// ZoneShare represents a shared DNS zone in OpenStack.
type ZoneShare struct {
	ID              string `json:"id"`
	TargetProjectID string `json:"target_project_id"`
}

func resourceDNSZoneShareV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSZoneShareV2Create,
		ReadContext:   resourceDNSZoneShareV2Read,
		DeleteContext: resourceDNSZoneShareV2Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDNSZoneShareV2Importer,
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the DNS zone to share",
			},
			"target_project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The target project ID with which to share the DNS zone",
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The owner project ID. If omitted, it is derived from the zone details.",
			},
			"share_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the created share",
			},
		},
	}
}

// resourceDNSZoneShareV2Importer supports both full-format (<zone_id>/<project_id>/<target_project_id>/<share_id>)
// and simplified (<zone_id>/<share_id>) imports.
func resourceDNSZoneShareV2Importer(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	client, err := config.DNSV2Client(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("error creating OpenStack DNS client: %s", err)
	}

	parts := strings.Split(d.Id(), "/")
	switch len(parts) {
	case 4:
		// Full format: <zone_id>/<project_id>/<target_project_id>/<share_id>
		zoneID := parts[0]
		projectID := parts[1]
		targetProjectID := parts[2]
		shareID := parts[3]

		d.Set("zone_id", zoneID)
		d.Set("project_id", projectID)
		d.Set("target_project_id", targetProjectID)
		d.Set("share_id", shareID)
		// Canonicalize the ID as <zone_id>/<share_id>
		d.SetId(fmt.Sprintf("%s/%s", zoneID, shareID))
	case 2:
		zoneID := parts[0]
		shareID := parts[1]
		d.Set("zone_id", zoneID)
		d.Set("share_id", shareID)

		// Determine the owner project ID.
		var ownerProjectID string
		if v, ok := d.GetOk("project_id"); ok && v.(string) != "" {
			ownerProjectID = v.(string)
		} else {
			zone, err := zones.Get(ctx, client, zoneID).Extract()
			if err != nil {
				return nil, fmt.Errorf("error retrieving zone details for zone %s: %s", zoneID, err)
			}
			ownerProjectID = zone.ProjectID
			d.Set("project_id", ownerProjectID)
		}

		// Determine target_project_id.
		targetProjectID := d.Get("target_project_id").(string)
		if targetProjectID == "" {
			shares, err := listZoneShares(ctx, client, zoneID, ownerProjectID)
			if err != nil {
				return nil, fmt.Errorf("error listing shares for DNS zone %s: %s", zoneID, err)
			}
			found := false
			for _, s := range shares {
				if s.ID == shareID {
					targetProjectID = s.TargetProjectID
					found = true
					break
				}
			}
			if !found {
				return nil, fmt.Errorf("could not find share with id %s for zone %s", shareID, zoneID)
			}
		}
		d.Set("target_project_id", targetProjectID)
		d.SetId(fmt.Sprintf("%s/%s", zoneID, shareID))
	default:
		return nil, fmt.Errorf("unexpected format (%s); expected either <zone_id>/<share_id> or <zone_id>/<project_id>/<target_project_id>/<share_id>", d.Id())
	}

	diags := resourceDNSZoneShareV2Read(ctx, d, meta)
	if diags.HasError() {
		return nil, fmt.Errorf("error reading DNS zone share: %v", diags)
	}
	return []*schema.ResourceData{d}, nil
}

func resourceDNSZoneShareV2Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client, err := config.DNSV2Client(ctx, "")
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating OpenStack DNS client: %s", err))
	}

	zoneID := d.Get("zone_id").(string)
	targetProjectID := d.Get("target_project_id").(string)

	// Use the provided project_id if set; otherwise, retrieve it from the zone.
	var sudoProjectID string
	if v, ok := d.GetOk("project_id"); ok && v.(string) != "" {
		sudoProjectID = v.(string)
	} else {
		zone, err := zones.Get(ctx, client, zoneID).Extract()
		if err != nil {
			return diag.FromErr(fmt.Errorf("error retrieving zone details for zone %s: %s", zoneID, err))
		}
		sudoProjectID = zone.ProjectID
		d.Set("project_id", sudoProjectID)
	}

	shareOpts := zones.ShareZoneOpts{
		TargetProjectID: targetProjectID,
	}

	client.Microversion = "2.0"
	client.MoreHeaders = map[string]string{
		"X-Auth-Sudo-Project-Id": sudoProjectID,
	}

	if err := zones.Share(ctx, client, zoneID, shareOpts).ExtractErr(); err != nil {
		return diag.FromErr(fmt.Errorf("error sharing DNS zone %s with project %s: %s", zoneID, targetProjectID, err))
	}

	shares, err := listZoneShares(ctx, client, zoneID, sudoProjectID)
	if err != nil {
		return diag.FromErr(err)
	}
	var shareID string
	for _, s := range shares {
		if s.TargetProjectID == targetProjectID {
			shareID = s.ID
			break
		}
	}
	if shareID == "" {
		return diag.FromErr(fmt.Errorf("failed to locate share for DNS zone %s with project %s", zoneID, targetProjectID))
	}

	d.SetId(fmt.Sprintf("%s/%s", zoneID, shareID))
	d.Set("share_id", shareID)
	return resourceDNSZoneShareV2Read(ctx, d, meta)
}

func resourceDNSZoneShareV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client, err := config.DNSV2Client(ctx, "")
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating OpenStack DNS client: %s", err))
	}

	zoneID, shareID, err := parseDNSSharedZoneID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	var ownerProjectID string
	if v, ok := d.GetOk("project_id"); ok && v.(string) != "" {
		ownerProjectID = v.(string)
	} else {
		zone, err := zones.Get(ctx, client, zoneID).Extract()
		if err != nil {
			return diag.FromErr(fmt.Errorf("error retrieving zone details for zone %s: %s", zoneID, err))
		}
		ownerProjectID = zone.ProjectID
		d.Set("project_id", ownerProjectID)
	}

	targetProjectID := d.Get("target_project_id").(string)

	shares, err := listZoneShares(ctx, client, zoneID, ownerProjectID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error listing shares for DNS zone %s: %s", zoneID, err))
	}

	found := false
	for _, s := range shares {
		if s.ID == shareID && s.TargetProjectID == targetProjectID {
			found = true
			break
		}
	}
	if !found {
		return diag.Errorf("DNS zone share not found: zone_id %s, share_id %s", zoneID, shareID)
	}
	return nil
}

func resourceDNSZoneShareV2Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client, err := config.DNSV2Client(ctx, "")
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating OpenStack DNS client: %s", err))
	}

	zoneID, shareID, err := parseDNSSharedZoneID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	var sudoProjectID string
	if v, ok := d.GetOk("project_id"); ok && v.(string) != "" {
		sudoProjectID = v.(string)
	} else {
		zone, err := zones.Get(ctx, client, zoneID).Extract()
		if err != nil {
			return diag.FromErr(fmt.Errorf("error retrieving zone details for zone %s: %s", zoneID, err))
		}
		sudoProjectID = zone.ProjectID
		d.Set("project_id", sudoProjectID)
	}

	url := client.ServiceURL("zones", zoneID, "shares", shareID)
	reqOpts := &gophercloud.RequestOpts{
		MoreHeaders: map[string]string{
			"X-Auth-Sudo-Project-Id": sudoProjectID,
		},
	}

	resp, err := client.Delete(ctx, url, reqOpts)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error unsharing DNS zone %s: %s", zoneID, err))
	}
	defer resp.Body.Close()

	d.SetId("")
	return nil
}

func parseDNSSharedZoneID(id string) (string, string, error) {
	if !strings.Contains(id, "/") {
		return "", id, nil
	}
	parts := strings.Split(id, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("unexpected ID format (%s), expected <zone_id>/<share_id>", id)
	}
	return parts[0], parts[1], nil
}

func listZoneShares(ctx context.Context, client *gophercloud.ServiceClient, zoneID string, ownerProjectID string) ([]ZoneShare, error) {
	url := client.ServiceURL("zones", zoneID, "shares")
	var result struct {
		Shares []ZoneShare `json:"shared_zones"`
	}
	reqOpts := &gophercloud.RequestOpts{}
	if ownerProjectID != "" {
		reqOpts.MoreHeaders = map[string]string{
			"X-Auth-Sudo-Project-Id": ownerProjectID,
		}
	}

	resp, err := client.Get(ctx, url, &result, reqOpts)
	if err != nil {
		return nil, fmt.Errorf("error getting shares for zone %s: %s", zoneID, err)
	}
	defer resp.Body.Close()
	return result.Shares, nil
}
