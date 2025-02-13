package openstack

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDNSZoneShareV2() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceDNSZoneShareV2Read,

		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the DNS zone",
			},
			"project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The owner project ID required to authorize sharing",
			},
			"shares": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"share_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the share",
						},
						"project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The project ID of the share",
						},
					},
				},
			},
		},
	}
}

func dataSourceDNSZoneShareV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.DNSV2Client(context.Background(), "")
	if err != nil {
		return fmt.Errorf("error creating OpenStack DNS client: %s", err)
	}

	zoneID := d.Get("zone_id").(string)

	// Get the sudo project ID if available
	var sudoProjectID string
	if v, ok := d.GetOk("project_id"); ok {
		sudoProjectID = v.(string)
	}

	// Fetch shared zones
	shares, err := listZoneShares(client, zoneID, sudoProjectID)
	if err != nil {
		return fmt.Errorf("error retrieving shared zones for DNS zone %s: %s", zoneID, err)
	}

	// Ensure we are parsing the correct JSON field
	var results []map[string]interface{}
	for _, share := range shares {
		results = append(results, map[string]interface{}{
			"share_id":   share.ID,
			"project_id": share.TargetProjectID,
		})
	}

	// 🔥 DEBUGGING: Print to confirm parsing is correct
	fmt.Printf("Retrieved %d shared zones for %s: %+v\n", len(results), zoneID, results)

	// Set the parsed data in Terraform
	if err := d.Set("shares", results); err != nil {
		return fmt.Errorf("error setting shared zones for DNS zone %s: %s", zoneID, err)
	}

	d.SetId(zoneID)
	return nil
}
