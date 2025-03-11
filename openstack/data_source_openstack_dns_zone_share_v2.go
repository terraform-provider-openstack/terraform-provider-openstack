package openstack

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/gophercloud/gophercloud/v2/openstack/dns/v2/zones"
)

// dataSourceDNSZoneShareV2 defines the schema for the DNS Zone Share data source.
func dataSourceDNSZoneShareV2() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDNSZoneShareV2Read,
		Schema: map[string]*schema.Schema{
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
			"shares": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of zone shares matching the filter.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"share_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The share ID.",
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The target project ID for this share.",
						},
					},
				},
			},
		},
	}
}

// flattenZoneShares converts a slice of ZoneShare objects into a slice of maps.
// If targetFilter is provided and non-empty, only shares with a matching TargetProjectID are included.
func flattenZoneShares(shares []ZoneShare, targetFilter interface{}) []map[string]interface{} {
	results := make([]map[string]interface{}, 0)
	filter := ""
	if targetFilter != nil {
		if s, ok := targetFilter.(string); ok {
			filter = s
		}
	}
	for _, s := range shares {
		if filter != "" && s.TargetProjectID != filter {
			continue
		}
		results = append(results, map[string]interface{}{
			"share_id":   s.ID,
			"project_id": s.TargetProjectID,
		})
	}
	return results
}

// dataSourceDNSZoneShareV2Read fetches the zone details and associated shares,
// updates the zone_id to the zone's FQDN, and sets the "shares" attribute using the flatten helper.
func dataSourceDNSZoneShareV2Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	client, err := config.DNSV2Client(ctx, "")
	if err != nil {
		return diag.FromErr(fmt.Errorf("error creating DNS client: %s", err))
	}

	// Retrieve the zone ID (as provided in the configuration)
	zoneID := d.Get("zone_id").(string)

	// Fetch the zone details so that we can update the zone_id to its FQDN.
	zone, err := zones.Get(ctx, client, zoneID).Extract()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving zone details for zone %s: %s", zoneID, err))
	}

	// Determine the owner project ID either from configuration or from the zone details.
	var ownerProjectID string
	if v, ok := d.GetOk("project_id"); ok {
		ownerProjectID = v.(string)
	} else {
		ownerProjectID = zone.ProjectID
	}

	// Update the zone_id attribute to the zone's FQDN.
	if err := d.Set("zone_id", zone.Name); err != nil {
		return diag.FromErr(err)
	}

	// Retrieve the filter value (if any) for target_project_id.
	targetFilter, _ := d.GetOk("target_project_id")

	// List all shares for the given zone.
	shares, err := listZoneShares(ctx, client, zoneID, ownerProjectID)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error listing shares for DNS zone %s: %s", zoneID, err))
	}

	// Flatten the list of shares using the helper function.
	results := flattenZoneShares(shares, targetFilter)

	// Set the "shares" attribute with the flattened list.
	if err := d.Set("shares", results); err != nil {
		return diag.FromErr(fmt.Errorf("error setting shares: %s", err))
	}

	// If target_project_id is not set, use the first share's target_project_id (if available)
	if d.Get("target_project_id") == "" && len(results) > 0 {
		if err := d.Set("target_project_id", results[0]["project_id"]); err != nil {
			return diag.FromErr(err)
		}
	}

	// Use the original zoneID as the data source ID.
	d.SetId(zoneID)
	return nil
}
