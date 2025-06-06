package openstack

import (
	"context"
	"log"

	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/zones"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// dataSourceDNSZoneShareV2 defines the schema for the DNS Zone Share data source.
func dataSourceDNSZoneShareV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDNSZoneShareV2Read,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"all_projects": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"share_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the DNS zone.",
			},

			"target_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional: If provided, filter shares by target_project_id.",
			},

			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Optional: The owner project ID. If omitted, it is derived from the zone details.",
			},
		},
	}
}

// dataSourceDNSZoneShareV2Read fetches the zone share details based on the provided parameters.
func dataSourceDNSZoneShareV2Read(ctx context.Context, d *schema.ResourceData, meta any) diag.Diagnostics {
	config := meta.(*Config)

	dnsClient, err := config.DNSV2Client(ctx, GetRegion(d, config))
	if err != nil {
		return diag.Errorf("error creating DNS client: %s", err)
	}

	if err := dnsClientSetAuthHeader(ctx, d, dnsClient); err != nil {
		return diag.Errorf("Error setting dns client auth headers: %s", err)
	}

	dnsClient.Microversion = "2.1"

	zoneID := d.Get("zone_id").(string)
	targetProjectID := d.Get("target_project_id").(string)

	shareID := d.Get("share_id").(string)
	if shareID != "" {
		// If share_id is provided, fetch the specific share.
		log.Printf("[DEBUG] Fetching specific DNS Zone share with ID: %s", shareID)

		share, err := zones.GetShare(ctx, dnsClient, zoneID, shareID).Extract()
		if err != nil {
			return diag.Errorf("error retrieving DNS Zone share %s for zone %s: %s", shareID, zoneID, err)
		}

		if targetProjectID != "" && share.TargetProjectID == targetProjectID {
			log.Printf("[DEBUG] Retrieved DNS Zone share %s: %+v", share.ID, share)

			d.SetId(share.ID)
			d.Set("share_id", share.ID)
			d.Set("project_id", share.ProjectID)
			d.Set("target_project_id", share.TargetProjectID)

			return nil
		}
	}

	// "all_projects" and "project_id" are handled in the dnsClientSetAuthHeader function.
	allPages, err := zones.ListShares(dnsClient, zoneID, nil).AllPages(ctx)
	if err != nil {
		return diag.Errorf("error listing shares for DNS zone %s: %s", zoneID, err)
	}

	// Extract the shares from the response.
	shares, err := zones.ExtractZoneShares(allPages)
	if err != nil {
		return diag.Errorf("error extracting shares for DNS zone %s: %s", zoneID, err)
	}

	var filteredShares []zones.ZoneShare

	if targetProjectID != "" {
		for _, share := range shares {
			if share.TargetProjectID == targetProjectID {
				filteredShares = append(filteredShares, share)
			}
		}
	} else {
		filteredShares = shares
	}

	if len(filteredShares) == 0 {
		return diag.Errorf("no shares found for DNS zone %s", zoneID)
	}

	if len(filteredShares) > 1 {
		return diag.Errorf("multiple shares found for DNS zone %s", zoneID)
	}

	share := filteredShares[0]

	log.Printf("[DEBUG] Retrieved DNS Zone share %s: %+v", share.ID, share)
	d.SetId(share.ID)
	d.Set("share_id", share.ID)
	d.Set("project_id", share.ProjectID)
	d.Set("target_project_id", share.TargetProjectID)
	d.Set("region", GetRegion(d, config))

	return nil
}
